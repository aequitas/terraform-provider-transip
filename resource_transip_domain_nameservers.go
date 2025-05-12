package main

import (
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/transip/gotransip/v6/domain"
	"github.com/transip/gotransip/v6/repository"
)

func resourceDomainNameservers() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainNameserversUpdate,
		Read:   resourceDomainNameserversRead,
		Update: resourceDomainNameserversUpdate,
		Delete: resourceDomainNameserversDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDomainNameserversImport,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Description: "The domain, including the tld",
				Required:    true,
			},
			"nameserver": {
				Type:        schema.TypeList,
				Description: "List of nameservers associated with domain",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:        schema.TypeString,
							Description: "The hostname of this nameserver",
							Required:    true,
						},
						"ipv4": {
							Type:        schema.TypeString,
							Description: "Optional ipv4 glue record for this nameserver",
							Optional:    true,
						},
						"ipv6": {
							Type:        schema.TypeString,
							Description: "Optional ipv6 glue record for this nameserver",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceDomainNameserversImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	d.Set("domain", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceDomainNameserversRead(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	domain := d.Get("domain").(string)
	nameservers, err := repository.GetNameservers(domain)
	if err != nil {
		return fmt.Errorf("failed to get nameservers of domain %q: %s", domain, err)
	}
	err = d.Set("nameserver", nameserversToMaps(nameservers))
	if err != nil {
		return fmt.Errorf("failed to parse nameservers of domain %q: %s", domain, err)
	}

	d.SetId(domain)
	return nil
}

func resourceDomainNameserversUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	domain := d.Get("domain").(string)
	nameservers := interfacesToNameservers(d.Get("nameserver").([]interface{}))
	err := repository.UpdateNameservers(domain, nameservers)
	if err != nil {
		return fmt.Errorf("failed to update nameservers of domain %q: %s", domain, err)
	}

	d.SetId(domain)
	return nil
}

func resourceDomainNameserversDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	nameservers := make([]domain.Nameserver, 3)
	nameservers[0] = domain.Nameserver{Hostname: "ns0.transip.net"}
	nameservers[1] = domain.Nameserver{Hostname: "ns1.transip.nl"}
	nameservers[2] = domain.Nameserver{Hostname: "ns2.transip.eu"}

	domain := d.Get("domain").(string)
	err := repository.UpdateNameservers(domain, nameservers)
	if err != nil {
		return fmt.Errorf("failed to reset nameservers of domain %q: %s", domain, err)
	}

	d.SetId("")
	return nil
}

func nameserversToMaps(nameservers []domain.Nameserver) []map[string]interface{} {
	maps := make([]map[string]interface{}, len(nameservers))
	for i, v := range nameservers {
		maps[i] = make(map[string]interface{})
		maps[i]["hostname"] = v.Hostname
		if v.IPv4 != nil {
			maps[i]["ipv4"] = v.IPv4.String()
		}
		if v.IPv4 != nil {
			maps[i]["ipv6"] = v.IPv6.String()
		}
	}
	return maps
}

func interfacesToNameservers(interfaces []interface{}) []domain.Nameserver {
	nameservers := make([]domain.Nameserver, len(interfaces))
	for i, v := range interfaces {
		map_ := v.(map[string]interface{})
		nameservers[i] = domain.Nameserver{
			Hostname: map_["hostname"].(string),
			IPv4:     net.ParseIP(map_["ipv4"].(string)),
			IPv6:     net.ParseIP(map_["ipv6"].(string)),
		}
	}
	return nameservers
}
