output "firewall_ids" {
  description = "List of firewall IDs"
  value       = [hcloud_firewall.main.id]
}

output "firewall_name" {
  description = "Name of the firewall"
  value       = hcloud_firewall.main.name
}
