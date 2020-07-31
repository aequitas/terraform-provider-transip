package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
)

func retryableVpsFirewallErrorf(err error, format string, a ...interface{}) *resource.RetryError {
	// Format the error
	e := fmt.Errorf(format+": %s", append(a, err)...)

	// Return the retryable error (retry or not)
	if strings.Contains(err.Error(), "has an action running, no modification is allowed") {
		return resource.RetryableError(e)
	} else {
		return resource.NonRetryableError(e)
	}
}

func resourceVpsFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpsFirewallCreate,
		Read:   resourceVpsFirewallRead,
		Delete: resourceVpsFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vps_name": {
				Type:        schema.TypeString,
				Description: "Name of the VPS",
				Required:    true,
				ForceNew:    true,
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Default:     true,
				Description: "Whether the firewall is enabled for this VPS",
				Optional:    true,
				ForceNew:    true,
			},
			"inbound_rule": {
				Type:        schema.TypeSet,
				Description: "Firewall rules",
				Optional:    true,
				ForceNew:    true,
				Elem:        vpsFirewallRuleSchema(),
			},
		},
	}
}

func vpsFirewallRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Description: "The rule name.",
				Required:    true,
				ForceNew:    true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Default:      "tcp",
				Description:  "The protocol `tcp`, `udp` or `tcp_udp`.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "tcp_udp"}, false),
			},
			"port": {
				Type:        schema.TypeString,
				Description: "Network port for this rule",
				Required:    true,
				ForceNew:    true,
			},
			"whitelist": {
				Type:        schema.TypeList,
				Description: "Whitelisted IP's or ranges that are allowed to connect, empty to allow all.",
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.CIDRNetwork(0, 128),
				},
			},
		},
	}
}

func resourceVpsFirewallRead(d *schema.ResourceData, m interface{}) error {
	vpsName := d.Id()

	// Obtain the firewall for the VPS
	client := m.(repository.Client)
	repository := vps.FirewallRepository{Client: client}

	log.Printf("[DEBUG] terraform-provider-transip reading firewall for VPS %s\n", vpsName)
	firewall, err := repository.GetFirewall(vpsName)
	if err != nil {
		return fmt.Errorf("failed to lookup vps firewall %q: %s", vpsName, err)
	}

	// Check if we a firewall exist
	log.Printf("[DEBUG] terraform-provider-transip firewall for VPS %s (enabled: %t) has %d inbound rules\n", vpsName, firewall.IsEnabled, len(firewall.RuleSet))
	if len(firewall.RuleSet) == 0 && !firewall.IsEnabled {
		d.SetId("")
		return nil
	}

	// Load information
	d.SetId(vpsName)
	d.Set("vps_name", vpsName)
	d.Set("is_enabled", firewall.IsEnabled)
	d.Set("inbound_rule", vpsFirewallRulesFlatten(firewall.RuleSet))

	return nil
}

func resourceVpsFirewallCreate(d *schema.ResourceData, m interface{}) error {
	vpsName := d.Get("vps_name").(string)
	inboundRules := d.Get("inbound_rule").(*schema.Set)

	// Create the API firewall
	var firewall vps.Firewall
	firewall.IsEnabled = d.Get("is_enabled").(bool)
	firewall.RuleSet, _ = vpsFirewallRulesExpand(inboundRules.List())

	// Set the firewall for the VPS
	client := m.(repository.Client)
	repository := vps.FirewallRepository{Client: client}

	// Try the update
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		log.Printf("[DEBUG] terraform-provider-transip updating firewall for VPS %s (%v)\n", vpsName, firewall.RuleSet)
		err := repository.UpdateFirewall(vpsName, firewall)
		if err != nil {
			return retryableVpsFirewallErrorf(err, "failed to update firewall for VPS %q", vpsName)
		}
		d.SetId(vpsName)

		return resource.NonRetryableError(resourceVpsFirewallRead(d, m))
	})
}

func resourceVpsFirewallDelete(d *schema.ResourceData, m interface{}) error {
	vpsName := d.Id()

	// Create an empty firewall which is also disabled
	var firewall vps.Firewall
	firewall.IsEnabled = false
	firewall.RuleSet = make([]vps.FirewallRule, 0)

	// Set the empty firewall for the VPS
	client := m.(repository.Client)
	repository := vps.FirewallRepository{Client: client}

	// Try the delete
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		log.Printf("[DEBUG] terraform-provider-transip removing firewall for VPS %s\n", vpsName)
		err := repository.UpdateFirewall(vpsName, firewall)
		if err != nil {
			return retryableVpsFirewallErrorf(err, "failed to delete firewall for VPS %q", vpsName)
		}
		return nil
	})
}
