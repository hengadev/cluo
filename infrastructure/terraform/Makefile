# =============================================================================
# Cluo Terraform Makefile
# Single VPS Architecture: Staging + Production
# =============================================================================

.PHONY: help init validate plan apply destroy refresh fmt output upgrade clean graph show

# Default target
.DEFAULT_GOAL := help

# Variables
TF_ARGS ?=
AUTO_APPROVE ?= 0

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)CLUO Terraform Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Usage:$(NC) make [target]"
	@echo "$(YELLOW)Example:$(NC) make plan"

# =============================================================================
# Core Terraform Commands
# =============================================================================

init: ## Initialize Terraform working directory
	@echo "$(BLUE)Initializing Terraform...$(NC)"
	terraform init $(TF_ARGS)
	@echo "$(GREEN)✓ Init complete$(NC)"

validate: ## Validate Terraform configuration
	@echo "$(BLUE)Validating configuration...$(NC)"
	terraform validate
	@echo "$(GREEN)✓ Configuration is valid$(NC)"

fmt: ## Format Terraform files
	@echo "$(BLUE)Formatting Terraform files...$(NC)"
	terraform fmt -recursive
	@echo "$(GREEN)✓ Files formatted$(NC)"

fmt-check: ## Check if Terraform files are formatted
	@terraform fmt -recursive -check
	@echo "$(GREEN)✓ Files are properly formatted$(NC)"

plan: validate ## Show Terraform execution plan
	@echo "$(BLUE)Creating execution plan...$(NC)"
	terraform plan $(TF_ARGS)

plan-out: validate ## Create saved execution plan
	@echo "$(BLUE)Creating execution plan...$(NC)"
	terraform plan -out=tfplan $(TF_ARGS)
	@echo "$(GREEN)✓ Plan saved to tfplan$(NC)"

apply: plan ## Apply Terraform changes (interactive)
	@echo "$(YELLOW)Applying Terraform changes...$(NC)"
	terraform apply $(TF_ARGS)
	@echo "$(GREEN)✓ Apply complete$(NC)"

apply-auto: validate ## Apply Terraform changes (auto-approve)
	@echo "$(RED)Auto-applying Terraform changes...$(NC)"
	@if [ "$(AUTO_APPROVE)" = "1" ]; then \
		terraform apply -auto-approve $(TF_ARGS); \
	else \
		read -p "Are you sure? [y/N] " confirm; \
		if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
			terraform apply -auto-approve $(TF_ARGS); \
		else \
			echo "Aborted"; \
		fi \
	fi
	@echo "$(GREEN)✓ Apply complete$(NC)"

apply-plan: plan-out ## Apply saved execution plan
	@echo "$(YELLOW)Applying saved plan...$(NC)"
	terraform apply tfplan
	@echo "$(GREEN)✓ Apply complete$(NC)"

destroy: ## Destroy Terraform resources (interactive)
	@echo "$(RED)Destroying Terraform resources...$(NC)"
	terraform destroy $(TF_ARGS)
	@echo "$(GREEN)✓ Destroy complete$(NC)"

destroy-auto: ## Destroy Terraform resources (auto-approve)
	@echo "$(RED)Auto-destroying Terraform resources...$(NC)"
	@if [ "$(AUTO_APPROVE)" = "1" ]; then \
		terraform destroy -auto-approve $(TF_ARGS); \
	else \
		read -p "Are you sure? This will delete all resources! [y/N] " confirm; \
		if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
			terraform destroy -auto-approve $(TF_ARGS); \
		else \
			echo "Aborted"; \
		fi \
	fi
	@echo "$(GREEN)✓ Destroy complete$(NC)"

refresh: ## Refresh Terraform state
	@echo "$(BLUE)Refreshing state...$(NC)"
	terraform refresh $(TF_ARGS)
	@echo "$(GREEN)✓ Refresh complete$(NC)"

output: ## Show Terraform outputs
	@echo "$(BLUE)Terraform outputs:$(NC)"
	@terraform output $(TF_ARGS)

output-json: ## Show Terraform outputs as JSON
	@echo "$(BLUE)Terraform outputs (JSON):$(NC)"
	@terraform output -json $(TF_ARGS)

show: ## Show Terraform state
	@echo "$(BLUE)Terraform state:$(NC)"
	terraform show $(TF_ARGS)

graph: ## Generate Terraform dependency graph
	@echo "$(BLUE)Generating graph...$(NC)"
	terraform graph -type=terraform | dot -Tpng > graph.png
	@echo "$(GREEN)✓ Graph saved to graph.png$(NC)"

upgrade: ## Upgrade Terraform providers and modules
	@echo "$(BLUE)Upgrading providers and modules...$(NC)"
	terraform init -upgrade $(TF_ARGS)
	@echo "$(GREEN)✓ Upgrade complete$(NC)"

providers: ## Show Terraform providers
	@echo "$(BLUE)Terraform providers:$(NC)"
	terraform providers

state-list: ## List resources in state
	@echo "$(BLUE)Resources in state:$(NC)"
	terraform state list

clean: ## Remove Terraform cache and temporary files
	@echo "$(YELLOW)Cleaning Terraform files...$(NC)"
	@rm -rf .terraform
	@rm -rf .terraform.lock.hcl
	@rm -f tfplan
	@rm -f graph.png
	@rm -f *.tfstate.backup
	@rm -rf .terraform-backup
	@echo "$(GREEN)✓ Clean complete$(NC)"

vars: ## Show example .tfvars file location
	@echo "$(BLUE)Configuration file:$(NC)"
	@if [ -f terraform.tfvars ]; then \
		echo "$(GREEN)✓ terraform.tfvars exists$(NC)"; \
	else \
		echo "$(YELLOW)✗ terraform.tfvars not found$(NC)"; \
		echo "Copy terraform.tfvars.example to get started:"; \
		echo "  cp terraform.tfvars.example terraform.tfvars"; \
	fi

version: ## Show Terraform version
	@terraform version

test: ## Run Terraform tests (if configured)
	@echo "$(BLUE)Running Terraform tests...$(NC)"
	terraform test $(TF_ARGS)

ci: fmt validate ## Run CI checks (format + validate)

# =============================================================================
# Single VPS Architecture Targets
# =============================================================================

ssh: ## SSH into the VPS
	@echo "$(BLUE)Connecting to VPS...$(NC)"
	@ssh root@$$(terraform output -raw vps_ipv4_address)

ssh-staging: ## SSH into VPS for staging environment
	@echo "$(BLUE)Connecting to VPS (staging)...$(NC)"
	@ssh root@$$(terraform output -raw vps_ipv4_address) -t "cd /opt/cluo && bash"

ssh-production: ## SSH into VPS for production environment
	@echo "$(BLUE)Connecting to VPS (production)...$(NC)"
	@ssh root@$$(terraform output -raw vps_ipv4_address) -t "cd /opt/cluo && bash"

urls: ## Show all environment URLs
	@echo "$(BLUE)Staging URLs:$(NC)"
	@terraform output -json staging_urls | jq -r '"  API:    \(.api)\n  Web:    \(.web)\n  Mobile: \(.mobile)\n  Assets: \(.assets)"'
	@echo ""
	@echo "$(BLUE)Production URLs:$(NC)"
	@terraform output -json production_urls | jq -r '"  API:    \(.api)\n  Web:    \(.web)\n  WWW:    \(.www)\n  Mobile: \(.mobile)\n  Assets: \(.assets)"'

dkim: ## Show SES DKIM tokens for DNS verification
	@echo "$(BLUE)SES DKIM Tokens (add to Cloudflare DNS):$(NC)"
	@terraform output -json ses_dkim_tokens | jq -r '.[] | "  \(.)._domainkey CNAME \(.).dkim.amazonses.com"'

buckets: ## Show all S3 bucket names
	@echo "$(BLUE)S3 Buckets:$(NC)"
	@echo "  Assets Staging:    $$(terraform output -json assets_buckets | jq -r '.staging')"
	@echo "  Assets Production: $$(terraform output -json assets_buckets | jq -r '.production')"
	@echo "  Backups Staging:    $$(terraform output -json backup_buckets | jq -r '.staging')"
	@echo "  Backups Production: $$(terraform output -json backup_buckets | jq -r '.production')"
	@echo "  Terraform State:    $$(terraform output -json terraform_state_bucket | jq -r '.')"
	@echo "  Vault Storage:      $$(terraform output -json vault_storage_bucket | jq -r '.')"

vault-keys: ## Show Vault IAM credentials (sensitive)
	@echo "$(RED)Vault S3 Access Credentials:$(NC)"
	@echo "  KMS Key ID:    $$(terraform output -json vault_kms_key_id | jq -r '.')"
	@echo "  Access Key ID: $$(terraform output -json vault_access_key_id | jq -r '.')"
	@echo "  Secret Key:    $$(terraform output -json vault_secret_access_key | jq -r '.')"

status: ## Show server and service status
	@echo "$(BLUE)Checking server status...$(NC)"
	@ssh root@$$(terraform output -raw vps_ipv4_address) "systemctl status caddy --no-pager | head -5"
	@echo ""
	@echo "$(BLUE)Docker containers:$(NC)"
	@ssh root@$$(terraform output -raw vps_ipv4_address) "docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null || echo 'Docker not accessible'"

logs-caddy: ## Show Caddy reverse proxy logs
	@echo "$(BLUE)Caddy logs:$(NC)"
	@ssh root@$$(terraform output -raw vps_ipv4_address) "journalctl -u caddy -n 50 --no-pager"

logs-staging: ## Show staging Docker logs
	@echo "$(BLUE)Staging logs:$(NC)"
	@ssh root@$$(terraform output -raw vps_ipv4_address) "cd /opt/cluo && docker-compose logs --tail=50 -f 2>/dev/null || echo 'Docker Compose not accessible'"

restart-all: ## Restart all Docker services
	@echo "$(YELLOW)Restarting all services...$(NC)"
	@ssh root@$$(terraform output -raw vps_ipv4_address) "cd /opt/cluo && docker-compose restart 2>/dev/null || echo 'Docker Compose not accessible'"
	@echo "$(GREEN)✓ Services restarted$(NC)"

# =============================================================================
# Deployment Targets
# =============================================================================

deploy: ## Deploy application to VPS (both environments)
	@echo "$(YELLOW)Deploying to VPS...$(NC)"
	@ssh root@$$(terraform output -raw vps_ipv4_address) "cd /opt/cluo && git pull && docker-compose up -d --build 2>/dev/null || echo 'Deployment failed'"
	@echo "$(GREEN)✓ Deployment complete$(NC)"

deploy-staging: ## Deploy staging environment only
	@echo "$(YELLOW)Deploying staging...$(NC)"
	@ssh root@$$(terraform output -raw vps_ipv4_address) "cd /opt/cluo && git pull && docker-compose up -d --build api_staging web_staging mobile_staging 2>/dev/null || echo 'Deployment failed'"
	@echo "$(GREEN)✓ Staging deployment complete$(NC)"

deploy-production: ## Deploy production environment only
	@echo "$(YELLOW)Deploying production...$(NC)"
	@read -p "Confirming PRODUCTION deployment. Continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		ssh root@$$(terraform output -raw vps_ipv4_address) "cd /opt/cluo && git pull && docker-compose up -d --build api_production web_production mobile_production 2>/dev/null || echo 'Deployment failed'"; \
		echo "$(GREEN)✓ Production deployment complete$(NC)"; \
	else \
		echo "Aborted"; \
	fi

# =============================================================================
# Backup Targets
# =============================================================================

backup-staging: ## Trigger staging database backup
	@echo "$(YELLOW)Backing up staging database...$(NC)"
	@ssh root@$$(terraform output -raw vps_ipv4_address) "cd /opt/cluo && docker-compose exec postgres_staging pg_dump -U cluo cluo_staging | gzip > /tmp/staging_backup_$$(date +%Y%m%d_%H%M%S).sql.gz 2>/dev/null || echo 'Backup failed'"
	@echo "$(GREEN)✓ Staging backup complete$(NC)"

backup-production: ## Trigger production database backup
	@echo "$(YELLOW)Backing up production database...$(NC)"
	@read -p "Confirming production backup. Continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		ssh root@$$(terraform output -raw vps_ipv4_address) "cd /opt/cluo && docker-compose exec postgres_production pg_dump -U cluo cluo_production | gzip > /tmp/production_backup_$$(date +%Y%m%d_%H%M%S).sql.gz 2>/dev/null || echo 'Backup failed'"; \
		echo "$(GREEN)✓ Production backup complete$(NC)"; \
	else \
		echo "Aborted"; \
	fi

# =============================================================================
# Production Targets
# =============================================================================

prod-plan: ## Production plan (with confirmation)
	@read -p "Confirming PRODUCTION infrastructure changes. Continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		$(MAKE) plan; \
	else \
		echo "Aborted"; \
	fi

prod-apply: ## Production apply (with confirmation)
	@read -p "Confirming PRODUCTION infrastructure changes. Continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		$(MAKE) apply; \
	else \
		echo "Aborted"; \
	fi
