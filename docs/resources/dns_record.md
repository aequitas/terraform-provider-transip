# Dns Record Resource



## Argument Reference

* `content` - (Required) The content of of the dns entry, for example '10 mail', '127.0.0.1' or 'www'.
* `domain` - (Required) The name, including the tld of the domain.
* `expire` - (Optional) The expiration period of the dns entry, in seconds. For example 86400 for a day of expiration.
* `name` - (Required) The name of the dns entry, for example '@' or 'www'.
* `type` - (Required) The type of dns entry. Possible types are 'A', 'AAAA', 'CAA', 'CNAME', 'DS', 'MX', 'NS', 'TXT', 'SRV', 'SSHFP', 'TLSA' and 'ALIAS'.

## Attribute Reference

* `id` - n/a

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import dns records using the ID format `domainname/type/name`. For example:

```terraform
import {
  to = transip_dns_record.example
  id = "example.com/A/@"
}
```

Using `terraform import`, import dns records using the ID format `domainname/type/name`. For example:

```console
% terraform import transip_dns_record.example "example.com/A/@"
```
