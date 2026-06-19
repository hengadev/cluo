# =============================================================================
# Hetzner Cloud Resources
# Single VPS hosting both staging and production environments
# =============================================================================

# -----------------------------------------------------------------------------
# SSH Key
# -----------------------------------------------------------------------------

# Reference existing SSH key created in Hetzner Cloud panel
data "hcloud_ssh_key" "default" {
  name = var.ssh_key_name
}

# -----------------------------------------------------------------------------
# Single VPS for Both Environments
# -----------------------------------------------------------------------------

resource "hcloud_server" "main" {
  name        = "cluo-vps"
  image       = var.server_image
  server_type = var.server_type  # cpx31 recommended (8GB RAM for both envs)
  location    = var.server_location
  ssh_keys    = [data.hcloud_ssh_key.default.id]
  backups     = var.enable_backups

  labels = {
    project     = var.project_name
    environment = "multi"  # staging+production on single VPS
    managed_by  = "terraform"
  }

  # Minimal bootstrap only — Ansible (site.yml) handles Docker, Caddy, the
  # app, and all secrets. See cloud-init.yml.tftpl.
  user_data = file("${path.module}/cloud-init.yml.tftpl")
}

# -----------------------------------------------------------------------------
# Firewall (Optional - Recommended for Production)
# -----------------------------------------------------------------------------

resource "hcloud_firewall" "main" {
  name = "cluo-firewall"
  labels = {
    project = var.project_name
  }

  # Allow SSH from anywhere (restrict to your IP in production)
  rule {
    direction = "in"
    protocol  = "tcp"
    port      = "22"
    source_ips = ["0.0.0.0/0"]
  }

  # Allow HTTP from Cloudflare IPs
  rule {
    direction = "in"
    protocol  = "tcp"
    port      = "80"
    source_ips = [
      "173.245.48.0/20",
      "103.21.244.0/22",
      "103.22.200.0/22",
      "103.31.4.0/22",
      "141.101.64.0/18",
      "108.162.192.0/18",
      "190.93.240.0/20",
      "188.114.96.0/20",
      "197.234.240.0/22",
      "198.41.128.0/17",
      "162.158.0.0/15",
      "104.16.0.0/13",
      "104.24.0.0/14",
      "172.64.0.0/13",
      "131.0.72.0/22"
    ]
  }

  # Allow HTTPS from Cloudflare IPs
  rule {
    direction = "in"
    protocol  = "tcp"
    port      = "443"
    source_ips = [
      "173.245.48.0/20",
      "103.21.244.0/22",
      "103.22.200.0/22",
      "103.31.4.0/22",
      "141.101.64.0/18",
      "108.162.192.0/18",
      "190.93.240.0/20",
      "188.114.96.0/20",
      "197.234.240.0/22",
      "198.41.128.0/17",
      "162.158.0.0/15",
      "104.16.0.0/13",
      "104.24.0.0/14",
      "172.64.0.0/13",
      "131.0.72.0/22"
    ]
  }

  # Allow outbound traffic
  rule {
    direction = "out"
    protocol  = "tcp"
    port      = "any"
    destination_ips = ["0.0.0.0/0"]
  }

  rule {
    direction = "out"
    protocol  = "udp"
    port      = "any"
    destination_ips = ["0.0.0.0/0"]
  }
}

# Attach firewall to server
resource "hcloud_firewall_attachment" "main" {
  firewall_id = hcloud_firewall.main.id
  server_ids  = [hcloud_server.main.id]
}
