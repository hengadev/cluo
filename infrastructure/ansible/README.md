# CLUO Infrastructure - Ansible

This Ansible playbook configures and secures a VPS that hosts both **staging** and **production** environments for the CLUO application, mirroring the deployment that currently runs on homelab's shared VPS (see `homelab/ansible/deploy.yml` and `homelab/docker/cluo/`).

## Architecture

**Production and staging are separate Docker Compose projects** (`/opt/cluo` and `/opt/cluo-staging`), not one shared compose file:

- **Production** (`/opt/cluo`): `cluo-prod-postgres`, `cluo-prod-redis`, `cluo-prod-api`, `cluo-prod-web`, `cluo-prod-mobile`
- **Staging** (`/opt/cluo-staging`): `cluo-staging-api`, `cluo-staging-web`, `cluo-staging-mobile`, `cluo-staging-backup`

Staging has **no Postgres or Redis containers of its own** — it joins production's `cluo-prod-internal` network and connects to `cluo-prod-postgres`/`cluo-prod-redis` directly, using a separate database (`cluo_staging`) and Redis DB index. This is why `app_deploy` creates the staging database role/database on the production Postgres container as part of deployment.

Both projects also join an external `proxy` network, reverse-proxied by Caddy (see `files/Caddyfile.snippet` — currently spliced into homelab's shared Caddyfile, not a standalone Caddy instance on this VPS yet).

Object storage (MinIO) and HashiCorp Vault (transit encryption + pepper secrets) are shared, external dependencies referenced by hostname in the env templates (`homelab-minio`, `homelab-vault` on the homelab box today).

## Prerequisites

- Ansible >= 2.14 installed on your local machine
- SSH access to the target VPS
- VPS already provisioned via Terraform
- A reachable Postgres-capable container for `cluo-prod-postgres` and Vault/MinIO endpoints reachable from the VPS

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

### 2. Configure Inventory

```bash
cd infrastructure/ansible

# Copy the example inventory
cp inventory.yml.example inventory.yml

# Edit with your server details and secrets
nano inventory.yml
```

See `inventory.yml.example` for the full variable list (`cluo_prod_db_*`, `cluo_staging_db_*`, `minio_root_*`, `vault_root_token`, `cluo_backup_*`, etc.) — these map directly to `cluo-prod.env.j2` / `cluo-staging.env.j2`.

### 3. Run the Playbook

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
| `app_deploy` | Production + staging application deployment | `app`, `deploy` |
| `monitoring` | cAdvisor and Node Exporter (optional, opt-in) | `monitoring` |
| `automatic_updates` | Unattended security updates | `updates` |
| `uptime_monitoring` | UptimeRobot health checks (optional, opt-in) | `uptime` |

`monitoring` and `uptime_monitoring` are forward-looking scaffolding for this dedicated VPS — they have no equivalent in homelab's current deployment and stay disabled (`enable_monitoring: false`, `enable_uptime_monitoring: false`) until explicitly turned on.

## Environment Files

Two environment files are templated, one per compose project:

**`/opt/cluo/.env`** (from `cluo-prod.env.j2`):
- Database: `cluo_production` on `cluo-prod-postgres`
- Public URL: `https://{{ cluo_domain }}`
- MinIO bucket: `cluo-prod`

**`/opt/cluo-staging/.env`** (from `cluo-staging.env.j2`):
- Database: `cluo_staging`, same Postgres container, separate role
- Public URL: `https://staging-api.{{ cluo_domain }}`
- MinIO bucket: `cluo-staging`
- Also carries the `cluo_backup_*` vars consumed by the `cluo-staging-backup` container

## Backups

Backups are **not** an Ansible role — they run inside the `cluo-staging-backup` container defined in `staging.docker-compose.yml`, which crons `files/backup.sh` to tar+GPG-encrypt MinIO data and upload it to S3. Nothing on the Ansible side needs to schedule or manage this.

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

### Fail2ban
- SSH brute-force protection
- Custom ban times
- Recidive handling

## Inventory Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `app_name` | Application name | `cluo` |
| `app_user` | Deploy user (no dedicated `cluo` system user is created) | `deploy` |
| `app_dir` | Production compose directory | `/opt/cluo` |
| `cluo_domain` | Root domain | `clientvault.fr` |
| `cluo_registry` | Docker image registry | `hengadev` |
| `cluo_prod_tag` / `cluo_staging_tag` | Image tags to deploy | `latest` |
| `cluo_prod_db_name` / `cluo_prod_db_user` / `cluo_prod_db_password` | Production database | `cluo_production` / `cluo` |
| `cluo_staging_db_name` / `cluo_staging_db_user` / `cluo_staging_db_password` | Staging database (same Postgres container) | `cluo_staging` / `cluo_staging` |
| `minio_root_user` / `minio_root_password` | MinIO credentials | - |
| `vault_root_token` | Vault root token for transit keys + pepper secrets | - |
| `cluo_backup_*` | AWS creds, region, bucket, GPG passphrase for the backup container | - |

## Post-Deployment Checklist

- [ ] Verify SSH access with key only
- [ ] Check firewall status: `sudo ufw status`
- [ ] Verify Fail2ban: `sudo fail2ban-client status`
- [ ] Check production containers: `cd /opt/cluo && docker compose ps`
- [ ] Check staging containers: `cd /opt/cluo-staging && docker compose ps`
- [ ] Test staging endpoints
- [ ] Test production endpoints
- [ ] Verify Caddy is routing correctly
- [ ] Check MinIO/S3 connectivity

## Maintenance

### View Containers

```bash
ssh deploy@your-server-ip

cd /opt/cluo && docker compose ps
cd /opt/cluo-staging && docker compose ps
```

### View Logs

```bash
cd /opt/cluo && docker compose logs -f
cd /opt/cluo-staging && docker compose logs -f
```

### Rebuild + Restart

```bash
cd /opt/cluo && docker compose pull && docker compose up -d --remove-orphans
cd /opt/cluo-staging && docker compose pull && docker compose up -d --remove-orphans
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

1. **Never commit** `inventory.yml` or `group_vars/*/vault.yml` with real IPs or credentials — both are gitignored
2. **Use Ansible Vault** for secrets in production
3. **Rotate credentials** regularly
4. **Keep Ansible updated** for security patches
5. **Test in staging** before production changes
6. **Monitor logs** for suspicious activity
