package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:        schema.TypeString,
				Description: "The unique VPS name.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The name that can be set by customer.",
				Optional:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if len(val.(string)) > 32 {
						errs = append(errs, fmt.Errorf("%q must be less than 33 characters", key))
					}
					return
				},
			},
			"product_name": {
				Type:        schema.TypeString,
				Description: "The product name.",
				Required:    true,
				ForceNew:    true,
			},
			"operating_system": {
				Type:        schema.TypeString,
				Description: "The VPS OperatingSystem.",
				Required:    true,
				ForceNew:    true,
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
				Default:     "ams0",
				Description: "The name of the availability zone the VPS is in.",
				Optional:    true,
				ForceNew:    true,
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "The custom tags added to this VPS.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"install_text": {
				Type:        schema.TypeString,
				Default:     "",
				Description: "Base64 encoded preseed / kickstart / cloudinit instructions, when installing unattended.",
				Optional:    true,
				ForceNew:    true,
			},
			"install_flavour": {
				Type:        schema.TypeString,
				Default:     "",
				Description: "The flavour of OS installation: 'installer', 'preinstallable' or 'cloudinit'.",
				Optional:    true,
				ForceNew:    true,
			},
			"ipv4_addresses": {
				Type:        schema.TypeList,
				Description: "All IPV4 addresses associated with this VPS.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ipv6_addresses": {
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

func resourceVpsCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	productName := d.Get("product_name").(string)
	operatingSystem := d.Get("operating_system").(string)
	availabilityZone := d.Get("availability_zone").(string)
	description := d.Get("description").(string)
	addons := []string{}
	installText := d.Get("install_text").(string)
	installFlavour := vps.InstallFlavour(d.Get("install_flavour").(string))

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

	// query using FilterOperatingSystems function, which requires the productName and optional addons

	availableOss, err := repository.FilterOperatingSystems(productName, addons)
	if err != nil {
		return fmt.Errorf("Failed to get available operating systems: %s", err)
	}
	validOs := false
	var osNames []string
	for _, os := range availableOss {
		osNames = append(osNames, os.Name)
	}
	for _, osName := range osNames {
		if osName == operatingSystem {
			validOs = true
		}
	}
	if !validOs {
		return fmt.Errorf("Operating system %s is invalid. Valid operating system names are: %v", operatingSystem, osNames)
	}

	base64InstallText := base64.StdEncoding.EncodeToString([]byte(installText))

	// Generate a unique temporary description to assign during VPS creation so we can reference
	// the VPS afterwards and get a Transip generated unique name.
	// Must be no more than 32 characters, hence MD5 hash.
	tempDescription := fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))

	vpsOrder := vps.Order{
		ProductName:       productName,
		OperatingSystem:   operatingSystem,
		AvailabilityZone:  availabilityZone,
		Description:       tempDescription,
		Hostname:          name,
		Addons:            addons,
		Base64InstallText: base64InstallText,
		InstallFlavour:    installFlavour,
	}

	err = repository.Order(vpsOrder)
	if err != nil {
		return fmt.Errorf("failed to order VPS %s: %s", name, err)

	}

	d.Set("install_text", installText)
	d.Set("install_flavour", installFlavour)

	// Wait for VPS to be in a "running" state
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		// The set name in the Terraform resource is not the same as the name used to query details about a VPS.
		// You'll need the unique name Transip generates to get the VPS details.
		all, err := repository.GetAll()
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("failed to get all VPS's: %s", err))
		}
		for _, v := range all {
			if v.Description != tempDescription {
				continue
			}
			if v.Status != vps.VpsStatusRunning {
				return resource.RetryableError(fmt.Errorf("VPS %s, not yet running.", d.Id()))
			}
			// replace temporary description with actual one
			v.Description = description
			err := repository.Update(v)
			if err != nil {
				return resource.RetryableError(fmt.Errorf("Failed to update description for VPS %s.", d.Id()))
			}

			d.SetId(v.Name)
		}
		if d.Id() == "" {
			return resource.RetryableError(fmt.Errorf("Failed to set ID for VPS %s", d.Id()))
		}
		return resource.NonRetryableError(resourceVpsRead(d, m))
	})
}

func resourceVpsRead(d *schema.ResourceData, m interface{}) error {
	name := d.Id()

	client := m.(repository.Client)
	repository := vps.Repository{Client: client}

	v, err := repository.GetByName(d.Id())
	if err != nil {
		return fmt.Errorf("failed to lookup vps %q: %s", name, err)
	}

	ipAddresses, err := repository.GetIPAddresses(d.Id())
	if err != nil {
		return fmt.Errorf("failed to lookup vps %q: %s", name, err)
	}

	var ipv4Addresses []string
	var ipv6Addresses []string
	for _, address := range ipAddresses {
		if address.Address.To4() != nil {
			ipv4Addresses = append(ipv4Addresses, address.Address.String())
		} else {
			ipv6Addresses = append(ipv6Addresses, address.Address.String())
		}
	}

	// Get all Active Addons
	addons, err := repository.GetAddons(d.Id())
	if err != nil {
		return fmt.Errorf("failed to lookup addons for vps %q: %s", name, err)
	}
	activeAddons := make([]string, len(addons.Active))
	for i, addon := range addons.Active {
		activeAddons[i] = addon.Name
	}

	d.Set("name", name)
	// Description returned by TransIP API == user defined name.
	d.Set("description", v.Description)
	d.Set("product_name", v.ProductName)
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
	d.Set("addons", activeAddons)

	operatingSystems, err := repository.FilterOperatingSystems(v.ProductName, []string{})

	if err != nil {
		return fmt.Errorf("Failed to get available operating systems: %s", err)
	}
	for _, os := range operatingSystems {
		if os.Description == v.OperatingSystem || os.Name == v.OperatingSystem {
			d.Set("operating_system", os.Name)
		}
	}

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
