# CLUO Infrastructure Makefile
# Orchestrates Terraform (provisioning) and Ansible (configuration)
# Usage: make <target>

.PHONY: help provision configure setup teardown status clean

# Default target
.DEFAULT_GOAL := help

# Directories
TERRAFORM_DIR := terraform
ANSIBLE_DIR := ansible

# Variables (passed through to sub-Makefiles)
AUTO_APPROVE ?= 0
TAGS ?=
LIMIT ?=

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)CLUO Infrastructure Commands$(NC)"
	@echo ""
	@echo "$(YELLOW)Workflow:$(NC) provision -> configure -> deploy"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Sub-commands:$(NC)"
	@echo "  $(GREEN)make terraform-help $(NC) Show Terraform targets"
	@echo "  $(GREEN)make ansible-help   $(NC) Show Ansible targets"
	@echo ""
	@echo "$(YELLOW)Full setup:$(NC) make setup"
	@echo "$(YELLOW)Quick deploy:$(NC) make deploy"

# ============================================================================
# FULL WORKFLOW
# ============================================================================

setup: provision configure ## Full setup: provision infrastructure + configure server
	@echo "$(GREEN)✓ Full infrastructure setup complete$(NC)"

teardown: ## Destroy all infrastructure (with confirmation)
	@echo "$(RED)This will destroy ALL infrastructure resources!$(NC)"
	@read -p "Are you absolutely sure? Type 'yes' to confirm: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		$(MAKE) -C $(TERRAFORM_DIR) destroy; \
	else \
		echo "Aborted"; \
	fi

# ============================================================================
# TERRAFORM (PROVISIONING)
# ============================================================================

provision: ## Provision cloud resources (Terraform apply)
	@echo "$(BLUE)Provisioning infrastructure...$(NC)"
	@$(MAKE) -C $(TERRAFORM_DIR) init
	@$(MAKE) -C $(TERRAFORM_DIR) apply
	@$(MAKE) -C $(TERRAFORM_DIR) save-outputs
	@echo "$(GREEN)✓ Infrastructure provisioned$(NC)"

provision-plan: ## Preview infrastructure changes (Terraform plan)
	@$(MAKE) -C $(TERRAFORM_DIR) init
	@$(MAKE) -C $(TERRAFORM_DIR) plan

provision-destroy: teardown ## Alias for teardown

terraform-%: ## Run any Terraform target (e.g., make terraform-plan)
	@$(MAKE) -C $(TERRAFORM_DIR) $*

terraform-help: ## Show Terraform help
	@$(MAKE) -C $(TERRAFORM_DIR) help

# ============================================================================
# ANSIBLE (CONFIGURATION)
# ============================================================================

configure: ## Configure server (Ansible full run)
	@echo "$(BLUE)Configuring server...$(NC)"
	@$(MAKE) -C $(ANSIBLE_DIR) deploy
	@echo "$(GREEN)✓ Server configured$(NC)"

configure-check: ## Dry run server configuration
	@$(MAKE) -C $(ANSIBLE_DIR) check

deploy: ## Deploy application only (Ansible app role)
	@$(MAKE) -C $(ANSIBLE_DIR) deploy-app

deploy-security: ## Deploy security configuration only
	@$(MAKE) -C $(ANSIBLE_DIR) deploy-security

ansible-%: ## Run any Ansible target (e.g., make ansible-ping)
	@$(MAKE) -C $(ANSIBLE_DIR) $*

ansible-help: ## Show Ansible help
	@$(MAKE) -C $(ANSIBLE_DIR) help

# ============================================================================
# OPERATIONS
# ============================================================================

status: ## Check infrastructure and application status
	@echo "$(BLUE)Infrastructure status:$(NC)"
	@echo ""
	@echo "$(YELLOW)Terraform state:$(NC)"
	@$(MAKE) -C $(TERRAFORM_DIR) state-list 2>/dev/null || echo "  No state found (run 'make provision' first)"
	@echo ""
	@echo "$(YELLOW)Server status:$(NC)"
	@$(MAKE) -C $(ANSIBLE_DIR) status 2>/dev/null || echo "  Cannot reach server (check inventory)"

ping: ## Test connectivity to hosts
	@$(MAKE) -C $(ANSIBLE_DIR) ping

logs: ## View application logs
	@$(MAKE) -C $(ANSIBLE_DIR) logs

restart: ## Restart application services
	@$(MAKE) -C $(ANSIBLE_DIR) restart

backup: ## Trigger backup on remote host
	@$(MAKE) -C $(ANSIBLE_DIR) backup

# ============================================================================
# VALIDATION & CI
# ============================================================================

validate: ## Validate both Terraform and Ansible configs
	@echo "$(BLUE)Validating all configurations...$(NC)"
	@$(MAKE) -C $(TERRAFORM_DIR) validate
	@$(MAKE) -C $(ANSIBLE_DIR) validate
	@echo "$(GREEN)✓ All configurations valid$(NC)"

lint: ## Lint both Terraform and Ansible
	@echo "$(BLUE)Linting all configurations...$(NC)"
	@$(MAKE) -C $(TERRAFORM_DIR) fmt-check
	@$(MAKE) -C $(ANSIBLE_DIR) lint
	@echo "$(GREEN)✓ Lint complete$(NC)"

ci: validate lint ## Run all CI checks

# ============================================================================
# INIT & CLEANUP
# ============================================================================

init: ## Initialize both Terraform and Ansible
	@echo "$(BLUE)Initializing infrastructure tools...$(NC)"
	@$(MAKE) -C $(TERRAFORM_DIR) init
	@$(MAKE) -C $(ANSIBLE_DIR) init
	@echo "$(GREEN)✓ Initialization complete$(NC)"

clean: ## Clean temporary files in both directories
	@$(MAKE) -C $(TERRAFORM_DIR) clean
	@$(MAKE) -C $(ANSIBLE_DIR) clean
	@echo "$(GREEN)✓ Clean complete$(NC)"

version: ## Show tool versions
	@$(MAKE) -C $(TERRAFORM_DIR) version
	@echo ""
	@$(MAKE) -C $(ANSIBLE_DIR) version
