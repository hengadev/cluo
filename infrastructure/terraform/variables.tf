variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "cluo"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "domain_name" {
  description = "Root domain name managed in Cloudflare"
  type        = string
}

# Hetzner variables
variable "hetzner_token" {
  description = "Hetzner Cloud API token"
  type        = string
  sensitive   = true
}

variable "hetzner_server_type" {
  description = "Hetzner server type (e.g., cx22, cpx11, etc.)"
  type        = string
  default     = "cpx11" # 2 vCPU, 2 GB RAM
}

variable "hetzner_server_location" {
  description = "Hetzner server location (e.g., nbg1, fsn1, hel1)"
  type        = string
  default     = "nbg1"
}

variable "hetzner_ssh_keys" {
  description = "List of SSH key names to add to the server"
  type        = list(string)
  default     = []
}

variable "hetzner_enable_backups" {
  description = "Whether to enable Hetzner automatic server backups"
  type        = bool
  default     = true
}

# Cloudflare variables
variable "cloudflare_api_token" {
  description = "Cloudflare API token with Zone Settings:Edit, Zone:Edit, and DNS:Edit permissions"
  type        = string
  sensitive   = true
}

variable "cloudflare_zone_id" {
  description = "Cloudflare zone ID for the domain"
  type        = string
}

# AWS variables
variable "aws_region" {
  description = "AWS region for S3 bucket"
  type        = string
  default     = "eu-central-1"
}

variable "s3_cors_allowed_origins" {
  description = "Allowed origins for S3 CORS on the media bucket (e.g., [\"https://app.example.com\"])"
  type        = list(string)
  default     = []
}
