# =============================================================================
# Cluo Terraform Variables
# Single VPS Architecture: Staging + Production on one server
# =============================================================================

# -----------------------------------------------------------------------------
# Core Variables
# -----------------------------------------------------------------------------

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "cluo"
}

variable "domain_name" {
  description = "Primary domain name (e.g., clientvault.fr)"
  type        = string
}

variable "aws_region" {
  description = "AWS region for S3, SES, KMS, and CloudFront resources"
  type        = string
  default     = "eu-central-1"
}

# -----------------------------------------------------------------------------
# Hetzner Cloud Variables
# -----------------------------------------------------------------------------

variable "hcloud_token" {
  description = "Hetzner Cloud API token"
  type        = string
  sensitive   = true
}

variable "server_type" {
  description = "Hetzner Cloud server type (cpx22 recommended for both envs: 8GB RAM)"
  type        = string
  default     = "cpx22"

  validation {
    condition     = can(regex("^cpx[0-9]{2}$", var.server_type))
    error_message = "Server type must be a valid CX plan (e.g., cpx22, cpx31)."
  }
}

variable "server_location" {
  description = "Hetzner Cloud datacenter location"
  type        = string
  default     = "nbg1"

  validation {
    condition     = contains(["nbg1", "fsn1", "hel1", "hil", "ash", "sin"], var.server_location)
    error_message = "Location must be a valid Hetzner datacenter code."
  }
}

variable "server_image" {
  description = "Server OS image"
  type        = string
  default     = "ubuntu-24.04"
}

variable "ssh_key_name" {
  description = "Name of existing SSH key in Hetzner Cloud"
  type        = string
  default     = "terraform-cluo"
}

variable "enable_backups" {
  description = "Enable automatic backups (additional 20% of server price)"
  type        = bool
  default     = true
}

# -----------------------------------------------------------------------------
# Cloudflare Variables
# -----------------------------------------------------------------------------

variable "cloudflare_token" {
  description = "Cloudflare API token with Zone:Edit and DNS:Edit permissions"
  type        = string
  sensitive   = true
}

variable "zone_id" {
  description = "Cloudflare zone ID for the domain"
  type        = string
}

variable "contact_email" {
  description = "Contact email for DMARC reports"
  type        = string
}

# -----------------------------------------------------------------------------
# Email (SES) Variables
# -----------------------------------------------------------------------------

variable "email_dkim_records" {
  description = "Mailbox provider DKIM records for inbound email forwarding"
  type = map(object({
    type    = string
    content = string
  }))
  default = {}
}

variable "email_dmarc_policy" {
  description = "DMARC policy (none, quarantine, reject)"
  type        = string
  default     = "quarantine"

  validation {
    condition     = contains(["none", "quarantine", "reject"], var.email_dmarc_policy)
    error_message = "DMARC policy must be one of: none, quarantine, reject."
  }
}

variable "mx_servers" {
  description = "Mail exchange servers for inbound email"
  type = list(object({
    server   = string
    priority = number
  }))
  default = [
    { server = "mx1.example.com", priority = 10 },
    { server = "mx2.example.com", priority = 20 }
  ]
}

# -----------------------------------------------------------------------------
# Database Variables
# -----------------------------------------------------------------------------

variable "staging_db_password" {
  description = "PostgreSQL password for staging database"
  type        = string
  sensitive   = true
}

variable "production_db_password" {
  description = "PostgreSQL password for production database"
  type        = string
  sensitive   = true
}
