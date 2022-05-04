package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/sshkey"
	"strconv"
)

func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSHKeyCreate,
		Read:   resourceSSHKeyRead,
		// Update: resourceSSHKeyUpdate,
		Delete: resourceSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Description: "The name that can be set by customer.",
				Required:    true,
				ForceNew:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "The public part of the SSH key.",
				Required:    true,
				ForceNew:    true,
			},
			"md5_fingerprint": {
				Type:        schema.TypeString,
				Description: "SSH key fingerprint.",
				Computed:    true,
			},
			"creation_date": {
				Type:        schema.TypeString,
				Description: "Creation date of the SSH key.",
				Computed:    true,
			},
		},
	}
}

func resourceSSHKeyCreate(d *schema.ResourceData, m interface{}) error {
	key := d.Get("key").(string)
	description := d.Get("description").(string)

	client := m.(repository.Client)
	repository := sshkey.Repository{Client: client}

	err := repository.Add(key, description)
	if err != nil {
		return fmt.Errorf("failed to add SSH key %q: %s", key, err)
	}

	sshKeys, err := repository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get SSH keys for determining recently added key id: %s", err)
	}

	for _, sshKey := range sshKeys {
		if sshKey.Key == key {
			d.SetId(strconv.FormatInt(sshKey.ID, 10))
			break
		}
	}

	return resourceSSHKeyRead(d, m)
}

func resourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse sshkey id %s: %s", d.Id(), err)
	}

	client := m.(repository.Client)
	repository := sshkey.Repository{Client: client}

	v, err := repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to lookup SSH key %q: %s", id, err)
	}

	// Description returned by TransIP API == user defined id.
	d.Set("description", v.Description)
	d.Set("key", v.Key)
	d.Set("md5_fingerprint", v.MD5Fingerprint)
	d.Set("creation_date", v.CreationDate)

	return nil
}

func resourceSSHKeyDelete(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse sshkey id %s: %s", d.Id(), err)
	}

	client := m.(repository.Client)
	repository := sshkey.Repository{Client: client}

	err = repository.Remove(id)
	if err != nil {
		return fmt.Errorf("failed to remove SSH key %q: %s", id, err)
	}

	return nil
}
