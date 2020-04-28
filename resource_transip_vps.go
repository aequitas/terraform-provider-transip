package main

import (
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/product"
	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
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
			"description": {
				Optional: true,
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
			"availability_zone": {
				Optional: true,
				Default:  "ams0",
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"install_text": {
				Optional: true,
				Type:     schema.TypeString,
				Default:  "",
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
			"is_blocked": {
				Computed: true,
				Type:     schema.TypeBool,
			},
			"is_customer_locked": {
				Computed: true,
				Type:     schema.TypeBool,
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

func resourceVpsCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	productName := d.Get("product_name").(string)
	operatingSystem := d.Get("operating_system").(string)
	availabilityZone := d.Get("availability_zone").(string)
	addons := []string{}
	InstallText := d.Get("install_text").(string)

	client := m.(repository.Client)
	repository := vps.Repository{Client: client}

	productRepository := product.Repository{Client: client}
	availableProducts, err := productRepository.GetAll()
	if err != nil {
		return fmt.Errorf("Failed to get available products: %s", err)
	}
	validProduct := false
	for _, product := range availableProducts.Vps {
		if product.Name == productName {
			validProduct = true
		}
	}
	if !validProduct {
		return fmt.Errorf("Product %s is invalid. Valid product names are: %v", productName, availableProducts)
	}

	base64InstallText := base64.StdEncoding.EncodeToString([]byte(InstallText))

	vpsOrder := vps.Order{
		ProductName:       productName,
		OperatingSystem:   operatingSystem,
		AvailabilityZone:  availabilityZone,
		Hostname:          name,
		Addons:            addons,
		Base64InstallText: base64InstallText,
	}

	err = repository.Order(vpsOrder)
	if err != nil {
		return fmt.Errorf("failed to order VPS %s: %s", name, err)

	}

	d.SetId(name)

	d.Set("install_text", InstallText)

	return resourceVpsRead(d, m)
}

func resourceVpsRead(d *schema.ResourceData, m interface{}) error {
	name := d.Id()

	client := m.(repository.Client)
	repository := vps.Repository{Client: client}
	v, err := repository.GetByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup vps %q: %s", name, err)
	}

	d.SetId(name)

	d.Set("name", name)
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

func resourceVpsDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	client := m.(repository.Client)
	repository := vps.Repository{Client: client}

	err := repository.Cancel(name, gotransip.CancellationTimeImmediately)
	if err != nil {
		return fmt.Errorf("failed to cancel VPS %q: %s", name, err)
	}

	return nil
}
