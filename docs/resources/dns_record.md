# Dns Record Resource



## Argument Reference

* `content` - (Required) The content of of the dns entry, for example '10 mail', '127.0.0.1' or 'www'.
* `domain` - (Required) n/a
* `expire` - (Optional) The expiration period of the dns entry, in seconds. For example 86400 for a day of expiration.
* `name` - (Required) The name of the dns entry, for example '@' or 'www'.
* `type` - (Required) The type of dns entry. Possbible types are 'A', 'AAAA', 'CNAME', 'MX', 'NS', 'TXT', 'SRV', 'SSHFP' and 'TLSA'.

## Attribute Reference

* `id` - n/a