# Vps Data Source



## Argument Reference

* `name` - (Required) The unique VPS name.

## Attribute Reference

* `availability_zone` - The name of the availability zone the VPS is in.
* `cpus` - The VPS cpu count.
* `description` - The name that can be set by customer.
* `disk_size` - The VPS disk size in kB.
* `id` - n/a
* `ip_address` - The VPS main ipAddress.
* `ipv4_addresses` - All IPV4 addresses associated with this VPS.
* `ipv6_address` - All IPV6 addresses associated with this VPS.
* `is_blocked` - If the VPS is administratively blocked.
* `is_customer_locked` - If this VPS is locked by the customer.
* `is_locked` - Whether or not another process is already doing stuff with this VPS.
* `mac_address` - The VPS macaddress.
* `memory_size` - The VPS memory size in kB.
* `operating_system` - The VPS OperatingSystem.
* `product_name` - The product name.
* `status` - The VPS status, either 'created', 'installing', 'running', 'stopped' or 'paused'.
* `tags` - The custom tags added to this VPS.