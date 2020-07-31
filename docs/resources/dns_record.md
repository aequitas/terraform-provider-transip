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

* `content` - (Required) n/a
* `domain` - (Required) n/a
* `expire` - (Optional) n/a
* `name` - (Required) n/a
* `type` - (Required) n/a

## Attribute Reference

* `id` - n/a