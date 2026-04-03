# =============================================================================
# Cluo Terraform Configuration
# Project: Multi-platform productivity application
# Architecture: Single VPS hosting both staging and production environments
# =============================================================================

terraform {
  required_providers {
    # AWS provider for S3, SES, KMS, CloudFront, ACM
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    # Hetzner Cloud provider for VPS
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "~> 1.45"
    }
    # Cloudflare provider for DNS management
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.0"
    }
  }

  required_version = ">= 1.5.0"
}

# =============================================================================
# Provider Configurations
# =============================================================================

provider "aws" {
  region  = var.aws_region
  profile = "terraform-cluo"

  default_tags {
    tags = {
      Project     = var.project_name
      ManagedBy   = "terraform"
      Environment = "shared"
    }
  }
}

provider "hcloud" {
  token = var.hcloud_token
}

provider "cloudflare" {
  api_token = var.cloudflare_token
}
