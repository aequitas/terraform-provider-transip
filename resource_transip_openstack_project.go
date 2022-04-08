package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/transip/gotransip/v6/openstack"
	"github.com/transip/gotransip/v6/repository"
)

func resourceOpenstackProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpenstackProjectCreate,
		Read:   resourceOpenstackProjectRead,
		Update: resourceOpenstackprojectUpdate,
		Delete: resourceOpenstackProjectDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Project name",
				Required:    true,
				ForceNew:    false,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Describes this project",
				Required:    false,
				ForceNew:    false,
			},
			"locked": {
				Type:        schema.TypeBool,
				Description: "Set to true when an ongoing process blocks the project from being modified",
				Required:    false,
				ForceNew:    false,
			},
			"blocked": {
				Type:        schema.TypeBool,
				Description: "Set to true when a project has been administratively blocked",
				Required:    false,
				ForceNew:    false,
			},
		},
	}
}

func resourceOpenstackProjectCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	client := m.(repository.Client)
	repository := openstack.ProjectRepository{Client: client}

	err := repository.Create(openstack.Project{
		Name:        name,
		Description: description,
	})

	if err != nil {
		return fmt.Errorf("failed to create openstack project %q: %s", name, err)
	}

	projects, err := repository.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get openstack projects %q: %s", name, err)
	}

	// get ID of openstack project
	for _, project := range projects {
		if project.Name == name {
			d.SetId(project.ID)
			break
		}
	}
	return resourceOpenstackProjectRead(d, m)
}

func resourceOpenstackProjectRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(repository.Client)
	repository := openstack.ProjectRepository{Client: client}

	i, err := repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get openstack project %q: %s", id, err)
	}

	d.SetId(i.ID)

	d.Set("name", i.Name)
	d.Set("description", i.Description)
	d.Set("locked", i.IsLocked)
	d.Set("blocked", i.IsBlocked)

	return nil
}

func resourceOpenstackprojectUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(repository.Client)
	repository := openstack.ProjectRepository{Client: client}

	repository.Update(openstack.Project{
		ID:          d.Id(),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		IsLocked:    d.Get("islocked").(bool),
		IsBlocked:   d.Get("isblocked").(bool),
	})

	return resourceOpenstackProjectRead(d, m)
}

func resourceOpenstackProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repository := openstack.ProjectRepository{Client: client}
	err := repository.Cancel(d.Id())
	if err != nil {
		return fmt.Errorf("failed to delete openstack project %q: %s", d.Get("name"), err)
	}

	d.SetId("")
	return nil
}
