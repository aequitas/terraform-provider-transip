package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip/v5"
	"github.com/transip/gotransip/v5/vps"
)

func resourceVps() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpsCreate,
		Read:   resourceVpsRead,
		// Update: resourceVpsUpdate,
		Delete: resourceVpsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"product_name": {
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"operating_system": {
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
				ForceNew: true,
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

func resourceVpsCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	productName := d.Get("product_name").(string)
	operatingSystem := d.Get("operating_system").(string)
	addons := []string{}

	client := m.(gotransip.Client)

	availableProducts, err := vps.GetAvailableProducts(client)
	if err != nil {
		return fmt.Errorf("Failed to get available products: %s", err)
	}
	validProduct := false
	for _, product := range availableProducts {
		if product.Name == productName {
			validProduct = true
		}
	}
	if !validProduct {
		return fmt.Errorf("Product %s is invalid. Valid product names are: %v", productName, availableProducts)
	}

	availableOss, err := vps.GetOperatingSystems(client)
	if err != nil {
		return fmt.Errorf("Failed to get available operating systems: %s", err)
	}
	validOs := false
	for _, os := range availableOss {
		if os.Name == operatingSystem {
			validOs = true
		}
	}
	if !validOs {
		return fmt.Errorf("Product %s is invalid. Valid product names are: %v", productName, availableProducts)
	}

	err = vps.OrderVps(client, productName, addons, operatingSystem, name)
	if err != nil {
		return fmt.Errorf("failed to order VPS %s: %s", name, err)

	}
	d.SetId(name)

	return resourceVpsRead(d, m)
}

func resourceVpsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(gotransip.Client)
	name := d.Get("name").(string)

	v, err := vps.GetVps(client, name)
	if err != nil {
		return fmt.Errorf("failed to lookup VPS %q: %s", name, err)
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

func resourceVpsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(gotransip.Client)
	name := d.Get("name").(string)

	err := vps.CancelVps(client, name, gotransip.CancellationTimeImmediately)
	if err != nil {
		return fmt.Errorf("failed to cancel VPS %q: %s", name, err)
	}

	return nil
}
