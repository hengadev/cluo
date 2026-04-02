# Cluo Infrastructure - Terraform

Single VPS architecture hosting both **staging** and **production** environments for the Cluo multi-platform productivity application.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Cloudflare CDN                               │
│  (SSL/TLS, DDoS Protection, WAF, DNS)                               │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                    ┌───────────────┴───────────────┐
                    ▼                               ▼
        ┌───────────────────────┐   ┌───────────────────────┐
        │     Staging Env       │   │    Production Env     │
        │   (Ports 8000-8999)   │   │ (Ports 3000-3999,     │
        │                       │   │            5000-5999)  │
        │ • API:   port 8080    │   │ • API:   port 5000    │
        │ • Web:   port 8100    │   │ • Web:   port 3100    │
        │ • Mobile: port 8200   │   │ • Mobile: port 3200   │
        └───────────────────────┘   └───────────────────────┘
                    │                               │
                    └───────────────┬───────────────┘
                                    ▼
                    ┌───────────────────────────┐
                    │  Single Hetzner VPS       │
                    │  (cpx31: 8GB RAM)         │
                    │                           │
                    │  • Docker Compose         │
                    │  • Caddy Reverse Proxy    │
                    │  • 2x PostgreSQL DBs      │
                    │  • 2x Redis instances     │
                    │  • Vault (auto-unseal)    │
                    └───────────────────────────┘
                                    │
                    ┌───────────────┼───────────────┐
                    ▼               ▼               ▼
            ┌───────────┐   ┌───────────┐   ┌───────────┐
            │   AWS S3  │   │   AWS SES │   │  AWS KMS  │
            │  Assets   │   │   Email   │   │  Vault    │
            │  Backups  │   │           │   │           │
            └───────────┘   └───────────┘   └───────────┘
```

## Subdomain Structure

### Staging Environment
| Subdomain | Type | Target | Port |
|-----------|------|--------|------|
| staging.clientvault.fr | A | VPS IP | 8100 |
| staging-api.clientvault.fr | A | VPS IP | 8080 |
| staging-mobile.clientvault.fr | A | VPS IP | 8200 |
| assets-staging.clientvault.fr | CNAME | S3 Bucket | - |

### Production Environment
| Subdomain | Type | Target | Port |
|-----------|------|--------|------|
| clientvault.fr | A | VPS IP | 3100 |
| www.clientvault.fr | CNAME | clientvault.fr | 3100 |
| api.clientvault.fr | A | VPS IP | 5000 |
| mobile.clientvault.fr | A | VPS IP | 3200 |
| assets.clientvault.fr | CNAME | CloudFront | - |

## Port Allocation

```
Staging Environment (8000-8999):
  8080 → API
  8100 → Web App
  8200 → Mobile App (PWA)

Production Environment (3000-3999, 5000-5999):
  3100 → Web App
  3200 → Mobile App (PWA)
  5000 → API
```

## Prerequisites

1. **Terraform** >= 1.5.0
2. **Hetzner Cloud** account with API token
3. **Cloudflare** account with API token
4. **AWS** account with IAM permissions
5. **jq** for JSON parsing (required by some make targets)

### 1. Hetzner Cloud Setup

1. Go to https://console.hetzner.cloud
2. Create an API token (Read & Write permissions)
3. Add your SSH key:
   - Navigate to Security → SSH Keys
   - Add your public key with name `terraform-cluo`

### 2. Cloudflare Setup

1. Get your Zone ID from the dashboard
2. Create an API token with permissions:
   - Zone:Edit
   - DNS:Edit

### 3. AWS Setup

Configure AWS credentials via one of these methods:

```bash
# Option 1: Environment variables
export AWS_ACCESS_KEY_ID="your_access_key"
export AWS_SECRET_ACCESS_KEY="your_secret_key"

# Option 2: AWS CLI profile
aws configure --profile cluo
export AWS_PROFILE=cluo
```

## Getting Started

### 1. Create Configuration File

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with your values:

```hcl
# Core
project_name = "cluo"
domain_name  = "clientvault.fr"
aws_region   = "eu-central-1"

# Hetzner
server_type    = "cpx31"  # 8GB RAM for both envs
server_location = "nbg1"

# Cloudflare
zone_id       = "your_zone_id"
contact_email = "admin@clientvault.fr"
```

### 2. Set Sensitive Values via Environment

```bash
export TF_VAR_hcloud_token="your_hetzner_token"
export TF_VAR_cloudflare_token="your_cloudflare_token"
export TF_VAR_staging_db_password="strong_password_1"
export TF_VAR_production_db_password="strong_password_2"
```

### 3. Initialize and Apply

```bash
# Initialize Terraform
make init

# Review the plan
make plan

# Apply the configuration
make apply
```

### 4. Configure SES DKIM Records

After Terraform apply, get the DKIM tokens:

```bash
make dkim
```

Add these CNAME records to Cloudflare DNS (Terraform will create them automatically).

### 5. Verify Deployment

```bash
# Show all URLs
make urls

# SSH into the server
make ssh

# Check service status
make status
```

## Common Commands

### Terraform Operations

```bash
make init          # Initialize Terraform
make plan          # Show execution plan
make apply         # Apply changes
make destroy       # Destroy all resources
make output        # Show outputs
make fmt           # Format files
make validate      # Validate configuration
```

### Server Operations

```bash
make ssh                # SSH into VPS
make status             # Show service status
make logs-caddy         # Show Caddy logs
make logs-staging       # Show staging logs
make restart-all        # Restart all services
```

### Deployment

```bash
make deploy             # Deploy both environments
make deploy-staging     # Deploy staging only
make deploy-production  # Deploy production only
```

### Backups

```bash
make backup-staging     # Backup staging database
make backup-production  # Backup production database
```

### Information

```bash
make urls          # Show all environment URLs
make dkim          # Show SES DKIM tokens
make buckets       # Show all S3 buckets
make vault-keys    # Show Vault IAM credentials
```

## Post-Provisioning Steps

### 1. Configure Vault

```bash
# SSH into the server
make ssh

# Initialize Vault
vault operator init

# Unseal Vault (3 times)
vault operator unseal <key1>
vault operator unseal <key2>
vault operator unseal <key3>

# Log in
vault login <root_token>
```

### 2. Configure Application Secrets

Store application secrets in Vault:

```bash
# Staging secrets
vault kv put secret/cluo/staging \
  database_url="postgres://cluo:password@localhost:5432/cluo_staging" \
  redis_url="redis://localhost:6379"

# Production secrets
vault kv put secret/cluo/production \
  database_url="postgres://cluo:password@localhost:5432/cluo_production" \
  redis_url="redis://localhost:6379"
```

### 3. Deploy Application

The cloud-init script will have already:
- Installed Docker and Docker Compose
- Installed Caddy reverse proxy
- Configured Caddy for multi-environment routing
- Cloned the repository to `/opt/cluo`
- Started Docker Compose services

Verify everything is running:

```bash
make status
```

## Docker Compose Structure

```yaml
services:
  # Staging Infrastructure
  postgres_staging:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: cluo_staging
      POSTGRES_USER: cluo
    volumes:
      - postgres_staging_data:/var/lib/postgresql/data

  redis_staging:
    image: redis:7-alpine

  # Staging Apps
  api_staging:
    environment:
      CLUO_ENVIRONMENT: staging
      CLUO_SERVER_PORT: 8080
    ports:
      - "8080:8080"

  web_staging:
    environment:
      PUBLIC_API_URL: https://staging-api.clientvault.fr
    ports:
      - "8100:3000"

  # Production Infrastructure
  postgres_production:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: cluo_production
      POSTGRES_USER: cluo
    volumes:
      - postgres_production_data:/var/lib/postgresql/data

  redis_production:
    image: redis:7-alpine

  # Production Apps
  api_production:
    environment:
      CLUO_ENVIRONMENT: production
      CLUO_SERVER_PORT: 5000
    ports:
      - "5000:5000"

  web_production:
    environment:
      PUBLIC_API_URL: https://api.clientvault.fr
    ports:
      - "3100:3000"
```

## Hetzner Server Types

| Type | vCPU | RAM | Storage | Price/month |
|------|------|-----|---------|-------------|
| cpx21 | 3 | 4 GB | 80 GB | ~€8 |
| cpx31 | 4 | 8 GB | 160 GB | ~€15 |
| cpx41 | 8 | 16 GB | 320 GB | ~€29 |

**Recommended:** `cpx31` for hosting both staging and production.

## Destroy Infrastructure

```bash
make destroy
```

**Warning:** This will delete the VPS, DNS records, and all AWS resources. Backup your data first!

## Security Notes

1. **SSH Keys**: Use SSH keys instead of passwords
2. **Firewall**: Cloudflare IPs are whitelisted for HTTP/HTTPS
3. **S3**: Production assets bucket is private, accessed via CloudFront OAC
4. **Secrets**: Never commit `terraform.tfvars` to git
5. **Vault**: Uses KMS auto-unseal for high availability
6. **Passwords**: Use strong, unique passwords for each environment

## Troubleshooting

### Cloud-Init Issues

Check cloud-init status:

```bash
make ssh
cloud-init status --wait
journalctl -u cloud-init
```

### Caddy Not Starting

```bash
make ssh
systemctl status caddy
journalctl -u caddy -f
```

### Docker Container Issues

```bash
make ssh
cd /opt/cluo
docker-compose logs -f
```

### DNS Propagation

DNS changes may take up to 24 hours to propagate globally. Use `dig` to check:

```bash
dig staging.clientvault.fr
dig api.clientvault.fr
```

## License

MIT
