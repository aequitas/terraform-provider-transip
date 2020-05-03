package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
)

func dataSourcePrivateNetwork() *schema.Resource {
	return &schema.Resource{
		Read: resourcePrivateNetworkRead,
		Schema: map[string]*schema.Schema{
			"description": {
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"name": {
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
	return nil
}

func dataSourcePrivateNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := vps.PrivateNetworkRepository{Client: client}

	getPrivateNetworkID(d, m)

	p, err := repository.GetByName(d.Id())
	if err != nil {
		return fmt.Errorf("failed to lookup private network %q: %s", d.Id(), err)
	}

	var vpsNames []string
	for _, vpsName := range p.VpsNames {
		vpsNames = append(vpsNames, vpsName)
	}

	d.Set("name", d.Id())
	d.Set("description", p.Description)
	d.Set("is_blocked", p.IsBlocked)
	d.Set("is_locked", p.IsLocked)
	d.Set("vps_names", vpsNames)

	return nil
}

func getPrivateNetworkID(d *schema.ResourceData, m interface{}) error {
	description := d.Get("description").(string)
	client := m.(repository.Client)
	repository := vps.PrivateNetworkRepository{Client: client}

	all, err := repository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get all private networks: %s", err)
	}

	found := false
	for _, privateNetwork := range all {
		if privateNetwork.Description == description {
			d.SetId(privateNetwork.Name)
			found = true
			return nil
		}
	}
	if !found {
		return (fmt.Errorf("Private network with description %s not found", description))
	}
	return nil
}
