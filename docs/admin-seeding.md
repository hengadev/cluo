# Admin user seeding

## Context

Cluo uses field-level encryption (`encx`) backed by HashiCorp Vault for all sensitive user data (email, role, password hash). Every column in `auth.users` is encrypted or hashed using Vault-managed keys. Raw SQL `INSERT` statements cannot seed a user because the ciphertext must be produced by the application's crypto pipeline.

A dedicated binary, `seed-admin`, handles first-time admin creation. It is compiled into the API Docker image alongside the main `cluo` binary.

## What the seed-admin binary does

- Reads admin credentials from `--email`/`--password` flags or `SEED_ADMIN_EMAIL`/`SEED_ADMIN_PASSWORD` env vars.
- Connects to PostgreSQL and Vault using the same environment variables as the API.
- Hashes the email and checks whether the user already exists — the command is **idempotent** (safe to run multiple times).
- If absent, creates a new user with role `administrator` using the full `encx` encryption pipeline.
- Exits 0 on success or if the user already exists; exits non-zero on any error.

Source: `cluo_api/cmd/seed-admin/main.go`

## Environment variables required

These are the same vars already present in the API container's environment:

| Variable | Purpose |
|---|---|
| `CLUO_DB_HOST` | PostgreSQL host |
| `CLUO_DB_PORT` | PostgreSQL port (default: `5432`) |
| `CLUO_DB_NAME` | Database name |
| `CLUO_DB_USER` | Database user |
| `CLUO_DB_PASSWORD` | Database password |
| `CLUO_DB_SSL_MODE` | SSL mode (default: `require`) |
| `VAULT_ADDR` | Vault address |
| `VAULT_APPROLE_ROLE_ID` | AppRole role ID (production auth method) |
| `VAULT_APPROLE_SECRET_ID` | AppRole secret ID (production auth method) |
| `VAULT_TOKEN` | Dev fallback if AppRole vars are absent |

## Running the command

### Current setup (homelab VPS, Docker Compose)

After the API container is up:

```bash
docker exec cluo-api ./seed-admin \
  --email admin@example.com \
  --password "changeme123"
```

Or pass credentials via env to avoid shell history:

```bash
docker exec \
  -e SEED_ADMIN_EMAIL=admin@example.com \
  -e SEED_ADMIN_PASSWORD=changeme123 \
  cluo-api ./seed-admin
```

The command inherits all other required env vars from the running container.

### Future infra (infrastructure/ Ansible)

Add a one-shot task to the `app_deploy` Ansible role, **after** the container is confirmed running. Use a stat/check guard so it only runs on first deploy.

Suggested task in `infrastructure/ansible/roles/app_deploy/tasks/main.yml`:

```yaml
- name: Check if admin user already exists
  community.docker.docker_container_exec:
    container: cluo-api
    command: ./seed-admin --email "{{ cluo_admin_email }}"
  register: seed_check
  ignore_errors: true
  changed_when: false

- name: Seed admin user on first deploy
  community.docker.docker_container_exec:
    container: cluo-api
    env:
      SEED_ADMIN_EMAIL: "{{ cluo_admin_email }}"
      SEED_ADMIN_PASSWORD: "{{ cluo_admin_password }}"
    command: ./seed-admin
  when: seed_check.rc != 0
```

Store `cluo_admin_email` and `cluo_admin_password` in Ansible Vault under `infrastructure/ansible/group_vars/all/vault.yml`. They should be rotated immediately after first login.

## Vault key names (hardcoded in seed-admin and the API)

Both `seed-admin` and the main `cluo` binary share these constants (defined in `internal/app/container/infrastructure.go` and mirrored in `cmd/seed-admin/main.go`):

| Config key | Value |
|---|---|
| `KEKAlias` | `cluo-encryption-key` |
| `PepperAlias` | `cluo` |
| `DBPath` | `data/.encx` |

The Vault transit key (`cluo-encryption-key`) and KV secret (`cluo` pepper) must exist in Vault before the seed command is run. These are provisioned by the Vault setup step in the infrastructure playbook (or manually before first deploy on the homelab VPS).

## Role string

The admin user is stored with role string `administrator` (from `identity.Administrator.String()`). This is the value decrypted at login and returned by `GET /auth/me`. The desktop app's dev mock uses the string `admin` — that is a local-dev shortcut only and does not match the production role string.
