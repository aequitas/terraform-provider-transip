package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Required:    true,
				Description: "Name of the Private Network.",
				Type:        schema.TypeString,
				ForceNew:    true,
			},
			"vps_id": {
				Required:    true,
				Description: "VPN name the Private Network is attached to.",
				Type:        schema.TypeString,
				ForceNew:    true,
			},
		},
	}
}

func resourcePrivateNetworkAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	errorStrings := []string{
		"has an action running, no modification is allowed",
		"is already locked to another action",
		"EOF"}

	privateNetworkID := d.Get("private_network_id").(string)
	vpsID := d.Get("vps_id").(string)

	client := m.(repository.Client)
	repository := vps.PrivateNetworkRepository{Client: client}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

		err := repository.AttachVps(vpsID, privateNetworkID)
		if err != nil {
			for _, errorString := range errorStrings {
				if strings.Contains(err.Error(), errorString) {
					return resource.RetryableError(fmt.Errorf("failed to attach VPS %s to private network %s, VPS busy: %s; retrying", vpsID, privateNetworkID, err))
				}
			}
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
			found = true
			break
		}
	}
	if !found {
		d.SetId("")
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
