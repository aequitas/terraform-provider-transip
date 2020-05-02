package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
)

func resourcePrivateNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrivateNetworkCreate,
		Read:   resourcePrivateNetworkRead,
		// Update: resourcePrivateNetworkUpdate,
		Delete: resourcePrivateNetworkDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
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
}

func resourcePrivateNetworkCreate(d *schema.ResourceData, m interface{}) error {
	description := d.Get("description").(string)
	client := m.(repository.Client)
	repository := vps.PrivateNetworkRepository{Client: client}

	err := repository.Order(description)
	if err != nil {
		return fmt.Errorf("failed to order private network %s: %s", description, err)

	}
	setPrivateNetworkID(d, m)

	return resourcePrivateNetworkRead(d, m)
}

func resourcePrivateNetworkRead(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := vps.PrivateNetworkRepository{Client: client}

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

func resourcePrivateNetworkDelete(d *schema.ResourceData, m interface{}) error {
	description := d.Get("description")
	client := m.(repository.Client)
	repository := vps.PrivateNetworkRepository{Client: client}

	err := repository.Cancel(d.Id(), gotransip.CancellationTimeImmediately)
	if err != nil {
		return fmt.Errorf("failed to cancel private network %s with id %q: %s", description, d.Id(), err)
	}

	return nil
}

func setPrivateNetworkID(d *schema.ResourceData, m interface{}) error {
	description := d.Get("description").(string)
	client := m.(repository.Client)
	repository := vps.PrivateNetworkRepository{Client: client}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		log.Printf("[DEBUG] terraform-provider-transip trying to get id for private network %s \n", description)

		// The set description in the Terraform resource is not the same as the name used to query details about a private network.
		// You'll need the unique name Transip generates to get the private network details.
		all, err := repository.GetAll()
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("failed to get all private networks: %s", err))
		}
		for _, privateNetwork := range all {
			if privateNetwork.Description == description {
				d.SetId(privateNetwork.Name)
				log.Printf("[DEBUG] terraform-provider-transip id found for private network %s:%s \n", description, d.Id())
			}
		}
		if d.Id() == "" {
			return resource.RetryableError(fmt.Errorf("Failed to set ID for private network %s", description))
		}
		return resource.NonRetryableError(resourcePrivateNetworkRead(d, m))
	})
}
