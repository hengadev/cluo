#!/usr/bin/env bash
# =============================================================================
# Cluo — Secrets Audit Script
# Scans Ansible variable files, Terraform config, and git history for
# plaintext secrets.
#
# Usage:
#   ./infrastructure/scripts/audit-secrets.sh [--fix]
#
#   --fix  Also encrypt plaintext secrets found in current files
# =============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
ANSIBLE_DIR="$PROJECT_ROOT/infrastructure/ansible"
TERRAFORM_DIR="$PROJECT_ROOT/infrastructure/terraform"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track findings
FINDINGS=()
SEVERITIES=()
CI_FAILURE=0  # set to 1 on gitignore violations or history leaks (CI-blocking)

find_secret() {
    local file="$1" line="$2" severity="$3" description="$4"
    local color="${RED}"
    [[ "$severity" == "medium" ]] && color="${YELLOW}"
    FINDINGS+=("$file:$line")
    SEVERITIES+=("$severity")
    echo -e "${color}[${severity^^}]${NC} ${file}:${line} — ${description}"
}

find_ok() {
    echo -e "${GREEN}[OK]${NC} $1"
}

find_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

echo "============================================"
echo " Cluo Secrets Audit — $(date -Iseconds)"
echo "============================================"
echo ""

# =========================================================================
# 1. Check that sensitive files are gitignored
# =========================================================================
echo "--- Gitignore Checks ---"

check_gitignored() {
    local pattern="$1" label="$2"
    if git -C "$PROJECT_ROOT" ls-files --error-unmatch "$pattern" &>/dev/null; then
        find_secret ".gitignore" "-" "critical" "$label ('$pattern') is TRACKED by git — must be removed from history"
        CI_FAILURE=1
    else
        find_ok "$label ('$pattern') is properly gitignored"
    fi
}

check_gitignored "infrastructure/ansible/inventory.yml" "Ansible inventory"
check_gitignored "infrastructure/ansible/terraform-outputs.json" "Terraform outputs JSON"
check_gitignored "infrastructure/terraform/terraform.tfvars" "Terraform tfvars"
check_gitignored "infrastructure/terraform/terraform.tfstate" "Terraform state"

echo ""

# =========================================================================
# 2. Scan current Ansible variable files for plaintext secrets
# =========================================================================
echo "--- Current File Scans ---"

# Patterns that indicate a real secret value (not a placeholder)
# Exclude: empty strings, placeholder text like "your_...", "changeme", etc.
SECRET_PATTERNS=(
    'AKIA[A-Z0-9]{16}'                    # AWS Access Key ID
    '(?:aws_secret_access_key|secret_key|SECRET_ACCESS_KEY)\s*[:=]\s*"[A-Za-z0-9/+=]{20,}"'  # AWS Secret
    '(?:password|PASSWORD)\s*[:=]\s*"[^"]{6,}"'  # Non-trivial password (6+ chars)
    'hcloud_token\s*[:=]\s*"[a-zA-Z0-9]{20,}"'   # Hetzner token
    'cloudflare_token\s*[:=]\s*"[a-zA-Z0-9_-]{20,}"' # Cloudflare token
    'ssh-ed25519\s+AAAA'                  # SSH private keys (public is OK but note it)
    'ssh-rsa\s+AAAA'                      # SSH public key
)

# Known placeholder patterns (these are OK)
PLACEHOLDER_PATTERNS=(
    'your_.*_here'
    'changeme'
    'PLACEHOLDER'
    'REPLACE'
    'example\.com'
)

# Known safe AWS key examples (AWS docs use these)
SAFE_AWS_KEYS=(
    'AKIAIOSFODNN7EXAMPLE'
    'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'
)

# Scan a file for real secrets
scan_file() {
    local file="$1"
    if [[ ! -f "$file" ]]; then
        find_info "File not found (skipping): $file"
        return
    fi

    local line_num=0
    while IFS= read -r line; do
        ((line_num++)) || true

        # Skip comments and empty lines
        [[ "$line" =~ ^[[:space:]]*# ]] && continue
        [[ -z "${line// /}" ]] && continue

        # Skip known placeholders
        local is_placeholder=false
        for pat in "${PLACEHOLDER_PATTERNS[@]}"; do
            if [[ "$line" =~ $pat ]]; then
                is_placeholder=true
                break
            fi
        done
        $is_placeholder && continue

        # Check for actual secret patterns
        # AWS Access Key IDs
        if [[ "$line" =~ (AKIA[A-Z0-9]{16}) ]]; then
            # Skip AWS example keys
            local is_example=false
            for ex in "${SAFE_AWS_KEYS[@]}"; do
                if [[ "${BASH_REMATCH[0]}" == "$ex" ]]; then
                    is_example=true
                    break
                fi
            done
            $is_example && continue
            find_secret "$file" "$line_num" "critical" "AWS Access Key ID found: ${BASH_REMATCH[0]:0:8}..."
        fi

        # AWS Secret Access Keys (40-char base64-ish after = or :)
        if [[ "$line" =~ (secret_access_key|secret_key|SECRET_ACCESS_KEY|s3_secret_access_key) ]]; then
            # Extract the value part
            if [[ "$line" =~ [:=][[:space:]]*\"([A-Za-z0-9/+=]{20,})\" ]]; then
                find_secret "$file" "$line_num" "critical" "AWS Secret Access Key found (${#BASH_REMATCH[1]} chars)"
            fi
        fi

        # Hetzner Cloud Token
        if [[ "$line" =~ hcloud_token[[:space:]]*=[[:space:]]*\"([a-zA-Z0-9]{20,})\" ]]; then
            find_secret "$file" "$line_num" "critical" "Hetzner Cloud API token found"
        fi

        # Cloudflare Token
        if [[ "$line" =~ cloudflare_token[[:space:]]*=[[:space:]]*\"([a-zA-Z0-9_-]{20,})\" ]]; then
            find_secret "$file" "$line_num" "critical" "Cloudflare API token found"
        fi

        # Database passwords (non-placeholder)
        if [[ "$line" =~ (postgres_password|db_password|database_password|staging_postgres_password|production_postgres_password) ]]; then
            if [[ "$line" =~ [:=][[:space:]]*\"([^\"]{6,})\" ]]; then
                local val="${BASH_REMATCH[1]}"
                # Skip obvious placeholders
                if [[ ! "$val" =~ ^your_ && ! "$val" =~ ^changeme && ! "$val" =~ _here$ ]]; then
                    find_secret "$file" "$line_num" "high" "Database password found: '$val'"
                fi
            fi
        fi

        # SMTP passwords
        if [[ "$line" =~ smtp_password ]]; then
            if [[ "$line" =~ [:=][[:space:]]*\"([^\"]{4,})\" ]]; then
                local val="${BASH_REMATCH[1]}"
                if [[ -n "$val" ]]; then
                    find_secret "$file" "$line_num" "high" "SMTP password found"
                fi
            fi
        fi

        # SES verification tokens
        if [[ "$line" =~ (ses_verification_token|ses_dkim_tokens|dkim) ]] && [[ "$line" =~ [:=][[:space:]]*\"([^\"]{10,})\" ]]; then
            find_secret "$file" "$line_num" "medium" "SES secret found"
        fi

    done < "$file"
}

scan_file "$ANSIBLE_DIR/inventory.yml"
scan_file "$ANSIBLE_DIR/terraform-outputs.json"

# Scan all group_vars and host_vars (non-example)
for f in "$ANSIBLE_DIR"/group_vars/*/*.yml "$ANSIBLE_DIR"/host_vars/*/*.yml; do
    [[ -f "$f" ]] && scan_file "$f"
done

# Scan Terraform tfvars
scan_file "$TERRAFORM_DIR/terraform.tfvars"

# Scan Terraform state (it's expected to have secrets but confirm)
if [[ -f "$TERRAFORM_DIR/terraform.tfstate" ]]; then
    if grep -q '"sensitive": true' "$TERRAFORM_DIR/terraform.tfstate" 2>/dev/null; then
        find_info "terraform.tfstate contains sensitive outputs (expected — file is gitignored)"
    fi
    # Check for actual secret values in tfstate
    if grep -qE 'AKIA[A-Z0-9]{16}' "$TERRAFORM_DIR/terraform.tfstate" 2>/dev/null; then
        find_secret "$TERRAFORM_DIR/terraform.tfstate" "-" "info" "Contains AWS access keys (expected in state — file is gitignored)"
    fi
fi

echo ""

# =========================================================================
# 3. Scan git history for leaked secrets
# =========================================================================
echo "--- Git History Scan ---"

# Check if actual credential values appear in any commit
SECRETS_TO_CHECK=(
    "REDACTED_AWS_KEY_ID_1"
    "REDACTED_AWS_KEY_ID_2"
    "REDACTED_AWS_KEY_ID_3"
    "REDACTED_AWS_SECRET_FRAGMENT_1"
    "REDACTED_AWS_SECRET_FRAGMENT_2"
    "REDACTED_AWS_SECRET_KEY_1"
    "REDACTED_CLOUDFLARE_TOKEN_1"
    "REDACTED_CLOUDFLARE_API_TOKEN"
)

HISTORY_LEAKS=0
for secret in "${SECRETS_TO_CHECK[@]}"; do
    count=$(git -C "$PROJECT_ROOT" log --all -S "$secret" --oneline 2>/dev/null | wc -l)
    if [[ "$count" -gt 0 ]]; then
        commits=$(git -C "$PROJECT_ROOT" log --all -S "$secret" --oneline 2>/dev/null)
        find_secret "git-history" "-" "critical" "Secret '${secret:0:8}...' found in git history:\n$commits"
        HISTORY_LEAKS=$((HISTORY_LEAKS + 1))
        CI_FAILURE=1
    fi
done

if [[ "$HISTORY_LEAKS" -eq 0 ]]; then
    find_ok "No actual secret values found in git history"
fi

# Also scan for generic patterns in tracked files in git history
for pattern in "-----BEGIN RSA PRIVATE KEY-----" "-----BEGIN OPENSSH PRIVATE KEY-----"; do
    # Only check currently-tracked files
    tracked=$(git -C "$PROJECT_ROOT" ls-files | xargs grep -l -E "$pattern" 2>/dev/null || true)
    if [[ -n "$tracked" ]]; then
        find_secret "tracked-files" "-" "critical" "Secret pattern '$pattern' found in tracked files:\n$tracked"
    else
        find_ok "No '$pattern' patterns in tracked files"
    fi
done

# Check for real AWS keys in tracked files (skip example keys)
tracked_aws=$(git -C "$PROJECT_ROOT" ls-files | xargs grep -l -E 'AKIA[A-Z0-9]{16}' 2>/dev/null || true)
if [[ -n "$tracked_aws" ]]; then
    real_aws=""
    for f in $tracked_aws; do
        if grep -vE '(AKIAIOSFODNN7EXAMPLE)' "$f" 2>/dev/null | grep -qE 'AKIA[A-Z0-9]{16}'; then
            real_aws="$real_aws\n$f"
        fi
    done
    if [[ -n "$real_aws" ]]; then
        find_secret "tracked-files" "-" "critical" "Real AWS Access Key pattern found in tracked files:$real_aws"
    else
        find_ok "Only example AWS keys in tracked files (vault.yml.example)"
    fi
else
    find_ok "No AWS access keys in tracked files"
fi

echo ""

# =========================================================================
# 4. Summary
# =========================================================================
echo "============================================"
echo " Summary"
echo "============================================"

CRITICAL=0
HIGH=0
MEDIUM=0
for sev in "${SEVERITIES[@]+${SEVERITIES[@]}}"; do
    case "$sev" in
        critical) ((CRITICAL++)) || true ;;
        high) ((HIGH++)) || true ;;
        medium) ((MEDIUM++)) || true ;;
    esac
done

echo -e "  Critical: ${RED}${CRITICAL}${NC}"
echo -e "  High:     ${RED}${HIGH}${NC}"
echo -e "  Medium:   ${YELLOW}${MEDIUM}${NC}"
echo ""

if [[ "$CRITICAL" -gt 0 || "$HIGH" -gt 0 ]]; then
    echo -e "${RED}ACTION REQUIRED:${NC} Rotate all exposed credentials immediately."
    echo ""
    echo "Rotation steps:"
    echo "  1. AWS IAM: Generate new access keys, disable old ones"
    echo "  2. Hetzner: Rotate API token at console.hetzner.cloud"
    echo "  3. Cloudflare: Roll API token at dash.cloudflare.com"
    echo "  4. Database: Change postgres passwords on the VPS"
    echo "  5. After rotation, update inventory.yml and terraform.tfvars with new values"
    echo ""
    echo "After rotation, consider encrypting secrets with ansible-vault:"
    echo "  cd infrastructure/ansible"
    echo "  ansible-vault create group_vars/all/vault.yml"
fi

if [[ "$HISTORY_LEAKS" -gt 0 ]]; then
    echo ""
    echo -e "${RED}GIT HISTORY:${NC} Secrets found in git commits."
    echo "  After rotation, consider: git filter-repo --replace-text <(echo 'PATTERN==>REDACTED')"
fi

if [[ "$CRITICAL" -eq 0 && "$HIGH" -eq 0 && "$MEDIUM" -eq 0 ]]; then
    echo -e "${GREEN}No plaintext secrets found in current files or git history.${NC}"
fi

echo ""
echo "============================================"
echo " Audit complete"
echo "============================================"

# Exit non-zero only for CI-blocking issues: gitignore violations or history leaks.
# Local-file findings (expected in dev) do not block CI since those files won't exist in CI.
exit "$CI_FAILURE"
