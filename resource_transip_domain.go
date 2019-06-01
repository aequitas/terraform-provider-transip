package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip"
	"github.com/transip/gotransip/domain"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		// Update: resourceDomainUpdate,
		Delete: resourceDomainDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// TODO: implement update
				ForceNew: true,
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(gotransip.Client)
	name := d.Get("name").(string)

	dom := domain.Domain{Name: name}

	_, err := domain.Register(client, dom)
	if err != nil {
		return fmt.Errorf("failed to register domain %q: %s", name, err)
	}

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		var err error
		_, err = domain.GetInfo(client, name)
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error creating domain: %s", err)
	}

	d.SetId(name)

	return resourceDomainRead(d, m)
}

func resourceDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(gotransip.Client)
	name := d.Get("name").(string)

	i, err := domain.GetInfo(client, name)

	// if domain does not exist, inform Terraform and gracefuly exit
	if err != nil && err.Error() == "SOAP Fault 102: One or more domains could not be found." {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to lookup domain %q: %s", name, err)
	}

	var nameservers []map[string]string
	for _, n := range i.Nameservers {
		nameservers = append(nameservers, map[string]string{
			"hostname":     n.Hostname,
			"ipv4_address": n.IPv4Address.String(),
			"ipv6_address": n.IPv6Address.String(),
		})
	}

	d.SetId(i.Name)
	d.Set("is_locked", i.IsLocked)
	d.Set("nameservers", nameservers)

	return nil
}

// func resourceDomainUpdate(d *schema.ResourceData, m interface{}) error {
// 	return resourceDomainRead(d, m)
// }

func resourceDomainDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(gotransip.Client)
	name := d.Get("name").(string)

	err := domain.Cancel(client, name, gotransip.CancellationTimeImmediately)
	if err != nil {
		return fmt.Errorf("failed to cancel domain %q: %s", name, err)
	}

	d.SetId("")
	return nil
}
