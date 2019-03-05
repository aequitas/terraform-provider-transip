package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Transip account.",
			},
			"private_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Contents of the private key file to be used to authenticate.",
			},
		},

		ConfigureFunc: providerConfigure,

		ResourcesMap: map[string]*schema.Resource{
			"transip_dns_record": resourceDNSRecord(),
			"transip_domain":     resourceDomain(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"transip_domain": dataSourceDomain(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client, err := gotransip.NewSOAPClient(gotransip.ClientConfig{
		AccountName:    d.Get("account_name").(string),
		PrivateKeyBody: []byte(d.Get("private_key").(string)),
		// Mode:           gotransip.APIModeReadOnly,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
