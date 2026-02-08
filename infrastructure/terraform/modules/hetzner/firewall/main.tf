locals {
  cloudflare_ips = concat(var.cloudflare_ipv4_ranges, var.cloudflare_ipv6_ranges)
}

resource "hcloud_firewall" "main" {
  name = "${var.project_name}-${var.environment}-firewall"

  # SSH - restricted to admin IPs
  rule {
    direction  = "in"
    protocol   = "tcp"
    port       = "22"
    source_ips = var.ssh_allowed_ips
  }

  # HTTP - Cloudflare only (for redirect to HTTPS)
  rule {
    direction  = "in"
    protocol   = "tcp"
    port       = "80"
    source_ips = local.cloudflare_ips
  }

  # HTTPS - Cloudflare only
  rule {
    direction  = "in"
    protocol   = "tcp"
    port       = "443"
    source_ips = local.cloudflare_ips
  }

  # Allow outbound traffic
  rule {
    direction       = "out"
    protocol        = "tcp"
    port            = "any"
    destination_ips = ["0.0.0.0/0", "::/0"]
  }

  rule {
    direction       = "out"
    protocol        = "udp"
    port            = "any"
    destination_ips = ["0.0.0.0/0", "::/0"]
  }

  labels = {
    project     = var.project_name
    environment = var.environment
  }
}
