package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/transip/gotransip/v6/domain"
	"github.com/transip/gotransip/v6/repository"
)

func dataSourceDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDomainsRead,
		Schema: map[string]*schema.Schema{
			"domains": {
				Type:        schema.TypeList,
				Description: "List of all domain names in your TransIP account.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceDomainsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := domain.Repository{Client: client}
	domains, err := repository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get all domains: %s", err)
	}
	var domainNames []string
	for i, domain := range domains {
		domainNames = append(domainNames, domain.Name)
		log.Printf(fmt.Sprintf("requesting all domains, %d/%d: %s", i+1, len(domains), domain.Name))
	}

	d.SetId(uuid.New().String())
	d.Set("domains", domainNames)

	return nil
}
