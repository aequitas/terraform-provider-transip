package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip"
	"github.com/transip/gotransip/vps"
)

func dataSourceVps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"description": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"operating_system": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"disk_size": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"memory_size": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"processors": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"ipv4_address": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"ipv6_address": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"mac_address": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"is_blocked": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"is_customer_locked": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"availability_zone": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func dataSourceVpsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(gotransip.Client)
	name := d.Get("name").(string)

	v, err := vps.GetVps(client, name)
	if err != nil {
		return fmt.Errorf("failed to lookup vps %q: %s", name, err)
	}

	d.SetId(v.Name)

	d.Set("description", v.Description)
	d.Set("operating_system", v.OperatingSystem)
	d.Set("disk_size", v.DiskSize)
	d.Set("memory_size", v.MemorySize)
	d.Set("processors", v.Processors)
	d.Set("status", v.Status)
	d.Set("ipv4_address", v.IPv4Address.String())
	d.Set("ipv6_address", v.IPv6Address.String())
	d.Set("mac_address", v.MACAddress)
	d.Set("is_blocked", v.IsBlocked)
	d.Set("is_customer_locked", v.IsCustomerLocked)
	d.Set("availability_zone", v.AvailabilityZone)

	return nil
}
