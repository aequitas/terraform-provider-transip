package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/transip/gotransip/v6/domain"
	"github.com/transip/gotransip/v6/repository"
)

// All errors that indicate a temporary failure
var errorStrings = []string{
	"Error setting Dns Entries",
	"Internal error occurred, please contact our support",
	"DNS Entries are currently being saved",
}

func retryableDNSRecordErrorf(err error, format string, a ...interface{}) *resource.RetryError {
	// Check if this is a retryable error
	isRetry := false
	for _, errorString := range errorStrings {
		if strings.Contains(err.Error(), errorString) {
			isRetry = true
			break
		}
	}

	// Format the error
	e := fmt.Errorf(format+": %s", append(a, err)...)

	// Return the retryable error (retry or not)
	if isRetry {
		return resource.RetryableError(e)
	} else {
		return resource.NonRetryableError(e)
	}
}

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
					"A", "AAAA", "CAA", "CNAME", "MX", "NS", "TXT", "SRV", "SSHFP", "TLSA",
				}, false),
			},
			"content": &schema.Schema{
				Type: schema.TypeSet,
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

	dnsEntries, err := repository.GetDNSEntries(domainName)
	if err != nil {
		return fmt.Errorf("failed to read DNS record entries for domain %s: %s", domainName, err)
	}

	for _, e := range dnsEntries {
		if e.Name == entryName && e.Type == entryType {
			return fmt.Errorf("DNS entries for %s record named %s already exist", entryType, entryName)
		}
	}

	// Note: as soon as we use SetId, we assume the resource has been created.
	// In this case that is not strictly true...
	id := fmt.Sprintf("%s/%s/%s", domainName, entryType, entryName)
	d.SetId(id)

	return resourceDNSRecordUpdate(d, m)
}

func resourceDNSRecordRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	// TODO: transip uniquely identifies the dnsentries using name, expire and type
	// https://github.com/transip/gotransip/blob/9defadb50daea3d11821aed85498078b9aff4986/domain/repository.go#L148
	// don't think it would hurt omitting the expire to keep compatible with older state files for now
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

	// We are now going to read from the domain (and retry)
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		dnsEntries, err := repository.GetDNSEntries(domainName)

		if err != nil {
			return retryableDNSRecordErrorf(err, "failed to read DNS record entries for domain %s", domainName)
		}

		if len(dnsEntries) == 0 {
			d.SetId("")
			return nil
		}

		var content []string
		var expire int
		for _, e := range dnsEntries {
			if e.Name == entryName && e.Type == entryType {
				expire = e.Expire
				content = append(content, e.Content)
			}
		}

		log.Printf("[DEBUG] terraform-provider-transip reading record %s, %d, %s, %v\n", entryName, expire, entryType, content)

		d.Set("name", entryName)
		d.Set("expire", expire)
		d.Set("type", entryType)
		d.Set("content", content)
		return nil
	})
}

func resourceDNSRecordUpdate(d *schema.ResourceData, m interface{}) error {
	domainName := d.Get("domain").(string)

	entryName := d.Get("name").(string)
	expire := d.Get("expire").(int)
	entryType := d.Get("type").(string)
	content := d.Get("content").(*schema.Set)

	client := m.(repository.Client)
	repository := domain.Repository{Client: client}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		// We lock resources because Transip only allows one change per domain
		// https://github.com/aequitas/terraform-provider-transip/issues/22
		dnsDomainMutexKV.Lock(domainName)
		defer dnsDomainMutexKV.Unlock(domainName)

		// Read current entries to figure out what needs to be changed
		log.Printf("[DEBUG] terraform-provider-transip update %s\n", entryName)
		dnsEntries, err := repository.GetDNSEntries(domainName)
		if err != nil {
			return retryableDNSRecordErrorf(err, "failed to get existing DNS record entries for domain %s", domainName)
		}

		// go through all dns entries in the zone (there is no way to read a single entry name)
		for _, existingEntry := range dnsEntries {
			// skip irrelevant entries (ie: the ones not being modified)
			if existingEntry.Name != entryName || existingEntry.Type != entryType {
				continue
			}

			// remove all entries for the current entry/expiry/type combination
			log.Printf("[DEBUG] terraform-provider-transip %s removing %v\n", entryName, existingEntry)
			err := repository.RemoveDNSEntry(domainName, existingEntry)
			if err != nil {
				return retryableDNSRecordErrorf(err, "failed to remove DNS record entry for domain %s (%v)", domainName, existingEntry)
			}
		}

		// add all desired entries for the current entry/expiry/type combination
		for _, c := range content.List() {
			dnsEntry := domain.DNSEntry{
				Name:    entryName,
				Expire:  expire,
				Type:    entryType,
				Content: c.(string),
			}

			log.Printf("[DEBUG] terraform-provider-transip: %s adding %v\n", entryName, dnsEntry)
			err := repository.AddDNSEntry(domainName, dnsEntry)

			if err != nil {
				return retryableDNSRecordErrorf(err, "failed to add DNS record entry for domain %s (%v)", domainName, dnsEntry)
			}
		}

		return resource.NonRetryableError(resourceDNSRecordRead(d, m))
	})
}

func resourceDNSRecordDelete(d *schema.ResourceData, m interface{}) error {
	d.Set("content", make([]interface{}, 0))

	return resourceDNSRecordUpdate(d, m)
}
