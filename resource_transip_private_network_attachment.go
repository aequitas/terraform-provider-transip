package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
)

func resourcePrivateNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrivateNetworkAttachmentCreate,
		Read:   resourcePrivateNetworkAttachmentRead,
		// Update: resourcePrivateNetworkUpdate,
		Delete: resourcePrivateNetworkAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"private_network_id": {
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"vps_id": {
				Required: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"action": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourcePrivateNetworkAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	privateNetworkID := d.Get("private_network_id").(string)
	vpsID := d.Get("vps_id").(string)

	client := m.(repository.Client)
	repository := vps.PrivateNetworkRepository{Client: client}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		log.Printf("[DEBUG] terraform-provider-transip trying to attach network %s to vps %s \n", privateNetworkID, vpsID)

		err := repository.AttachVps(vpsID, privateNetworkID)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("VPS '%s' has an action running, no modification is allowed", vpsID)) {
				log.Printf("[DEBUG] terraform-provider-transip VPS %s busy, retrying to attach to %s in a bit \n", vpsID, privateNetworkID)
				return resource.RetryableError(fmt.Errorf("failed to attach VPS %s to private network %s, VPS busy: %s", vpsID, privateNetworkID, err))
			}
			log.Printf("[DEBUG] terraform-provider-transip something went wrong while attaching VPS %s to %s \n", vpsID, privateNetworkID)
			return resource.NonRetryableError(fmt.Errorf("failed to attach private network %s to VPS %s: %s", privateNetworkID, vpsID, err))
		}
		return resource.NonRetryableError(resourcePrivateNetworkAttachmentRead(d, m))
	})
}

func resourcePrivateNetworkAttachmentRead(d *schema.ResourceData, m interface{}) error {
	privateNetworkID := d.Get("private_network_id").(string)
	vpsID := d.Get("vps_id").(string)
	client := m.(repository.Client)
	repository := vps.PrivateNetworkRepository{Client: client}

	p, err := repository.GetByName(privateNetworkID)
	if err != nil {
		return fmt.Errorf("failed to lookup private network %q: %s", d.Id(), err)
	}

	found := false
	for _, vpsName := range p.VpsNames {
		if vpsName == vpsID {
			log.Printf("[DEBUG] terraform-provider-transip success: vps %s is the same as %s in private network %s \n", vpsID, vpsName, privateNetworkID)
			found = true
			break
		}
		if !found {
			d.SetId("")
			log.Printf("[DEBUG] terraform-provider-transip vps %s not found in private network %s \n", vpsID, privateNetworkID)
		}
	}
	d.SetId(p.Name)
	return nil
}

func resourcePrivateNetworkAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	privateNetworkID := d.Get("private_network_id").(string)
	vpsID := d.Get("vps_id").(string)

	client := m.(repository.Client)
	repository := vps.PrivateNetworkRepository{Client: client}
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err := repository.DetachVps(vpsID, privateNetworkID)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("VPS '%s' has an action running, no modification is allowed", vpsID)) {
				return resource.RetryableError(fmt.Errorf("retrying to detach private network %s from VPS %s: %s", privateNetworkID, vpsID, err))
			}
		}
		return nil
	})
	return nil
}
