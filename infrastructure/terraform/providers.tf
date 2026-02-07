# Hetzner Provider
provider "hetzner" {
  token = var.hetzner_token
}

# Cloudflare Provider
provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

# AWS Provider (for S3 only)
provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "cluo"
      ManagedBy   = "terraform"
      Environment = var.environment
    }
  }
}
