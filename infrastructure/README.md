# CLUO Infrastructure

This directory contains all infrastructure-as-code and configuration for deploying CLUO with a **single VPS multi-environment architecture**.

## Architecture

**Single Hetzner VPS** hosting both staging and production environments:

| Environment | Services | Ports |
|-------------|----------|-------|
| **Staging** | API, Web, Mobile | 8080, 8100, 8200 |
| **Production** | API, Web, Mobile | 5000, 3100, 3200 |

Each environment has:
- Separate PostgreSQL database (`cluo_staging`, `cluo_production`)
- Separate Redis instance
- Separate S3 buckets for assets
- Hostname-based routing via Caddy reverse proxy

## Structure

```
infrastructure/
├── terraform/     # Cloud resource provisioning
│   ├── main.tf           # Core Terraform configuration
│   ├── variables.tf      # Input variables
│   ├── outputs.tf        # Output values for Ansible
│   ├── hetzner.tf        # Hetzner VPS resources
│   ├── cloudflare.tf     # DNS records
│   ├── s3.tf             # S3 buckets (assets, backups, vault)
│   ├── backups.tf        # IAM users for backups
│   ├── cloud-init.yml.tftpl  # Server initialization template
│   └── terraform.tfvars  # Your configuration values
│
├── ansible/       # Server configuration and deployment
│   ├── site.yml            # Main playbook
│   ├── inventory.yml       # Server inventory (not in git)
│   ├── inventory.yml.example
│   ├── roles/
│   │   ├── system_hardening/
│   │   ├── ssh_hardening/
│   │   ├── firewall/
│   │   ├── fail2ban/
│   │   ├── docker/
│   │   ├── app_user/
│   │   ├── app_deploy/     # Docker Compose deployment
│   │   │   └── files/
│   │   │       └── docker-compose.yml  # Multi-environment setup
│   │   ├── monitoring/
│   │   ├── uptime_monitoring/  # UptimeRobot /health check
│   │   ├── automatic_updates/
│   │   └── backup/
│   └── README.md
│
└── Makefile       # Unified commands for Terraform + Ansible
```

## Quick Start

### 1. Provision Infrastructure (Terraform)

```bash
cd infrastructure

# Configure (if not done already)
cd terraform
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your tokens

# Provision
make init
make apply

# Get server IP and IAM credentials
terraform output server_ipv4
terraform output staging_assets_iam_access_key
terraform output staging_assets_iam_secret_key
terraform output production_assets_iam_access_key
terraform output production_assets_iam_secret_key
```

### 2. Configure Server (Ansible)

```bash
cd infrastructure/ansible

# Configure inventory with Terraform outputs
cp inventory.yml.example inventory.yml
nano inventory.yml  # Add server IP, credentials from Terraform

# Deploy
cd ..
make configure
```

## Workflow

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│   Hetzner VPS   │──────│   Cloudflare    │──────│      AWS S3     │
│   Single Server │     │   DNS Routing   │      │   Assets/Backups│
│   Staging+Prod  │      │   (Terraform)   │      │   (Terraform)   │
└────────┬────────┘      └─────────────────┘      └─────────────────┘
         │
         │ Ansible Configuration
         ▼
┌─────────────────────────────────────────────────────────────────┐
│                  Single VPS Configuration                       │
│  ┌─────────────────┐  ┌─────────────────┐                      │
│  │ Staging Env     │  │ Production Env  │                      │
│  │ Ports: 8xxx     │  │ Ports: 3xxx/5xxx│                      │
│  │ - postgres_stg  │  │ - postgres_prod │                      │
│  │ - redis_stg     │  │ - redis_prod    │                      │
│  │ - api/web/mobile│  │ - api/web/mobile│                      │
│  └─────────────────┘  └─────────────────┘                      │
│                                                                  │
│  Caddy Reverse Proxy (hostname-based routing)                   │
│  System Hardening, Firewall, Fail2ban, Docker                   │
└─────────────────────────────────────────────────────────────────┘
```

## Domain Routing

All traffic flows through Cloudflare to the single VPS:

| Hostname | Route | Environment |
|----------|-------|-------------|
| `staging-api.clientvault.fr` | → port 8080 | Staging API |
| `staging.clientvault.fr` | → port 8100 | Staging Web |
| `staging-mobile.clientvault.fr` | → port 8200 | Staging Mobile |
| `api.clientvault.fr` | → port 5000 | Production API |
| `clientvault.fr` | → port 3100 | Production Web |
| `mobile.clientvault.fr` | → port 3200 | Production Mobile |

## Cost Summary

| Service | Est. Monthly Cost |
|---------|-------------------|
| Hetzner CPX22 (3 vCPU, 8GB RAM) | ~€11 |
| AWS S3 (assets + backups) | ~€1-2 |
| Cloudflare Free | €0 |
| **Total** | **~€12-13/month** |

## Security Checklist

### Terraform
- [ ] S3 buckets configured with proper CORS
- [ ] IAM users have minimal scoped permissions
- [ ] Cloudflare API token has least privilege
- [ ] Firewall restricts to Cloudflare IPs only

### Ansible
- [ ] SSH keys only (no passwords)
- [ ] Firewall configured (UFW)
- [ ] Fail2ban enabled
- [ ] Automatic updates enabled
- [ ] Caddy reverse proxy configured

### Operational
- [ ] Database backups configured
- [ ] Log rotation configured
- [ ] Container health checks enabled
- [ ] Systemd service for auto-start
- [ ] External uptime monitoring configured (UptimeRobot)

## Maintenance

### Uptime Monitoring

The production API `/health` endpoint is monitored by [UptimeRobot](https://uptimerobot.com) (free tier).

| Setting | Value |
|---|---|
| Monitor type | HTTP(S) |
| URL | `https://api.clientvault.fr/health` |
| Check interval | 5 minutes |
| Alert channel | Email to developer |
| Documentation | `/opt/cluo/UPTIME_MONITOR.txt` on the VPS |

**Setup:**
```bash
cd infrastructure
# Set in inventory.yml:
#   enable_uptime_monitoring: true
#   uptimerobot_api_key: "your-api-key"
make configure -- --tags uptime
```

**Recovery:** If the UptimeRobot account is lost, create a new free account and re-run the Ansible role.
The monitor ID and configuration are recorded in `/opt/cluo/UPTIME_MONITOR.txt`.

### Daily
- Monitor application logs
- Check container health

### Weekly
- Review security logs (Fail2ban)
- Check disk space usage
- Review backup status

### Monthly
- Update Docker images
- Review and rotate secrets
- Test backup restoration
- Review AWS S3 costs

## Useful Commands

### Using Make (Recommended)

```bash
cd infrastructure

make help              # Show all commands
make setup             # Full setup: Terraform + Ansible
make provision         # Run Terraform apply
make configure         # Run Ansible playbook
make status            # Check infrastructure status
make logs              # View application logs
make restart           # Restart services
make ping              # Test server connectivity
```

### Terraform

```bash
cd terraform
make init              # Initialize
make plan              # Preview changes
make apply             # Apply changes
make destroy           # Destroy resources
terraform output       # Show outputs
```

### Ansible

```bash
cd ansible
ansible-playbook -i inventory.yml site.yml  # Full deployment
ansible-playbook -i inventory.yml site.yml --tags app,deploy  # Specific roles
ansible all -i inventory.yml -m ping       # Test connection
```

### Server Operations

```bash
# SSH to server
ssh root@<server-ip>

# Check all containers
cd /opt/cluo && docker compose ps

# View logs
docker compose logs -f

# Restart services
docker compose restart

# Restart specific environment
docker compose restart api_staging web_staging mobile_staging
```

## Troubleshooting

### Terraform State Issues
```bash
# State locked
terraform force-unlock <LOCK_ID>

# Reconfigure backend
terraform init -migrate-state
```

### Ansible Connection Issues
```bash
# Debug mode
ansible-playbook -i inventory.yml site.yml -vvv

# Test SSH connection
ssh root@<server-ip> -i ~/.ssh/cluo
```

### Container Issues
```bash
# Check what's using a port
sudo lsof -i :8080

# View container logs
docker compose logs api_staging

# Rebuild containers
docker compose up -d --build
```

## Documentation

- [Terraform Documentation](terraform/README.md)
- [Ansible Documentation](ansible/README.md)
- [Hetzner Docs](https://docs.hetzner.com/)
- [Cloudflare Docs](https://developers.cloudflare.com/)
- [AWS S3 Docs](https://docs.aws.amazon.com/s3/)
