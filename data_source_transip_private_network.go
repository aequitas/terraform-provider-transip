package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
)

func dataSourcePrivateNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePrivateNetworkRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"description": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"is_blocked": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"is_locked": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"vps_names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourcePrivateNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := vps.PrivateNetworkRepository{Client: client}
	name := d.Get("name").(string)
	// getPrivateNetworkID(d, m)

	p, err := repository.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup private network %q: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] terraform-provider-transip success: private network %s found, getting details.\n", name)
	var vpsNames []string
	for _, vpsName := range p.VpsNames {
		log.Printf("[DEBUG] terraform-provider-transip success: private network %s has vps %s \n", name, vpsName)
		vpsNames = append(vpsNames, vpsName)
	}

	d.SetId(name)

	d.Set("name", name)
	d.Set("description", p.Description)
	log.Printf("[DEBUG] terraform-provider-transip success: private network %s has description %s \n", name, p.Description)
	d.Set("is_blocked", p.IsBlocked)
	d.Set("is_locked", p.IsLocked)
	d.Set("vps_names", vpsNames)

	return nil
}

// func getPrivateNetworkID(d *schema.ResourceData, m interface{}) error {
// 	description := d.Get("description").(string)
// 	client := m.(repository.Client)
// 	repository := vps.PrivateNetworkRepository{Client: client}

// 	all, err := repository.GetAll()
// 	if err != nil {
// 		return fmt.Errorf("failed to get all private networks: %s", err)
// 	}

// 	found := false
// 	for _, privateNetwork := range all {
// 		if privateNetwork.Description == description {
// 			d.SetId(privateNetwork.Name)
// 			found = true
// 			return nil
// 		}
// 	}
// 	if !found {
// 		return (fmt.Errorf("Private network with description %s not found", description))
// 	}
// 	return nil
// }
