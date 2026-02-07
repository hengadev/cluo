output "server_ip" {
  description = "Public IP address of the server"
  value       = hcloud_server.main.ipv4_address
}

output "server_name" {
  description = "Name of the server"
  value       = hcloud_server.main.name
}

output "server_id" {
  description = "ID of the server"
  value       = hcloud_server.main.id
}
