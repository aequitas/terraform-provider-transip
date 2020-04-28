package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
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
			"product_name": {
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
			"cpus": {
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
			"is_locked": {
				Computed: true,
				Type:     schema.TypeBool,
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
			"tags": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceVpsRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	client := m.(repository.Client)
	repository := vps.Repository{Client: client}
	v, err := repository.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup vps %q: %s", name, err)
	}

	d.SetId(v.Name)

	d.Set("description", v.Description)
	d.Set("product_name", v.ProductName)
	d.Set("operating_system", v.OperatingSystem)
	d.Set("disk_size", v.DiskSize)
	d.Set("memory_size", v.MemorySize)
	d.Set("cpus", v.CPUs)
	d.Set("status", v.Status)
	d.Set("ip_address", v.IPAddress)
	d.Set("mac_address", v.MacAddress)
	d.Set("is_locked", v.IsLocked)
	d.Set("is_blocked", v.IsBlocked)
	d.Set("is_customer_locked", v.IsCustomerLocked)
	d.Set("availability_zone", v.AvailabilityZone)
	d.Set("tags", v.Tags)

	return nil
}
