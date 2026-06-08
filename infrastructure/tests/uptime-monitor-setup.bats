#!/usr/bin/env bats
# =============================================================================
# Bats tests for the UptimeRobot monitor setup logic.
#
# These tests exercise the UptimeRobot API interaction using a stub HTTP
# server so the monitor creation and idempotency logic can be verified
# without a real API key.
# =============================================================================

setup() {
    TEST_DIR="$(mktemp -d)"

    # We'll test the Ansible task logic indirectly by verifying the
    # API call patterns.  For unit-testing the monitor script we create
    # a small shell wrapper that mirrors the Ansible URI calls.
    STUB_DIR="$TEST_DIR/bin"
    mkdir -p "$STUB_DIR"

    # Capture HTTP requests
    REQUEST_LOG="$TEST_DIR/requests.log"

    # --- Stub curl that mimics UptimeRobot API responses ---
    cat > "$STUB_DIR/curl" <<'CURL_STUB'
#!/bin/bash
set -euo pipefail

# Parse the request body to determine the action
BODY=""
URL=""
NEXT_IS_DATA=false
for arg in "$@"; do
    if $NEXT_IS_DATA; then
        BODY="$arg"
        NEXT_IS_DATA=false
    fi
    case "$arg" in
        --data|-d) NEXT_IS_DATA=true ;;
        https://*) URL="$arg" ;;
    esac
done

# Log the request
echo "URL=$URL BODY=$BODY" >> "${REQUEST_LOG:?}"

case "$URL" in
    */getMonitors)
        # If search URL matches a known monitor, return it
        if echo "$BODY" | grep -q "api.clientvault.fr/health"; then
            cat <<'JSON'
{"stat":"ok","monitors":[{"id":12345,"friendly_name":"Cluo API /health","url":"https://api.clientvault.fr/health","interval":300}]}
JSON
        else
            cat <<'JSON'
{"stat":"ok","monitors":[]}
JSON
        fi
        ;;
    */newMonitor)
        cat <<'JSON'
{"stat":"ok","monitor":{"id":67890,"friendly_name":"Cluo API /health","url":"https://api.clientvault.fr/health","interval":300}}
JSON
        ;;
    *)
        echo '{"stat":"fail","error":{"message":"Unknown endpoint"}}' >&2
        exit 1
        ;;
esac
CURL_STUB
    chmod +x "$STUB_DIR/curl"

    export PATH="$STUB_DIR:$PATH"
    export REQUEST_LOG

    # --- Monitor setup script (mirrors Ansible logic) ---
    SETUP_SCRIPT="$TEST_DIR/setup-monitor.sh"
    cat > "$SETUP_SCRIPT" <<'SCRIPT'
#!/bin/bash
set -euo pipefail

API_KEY="${UPTIMEROBOT_API_KEY:?UPTIMEROBOT_API_KEY must be set}"
HEALTH_URL="${HEALTH_URL:-https://api.clientvault.fr/health}"
FRIENDLY_NAME="${FRIENDLY_NAME:-Cluo API /health}"
INTERVAL="${INTERVAL:-300}"
ALERT_EMAIL="${ALERT_EMAIL:-test@example.com}"
STATE_DIR="${STATE_DIR:?STATE_DIR must be set}"
RESULT_FILE="$STATE_DIR/monitor_result.txt"

# Step 1: Check if monitor exists
RESPONSE=$(curl -s --data "{\"api_key\":\"$API_KEY\",\"format\":\"json\",\"search\":\"$HEALTH_URL\"}" -H "Content-Type: application/json" "https://api.uptimerobot.com/v2/getMonitors")

# Extract existing monitor ID (simple grep-based)
MONITOR_ID=$(echo "$RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2 || true)

if [ -n "$MONITOR_ID" ] && echo "$RESPONSE" | grep -q "$HEALTH_URL"; then
    echo "monitor_exists=true" >> "$RESULT_FILE"
    echo "monitor_id=$MONITOR_ID" >> "$RESULT_FILE"
    exit 0
fi

# Step 2: Create new monitor
RESPONSE=$(curl -s --data "{\"api_key\":\"$API_KEY\",\"format\":\"json\",\"type\":1,\"url\":\"$HEALTH_URL\",\"friendly_name\":\"$FRIENDLY_NAME\",\"interval\":$INTERVAL,\"http_method\":\"GET\",\"alert_type\":2}" -H "Content-Type: application/json" "https://api.uptimerobot.com/v2/newMonitor")

MONITOR_ID=$(echo "$RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2 || true)

if [ -n "$MONITOR_ID" ]; then
    echo "monitor_created=true" >> "$RESULT_FILE"
    echo "monitor_id=$MONITOR_ID" >> "$RESULT_FILE"
else
    echo "monitor_failed=true" >> "$RESULT_FILE"
    exit 1
fi
SCRIPT
    chmod +x "$SETUP_SCRIPT"
}

teardown() {
    rm -rf "$TEST_DIR"
}

# ---------------------------------------------------------------------------
# Tests
# ---------------------------------------------------------------------------

@test "creates a new monitor when none exists" {
    # Use a URL not in the stub's "existing" list
    export UPTIMEROBOT_API_KEY="test-key-123"
    export HEALTH_URL="https://staging-api.clientvault.fr/health"
    export ALERT_EMAIL="test@example.com"
    export STATE_DIR="$TEST_DIR/state"
    mkdir -p "$STATE_DIR"

    run "$SETUP_SCRIPT"
    [ "$status" -eq 0 ]

    # Result file should show creation
    grep -q "monitor_created=true" "$STATE_DIR/monitor_result.txt"
    grep -q "monitor_id=67890" "$STATE_DIR/monitor_result.txt"

    # Should have made exactly 2 API calls: getMonitors + newMonitor
    CALL_COUNT=$(grep -c "URL=" "$REQUEST_LOG")
    [ "$CALL_COUNT" -eq 2 ]
}

@test "skips creation when monitor already exists" {
    export UPTIMEROBOT_API_KEY="test-key-123"
    export HEALTH_URL="https://api.clientvault.fr/health"
    export ALERT_EMAIL="test@example.com"
    export STATE_DIR="$TEST_DIR/state"
    mkdir -p "$STATE_DIR"

    run "$SETUP_SCRIPT"
    [ "$status" -eq 0 ]

    # Result file should show existing monitor
    grep -q "monitor_exists=true" "$STATE_DIR/monitor_result.txt"
    grep -q "monitor_id=12345" "$STATE_DIR/monitor_result.txt"

    # Should have made exactly 1 API call: getMonitors only
    CALL_COUNT=$(grep -c "URL=" "$REQUEST_LOG")
    [ "$CALL_COUNT" -eq 1 ]
}

@test "getMonitors request includes correct search URL" {
    export UPTIMEROBOT_API_KEY="test-key-123"
    export HEALTH_URL="https://api.clientvault.fr/health"
    export ALERT_EMAIL="test@example.com"
    export STATE_DIR="$TEST_DIR/state"
    mkdir -p "$STATE_DIR"

    "$SETUP_SCRIPT" || true

    grep -q "getMonitors" "$REQUEST_LOG"
    grep -q "api.clientvault.fr/health" "$REQUEST_LOG"
}

@test "newMonitor request includes correct parameters" {
    export UPTIMEROBOT_API_KEY="test-key-123"
    export HEALTH_URL="https://new-endpoint.example.com/health"
    export FRIENDLY_NAME="Test Monitor"
    export INTERVAL="600"
    export ALERT_EMAIL="test@example.com"
    export STATE_DIR="$TEST_DIR/state"
    mkdir -p "$STATE_DIR"

    "$SETUP_SCRIPT" || true

    # The newMonitor call should be present with type=1 (HTTP) and interval
    grep -q "newMonitor" "$REQUEST_LOG"
    grep -q '"type":1' "$REQUEST_LOG"
    grep -q '"interval":600' "$REQUEST_LOG"
}
