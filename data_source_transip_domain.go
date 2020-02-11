package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip/v5"
	"github.com/transip/gotransip/v5/domain"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDomainRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"nameservers": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"is_locked": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			// "whoiscontact": {
			// 	Computed: true,
			// 	Type:     schema.TypeList,
			// 	Elem: &schema.Schema{
			// 		Type: schema.TypeMap,
			// 	},
			// },
			// "registration_data": {
			// 	Computed: true,
			// 	Type:     schema.TypeString,
			// },
			// "renewal_date": {
			// 	Computed: true,
			// 	Type:     schema.TypeString,
			// },
		},
	}
}

func dataSourceDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(gotransip.Client)
	name := d.Get("name").(string)

	i, err := domain.GetInfo(client, name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %q: %s", name, err)
	}

	var nameservers []map[string]string
	for _, n := range i.Nameservers {
		nameservers = append(nameservers, map[string]string{
			"hostname":     n.Hostname,
			"ipv4_address": n.IPv4Address.String(),
			"ipv6_address": n.IPv6Address.String(),
		})
	}

	d.SetId(i.Name)
	d.Set("is_locked", i.IsLocked)
	d.Set("nameservers", nameservers)

	return nil
}
