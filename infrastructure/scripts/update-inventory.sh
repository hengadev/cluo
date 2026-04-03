#!/bin/bash
# Update Ansible inventory from Terraform outputs

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INFRA_DIR="$(dirname "$SCRIPT_DIR")"
TF_DIR="$INFRA_DIR/terraform"
ANSIBLE_DIR="$INFRA_DIR/ansible"
OUTPUTS_FILE="$ANSIBLE_DIR/terraform-outputs.json"
INVENTORY_FILE="$ANSIBLE_DIR/inventory.yml"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}Updating Ansible inventory from Terraform outputs...${NC}"

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is not installed${NC}"
    echo "Install jq:"
    echo "  Ubuntu/Debian: sudo apt install jq"
    echo "  macOS: brew install jq"
    exit 1
fi

# Check if outputs file exists
if [ ! -f "$OUTPUTS_FILE" ]; then
    echo -e "${RED}Error: Terraform outputs file not found at $OUTPUTS_FILE${NC}"
    echo "Please run: cd $TF_DIR && terraform output -json > $OUTPUTS_FILE"
    exit 1
fi

# Check if inventory file exists
if [ ! -f "$INVENTORY_FILE" ]; then
    echo -e "${YELLOW}Warning: Inventory file not found. Copying from example...${NC}"
    cp "$ANSIBLE_DIR/inventory.yml.example" "$INVENTORY_FILE"
fi

# Extract values from Terraform outputs
VPS_IP=$(jq -r '.vps_ipv4_address.value // empty' "$OUTPUTS_FILE")
STAGING_ACCESS_KEY=$(jq -r '.staging_assets_iam_access_key.value // empty' "$OUTPUTS_FILE")
STAGING_SECRET_KEY=$(jq -r '.staging_assets_iam_secret_key.value // empty' "$OUTPUTS_FILE")
PRODUCTION_ACCESS_KEY=$(jq -r '.production_assets_iam_access_key.value // empty' "$OUTPUTS_FILE")
PRODUCTION_SECRET_KEY=$(jq -r '.production_assets_iam_secret_key.value // empty' "$OUTPUTS_FILE")

# Validate required values
if [ -z "$VPS_IP" ]; then
    echo -e "${RED}Error: Could not extract vps_ipv4_address from Terraform outputs${NC}"
    exit 1
fi

echo -e "${GREEN}Found values:${NC}"
echo "  Server IP: $VPS_IP"
echo "  Staging Access Key: ${STAGING_ACCESS_KEY:0:8}..."
echo "  Production Access Key: ${PRODUCTION_ACCESS_KEY:0:8}..."

# Update inventory.yml using sed (preserve SSH keys and other manual settings)
echo -e "${BLUE}Updating $INVENTORY_FILE${NC}"

# Update ansible_host
sed -i 's|ansible_host: ".*"  # Update with actual server IP from Terraform|ansible_host: "'"$VPS_IP"'"  # From Terraform output|' "$INVENTORY_FILE"

# Update S3 credentials
sed -i 's|staging_s3_access_key_id: ""  # From Terraform output: staging_assets_iam_access_key|staging_s3_access_key_id: "'"$STAGING_ACCESS_KEY"'"  # From Terraform output|' "$INVENTORY_FILE"
sed -i 's|staging_s3_secret_access_key: ""  # From Terraform output: staging_assets_iam_secret_key|staging_s3_secret_access_key: "'"$STAGING_SECRET_KEY"'"  # From Terraform output|' "$INVENTORY_FILE"
sed -i 's|production_s3_access_key_id: ""  # From Terraform output: production_assets_iam_access_key|production_s3_access_key_id: "'"$PRODUCTION_ACCESS_KEY"'"  # From Terraform output|' "$INVENTORY_FILE"
sed -i 's|production_s3_secret_access_key: ""  # From Terraform output: production_assets_iam_secret_key|production_s3_secret_access_key: "'"$PRODUCTION_SECRET_KEY"'"  # From Terraform output|' "$INVENTORY_FILE"

echo -e "${GREEN}✓ Inventory updated successfully${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "  1. Review the inventory: cat $INVENTORY_FILE"
echo "  2. Run Ansible configuration: make configure"
