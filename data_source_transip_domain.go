package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip/v6/domain"
	"github.com/transip/gotransip/v6/repository"
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
			"tags": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"authcode": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"cancellation_date": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"cancellation_status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"is_dns_only": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"is_transfer_locked": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"is_whitelabel": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"registration_date": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"renewal_date": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func dataSourceDomainRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	client := m.(repository.Client)
	repository := domain.Repository{Client: client}
	i, err := repository.GetByDomainName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %q: %s", name, err)
	}

	var nameservers []map[string]string
	ns, err := repository.GetNameservers(name)
	if err != nil {
		return fmt.Errorf("failed to lookup nameservers for domain %q: %s", name, err)
	}

	for _, n := range ns {
		nameservers = append(nameservers, map[string]string{
			"hostname":     n.Hostname,
			"ipv4_address": n.IPv4.String(),
			"ipv6_address": n.IPv6.String(),
		})
	}

	d.SetId(i.Name)
	d.Set("nameservers", nameservers)
	d.Set("tags", i.Tags)
	d.Set("authcode", i.AuthCode)
	d.Set("cancellation_date", i.CancellationDate)
	d.Set("cancellation_status", i.CancellationStatus)
	d.Set("is_dns_only", i.IsDNSOnly)
	d.Set("is_transfer_locked", i.IsTransferLocked)
	d.Set("is_whitelabel", i.IsWhitelabel)
	d.Set("registation_date", i.RegistrationDate)
	d.Set("renewal_date", i.RenewalDate)

	return nil
}
