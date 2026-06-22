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

### Current setup (cluo's own VPS, Makefile-driven deploy)

Seeding is automatic — no manual `docker exec` needed. Create a root `.env`
(gitignored, never committed) with:

```bash
SEED_ADMIN_EMAIL=admin@example.com
SEED_ADMIN_PASSWORD=changeme123
```

`make seed-admin-staging` and `make seed-admin-prod` read that file and run
`./seed-admin` inside `cluo-staging-api`/`cluo-prod-api` over SSH. These two
targets are wired into the end of `restart-staging`, `restart-staging-api`,
`restart-prod`, and `restart-prod-api` (and therefore the `deploy-staging*`
and `deploy-prod*` targets that depend on them), so seeding happens
automatically on every deploy. Because `seed-admin` itself checks whether
the user already exists before creating one, re-running it on every deploy
is a no-op after the first time.

If `.env` doesn't exist, the targets print a warning and skip rather than
failing the deploy.

To seed manually (e.g. to debug), the underlying command is:

```bash
docker exec \
  -e SEED_ADMIN_EMAIL=admin@example.com \
  -e SEED_ADMIN_PASSWORD=changeme123 \
  cluo-staging-api ./seed-admin   # or cluo-prod-api
```

The command inherits all other required env vars from the running container.

## Vault key names (hardcoded in seed-admin and the API)

Both `seed-admin` and the main `cluo` binary share these constants (defined in `internal/app/container/infrastructure.go` and mirrored in `cmd/seed-admin/main.go`):

| Config key | Value |
|---|---|
| `KEKAlias` | `cluo-encryption-key` |
| `PepperAlias` | `cluo` |
| `DBPath` | `data/.encx` |

The Vault transit key (`cluo-encryption-key`) and KV-v2 secrets engine must exist before the seed command is run — see step 6 of `infrastructure/DEPLOY_RUNBOOK.md` (one-time, manual, run once per VPS). The pepper itself is generated and persisted by `encx` automatically on first successful API boot once those engines exist.

## Role string

The admin user is stored with role string `administrator` (from `identity.Administrator.String()`). This is the value decrypted at login and returned by `GET /auth/me`. The desktop app's dev mock uses the string `admin` — that is a local-dev shortcut only and does not match the production role string.
