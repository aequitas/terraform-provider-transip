package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/transip/gotransip/v6/openstack"
	"github.com/transip/gotransip/v6/repository"
)

func resourceOpenstackUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpenstackUserCreate,
		Read:   resourceOpenstackUserRead,
		Update: resourceOpenstackUserUpdate,
		Delete: resourceOpenstackUserDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"projectid": {
				Type:        schema.TypeString,
				Description: "Grant user access to a project",
				Required:    true,
				ForceNew:    false,
			},
			"email": {
				Type:        schema.TypeString,
				Description: "Email address",
				Required:    true,
				ForceNew:    false,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username",
				Required:    true,
				ForceNew:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password",
				Required:    true,
				ForceNew:    false,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description",
				Required:    false,
				Optional:    true,
				ForceNew:    false,
			},
		},
	}
}

func resourceOpenstackUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := openstack.UserRepository{Client: client}

	err := repository.Create(openstack.CreateUserRequest{
		ProjectID:   d.Get("projectid").(string),
		Username:    d.Get("username").(string),
		Password:    d.Get("password").(string),
		Description: d.Get("description").(string),
		Email:       d.Get("email").(string),
	})

	if err != nil {
		return fmt.Errorf("failed to create openstack user %q: %s", d.Get("username").(string), err)
	}

	// Return with retryable error, as the user is not found instantly.
	return resource.Retry(d.Timeout(schema.TimeoutCreate)-15*time.Minute, func() *resource.RetryError {
		users, err := repository.GetAll()
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("failed to list openstack users: %s", err))
		}
		var found bool
		// get ID of openstack project
		log.Printf("[DEBUG] listing all users from openstack %v", users)
		for _, user := range users {
			if user.Username == d.Get("username") {
				log.Printf("[DEBUG] found id %v for user %v\n", user.Username, user.ID)
				d.SetId(user.ID)
				found = true
				break
			}
		}
		if !found {
			return resource.RetryableError(fmt.Errorf("UserId not available yet due to OpenStack delay, retrying"))
		}

		err = resourceOpenstackUserRead(d, m)
		if err != nil {
			return resource.NonRetryableError(err)
		} else {
			return nil
		}

	})
}

func resourceOpenstackUserRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(repository.Client)
	repository := openstack.UserRepository{Client: client}

	i, err := repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get openstack users %q: %s", id, err)
	}

	d.Set("email", i.Email)
	d.Set("username", i.Username)
	d.Set("description", i.Description)

	return nil
}

func resourceOpenstackUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := openstack.UserRepository{Client: client}

	repository.Update(openstack.User{
		ID:          d.Id(),
		Username:    d.Get("username").(string),
		Email:       d.Get("email").(string),
		Description: d.Get("description").(string),
	})

	return resourceOpenstackUserRead(d, m)
}

func resourceOpenstackUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := openstack.UserRepository{Client: client}
	err := repository.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("failed to delete openstack user %q: %s", d.Get("username"), err)
	}

	d.SetId("")
	return nil
}
