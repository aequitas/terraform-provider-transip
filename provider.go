package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/transip/gotransip/v6"
	"strings"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Transip account.",
				DefaultFunc: schema.EnvDefaultFunc("TRANSIP_ACCOUNT_NAME", nil),
			},
			"private_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Contents of the private key file to be used to authenticate.",
				DefaultFunc: schema.EnvDefaultFunc("TRANSIP_PRIVATE_KEY", nil),
			},
			// "access_token": &schema.Schema{
			//   Type:        schema.TypeString,
			//   Required:    true,
			//   Description: "Temporary access token used for authentication.",
			//   DefaultFunc: schema.EnvDefaultFunc("TRANSIP_ACCESS_TOKEN", nil),
			// },
			"read_only": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disable API write calls.",
				Default:     false,
			},
		},

		ConfigureFunc: providerConfigure,

		ResourcesMap: map[string]*schema.Resource{
			"transip_dns_record": resourceDNSRecord(),
			// "transip_domain":     resourceDomain(),
			// "transip_vps":        resourceVps(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"transip_domain": dataSourceDomain(),
			"transip_vps":    dataSourceVps(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiMode := gotransip.APIModeReadWrite
	if d.Get("read_only").(bool) {
		apiMode = gotransip.APIModeReadOnly
	}

	private_key_body := d.Get("private_key").(string)
	if private_key_body == "" {
		return nil, fmt.Errorf("private_key not provided")
	}

	private_key := strings.NewReader(private_key_body)

	client, err := gotransip.NewClient(gotransip.ClientConfiguration{
		AccountName:      d.Get("account_name").(string),
		PrivateKeyReader: private_key,
		Mode:             apiMode,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
