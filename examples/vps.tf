# Because concurrency is a bit finicky at the moment run this with -parallelism=1

variable "vps_name" {
  default = "test"
}

variable "vps_product" {
  default = "vps-bladevps-x1"
}

variable "vps_os" {
  default = "Debian 6"
}

resource "transip_vps" "test" {
  name = var.vps_name
  product_name = var.vps_product
  operating_system = var.vps_os
}