# Modules
module "hetzner" {
  source = "./modules/hetzner"

  server_name     = "${var.project_name}-${var.environment}"
  server_type     = var.hetzner_server_type
  server_location = var.hetzner_server_location
  ssh_keys        = var.hetzner_ssh_keys
  firewall_ids    = module.hetzner_firewall.firewall_ids
  enable_backups  = var.hetzner_enable_backups
}

module "hetzner_firewall" {
  source = "./modules/hetzner/firewall"

  project_name = var.project_name
  environment  = var.environment

  # SSH access restricted to admin IPs (override in tfvars)
  # ssh_allowed_ips = ["your.admin.ip/32"]

  # Cloudflare IP ranges are set as defaults in the module
  # and will restrict HTTP/HTTPS to Cloudflare-proxied traffic only
}

module "cloudflare" {
  source = "./modules/cloudflare"

  zone_id     = var.cloudflare_zone_id
  domain_name = var.domain_name
  server_ip   = module.hetzner.server_ip

  # Proxied records through Cloudflare
  proxy_api    = true
  proxy_web    = true
  proxy_mobile = true
}

module "aws_s3" {
  source = "./modules/aws_s3"

  bucket_name          = "${var.project_name}-media-${var.environment}"
  region               = var.aws_region
  cors_allowed_origins = var.s3_cors_allowed_origins

  # Create IAM user for the application
  create_iam_user = true
  iam_user_name   = "${var.project_name}-${var.environment}"
}

module "aws_s3_backup" {
  source = "./modules/aws_s3"

  bucket_name = "${var.project_name}-backups-${var.environment}"
  region      = var.aws_region

  # Separate IAM user for backup access
  create_iam_user = true
  iam_user_name   = "${var.project_name}-backup-${var.environment}"

  # No CORS needed for backup bucket
  cors_allowed_origins = []

  # Custom lifecycle: move to STANDARD_IA after 1 day, expire after 90 days
  lifecycle_rules = [
    {
      id                                 = "backup-lifecycle"
      status                             = "Enabled"
      noncurrent_version_expiration_days = null
      transition_days                    = 1
      transition_storage_class           = "STANDARD_IA"
      expiration_days                    = 90
    }
  ]
}
