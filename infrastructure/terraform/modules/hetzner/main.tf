# SSH Key resource (optional - if you want to manage keys with Terraform)
# resource "hcloud_ssh_key" "default" {
#   for_each   = { for idx, key in var.ssh_keys_public : idx => key }
#   name       = each.value.name
#   public_key = each.value.public_key
# }

# Server resource
resource "hcloud_server" "main" {
  name        = var.server_name
  server_type = var.server_type
  location    = var.server_location
  image       = "ubuntu-24.04"
  ssh_keys    = var.ssh_keys
  firewall_ids = var.firewall_ids

  # Enable backups (optional, costs extra)
  # backups = true

  # User data for cloud-init
  user_data = <<-EOF
    #cloud-config
    package_update: true
    package_upgrade: true
    packages:
      - curl
      - git
      - ufw
    runcmd:
      - ufw allow 22/tcp
      - ufw allow 80/tcp
      - ufw allow 443/tcp
      - ufw --force enable
      - curl -fsSL https://get.docker.com -o /tmp/get-docker.sh
      - sh /tmp/get-docker.sh
      - usermod -aG docker ubuntu
  EOF

  labels = {
    project     = "cluo"
    environment = "cluo"
  }
}

# Get the server network info
data "hcloud_server" "main" {
  id = hcloud_server.main.id
}
