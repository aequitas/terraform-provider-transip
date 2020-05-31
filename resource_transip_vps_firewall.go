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

func retryableErrorf(err error, format string, a ...interface{}) *resource.RetryError {
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the vps",
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     true,
				Description: "Whether the firewall is enabled for this VPS",
			},
			"inbound_rule": {
				Type:        schema.TypeSet,
				ForceNew:    true,
				Optional:    true,
				Description: "Ruleset of the VPS",
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
				ForceNew:    true,
				Required:    true,
				Description: "Description of the rule",
			},
			"protocol": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "tcp",
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "tcp_udp"}, false),
				Description:  "Protocol for this rule (tcp, udp or tcp_udp)",
			},
			"port": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Network port for this rule",
			},
			"whitelist": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.CIDRNetwork(0, 128),
				},
				Description: "Whitelisted IP's or ranges that are allowed to connect, empty to allow all",
			},
		},
	}
}

func resourceVpsFirewallRead(d *schema.ResourceData, m interface{}) error {
	name := d.Id()

	// Obtain the firewall for the VPS
	client := m.(repository.Client)
	repository := vps.FirewallRepository{Client: client}

	log.Printf("[DEBUG] terraform-provider-transip reading VPS firewall %s\n", name)
	firewall, err := repository.GetFirewall(name)
	if err != nil {
		return fmt.Errorf("failed to lookup vps firewall %q: %s", name, err)
	}

	// Check if we a firewall exist
	log.Printf("[DEBUG] terraform-provider-transip VPS firewall %s (enabled: %t) has %d inbound rules\n", name, firewall.IsEnabled, len(firewall.RuleSet))
	if len(firewall.RuleSet) == 0 && !firewall.IsEnabled {
		d.SetId("")
		return nil
	}

	// Load information
	d.SetId(name)
	d.Set("name", name)
	d.Set("is_enabled", firewall.IsEnabled)
	d.Set("inbound_rule", vpsFirewallRulesFlatten(firewall.RuleSet))

	return nil
}

func resourceVpsFirewallCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
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
		log.Printf("[DEBUG] terraform-provider-transip updating VPS firewall %s (%v)\n", name, firewall.RuleSet)
		err := repository.UpdateFirewall(name, firewall)
		if err != nil {
			return retryableErrorf(err, "failed to update vps firewall %q", name)
		}
		d.SetId(name)

		return resource.NonRetryableError(resourceVpsFirewallRead(d, m))
	})
}

func resourceVpsFirewallDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Id()

	// Create an empty firewall which is also disabled
	var firewall vps.Firewall
	firewall.IsEnabled = false
	firewall.RuleSet = make([]vps.FirewallRule, 0)

	// Set the empty firewall for the VPS
	client := m.(repository.Client)
	repository := vps.FirewallRepository{Client: client}

	// Try the delete
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		log.Printf("[DEBUG] terraform-provider-transip removing VPS firewall %s\n", name)
		err := repository.UpdateFirewall(name, firewall)
		if err != nil {
			return retryableErrorf(err, "failed to delete vps firewall %q", name)
		}
		return nil
	})
}
