package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/transip/gotransip/v6/openstack"
	"github.com/transip/gotransip/v6/repository"
)

func dataSourceOpenstackUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpenstackUserRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Description: "Username",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Describes this user",
			},
			"email": {
				Type:        schema.TypeString,
				Description: "Email",
			},
		},
	}
}

func dataSourceOpenstackUserRead(d *schema.ResourceData, m interface{}) error {
	username := d.Get("username").(string)

	client := m.(repository.Client)
	repository := openstack.UserRepository{Client: client}

	users, err := repository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get openstack users: %s", err)
	}

	// get ID of openstack user
	for _, user := range users {
		if user.Username == username {
			d.SetId(user.ID)
			break
		}
	}

	i, err := repository.GetByID(d.Id())
	if err != nil {
		return fmt.Errorf("failed to lookup openstack user %q: %s", username, err)
	}

	d.Set("username", i.Username)
	d.Set("description", i.Description)
	d.Set("email", i.Email)

	return nil
}
