package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
)

func dataSourcePrivateNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePrivateNetworkRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique private network name",
				ForceNew:    true,
			},
			"description": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The custom name that can be set by customer.",
			},
			"is_blocked": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "If the Private Network is administratively blocked.",
			},
			"is_locked": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "When locked, another process is already working with this private network.",
			},
			"vps_names": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The VPSes in this private network.",
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

	p, err := repository.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup private network %q: %s", name, err)
	}

	var vpsNames []string
	for _, vpsName := range p.VpsNames {
		vpsNames = append(vpsNames, vpsName)
	}

	d.SetId(name)

	d.Set("name", name)
	d.Set("description", p.Description)
	d.Set("is_blocked", p.IsBlocked)
	d.Set("is_locked", p.IsLocked)
	d.Set("vps_names", vpsNames)

	return nil
}
