# =============================================================================
# Cluo Terraform Outputs
# Single VPS Architecture: Staging + Production
# =============================================================================

# -----------------------------------------------------------------------------
# VPS Outputs
# -----------------------------------------------------------------------------

output "vps_ipv4_address" {
  description = "VPS IPv4 address (hosts both staging and production)"
  value       = hcloud_server.main.ipv4_address
}

output "ssh_connection_string" {
  description = "SSH connection string"
  value       = "ssh root@${hcloud_server.main.ipv4_address}"
}

# -----------------------------------------------------------------------------
# Environment URLs
# -----------------------------------------------------------------------------

output "staging_urls" {
  value = {
    api    = "https://staging-api.${var.domain_name}"
    web    = "https://staging.${var.domain_name}"
    mobile = "https://staging-mobile.${var.domain_name}"
    assets = "https://assets-staging.${var.domain_name}"
  }
  description = "Staging environment URLs"
}

output "production_urls" {
  value = {
    api    = "https://api.${var.domain_name}"
    web    = "https://${var.domain_name}"
    www    = "https://www.${var.domain_name} (managed outside Terraform)"
    mobile = "https://mobile.${var.domain_name}"
    assets = "https://assets.${var.domain_name}"
  }
  description = "Production environment URLs"
}

# -----------------------------------------------------------------------------
# S3 Buckets
# -----------------------------------------------------------------------------

output "assets_buckets" {
  value = {
    staging    = aws_s3_bucket.assets_staging.id
    production = aws_s3_bucket.assets_production.id
  }
  description = "S3 buckets for assets"
}

output "backup_buckets" {
  value = {
    staging    = aws_s3_bucket.backups_staging.id
    production = aws_s3_bucket.backups_production.id
  }
  description = "S3 buckets for PostgreSQL backups"
}

output "vault_storage_bucket" {
  description = "S3 bucket for Vault storage"
  value       = aws_s3_bucket.vault.id
}

# -----------------------------------------------------------------------------
# IAM Access Keys for Ansible/Application
# -----------------------------------------------------------------------------

output "staging_assets_iam_access_key" {
  description = "IAM access key ID for staging assets S3 access"
  value       = aws_iam_access_key.assets_staging_key.id
  sensitive   = true
}

output "staging_assets_iam_secret_key" {
  description = "IAM secret access key for staging assets S3 access"
  value       = aws_iam_access_key.assets_staging_key.secret
  sensitive   = true
}

output "production_assets_iam_access_key" {
  description = "IAM access key ID for production assets S3 access"
  value       = aws_iam_access_key.assets_production_key.id
  sensitive   = true
}

output "production_assets_iam_secret_key" {
  description = "IAM secret access key for production assets S3 access"
  value       = aws_iam_access_key.assets_production_key.secret
  sensitive   = true
}

output "cluo_app_access_key_id" {
  description = "IAM access key ID for cluo-app (least-privilege, cluo-assets-prod only)"
  value       = aws_iam_access_key.cluo_app.id
  sensitive   = true
}

output "cluo_app_secret_access_key" {
  description = "IAM secret access key for cluo-app (least-privilege, cluo-assets-prod only)"
  value       = aws_iam_access_key.cluo_app.secret
  sensitive   = true
}

# -----------------------------------------------------------------------------
# SES Email
# -----------------------------------------------------------------------------

output "ses_dkim_tokens" {
  value       = aws_ses_domain_dkim.main.dkim_tokens
  description = "SES DKIM tokens for DNS verification"
}

output "ses_verification_token" {
  value       = aws_ses_domain_identity.main.verification_token
  description = "SES domain verification token"
  sensitive   = true
}

# -----------------------------------------------------------------------------
# Vault
# -----------------------------------------------------------------------------

output "vault_kms_key_id" {
  description = "KMS key ID for Vault auto-unseal"
  value       = aws_kms_key.vault.id
}

output "vault_kms_key_arn" {
  description = "KMS key ARN for Vault auto-unseal"
  value       = aws_kms_key.vault.arn
}

output "vault_access_key_id" {
  description = "IAM access key ID for Vault S3 storage"
  value       = aws_iam_access_key.vault_key.id
  sensitive   = true
}

output "vault_secret_access_key" {
  description = "IAM secret access key for Vault S3 storage"
  value       = aws_iam_access_key.vault_key.secret
  sensitive   = true
}

# -----------------------------------------------------------------------------
# CloudFront CDN
# -----------------------------------------------------------------------------

# CloudFront outputs - commented out as CloudFront is disabled for now
# output "cloudfront_distribution_id" {
#   description = "CloudFront distribution ID for production assets"
#   value       = aws_cloudfront_distribution.assets.id
# }
#
# output "cloudfront_domain_name" {
#   description = "CloudFront distribution domain name"
#   value       = aws_cloudfront_distribution.assets.domain_name
# }

# -----------------------------------------------------------------------------
# Post-Provisioning Instructions
# -----------------------------------------------------------------------------

output "next_steps" {
  description = "Next steps after Terraform apply"
  value       = <<-EOT
    1. SSH into the server:
       ssh root@${hcloud_server.main.ipv4_address}

    2. Wait for cloud-init to complete:
       cloud-init status --wait

    3. Verify Docker services are running:
       docker ps

    4. Check Caddy is running:
       systemctl status caddy

    5. View logs:
       journalctl -u caddy -f

    Environment URLs:
      Staging:  https://staging.${var.domain_name}
      Production: https://${var.domain_name}
  EOT
}
