package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/transip/gotransip"
	"github.com/transip/gotransip/domain"
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
					string(domain.DNSEntryTypeA),
					string(domain.DNSEntryTypeAAAA),
					string(domain.DNSEntryTypeCNAME),
					string(domain.DNSEntryTypeMX),
					string(domain.DNSEntryTypeNS),
					string(domain.DNSEntryTypeTXT),
					string(domain.DNSEntryTypeSRV),
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
	entryType := domain.DNSEntryType(d.Get("type").(string))

	client := m.(gotransip.Client)
	dom, err := domain.GetInfo(client, domainName)
	if err != nil {
		return fmt.Errorf("failed to get domain %s for reading DNS record entries: %s", domainName, err)
	}
	for _, e := range dom.DNSEntries {
		if e.Name == entryName && e.Type == entryType {
			return fmt.Errorf("DNS entries for %s record named %s already exist", entryType, entryName)
		}
	}

	id := fmt.Sprintf("%s/%s/%s", domainName, entryType, entryName)
	d.SetId(id)

	return resourceDNSRecordUpdate(d, m)
}

func resourceDNSRecordRead(d *schema.ResourceData, m interface{}) error {
	client := m.(gotransip.Client)

	id := d.Id()
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
	entryType := domain.DNSEntryType(d.Get("type").(string))

	dom, err := domain.GetInfo(client, domainName)
	if err != nil {
		return fmt.Errorf("failed to get domain %s for reading DNS record entries: %s", domainName, err)
	}

	var content []string
	var ttl int64
	for _, e := range dom.DNSEntries {
		if e.Name == entryName && e.Type == entryType {
			ttl = e.TTL
			content = append(content, e.Content)
		}
	}
	if len(content) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("name", entryName)
	d.Set("expire", ttl)
	d.Set("type", entryType)
	d.Set("content", content)
	return nil
}

func resourceDNSRecordUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(gotransip.Client)
	domainName := d.Get("domain").(string)

	entryName := d.Get("name").(string)
	expire := int64(d.Get("expire").(int))
	entryType := domain.DNSEntryType(d.Get("type").(string))
	content := d.Get("content").([]interface{})

	dom, err := domain.GetInfo(client, domainName)
	if err != nil {
		return fmt.Errorf("failed to get domain %s for writing DNS record entries: %s", domainName, err)
	}

	var newEntries []domain.DNSEntry
	for _, e := range dom.DNSEntries {
		if e.Name == entryName && e.Type == entryType {
			continue
		}
		newEntries = append(newEntries, e)
	}

	for _, c := range content {
		newEntries = append(newEntries, domain.DNSEntry{
			Name:    entryName,
			TTL:     expire,
			Type:    entryType,
			Content: c.(string),
		})
	}

	err = domain.SetDNSEntries(client, domainName, newEntries)
	if err != nil {
		return fmt.Errorf("failed to update DNS entries for domain %s: %s", domainName, err)
	}

	return resourceDNSRecordRead(d, m)
}

func resourceDNSRecordDelete(d *schema.ResourceData, m interface{}) error {
	d.Set("content", []string{})

	return resourceDNSRecordUpdate(d, m)
}
