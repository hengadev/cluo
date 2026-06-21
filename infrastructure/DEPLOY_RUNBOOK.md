# Cluo Dedicated VPS — First Deploy Runbook

Steps that can't be automated by Terraform/Ansible alone, in order. Run these
once, the first time this VPS is provisioned.

## 1. `terraform apply`

Only after reviewing `terraform plan` and getting explicit go-ahead — this
creates real billed resources (VPS, DNS, AWS IAM/S3/KMS).

```bash
cd infrastructure/terraform
terraform apply
```

## 2. Pull the values Ansible needs out of Terraform's outputs

These don't exist before step 1, so they can't be pre-filled:

```bash
terraform output vault_kms_key_id
terraform output vault_storage_bucket
terraform output vault_access_key_id            # sensitive
terraform output vault_secret_access_key         # sensitive
terraform output backup_production_access_key_id      # sensitive
terraform output backup_production_secret_access_key  # sensitive
```

Put these into `infrastructure/ansible/group_vars/all/vault.yml`
(`ansible-vault edit group_vars/all/vault.yml`) under:
`vault_kms_key_id`, `vault_s3_bucket`, `vault_access_key_id`,
`vault_secret_access_key`, `cluo_backup_aws_access_key_id`,
`cluo_backup_aws_secret_access_key`.

Also copy `terraform output vps_ipv4_address` into
`infrastructure/ansible/inventory.yml` (from `inventory.yml.example`).

## 3. First Ansible run

```bash
cd infrastructure/ansible
ansible-playbook site.yml -i inventory.yml
```

`cluo-prod-vault` will come up **sealed and uninitialized** — that's expected,
its healthcheck will show unhealthy until step 4. Nothing else depends on it
being healthy, so this run completes fine.

## 4. Initialize Vault (one-time, manual — do not script this)

```bash
ssh deploy@<vps-ip>
docker exec -it cluo-prod-vault vault operator init
```

This prints 5 recovery key shares and an **Initial Root Token**. Save the
recovery key shares somewhere durable and offline (they're only needed for
rare operations like re-keying — AWS KMS handles unsealing automatically on
every restart after this point). Copy the root token.

## 5. Push the root token back in

```bash
cd infrastructure/ansible
ansible-vault edit group_vars/all/vault.yml --vault-password-file .vault_pass
# set vault_root_token to the Initial Root Token from step 4
ansible-playbook site.yml -i inventory.yml --tags env,deploy
```

This re-templates `.env` with the real token and restarts
`cluo-prod-api`/`cluo-staging-api` to pick it up.

## 6. Enable the transit engine + KV-v2 secrets Vault needs

The `encx` library (`cluo_api/internal/app/container/infrastructure.go`)
hardcodes a single **shared** `KEKAlias: "cluo-encryption-key"` and
`PepperAlias: "cluo"` — used by both prod and staging, not one per
environment. (The `CLUO_VAULT_KEY_NAME`/`CLUO_VAULT_PEPPER_PATH` vars in
`cluo-prod.env.j2`/`cluo-staging.env.j2` are templated but never read by the
app — ignore them.) Nothing creates the transit key or mounts the KV-v2
engine by default, so `cluo-prod-api` will crash-loop on a pepper-storage
404 until you do this:

```bash
docker exec -it cluo-prod-vault sh
export VAULT_TOKEN=<root token from step 4>

vault secrets enable -path=secret kv-v2
vault secrets enable transit
vault write -f transit/keys/cluo-encryption-key
```

The pepper itself doesn't need to be created manually — `encx` generates
and persists it under `secret/data/encx/cluo/pepper` on first successful
boot once the engines above exist. Restart both API containers afterwards
so they pick it up:

```bash
docker compose -f /opt/cluo/docker-compose.yml restart cluo-prod-api
docker compose -f /opt/cluo-staging/docker-compose.yml restart cluo-staging-api
```

Because the KEK/pepper are shared, prod and staging can decrypt each
other's encrypted fields — consistent with them already sharing one
Postgres instance and one MinIO instance.

## 7. Verify

```bash
docker compose -f /opt/cluo/docker-compose.yml ps        # all healthy
docker exec cluo-prod-vault vault status                  # Sealed: false
curl -s https://api.clientvault.fr/health
curl -s https://staging-api.clientvault.fr/health
```
