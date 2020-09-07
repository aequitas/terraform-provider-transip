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
		Update: resourceVpsUpdate,
		Delete: resourceVpsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Computed:    true,
				Description: "The unique VPS name.",
				Type:        schema.TypeString,
			},
			"description": {
				Optional:    true,
				Description: "The name that can be set by customer.",
				Type:        schema.TypeString,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if len(val.(string)) > 32 {
						errs = append(errs, fmt.Errorf("%q must be less than 33 characters", key))
					}
					return
				},
			},
			"product_name": {
				Required:    true,
				Description: "The product name.",
				Type:        schema.TypeString,
				ForceNew:    true,
			},
			"operating_system": {
				Required:    true,
				Description: "The VPS OperatingSystem.",
				Type:        schema.TypeString,
				ForceNew:    true,
			},
			"disk_size": {
				Computed:    true,
				Description: "The VPS disk size in kB.",
				Type:        schema.TypeInt,
			},
			"memory_size": {
				Computed:    true,
				Description: "The VPS memory size in kB.",
				Type:        schema.TypeInt,
			},
			"cpus": {
				Computed:    true,
				Description: "The VPS cpu count.",
				Type:        schema.TypeInt,
			},
			"status": {
				Computed:    true,
				Description: "The VPS status, either 'created', 'installing', 'running', 'stopped' or 'paused'.",
				Type:        schema.TypeString,
			},
			"ip_address": {
				Computed:    true,
				Description: "The VPS main ipAddress.",
				Type:        schema.TypeString,
			},
			"mac_address": {
				Computed:    true,
				Description: "The VPS macaddress.",
				Type:        schema.TypeString,
			},
			"is_blocked": {
				Computed:    true,
				Description: "If the VPS is administratively blocked.",
				Type:        schema.TypeBool,
			},
			"is_customer_locked": {
				Optional:    true,
				Description: "If this VPS is locked by the customer.",
				Type:        schema.TypeBool,
				Default:     false,
				ForceNew:    false,
			},
			"availability_zone": {
				Optional:    true,
				Description: "The name of the availability zone the VPS is in.",
				Default:     "ams0",
				Type:        schema.TypeString,
				ForceNew:    true,
			},
			"tags": {
				Computed:    true,
				Description: "The custom tags added to this VPS.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"install_text": {
				Optional:    true,
				Description: "Base64 encoded preseed / kickstart / cloudinit instructions, when installing unattended.",
				Type:        schema.TypeString,
				Default:     "",
				ForceNew:    true,
			},
			"ipv4_addresses": {
				Computed:    true,
				Description: "All IPV4 addresses associated with this VPS.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ipv6_addresses": {
				Computed:    true,
				Description: "All IPV6 addresses associated with this VPS.",
				Type:        schema.TypeList,
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

	// query with fake vps name "x"
	availableOss, err := repository.GetOperatingSystems("x")
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
	}

	err = repository.Order(vpsOrder)
	if err != nil {
		return fmt.Errorf("failed to order VPS %s: %s", name, err)

	}

	d.Set("install_text", installText)

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

	// Transip API requires OS Name for creating VPS but return OS Description on a VPS query.
	// So it needs to be translated to avoid Terraform detecting changes.
	operatingSystems, err := repository.GetOperatingSystems("x")
	if err != nil {
		return fmt.Errorf("Failed to get available operating systems: %s", err)
	}
	for _, os := range operatingSystems {
		if os.Description == v.OperatingSystem {
			d.Set("operating_system", os.Name)
		}
	}

	return nil
}

func resourceVpsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := vps.Repository{Client: client}
	description := d.Get("description").(string)
	diskSize := d.Get("disk_size").(int)
	memorySize := d.Get("memory_size").(int)

	vps := vps.Vps{
		Name:        d.Id(), // Unique ID provided by TransIP
		Description: description,
		// TransIP API expects int64, while Terraform Schema expects TypeInt.
		DiskSize:         int64(diskSize),
		MemorySize:       int64(memorySize),
		CPUs:             d.Get("cpus").(int),
		IsCustomerLocked: d.Get("is_customer_locked").(bool),
	}

	err := repository.Update(vps)

	if err != nil {
		return fmt.Errorf("failed to update vps %s with id %q: %s", description, d.Id(), err)
	}
	return resourceVpsRead(d, m)
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
