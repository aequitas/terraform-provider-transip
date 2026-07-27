package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/transip/gotransip/v6/ipaddress"
	"github.com/transip/gotransip/v6/repository"
	"github.com/transip/gotransip/v6/vps"
)

func resourceVpsReverseDNS() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpsReverseDNSCreate,
		Read:   resourceVpsReverseDNSRead,
		Update: resourceVpsReverseDNSCreate,
		Delete: resourceVpsReverseDNSDelete,

		Importer: &schema.ResourceImporter{
			State: resourceVpsReverseDNSImport,
		},

		Schema: map[string]*schema.Schema{
			"vps_name": {
				Type:        schema.TypeString,
				Description: "The VPS name.",
				Required:    true,
				ForceNew:    true,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Description: "The IP address to set reverse DNS for.",
				Required:    true,
				ForceNew:    true,
			},
			"reverse_dns": {
				Type:        schema.TypeString,
				Description: "The reverse DNS hostname (FQDN) for this IP address.",
				Required:    true,
			},
		},
	}
}

func resourceVpsReverseDNSImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format: expected vps_name/ip_address, got %s", d.Id())
	}
	d.Set("vps_name", parts[0])
	d.Set("ip_address", parts[1])
	return []*schema.ResourceData{d}, nil
}

func resourceVpsReverseDNSCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repo := vps.Repository{Client: client}

	vpsName := d.Get("vps_name").(string)
	ipStr := d.Get("ip_address").(string)
	reverseDNS := d.Get("reverse_dns").(string)

	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address: %s", ipStr)
	}

	ip := ipaddress.IPAddress{
		Address:    parsedIP,
		ReverseDNS: reverseDNS,
	}

	err := repo.UpdateReverseDNS(vpsName, ip)
	if err != nil {
		return fmt.Errorf("failed to set reverse DNS for %s on VPS %q: %s", ipStr, vpsName, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", vpsName, ipStr))
	return resourceVpsReverseDNSRead(d, m)
}

func resourceVpsReverseDNSRead(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repo := vps.Repository{Client: client}

	vpsName := d.Get("vps_name").(string)
	ipStr := d.Get("ip_address").(string)

	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address: %s", ipStr)
	}

	ip, err := repo.GetIPAddressByAddress(vpsName, parsedIP)
	if err != nil {
		return fmt.Errorf("failed to get IP address %s for VPS %q: %s", ipStr, vpsName, err)
	}

	d.Set("reverse_dns", ip.ReverseDNS)
	d.SetId(fmt.Sprintf("%s/%s", vpsName, ipStr))
	return nil
}

func resourceVpsReverseDNSDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(repository.Client)
	repo := vps.Repository{Client: client}

	vpsName := d.Get("vps_name").(string)
	ipStr := d.Get("ip_address").(string)

	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address: %s", ipStr)
	}

	ip := ipaddress.IPAddress{
		Address:    parsedIP,
		ReverseDNS: "",
	}

	err := repo.UpdateReverseDNS(vpsName, ip)
	if err != nil {
		return fmt.Errorf("failed to clear reverse DNS for %s on VPS %q: %s", ipStr, vpsName, err)
	}

	d.SetId("")
	return nil
}
