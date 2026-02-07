# CLUO Infrastructure

This directory contains all infrastructure-as-code and configuration for deploying CLUO.

## Structure

```
infrastructure/
├── terraform/     # Cloud resource provisioning
│   ├── modules/   # Reusable Terraform modules
│   └── README.md  # Terraform documentation
│
└── ansible/       # Server configuration and deployment
    ├── roles/     # Ansible roles for different tasks
    ├── group_vars/ # Variables for groups of hosts
    └── README.md  # Ansible documentation
```

## Quick Start

### 1. Provision Infrastructure (Terraform)

```bash
cd infrastructure/terraform

# Configure
cp terraform.tfvars.example terraform.tfvars
# Edit with your tokens

# Provision
terraform init
terraform apply

# Save outputs
terraform output -json > ../ansible/terraform-outputs.json
```

### 2. Configure Server (Ansible)

```bash
cd infrastructure/ansible

# Configure inventory
cp inventory.yml.example inventory.yml
# Add server IP from Terraform output

# Configure variables
cp group_vars/all.yml.example group_vars/all.yml
cp group_vars/vault.yml.example group_vars/vault.yml

# Add secrets and encrypt
nano group_vars/vault.yml
ansible-vault encrypt group_vars/vault.yml

# Deploy
ansible-playbook -i inventory.yml site.yml --ask-vault-pass
```

## Workflow

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│   Hetzner VPS   │──────│   Cloudflare    │──────│      AWS S3     │
│   (Terraform)   │      │   (Terraform)   │      │   (Terraform)   │
└────────┬────────┘      └─────────────────┘      └─────────────────┘
         │
         │ Ansible
         ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Server Configuration                        │
│  ├─ System Hardening  │  ├─ SSH Hardening  │  ├─ Firewall      │
│  ├─ Docker + Compose  │  ├─ App Deploy     │  ├─ Fail2ban      │
│  ├─ Automatic Updates │  ├─ Monitoring     │  └─ Backups       │
└─────────────────────────────────────────────────────────────────┘
```

## Security Checklist

### Terraform
- [ ] S3 bucket is private
- [ ] IAM user has minimal permissions
- [ ] Cloudflare API token has least privilege
- [ ] Terraform state is secured (backend configured)

### Ansible
- [ ] SSH keys only (no passwords)
- [ ] Firewall configured
- [ ] Fail2ban enabled
- [ ] Secrets encrypted with ansible-vault
- [ ] Automatic updates enabled
- [ ] Backups configured

### Operational
- [ ] Monitoring enabled (cAdvisor/Node Exporter)
- [ ] Logs configured and rotating
- [ ] Backup script tested
- [ ] Disaster recovery plan documented

## Cost Summary

| Service | Est. Monthly Cost |
|---------|-------------------|
| Hetzner CPX11 | ~€4 |
| AWS S3 (10GB) | ~€0.25 |
| Cloudflare Free | €0 |
| **Total** | **~€4.25/month** |

## Maintenance

### Daily
- Monitor backup logs
- Check application status

### Weekly
- Review security logs (Fail2ban)
- Check disk space
- Review Docker container health

### Monthly
- Update Docker images
- Review and rotate secrets
- Test backup restoration
- Review AWS S3 costs

### Quarterly
- Security audit
- Dependency updates
- Performance review
- Disaster recovery test

## Troubleshooting

### Terraform Issues
```bash
# State locked
terraform force-unlock <LOCK_ID>

# Reconfigure backend
terraform init -migrate-state
```

### Ansible Issues
```bash
# Debug mode
ansible-playbook -i inventory.yml site.yml -vvv

# Skip specific host
ansible-playbook -i inventory.yml site.yml --limit 'all:!problem-host'
```

## Useful Commands

### Terraform
```bash
terraform plan              # Preview changes
terraform apply             # Apply changes
terraform destroy           # Destroy resources
terraform output            # Show outputs
terraform refresh           # Refresh state
```

### Ansible
```bash
ansible all -i inventory.yml -m ping                    # Test connection
ansible all -i inventory.yml -m setup                   # Gather facts
ansible-playbook site.yml --check                       # Dry run
ansible-playbook site.yml --tags docker                 # Run specific role
```

### Server
```bash
ssh cluo@server                           # SSH to server
docker compose ps                         # Check containers
docker compose logs -f                    # View logs
sudo ufw status                           # Check firewall
sudo fail2ban-client status               # Check Fail2ban
```

## Documentation

- [Terraform Documentation](terraform/README.md)
- [Ansible Documentation](ansible/README.md)
- [Hetzner Docs](https://docs.hetzner.com/)
- [Cloudflare Docs](https://developers.cloudflare.com/)
- [AWS S3 Docs](https://docs.aws.amazon.com/s3/)
