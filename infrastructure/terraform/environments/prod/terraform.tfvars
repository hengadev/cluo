# Production environment configuration
# Sensitive values (tokens, API keys) should be set via environment variables:
#   export TF_VAR_hetzner_token="..."
#   export TF_VAR_cloudflare_api_token="..."
#   AWS credentials via: aws configure / env vars / ~/.aws/credentials

project_name = "cluo"
environment  = "prod"
domain_name  = "yourdomain.com"

# Hetzner
hetzner_server_type     = "cpx21" # 3 vCPU, 4 GB RAM
hetzner_server_location = "nbg1"
hetzner_ssh_keys        = ["your-ssh-key-name"]
hetzner_enable_backups  = true

# Cloudflare
cloudflare_zone_id = "your_cloudflare_zone_id"

# AWS
aws_region = "eu-central-1"

# S3 CORS — production domains only
s3_cors_allowed_origins = [
  "https://app.yourdomain.com",
  "https://mobile.yourdomain.com",
]
