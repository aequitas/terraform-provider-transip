# Add basic firewall to VPS
resource "transip_vps_firewall" "test" {
  vps_name = transip_vps.test.name

  inbound_rule {
    description = "HTTP"
    port        = "80"
    protocol    = "tcp"
  }

  inbound_rule {
    description = "HTTPS"
    port        = "443"
    protocol    = "tcp"
  }

  inbound_rule {
    description = "SSH"
    port        = "22"
    protocol    = "tcp"
    whitelist   = [
      "192.0.2.0/24"
    ]
  }
}