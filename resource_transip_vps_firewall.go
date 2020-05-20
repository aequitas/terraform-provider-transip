package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
)

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
				Optional:    true,
				ForceNew:    true,
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
				Required:    true,
				ForceNew:    true,
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
					ValidateFunc: validation.CIDRNetwork(0, 32),
				},
				Description: "Whitelisted IP’s or ranges that are allowed to connect, empty to allow all",
			},
		},
	}
}

func resourceVpsFirewallRead(d *schema.ResourceData, m interface{}) error {
	name := d.Id()

	// Obtain the firewall for the VPS
	client := m.(repository.Client)
	repository := vps.FirewallRepository{Client: client}
	v, err := repository.GetFirewall(name)
	if err != nil {
		return fmt.Errorf("failed to lookup vps firewall %q: %s", name, err)
	}

	// Convert all inbound API rules to state rules
	inboundRules := make([]interface{}, len(v.RuleSet))
	for i, rule := range v.RuleSet {
		inboundRules[i] = vpsFirewallRuleFlatten(&rule)
	}

	// Load information
	d.SetId(name)
	d.Set("name", name)
	d.Set("is_enabled", v.IsEnabled)
	d.Set("inbound_rule", inboundRules)

	return nil
}

func resourceVpsFirewallCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	// Create the API firewall
	var firewall vps.Firewall
	firewall.IsEnabled = d.Get("is_enabled").(bool)

	// Convert all inbound state rules to API rules
	inboundRules := d.Get("inbound_rule").(*schema.Set)
	for _, rule := range inboundRules.List() {
		r, err := vpsFirewallRuleExpand(rule)
		if err != nil {
			return err
		}
		firewall.RuleSet = append(firewall.RuleSet, *r)
	}

	// Set the firewall for the VPS
	client := m.(repository.Client)
	repository := vps.FirewallRepository{Client: client}
	err := repository.UpdateFirewall(name, firewall)
	if err != nil {
		return fmt.Errorf("failed to update vps firewall %q: %s", name, err)
	}

	d.SetId(name)

	return resourceVpsFirewallRead(d, m)
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
	err := repository.UpdateFirewall(name, firewall)
	if err != nil {
		return fmt.Errorf("failed to delete vps firewall %q: %s", name, err)
	}

	return nil
}
