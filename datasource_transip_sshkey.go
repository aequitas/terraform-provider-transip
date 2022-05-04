package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/sshkey"
	"strconv"
)

func datasourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: datasourceSSHKeyRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Description:  "The ID of the SSH key.",
				Optional:     true,
				AtLeastOneOf: []string{"id", "description", "key", "md5_fingerprint"},
			},
			"description": {
				Type:         schema.TypeString,
				Description:  "The name that can be set by customer.",
				Optional:     true,
				AtLeastOneOf: []string{"id", "description", "key", "md5_fingerprint"},
			},
			"key": {
				Type:         schema.TypeString,
				Description:  "The public part of the SSH key.",
				Optional:     true,
				AtLeastOneOf: []string{"id", "description", "key", "md5_fingerprint"},
			},
			"md5_fingerprint": {
				Type:         schema.TypeString,
				Description:  "SSH key fingerprint.",
				Optional:     true,
				AtLeastOneOf: []string{"id", "description", "key", "md5_fingerprint"},
			},
		},
	}
}

func datasourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := sshkey.Repository{Client: client}

	var id int64
	var err error
	if d.Id() != "" {
		id, err = strconv.ParseInt(d.Id(), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse sshkey id %s: %s", d.Id(), err)
		}
	} else {

		key := d.Get("key").(string)
		description := d.Get("description").(string)
		md5_fingerprint := d.Get("md5_fingerprint").(string)

		sshKeys, err := repository.GetAll()
		if err != nil {
			return fmt.Errorf("failed to get SSH keys: %s", err)
		}

		for _, sshKey := range sshKeys {
			if sshKey.Key == key || sshKey.Description == description || sshKey.MD5Fingerprint == md5_fingerprint {
				id = sshKey.ID
				break
			}
		}
	}

	d.SetId(strconv.FormatInt(id, 10))

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
