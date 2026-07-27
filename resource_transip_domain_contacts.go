package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/transip/gotransip/v6/domain"
	"github.com/transip/gotransip/v6/repository"
)

var validContactTypes = []string{"registrant", "administrative", "technical"}
var validCompanyTypes = []string{"BV", "BVI/O", "COOP", "CV", "EENMANSZAAK", "KERK", "NV", "OWM", "REDR", "STICHTING", "VERENIGING", "VOF", "BEG", "BRO", "EESV", "ANDERS"}

func resourceDomainContacts() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainContactsUpdate,
		Read:   resourceDomainContactsRead,
		Update: resourceDomainContactsUpdate,
		Delete: resourceDomainContactsDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDomainContactsImport,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Description: "The domain name, including the tld.",
				Required:    true,
				ForceNew:    true,
			},
			"contact": {
				Type:        schema.TypeList,
				Description: "List of whois contacts for this domain.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Description:  "The type of contact: 'registrant', 'administrative' or 'technical'.",
							Required:     true,
							ValidateFunc: validation.StringInSlice(validContactTypes, false),
						},
						"first_name": {
							Type:        schema.TypeString,
							Description: "The first name of this contact.",
							Required:    true,
						},
						"last_name": {
							Type:        schema.TypeString,
							Description: "The last name of this contact.",
							Required:    true,
						},
						"company_name": {
							Type:        schema.TypeString,
							Description: "The company name of this contact.",
							Optional:    true,
						},
						"company_kvk": {
							Type:        schema.TypeString,
							Description: "The kvk number of this contact.",
							Optional:    true,
						},
						"company_type": {
							Type:         schema.TypeString,
							Description:  "The company type. Possible values: 'BV', 'BVI/O', 'COOP', 'CV', 'EENMANSZAAK', 'KERK', 'NV', 'OWM', 'REDR', 'STICHTING', 'VERENIGING', 'VOF', 'BEG', 'BRO', 'EESV', 'ANDERS'.",
							Optional:     true,
							ValidateFunc: validation.StringInSlice(validCompanyTypes, false),
						},
						"street": {
							Type:        schema.TypeString,
							Description: "The street of the address of this contact.",
							Required:    true,
						},
						"number": {
							Type:        schema.TypeString,
							Description: "The street number of the address of this contact.",
							Required:    true,
						},
						"postal_code": {
							Type:        schema.TypeString,
							Description: "The postal code of the address of this contact.",
							Required:    true,
						},
						"city": {
							Type:        schema.TypeString,
							Description: "The city of the address of this contact.",
							Required:    true,
						},
						"phone_number": {
							Type:        schema.TypeString,
							Description: "The phone number of this contact.",
							Required:    true,
						},
						"fax_number": {
							Type:        schema.TypeString,
							Description: "The fax number of this contact.",
							Optional:    true,
						},
						"email": {
							Type:        schema.TypeString,
							Description: "The email address of this contact.",
							Required:    true,
						},
						"country": {
							Type:        schema.TypeString,
							Description: "The country code (ISO 3166-1 alpha-2, lowercase) of this contact.",
							Required:    true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								if len(val.(string)) != 2 {
									errs = append(errs, fmt.Errorf("%q must be a 2-letter ISO country code", key))
								}
								return
							},
						},
					},
				},
			},
		},
	}
}

func resourceDomainContactsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	d.Set("domain", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceDomainContactsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repo := domain.Repository{Client: client}

	domainName := d.Get("domain").(string)
	contacts, err := repo.GetContacts(domainName)
	if err != nil {
		return fmt.Errorf("failed to get contacts of domain %q: %s", domainName, err)
	}

	err = d.Set("contact", whoisContactsToMaps(contacts))
	if err != nil {
		return fmt.Errorf("failed to set contacts for domain %q: %s", domainName, err)
	}

	d.SetId(domainName)
	return nil
}

func resourceDomainContactsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repo := domain.Repository{Client: client}

	domainName := d.Get("domain").(string)
	contacts := interfacesToWhoisContacts(d.Get("contact").([]interface{}))
	err := repo.UpdateContacts(domainName, contacts)
	if err != nil {
		return fmt.Errorf("failed to update contacts of domain %q: %s", domainName, err)
	}

	d.SetId(domainName)
	return nil
}

func resourceDomainContactsDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func whoisContactsToMaps(contacts []domain.WhoisContact) []map[string]interface{} {
	maps := make([]map[string]interface{}, len(contacts))
	for i, c := range contacts {
		m := make(map[string]interface{})
		m["type"] = c.Type
		m["first_name"] = c.FirstName
		m["last_name"] = c.LastName
		m["company_name"] = c.CompanyName
		m["company_kvk"] = c.CompanyKvk
		m["company_type"] = c.CompanyType
		m["street"] = c.Street
		m["number"] = c.Number
		m["postal_code"] = c.PostalCode
		m["city"] = c.City
		m["phone_number"] = c.PhoneNumber
		m["fax_number"] = c.FaxNumber
		m["email"] = c.Email
		m["country"] = c.Country
		maps[i] = m
	}
	return maps
}

func interfacesToWhoisContacts(interfaces []interface{}) []domain.WhoisContact {
	contacts := make([]domain.WhoisContact, len(interfaces))
	for i, v := range interfaces {
		m := v.(map[string]interface{})
		contacts[i] = domain.WhoisContact{
			Type:        m["type"].(string),
			FirstName:   m["first_name"].(string),
			LastName:    m["last_name"].(string),
			CompanyName: m["company_name"].(string),
			CompanyKvk:  m["company_kvk"].(string),
			CompanyType: m["company_type"].(string),
			Street:      m["street"].(string),
			Number:      m["number"].(string),
			PostalCode:  m["postal_code"].(string),
			City:        m["city"].(string),
			PhoneNumber: m["phone_number"].(string),
			FaxNumber:   m["fax_number"].(string),
			Email:       m["email"].(string),
			Country:     m["country"].(string),
		}
	}
	return contacts
}
