package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/transip/gotransip/v6/domain"
	"github.com/transip/gotransip/v6/repository"
)

func resourceDomainDNSSec() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainDNSSecUpdate,
		Read:   resourceDomainDNSSecRead,
		Update: resourceDomainDNSSecUpdate,
		Delete: resourceDomainDNSSecDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDomainDNSSecImport,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Description: "The domain, including the tld",
				Required:    true,
			},
			"dnssec": {
				Type:        schema.TypeList,
				Description: "List of dnssec entries associated with domain",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_tag": {
							Type:        schema.TypeInt,
							Description: "A 5-digit key of the Zonesigner",
							Required:    true,
						},
						"flags": {
							Type:        schema.TypeInt,
							Description: "The signing key number, either 256 (Zone Signing Key) or 257 (Key Signing Key)",
							Required:    true,
						},
						"algorithm": {
							Type:        schema.TypeInt,
							Description: "The algorithm type that is used, see: https://www.transip.nl/vragen/461-domeinnaam-nameservers-gebruikt-beveiligen-dnssec/ for the possible options.",
							Required:    true,
						},
						"public_key": {
							Type:        schema.TypeString,
							Description: "The public key",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceDomainDNSSecImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	d.Set("domain", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceDomainDNSSecRead(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	domain := d.Get("domain").(string)
	entries, err := repository.GetDNSSecEntries(domain)
	if err != nil {
		return fmt.Errorf("failed to get dnssec entries of domain %q: %s", domain, err)
	}
	err = d.Set("dnssec", entriesToMaps(entries))
	if err != nil {
		return fmt.Errorf("failed to parse dnssec entries of domain %q: %s", domain, err)
	}

	d.SetId(domain)
	return nil
}

func resourceDomainDNSSecUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	domain := d.Get("domain").(string)
	entries := interfacesToDNSSec(d.Get("dnssec").([]interface{}))
	err := repository.ReplaceDNSSecEntries(domain, entries)
	if err != nil {
		return fmt.Errorf("failed to update dnssec entries of domain %q: %s", domain, err)
	}

	d.SetId(domain)
	return nil
}

func resourceDomainDNSSecDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	entries := make([]domain.DNSSecEntry, 0)
	domain := d.Get("domain").(string)
	err := repository.ReplaceDNSSecEntries(domain, entries)
	if err != nil {
		return fmt.Errorf("failed to reset dnssec of domain %q: %s", domain, err)
	}

	d.SetId("")
	return nil
}

func entriesToMaps(entries []domain.DNSSecEntry) []map[string]interface{} {
	maps := make([]map[string]interface{}, len(entries))
	for i, v := range entries {
		maps[i] = make(map[string]interface{})
		maps[i]["key_tag"] = v.KeyTag
		maps[i]["flags"] = v.Flags
		maps[i]["algorithm"] = v.Algorithm
		maps[i]["public_key"] = v.PublicKey
	}
	return maps
}

func interfacesToDNSSec(interfaces []interface{}) []domain.DNSSecEntry {
	entries := make([]domain.DNSSecEntry, len(interfaces))
	for i, v := range interfaces {
		map_ := v.(map[string]interface{})
		entries[i] = domain.DNSSecEntry{
			KeyTag:    map_["key_tag"].(int),
			Flags:     map_["flags"].(int),
			Algorithm: map_["algorithm"].(int),
			PublicKey: map_["public_key"].(string),
		}
	}
	return entries
}
