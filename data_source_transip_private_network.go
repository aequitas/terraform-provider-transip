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
				Type:        schema.TypeString,
				Description: "The unique private network name",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The custom name that can be set by customer.",
				Computed:    true,
			},
			"is_blocked": {
				Type:        schema.TypeString,
				Description: "If the Private Network is administratively blocked.",
				Computed:    true,
			},
			"is_locked": {
				Type:        schema.TypeString,
				Description: "When locked, another process is already working with this privatenetwork.",
				Computed:    true,
			},
			"vps_names": {
				Type:        schema.TypeList,
				Description: "The VPSes in this private network.",
				Computed:    true,
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
