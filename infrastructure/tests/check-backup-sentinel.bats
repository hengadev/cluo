#!/usr/bin/env bats
# =============================================================================
# Bats tests for the backup sentinel alert logic.
#
# These tests exercise the check-backup.sh script's sentinel-detection
# behaviour in isolation.  The msmtp binary is replaced with a stub so
# alerts are captured to a file rather than actually sent.
# =============================================================================

setup() {
    # Isolated temp directories
    TEST_DIR="$(mktemp -d)"
    BACKUP_DIR="$TEST_DIR/backups"
    mkdir -p "$BACKUP_DIR"

    # Stub msmtp: capture mail to a file so we can assert on it
    STUB_DIR="$TEST_DIR/bin"
    mkdir -p "$STUB_DIR"
    MAIL_LOG="$TEST_DIR/mail.log"
    printf '#!/bin/sh\necho "$@" >> "%s"\ncat >> "%s"\n' "$MAIL_LOG" "$MAIL_LOG" > "$STUB_DIR/msmtp"
    chmod +x "$STUB_DIR/msmtp"
    export PATH="$STUB_DIR:$PATH"

    # Build a testable check-backup script with injected variables
    CHECK_SCRIPT="$TEST_DIR/check-backup.sh"
    cat > "$CHECK_SCRIPT" <<'SCRIPT'
#!/bin/bash
set -euo pipefail

# Backup alert check — testable version.
# Paths injected via env vars with sensible defaults.

BACKUP_DIR="${TEST_BACKUP_DIR:?TEST_BACKUP_DIR must be set}"
SENTINEL_FILE="$BACKUP_DIR/.last_success"
FAILURE_SENTINEL="$BACKUP_DIR/.last_failure"
MAX_AGE_SECONDS=$((25 * 3600))
ALERT_EMAIL="${TEST_ALERT_EMAIL:-root@localhost}"
HOSTNAME_STR="$(hostname -f 2>/dev/null || hostname)"

send_alert() {
  local subject="$1"
  local body="$2"
  printf "To: %s\nSubject: %s\n\n%s\n" "$ALERT_EMAIL" "$subject" "$body" | msmtp "$ALERT_EMAIL"
  echo "ALERT: $subject"
}

# Check for a recent failure sentinel first
if [ -f "$FAILURE_SENTINEL" ]; then
  FAILURE_INFO=$(cat "$FAILURE_SENTINEL")
  SUBJECT="[CLUO ALERT] Backup failed on $HOSTNAME_STR"
  BODY="A backup failure was detected on $HOSTNAME_STR.

Failure sentinel : $FAILURE_SENTINEL
Contents         :
$FAILURE_INFO

The last backup run exited with an error. Please investigate.

Timestamp (UTC): $(date -u +'%Y-%m-%dT%H:%M:%SZ')"
  send_alert "$SUBJECT" "$BODY"
  exit 1
fi

if [ ! -f "$SENTINEL_FILE" ]; then
  SUBJECT="[CLUO ALERT] Backup sentinel missing on $HOSTNAME_STR"
  BODY="The backup sentinel file was not found at $SENTINEL_FILE on host $HOSTNAME_STR.

This means either:
  - No backup has ever completed successfully, or
  - The sentinel file was deleted.

Please investigate immediately.

Timestamp (UTC): $(date -u +'%Y-%m-%dT%H:%M:%SZ')"
  send_alert "$SUBJECT" "$BODY"
  exit 1
fi

SENTINEL_TS=$(head -1 "$SENTINEL_FILE")
SENTINEL_EPOCH=$(date -d "$SENTINEL_TS" +%s 2>/dev/null || echo 0)
NOW_EPOCH=$(date +%s)
AGE_SECONDS=$(( NOW_EPOCH - SENTINEL_EPOCH ))

if [ "$AGE_SECONDS" -gt "$MAX_AGE_SECONDS" ]; then
  SUBJECT="[CLUO ALERT] Backup is stale on $HOSTNAME_STR"
  BODY="The last successful backup on $HOSTNAME_STR is stale.

Sentinel file : $SENTINEL_FILE
Last success  : $SENTINEL_TS
Age           : $(( AGE_SECONDS / 3600 )) hours
Threshold     : 25 hours

The backup may have failed silently. Please investigate.

Timestamp (UTC): $(date -u +'%Y-%m-%dT%H:%M:%SZ')"
  send_alert "$SUBJECT" "$BODY"
  exit 1
fi

# Sentinel is fresh — exit silently
exit 0
SCRIPT
    chmod +x "$CHECK_SCRIPT"
}

teardown() {
    rm -rf "$TEST_DIR"
}

# ---------------------------------------------------------------------------
# Tests
# ---------------------------------------------------------------------------

@test "missing sentinel triggers alert" {
    run env TEST_BACKUP_DIR="$BACKUP_DIR" "$CHECK_SCRIPT"
    [ "$status" -ne 0 ]
    grep -q "Backup sentinel missing" "$MAIL_LOG"
}

@test "fresh success sentinel — no alert" {
    date -u +'%Y-%m-%dT%H:%M:%SZ' > "$BACKUP_DIR/.last_success"
    run env TEST_BACKUP_DIR="$BACKUP_DIR" "$CHECK_SCRIPT"
    [ "$status" -eq 0 ]
    [ ! -s "$MAIL_LOG" ]
}

@test "stale success sentinel (26 hours old) triggers alert" {
    # Write a sentinel timestamp 26 hours in the past
    STALE_TS=$(date -u -d '26 hours ago' +'%Y-%m-%dT%H:%M:%SZ')
    echo "$STALE_TS" > "$BACKUP_DIR/.last_success"
    run env TEST_BACKUP_DIR="$BACKUP_DIR" "$CHECK_SCRIPT"
    [ "$status" -ne 0 ]
    grep -q "Backup is stale" "$MAIL_LOG"
}

@test "failure sentinel triggers alert even when success sentinel is fresh" {
    # Fresh success sentinel
    date -u +'%Y-%m-%dT%H:%M:%SZ' > "$BACKUP_DIR/.last_success"
    # But a failure sentinel also exists
    printf '%s\n%s\n' "$(date -u +'%Y-%m-%dT%H:%M:%SZ')" "exit_code=1" > "$BACKUP_DIR/.last_failure"

    run env TEST_BACKUP_DIR="$BACKUP_DIR" "$CHECK_SCRIPT"
    [ "$status" -ne 0 ]
    grep -q "Backup failed" "$MAIL_LOG"
}

@test "failure sentinel without success sentinel triggers failure alert" {
    printf '%s\n%s\n' "$(date -u +'%Y-%m-%dT%H:%M:%SZ')" "exit_code=1" > "$BACKUP_DIR/.last_failure"

    run env TEST_BACKUP_DIR="$BACKUP_DIR" "$CHECK_SCRIPT"
    [ "$status" -ne 0 ]
    # Should get the failure alert, not the missing-sentinel alert
    grep -q "Backup failed" "$MAIL_LOG"
}

@test "barely-fresh sentinel (24 hours old) does not trigger alert" {
    FRESH_TS=$(date -u -d '24 hours ago' +'%Y-%m-%dT%H:%M:%SZ')
    echo "$FRESH_TS" > "$BACKUP_DIR/.last_success"
    run env TEST_BACKUP_DIR="$BACKUP_DIR" "$CHECK_SCRIPT"
    [ "$status" -eq 0 ]
    [ ! -s "$MAIL_LOG" ]
}
