# CLUO Infrastructure - Ansible

This Ansible playbook configures and secures a VPS for the CLUO application deployment.

## Prerequisites

- Ansible >= 2.14 installed on your local machine
- SSH access to the target VPS
- VPS already provisioned (via Terraform or manually)

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
# Copy the example inventory
cp inventory.yml.example inventory.yml

# Edit with your server details
nano inventory.yml
```

### 3. Run the Playbook

```bash
# Full deployment
ansible-playbook -i inventory.yml site.yml

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
| `app_deploy` | Application deployment | `app`, `deploy` |
| `monitoring` | cAdvisor and Node Exporter | `monitoring` |
| `automatic_updates` | Unattended security updates | `updates` |
| `backup` | Automated backups | `backup` |

## Security Features Implemented

### System Hardening
- Kernel parameters tuned for security
- Strong password policies
- Secure file permissions
- Log rotation configured
- Audit logging enabled

### SSH Hardening
- Key-based authentication only (no passwords)
- Restricted cipher suites and algorithms
- Login banners
- Session limits
- Root login disabled (prohibit-password)

### Firewall (UFW)
- Default deny incoming, allow outgoing
- SSH rate limiting
- HTTP/HTTPS allowed
- Application ports restricted to localhost
- Cloudflare IP support

### Fail2ban
- SSH brute-force protection
- Nginx HTTP auth protection
- Recidive (repeat offender) handling
- Custom ban times

### Automatic Updates
- Unattended security updates
- Automatic package updates
- Update notifications
- Docker packages excluded (managed separately)

## Secrets Management

For production, use Ansible Vault to encrypt sensitive data:

```bash
# Encrypt a string
ansible-vault encrypt_string 'my_secret_password' --name 'postgres_password'

# Create encrypted vault file
ansible-vault create group_vars/all/vault.yml

# Edit encrypted vault
ansible-vault edit group_vars/all/vault.yml

# Run playbook with vault
ansible-playbook -i inventory.yml site.yml --ask-vault-pass
```

### Example vault.yml

```yaml
# Encrypted with: ansible-vault encrypt group_vars/all/vault.yml
postgres_password: "your_secure_password_here"
redis_password: "your_redis_password_here"
jwt_secret: "your_jwt_secret_here"
session_secret: "your_session_secret_here"
aws_access_key_id: "your_aws_access_key"
aws_secret_access_key: "your_aws_secret_key"
```

## Inventory Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `app_name` | Application name | `cluo` |
| `app_user` | Application user | `cluo` |
| `app_dir` | Application directory | `/opt/cluo` |
| `domain` | Root domain | - |
| `api_domain` | API subdomain | `api.yourdomain.com` |
| `web_domain` | Web subdomain | `app.yourdomain.com` |
| `mobile_domain` | Mobile subdomain | `mobile.yourdomain.com` |
| `ssh_port` | SSH port | `22` |
| `enable_fail2ban` | Enable Fail2ban | `true` |
| `enable_monitoring` | Enable monitoring | `false` |
| `enable_automatic_updates` | Enable auto updates | `true` |
| `enable_backups` | Enable backups | `false` |

## Post-Deployment Checklist

- [ ] Verify SSH access with key only
- [ ] Check firewall status: `sudo ufw status`
- [ ] Verify Fail2ban: `sudo fail2ban-client status`
- [ ] Check Docker containers: `docker compose ps`
- [ ] Test application endpoints
- [ ] Verify backups: `ls -lh /opt/cluo/backups/`
- [ ] Check monitoring: `curl http://localhost:8081` (cAdvisor)

## Maintenance

### Update Application

```bash
# SSH into server
ssh cluo@your-server-ip

# Pull latest changes
cd /opt/cluo
git pull

# Rebuild and restart
docker compose up -d --build
```

### View Logs

```bash
# All logs
docker compose logs -f

# Specific service
docker compose logs -f api

# System logs
sudo journalctl -u cluo -f
```

### Restore from Backup

```bash
cd /opt/cluo
./restore.sh backups/postgres_20240101_020000.sql.gz
```

## Troubleshooting

### SSH Access Issues

If you get locked out of SSH, you can use the Hetzner console:

1. Log in to Hetzner Cloud Console
2. Use VNC console to access the server
3. Check SSH logs: `sudo journalctl -u sshd`

### Docker Issues

```bash
# Check Docker status
sudo systemctl status docker

# View Docker logs
sudo journalctl -u docker -n 50

# Restart Docker
sudo systemctl restart docker
```

### Rollback Changes

```bash
# Re-run with specific tags
ansible-playbook -i inventory.yml site.yml --tags docker
```

## Development

### Testing Changes

```bash
# Check mode (no changes made)
ansible-playbook -i inventory.yml site.yml --check --diff

# Limit to specific host
ansible-playbook -i inventory.yml site.yml --limit cluo-prod

# Run with specific tags
ansible-playbook -i inventory.yml site.yml --tags security
```

### Adding New Roles

1. Create role directory: `ansible-galaxy init roles/new_role`
2. Add tasks to `roles/new_role/tasks/main.yml`
3. Reference in `site.yml`

## Security Notes

1. **Never commit** `inventory.yml` with real IPs
2. **Use Ansible Vault** for secrets in production
3. **Rotate credentials** regularly
4. **Keep Ansible updated** for security patches
5. **Review logs** for suspicious activity
6. **Test in staging** before production changes
