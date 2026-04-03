# CLUO Infrastructure - Ansible

This Ansible playbook configures and secures a single VPS that hosts both **staging** and **production** environments for the CLUO application.

## Architecture

**Single VPS Multi-Environment Setup:**
- **Staging Environment** (Ports 8000-8999)
  - API: `staging-api.clientvault.fr` → port 8080
  - Web: `staging.clientvault.fr` → port 8100
  - Mobile: `staging-mobile.clientvault.fr` → port 8200

- **Production Environment** (Ports 3000-3999, 5000-5999)
  - API: `api.clientvault.fr` → port 5000
  - Web: `clientvault.fr` → port 3100
  - Mobile: `mobile.clientvault.fr` → port 3200

Each environment has:
- Separate PostgreSQL database (`cluo_staging`, `cluo_production`)
- Separate Redis instance
- Separate S3 buckets for assets
- Separate environment files (`.env.staging`, `.env.production`)

## Prerequisites

- Ansible >= 2.14 installed on your local machine
- SSH access to the target VPS
- VPS already provisioned via Terraform
- Terraform outputs available for IAM credentials

## Quick Start

### 1. Install Ansible

```bash
# macOS
brew install ansible

# Ubuntu/Debian
sudo apt update
sudo apt install ansible -y

# Using pip
pip install ansible
```

### 2. Get Terraform Outputs

First, get the required values from your Terraform deployment:

```bash
cd infrastructure/terraform
terraform output server_ipv4          # Server IP
terraform output staging_assets_iam_access_key
terraform output staging_assets_iam_secret_key
terraform output production_assets_iam_access_key
terraform output production_assets_iam_secret_key
```

### 3. Configure Inventory

```bash
cd infrastructure/ansible

# Copy the example inventory
cp inventory.yml.example inventory.yml

# Edit with your server details
nano inventory.yml
```

**Required values in inventory.yml:**
- `ansible_host`: Server IP from Terraform
- `ansible_private_key_file`: Path to your SSH key (~/.ssh/cluo)
- `staging_postgres_password`: From terraform.tfvars
- `production_postgres_password`: From terraform.tfvars
- `staging_s3_access_key_id`: From Terraform output
- `staging_s3_secret_access_key`: From Terraform output
- `production_s3_access_key_id`: From Terraform output
- `production_s3_secret_access_key`: From Terraform output

### 4. Run the Playbook

```bash
# Full deployment
ansible-playbook -i inventory.yml site.yml

# Using Make (from infrastructure directory)
make configure

# Dry run (check mode)
ansible-playbook -i inventory.yml site.yml --check

# Run specific roles with tags
ansible-playbook -i inventory.yml site.yml --tags security,docker
```

## Roles

| Role | Description | Tags |
|------|-------------|------|
| `system_hardening` | Baseline security configurations | `security`, `hardening` |
| `ssh_hardening` | SSH security best practices | `security`, `ssh` |
| `firewall` | UFW firewall configuration | `security`, `firewall` |
| `fail2ban` | Intrusion prevention | `security`, `fail2ban` |
| `docker` | Docker and Docker Compose installation | `docker` |
| `app_user` | Application user and directories | `app`, `user` |
| `app_deploy` | Multi-environment application deployment | `app`, `deploy` |
| `monitoring` | cAdvisor and Node Exporter | `monitoring` |
| `automatic_updates` | Unattended security updates | `updates` |
| `backup` | Automated backups | `backup` |

## Docker Compose Services

The deployment includes **12 containers** (6 per environment):

### Staging Services
- `postgres_staging`: PostgreSQL 16 for staging
- `redis_staging`: Redis 7 for staging
- `api_staging`: API backend (port 8080)
- `web_staging`: Web frontend (port 8100)
- `mobile_staging`: Mobile frontend (port 8200)

### Production Services
- `postgres_production`: PostgreSQL 16 for production
- `redis_production`: Redis 7 for production
- `api_production`: API backend (port 5000)
- `web_production`: Web frontend (port 3100)
- `mobile_production`: Mobile frontend (port 3200)

## Reverse Proxy (Caddy)

Caddy is configured to route traffic based on hostname:

```
staging-api.clientvault.fr    → localhost:8080 (staging API)
staging.clientvault.fr        → localhost:8100 (staging Web)
staging-mobile.clientvault.fr → localhost:8200 (staging Mobile)

api.clientvault.fr    → localhost:5000 (production API)
clientvault.fr        → localhost:3100 (production Web)
mobile.clientvault.fr → localhost:3200 (production Mobile)
```

## Environment Files

Two environment files are created:

**.env.staging:**
- Database: `cluo_staging`
- API URL: `https://staging-api.clientvault.fr`
- Assets bucket: `cluo-assets-staging`

**.env.production:**
- Database: `cluo_production`
- API URL: `https://api.clientvault.fr`
- Assets bucket: `cluo-assets-production`

## Security Features

### System Hardening
- Kernel parameters tuned for security
- Strong password policies
- Secure file permissions
- Log rotation configured

### SSH Hardening
- Key-based authentication only
- Restricted cipher suites
- Root login disabled (prohibit-password)

### Firewall (UFW)
- Default deny incoming, allow outgoing
- SSH rate limiting
- HTTP/HTTPS allowed
- Application ports restricted to localhost

### Fail2ban
- SSH brute-force protection
- Custom ban times
- Recidive handling

## Inventory Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `app_name` | Application name | `cluo` |
| `app_user` | Application user | `cluo` |
| `app_dir` | Application directory | `/opt/cluo` |
| `domain` | Root domain | - |
| `staging_api_domain` | Staging API subdomain | `staging-api.domain` |
| `staging_web_domain` | Staging web subdomain | `staging.domain` |
| `staging_mobile_domain` | Staging mobile subdomain | `staging-mobile.domain` |
| `production_api_domain` | Production API subdomain | `api.domain` |
| `production_web_domain` | Production web subdomain | `domain` |
| `production_mobile_domain` | Production mobile subdomain | `mobile.domain` |
| `staging_postgres_db` | Staging database name | `cluo_staging` |
| `production_postgres_db` | Production database name | `cluo_production` |
| `staging_assets_bucket` | Staging S3 bucket | `cluo-assets-staging` |
| `production_assets_bucket` | Production S3 bucket | `cluo-assets-production` |

## Post-Deployment Checklist

- [ ] Verify SSH access with key only
- [ ] Check firewall status: `sudo ufw status`
- [ ] Verify Fail2ban: `sudo fail2ban-client status`
- [ ] Check all Docker containers: `docker compose ps`
- [ ] Test staging endpoints
- [ ] Test production endpoints
- [ ] Verify Caddy is routing correctly
- [ ] Check S3 connectivity

## Maintenance

### View All Containers

```bash
# SSH into server
ssh root@your-server-ip

# Check containers
cd /opt/cluo
docker compose ps
```

### View Logs

```bash
# All logs
docker compose logs -f

# Staging API logs
docker compose logs -f api_staging

# Production API logs
docker compose logs -f api_production

# Caddy logs
sudo journalctl -u caddy -f
```

### Restart Services

```bash
# Restart all services
docker compose restart

# Restart staging only
docker compose restart api_staging web_staging mobile_staging

# Restart production only
docker compose restart api_production web_production mobile_production
```

### Update Application

```bash
# SSH into server
ssh root@your-server-ip

# Pull latest changes
cd /opt/cluo
git pull

# Rebuild and restart
docker compose up -d --build
```

## Troubleshooting

### Port Conflicts

If you get port conflicts:
```bash
# Check what's using the port
sudo lsof -i :8080

# Kill the process if needed
sudo kill -9 <PID>
```

### Database Connection Issues

```bash
# Check PostgreSQL is running
docker compose ps postgres_staging
docker compose ps postgres_production

# View database logs
docker compose logs postgres_staging
docker compose logs postgres_production
```

### Caddy Routing Issues

```bash
# Check Caddy status
sudo systemctl status caddy

# View Caddy logs
sudo journalctl -u caddy -n 50

# Validate Caddyfile
sudo caddy validate --config /etc/caddy/Caddyfile

# Reload Caddy
sudo systemctl reload caddy
```

## Development

### Testing Changes

```bash
# Check mode (no changes made)
ansible-playbook -i inventory.yml site.yml --check --diff

# Run with specific tags
ansible-playbook -i inventory.yml site.yml --tags app,deploy

# Limit to specific host
ansible-playbook -i inventory.yml site.yml --limit cluo-vps
```

## Security Notes

1. **Never commit** `inventory.yml` with real IPs or credentials
2. **Use Ansible Vault** for secrets in production
3. **Rotate credentials** regularly
4. **Keep Ansible updated** for security patches
5. **Test in staging** before production changes
6. **Monitor logs** for suspicious activity
