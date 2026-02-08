variable "zone_id" {
  description = "Cloudflare zone ID"
  type        = string
}

variable "domain_name" {
  description = "Root domain name"
  type        = string
}

variable "server_ip" {
  description = "IP address of the server"
  type        = string
}

variable "proxy_api" {
  description = "Whether to proxy the API record through Cloudflare"
  type        = bool
  default     = false # Set to false for API to avoid issues with WebSocket
}

variable "proxy_web" {
  description = "Whether to proxy the web record through Cloudflare"
  type        = bool
  default     = true
}

variable "proxy_mobile" {
  description = "Whether to proxy the mobile record through Cloudflare"
  type        = bool
  default     = true
}

variable "enable_cache_rules" {
  description = "Whether to create Cloudflare cache rules (bypass API, cache static assets)"
  type        = bool
  default     = false
}
