package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/authenticator"
)

var dnsDomainMutexKV = mutexkv.NewMutexKV()

func envBoolFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v == "1" {
			return true, nil
		}
		return false, nil
	}
}

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
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Contents of the private key file to be used to authenticate.",
				DefaultFunc:   schema.EnvDefaultFunc("TRANSIP_PRIVATE_KEY", nil),
				ConflictsWith: []string{"access_token"},
			},
			"access_token": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Temporary access token used for authentication.",
				DefaultFunc:   schema.EnvDefaultFunc("TRANSIP_ACCESS_TOKEN", nil),
				ConflictsWith: []string{"private_key"},
			},
			"read_only": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disable API write calls.",
				Default:     false,
			},
			"test_mode": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use API test mode.",
				DefaultFunc: envBoolFunc("TRANSIP_TEST_MODE"),
			},
		},

		ConfigureFunc: providerConfigure,

		ResourcesMap: map[string]*schema.Resource{
			"transip_dns_record":                 resourceDNSRecord(),
			"transip_domain":                     resourceDomain(),
			"transip_vps":                        resourceVps(),
			"transip_vps_firewall":               resourceVpsFirewall(),
			"transip_private_network":            resourcePrivateNetwork(),
			"transip_private_network_attachment": resourcePrivateNetworkAttachment(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"transip_domain":          dataSourceDomain(),
			"transip_vps":             dataSourceVps(),
			"transip_private_network": dataSourcePrivateNetwork(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiMode := gotransip.APIModeReadWrite
	if d.Get("read_only").(bool) {
		apiMode = gotransip.APIModeReadOnly
	}

	testMode := d.Get("test_mode").(bool)

	private_key_body := d.Get("private_key").(string)
	access_token := d.Get("access_token").(string)
	if private_key_body == "" && access_token == "" {
		return nil, fmt.Errorf("either private_key or access_token must be provided")
	}

	var cacheDir string
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	cacheFile := filepath.Join(cacheDir, "gotransip_test_token_cache")
	// create cachefile with restricted permissions
	_, err = os.OpenFile(cacheFile, os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to create token cache dir")
	}
	cache, err := authenticator.NewFileTokenCache(cacheFile)
	if err != nil {
		panic(err.Error())
	}

	var client_configuration gotransip.ClientConfiguration

	if private_key_body != "" {
		private_key := strings.NewReader(private_key_body)

		client_configuration = gotransip.ClientConfiguration{
			AccountName:      d.Get("account_name").(string),
			PrivateKeyReader: private_key,
			Mode:             apiMode,
			TestMode:         testMode,
			TokenCache:       cache,
		}
	} else {
		client_configuration = gotransip.ClientConfiguration{
			AccountName: d.Get("account_name").(string),
			Mode:        apiMode,
			TestMode:    testMode,
			Token:       access_token,
		}
	}

	client, err := gotransip.NewClient(client_configuration)
	if err != nil {
		return nil, err
	}

	return client, nil
}
