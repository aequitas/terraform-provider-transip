package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/domain"
	"github.com/transip/gotransip/v6/repository"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		// Update: resourceDomainUpdate,
		Delete: resourceDomainDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name, including the tld of this domain",
				Required:    true,
				// TODO: implement update
				ForceNew: true,
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	register := domain.Register{
		DomainName: name,
	}
	err := repository.Register(register)
	if err != nil {
		return fmt.Errorf("failed to register domain %q: %s", name, err)
	}

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		var err error
		_, err = repository.GetByDomainName(name)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error waiting for domain to be created: %s", err)
	}

	d.SetId(name)

	return resourceDomainRead(d, m)
}

func resourceDomainRead(d *schema.ResourceData, m interface{}) error {
	name := d.Id()

	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	i, err := repository.GetByDomainName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %q: %s", name, err)
	}

	d.SetId(i.Name)

	d.Set("name", name)

	return nil
}

// func resourceDomainUpdate(d *schema.ResourceData, m interface{}) error {
// 	return resourceDomainRead(d, m)
// }

func resourceDomainDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)

	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	err := repository.Cancel(name, gotransip.CancellationTimeImmediately)
	if err != nil {
		return fmt.Errorf("failed to cancel domain %q: %s", name, err)
	}

	d.SetId("")
	return nil
}
