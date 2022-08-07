package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/transip/gotransip/v6/openstack"
	"github.com/transip/gotransip/v6/repository"
)

func dataSourceOpenstackProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpenstackProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Project name",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Describes this project",
				Optional:    true,
			},
			"locked": {
				Type:        schema.TypeBool,
				Description: "Set to true when an ongoing process blocks the project from being modified",
				Optional:    true,
			},
			"blocked": {
				Type:        schema.TypeBool,
				Description: "Set to true when a project has been administratively blocked",
				Optional:    true,
			},
		},
	}
}

func dataSourceOpenstackProjectRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	client := m.(repository.Client)
	repository := openstack.ProjectRepository{Client: client}

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

	i, err := repository.GetByID(d.Id())
	if err != nil {
		return fmt.Errorf("failed to lookup openstack project %q: %s", name, err)
	}

	d.Set("name", i.Name)
	d.Set("description", i.Description)
	d.Set("locked", i.IsLocked)
	d.Set("blocked", i.IsBlocked)

	return nil
}
