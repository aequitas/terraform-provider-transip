# VPS Reverse DNS Resource

Manages reverse DNS (PTR) records for VPS IP addresses.

## Example Usage

```hcl
resource "transip_vps_reverse_dns" "example" {
  vps_name    = "example-vps"
  ip_address  = transip_vps.example.ip_address
  reverse_dns = "server.example.com."
}
```

## Argument Reference

* `vps_name` - (Required) The VPS name.
* `ip_address` - (Required) The IP address to set reverse DNS for.
* `reverse_dns` - (Required) The reverse DNS hostname (FQDN) for this IP address.

## Attribute Reference

* `id` - The VPS name and IP address in `vps_name/ip_address` format.
