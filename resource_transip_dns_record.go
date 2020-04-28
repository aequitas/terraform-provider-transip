package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/transip/gotransip/v6/domain"
	"github.com/transip/gotransip/v6/repository"
)

func resourceDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecordCreate,
		Read:   resourceDNSRecordRead,
		Update: resourceDNSRecordUpdate,
		Delete: resourceDNSRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// TODO: true for transip?
				StateFunc: func(v interface{}) string {
					value := strings.TrimSuffix(v.(string), ".")
					return strings.ToLower(value)
				},
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"expire": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  86400,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A", "AAAA", "CNAME", "MX", "NS", "TXT", "SRV", "SSHFP", "TLSA",
				}, false),
			},
			"content": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceDNSRecordCreate(d *schema.ResourceData, m interface{}) error {
	domainName := d.Get("domain").(string)
	entryName := d.Get("name").(string)
	entryType := d.Get("type").(string)

	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	// dom, err := domain.GetInfo(client, domainName)
	// if err != nil {
	//   return fmt.Errorf("failed to get domain %s for reading DNS record entries: %s", domainName, err)
	// }

	dnsEntries, err := repository.GetDNSEntries(domainName)
	if err != nil {
		return fmt.Errorf("failed to read DNS record entries for domain %s: %s", domainName, err)
	}

	for _, e := range dnsEntries {
		if e.Name == entryName && e.Type == entryType {
			return fmt.Errorf("DNS entries for %s record named %s already exist", entryType, entryName)
		}
	}

	id := fmt.Sprintf("%s/%s/%s", domainName, entryType, entryName)
	d.SetId(id)

	return resourceDNSRecordUpdate(d, m)
}

func resourceDNSRecordRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	// TODO: transip uniquely identifies the dnsentries using name, expire and type
	// https://github.com/transip/gotransip/blob/9defadb50daea3d11821aed85498078b9aff4986/domain/repository.go#L148
	// don't think it would hurt omiting the expire to keep compatible with older state files for now
	if id != "" {
		idparts := strings.Split(id, "/")
		if len(idparts) == 3 {
			d.Set("domain", idparts[0])
			d.Set("type", idparts[1])
			d.Set("name", idparts[2])
		} else {
			return fmt.Errorf("Incorrect ID format, should match `domainname/type/name`")
		}
	}

	domainName := d.Get("domain").(string)
	entryName := d.Get("name").(string)
	entryType := d.Get("type").(string)

	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	dnsEntries, err := repository.GetDNSEntries(domainName)
	if err != nil {
		return fmt.Errorf("failed to read DNS record entries for domain %s: %s", domainName, err)
	}

	var content []string
	var expire int
	for _, e := range dnsEntries {
		if e.Name == entryName && e.Type == entryType {
			expire = e.Expire
			content = append(content, e.Content)
		}
	}
	if len(content) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("name", entryName)
	d.Set("expire", expire)
	d.Set("type", entryType)
	d.Set("content", content)
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func resourceDNSRecordUpdate(d *schema.ResourceData, m interface{}) error {
	domainName := d.Get("domain").(string)

	entryName := d.Get("name").(string)
	expire := d.Get("expire").(int)
	entryType := d.Get("type").(string)
	content := d.Get("content").([]interface{})

	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		log.Printf("[DEBUG] terraform-provider-transip update %s\n", entryName)

		dnsEntries, err := repository.GetDNSEntries(domainName)
		if err != nil {
			// This error props up often when doing concurrent requests, so probably related to
			// internal locking of the database? Similar as with the SOAP API.
			if strings.Contains(err.Error(), "Internal error occurred, please contact our support") {
				return resource.RetryableError(fmt.Errorf("failed to get existing DNS record entries for domain %s: %s", domainName, err))
			}
			return resource.NonRetryableError(fmt.Errorf("failed to get existing DNS record entries for domain %s: %s", domainName, err))
		}

		// go through all dns entries in the zone (there is no way to read a single entry name)
		for _, existingEntry := range dnsEntries {
			// skip irrelevant entries (ie: the ones not being modified)
			if existingEntry.Name != entryName || existingEntry.Type != entryType {
				continue
			}
			log.Printf("[DEBUG] terraform-provider-transip %s removing %v\n", entryName, existingEntry)
			// remove all entries for the current entry/expiry/type combination
			err := repository.RemoveDNSEntry(domainName, existingEntry)
			if err != nil {
				// This error props up often when doing concurrent requests, so probably related to
				// internal locking of the database? Similar as with the SOAP API.
				if strings.Contains(err.Error(), "Internal error occurred, please contact our support") {
					return resource.RetryableError(fmt.Errorf("failed to remove DNS record entry for domain %s (%v): %s", domainName, existingEntry, err))
				}
				return resource.NonRetryableError(fmt.Errorf("failed to remove DNS record entry for domain %s (%v): %s", domainName, existingEntry, err))
			}
		}

		// add all desired entries for the current entry/expiry/type combination
		for _, c := range content {
			dnsEntry := domain.DNSEntry{
				Name:    entryName,
				Expire:  expire,
				Type:    entryType,
				Content: c.(string),
			}
			log.Printf("[DEBUG] terraform-provider-transip: %s adding %v\n", entryName, dnsEntry)
			err := repository.AddDNSEntry(domainName, dnsEntry)
			if err != nil {
				if strings.Contains(err.Error(), "Internal error occurred, please contact our support") {
					return resource.RetryableError(fmt.Errorf("failed to add DNS record entry for domain %s (%v): %s", domainName, dnsEntry, err))
				}
				return resource.NonRetryableError(fmt.Errorf("failed to add DNS record entry for domain %s (%v): %s", domainName, dnsEntry, err))
			}
		}

		return resource.NonRetryableError(resourceDNSRecordRead(d, m))
	})
}

func resourceDNSRecordDelete(d *schema.ResourceData, m interface{}) error {
	d.Set("content", []string{})

	return resourceDNSRecordUpdate(d, m)
}
