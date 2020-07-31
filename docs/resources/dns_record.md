# DNS Record Resource

DNS record for zones managed by Transip.

## Example Usage

```hcl
data "transip_domain" "demo" {
  name = var.domain
}

resource "transip_dns_record" "demo" {
  domain  = data.transip_domain.demo.id
  name    = "demo"
  type    = "A"
  content = [
    "192.0.2.0",
  ]
}
```

## Argument Reference

* `domain` - (Required) Domain name to register the record under.
* `name` - (Required) Name of the record.
* `expire` - (Required) TTL/expiry time in seconds.
* `type` - (Required) Type of the record (eg: A, AAAA, CNAME, TXT).
* `content` - (Required) List of record contents (eg: IP address).

## Attribute Reference

* N/A