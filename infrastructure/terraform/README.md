# CLUO Infrastructure - Terraform

This Terraform configuration provisions the infrastructure for the CLUO application across Hetzner, Cloudflare, and AWS.

## Architecture

- **Hetzner**: Single VPS running Docker Compose with all services
- **Cloudflare**: DNS management and SSL/TLS termination
- **AWS S3**: Media file storage only

## Prerequisites

1. **Terraform** >= 1.5.0 installed
2. **Hetzner Cloud** account with API token
3. **Cloudflare** account with API token
4. **AWS** account with IAM permissions

## Getting Started

### 1. Get Required Tokens & IDs

**Hetzner:**
- Go to https://console.hetzner.cloud/projects/YOUR_PROJECT/security/tokens
- Create an API token with Read & Write permissions

**Cloudflare:**
- Get your Zone ID from the dashboard
- Create an API token at https://dash.cloudflare.com/profile/api-tokens
- Required permissions: `Zone:Edit` and `DNS:Edit`

**AWS:**
- You'll need credentials with permissions to create:
  - S3 buckets
  - IAM users and access keys

### 2. Configure Terraform

```bash
# Copy example terraform.tfvars
cp terraform.tfvars.example terraform.tfvars

# Edit with your values
nano terraform.tfvars
```

### 3. Initialize and Apply

```bash
# Initialize Terraform
terraform init

# Review the plan
terraform plan

# Apply the configuration
terraform apply
```

### 4. Get Outputs

After applying, get your credentials:

```bash
# Show all outputs (including sensitive)
terraform output -json

# Or show specific output
terraform output iam_access_key_id
terraform output iam_secret_access_key
```

## Post-Provisioning Steps

### 1. Connect to the Server

```bash
ssh root@$(terraform output -raw server_ip)
```

### 2. Set Up Docker Compose

```bash
# Clone your repository or copy files
git clone <your-repo> /opt/cluo
cd /opt/cluo

# Copy and configure .env with your S3 credentials
cp .env.example .env
nano .env  # Add AWS credentials from terraform output

# Start services
docker compose up -d
```

### 3. Configure Reverse Proxy (Optional)

If you want to use Cloudflare proxy with multiple domains, you may need a reverse proxy like nginx:

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
    }
}

server {
    listen 80;
    server_name app.yourdomain.com;

    location / {
        proxy_pass http://localhost:3001;
    }
}

server {
    listen 80;
    server_name mobile.yourdomain.com;

    location / {
        proxy_pass http://localhost:3000;
    }
}
```

## DNS Records Created

| Subdomain | Type | Target |
|-----------|------|--------|
| api | A | Server IP (DNS-only by default) |
| app | A | Server IP (Proxied) |
| mobile | A | Server IP (Proxied) |
| * | A | Server IP (Proxied) |

## Variables Reference

| Variable | Description | Default |
|----------|-------------|---------|
| `project_name` | Project name | `cluo` |
| `environment` | Environment name | `dev` |
| `domain_name` | Root domain | - |
| `hetzner_token` | Hetzner API token | - |
| `hetzner_server_type` | Server type | `cpx11` |
| `hetzner_server_location` | Server location | `nbg1` |
| `hetzner_ssh_keys` | SSH key names | `[]` |
| `cloudflare_api_token` | Cloudflare token | - |
| `cloudflare_zone_id` | Zone ID | - |
| `aws_region` | AWS region | `eu-central-1` |

## Hetzner Server Types

| Type | vCPU | RAM | Storage | Price/month |
|------|------|-----|---------|-------------|
| cpx11 | 2 | 2 GB | 40 GB | ~€4 |
| cpx21 | 3 | 4 GB | 80 GB | ~€8 |
| cpx31 | 4 | 8 GB | 160 GB | ~€15 |

## Destroy Infrastructure

```bash
terraform destroy
```

**Warning:** This will delete the VPS, DNS records, and S3 bucket. Backup your data first!

## Security Notes

1. **SSH Keys**: Use SSH keys instead of passwords
2. **Firewall**: Configure allowed IPs appropriately
3. **S3**: The bucket is private by default
4. **Secrets**: Never commit `terraform.tfvars` to git
5. **IAM**: The created IAM user has minimal S3 permissions only
