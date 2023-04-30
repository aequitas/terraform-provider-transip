package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
)

func dataSourceVps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The unique VPS name.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The name that can be set by customer.",
				Computed:    true,
			},
			"product_name": {
				Type:        schema.TypeString,
				Description: "The product name.",
				Computed:    true,
			},
			"operating_system": {
				Type:        schema.TypeString,
				Description: "The VPS OperatingSystem.",
				Computed:    true,
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Description: "The VPS disk size in kB.",
				Computed:    true,
			},
			"memory_size": {
				Type:        schema.TypeInt,
				Description: "The VPS memory size in kB.",
				Computed:    true,
			},
			"cpus": {
				Type:        schema.TypeInt,
				Description: "The VPS cpu count.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "The VPS status, either 'created', 'installing', 'running', 'stopped' or 'paused'.",
				Computed:    true,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Description: "The VPS main ipAddress.",
				Computed:    true,
			},
			"mac_address": {
				Type:        schema.TypeString,
				Description: "The VPS macaddress.",
				Computed:    true,
			},
			"is_locked": {
				Type:        schema.TypeBool,
				Description: "Whether or not another process is already doing stuff with this VPS.",
				Computed:    true,
			},
			"is_blocked": {
				Type:        schema.TypeBool,
				Description: "If the VPS is administratively blocked.",
				Computed:    true,
			},
			"is_customer_locked": {
				Type:        schema.TypeBool,
				Description: "If this VPS is locked by the customer.",
				Computed:    true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "The name of the availability zone the VPS is in.",
				Computed:    true,
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "The custom tags added to this VPS.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ipv4_addresses": {
				Type:        schema.TypeList,
				Description: "All IPV4 addresses associated with this VPS.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ipv6_address": {
				Type:        schema.TypeList,
				Description: "All IPV6 addresses associated with this VPS.",
				Computed:    true,
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

	ipAddresses, err := repository.GetIPAddresses(name)
	if err != nil {
		return fmt.Errorf("failed to lookup vps %q: %s", name, err)
	}

	var ipv4Addresses []string
	var ipv6Addresses []string
	for _, address := range ipAddresses {
		if len(address.Address) == 4 {
			ipv4Addresses = append(ipv4Addresses, address.Address.String())
		}
		if len(address.Address) == 16 {
			ipv6Addresses = append(ipv6Addresses, address.Address.String())
		}
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
	d.Set("ipv4_addresses", ipv4Addresses)
	d.Set("ipv6_addresses", ipv6Addresses)

	return nil
}
