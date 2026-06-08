# Secrets Audit — Ansible Infrastructure

**Date:** 2026-06-08
**Status:** Completed — Secrets migrated to ansible-vault

## Scope

All Ansible inventory, variable, and role files under `infrastructure/ansible/`.

## Findings

### Plaintext secrets found on disk (never committed to git)

| File | Secrets found |
|------|--------------|
| `inventory.yml` | AWS IAM access keys + secret keys (staging & production), PostgreSQL passwords, SSH public key, admin email, server IP |
| `terraform-outputs.json` | AWS IAM keys (staging, production, vault user), SES verification token, vault KMS key ARN, server IP |

**Both files were properly gitignored and had never been committed to git history.** No real secrets exist in any tracked file. The `.example` files and role templates contain only placeholder values.

### Committed files — clean

All tracked files (`inventory.yml.example`, `group_vars/*.example`, `roles/**/*.yml`, `site.yml`, `ansible.cfg`, `Makefile`, `README.md`) contain only:
- Placeholder/example strings (e.g., `your_staging_db_password_here`)
- Jinja2 template references (e.g., `{{ staging_postgres_password }}`)

No real credentials were ever committed.

## Actions taken

1. **Created `group_vars/all/vault.yml`** — ansible-vault encrypted file containing all sensitive values:
   - Database passwords (staging & production)
   - AWS IAM credentials (staging, production, vault user)
   - AWS SES verification token
   - SMTP credentials

2. **Refactored `inventory.yml`** — replaced all plaintext secrets with Jinja2 references to vault variables (`{{ vault_staging_postgres_password }}`, etc.)

3. **Updated `ansible.cfg`** — added `vault_password_file = .vault_pass` for automatic vault decryption

4. **Updated `.gitignore`** — confirmed coverage of `group_vars/*/vault.yml`, `inventory.yml`, `terraform-outputs.json`, `.vault_pass`

5. **Fixed pre-existing YAML parse error** in `roles/backup/tasks/main.yml` — the backup alert script contained multiline shell strings that broke YAML parsing. Rewrote using `printf` to avoid the issue.

6. **Updated `Makefile`** — corrected `VAULT_FILE` path from `group_vars/vault.yml` to `group_vars/all/vault.yml`

## Verification

- `ansible-playbook site.yml --syntax-check` passes ✅
- All vault variables resolve correctly via `ansible -m debug` ✅
- No real secrets in any git-tracked file (`git ls-files | xargs grep`) ✅
- All sensitive files confirmed gitignored (`git check-ignore`) ✅

## Vault contents reference

| Variable | Description |
|----------|-------------|
| `vault_staging_postgres_password` | Staging DB password |
| `vault_production_postgres_password` | Production DB password |
| `vault_staging_s3_access_key_id` | Staging AWS IAM access key |
| `vault_staging_s3_secret_access_key` | Staging AWS IAM secret key |
| `vault_production_s3_access_key_id` | Production AWS IAM access key |
| `vault_production_s3_secret_access_key` | Production AWS IAM secret key |
| `vault_aws_vault_access_key_id` | Vault user AWS IAM access key |
| `vault_aws_vault_secret_access_key` | Vault user AWS IAM secret key |
| `vault_ses_verification_token` | AWS SES domain verification token |
| `vault_smtp_user` | SMTP username (empty if unused) |
| `vault_smtp_password` | SMTP password (empty if unused) |

## Recommendations

1. **Rotate credentials** — The real IAM keys and passwords were present in `inventory.yml` on disk. While never committed, anyone with local access could have read them. Consider rotating:
   - AWS IAM access keys (staging, production, vault user)
   - PostgreSQL passwords
   - SES verification token is non-sensitive (only used for DNS validation)

2. **Remove `terraform-outputs.json`** from disk after rotating, since its secrets are now in the vault.

3. **Back up `.vault_pass`** — The vault password file (`.vault_pass`) must be preserved. Consider storing it in a password manager.
