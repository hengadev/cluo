variable "server_name" {
  description = "Name of the server"
  type        = string
}

variable "server_type" {
  description = "Hetzner server type"
  type        = string
}

variable "server_location" {
  description = "Hetzner server location"
  type        = string
}

variable "ssh_keys" {
  description = "List of SSH key names to add to the server"
  type        = list(string)
  default     = []
}

variable "firewall_ids" {
  description = "List of firewall IDs to attach to the server"
  type        = list(string)
  default     = []
}

variable "enable_backups" {
  description = "Whether to enable Hetzner automatic backups"
  type        = bool
  default     = true
}
