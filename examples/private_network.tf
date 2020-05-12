variable "private_network_name" {
    default = "test"
}

# the private network data resource can be created by using the unique name TransIP assigns to the resource. (You can view this in the control panel)
data "transip_private_network" "test" {
    name = var.private_network_name
}

# order a new private network, or use `terraform import transip_private_netwerk.test test` to import existing one
resource "transip_private_network" "test" {
  description = "test" # Description is the user defined name, as the name attribute is set by TransIP as a unique ID.
}