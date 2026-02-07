# CLUO Terraform Makefile
# Usage: make <target>

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
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Usage:$(NC) make [target]"
	@echo "$(YELLOW)Example:$(NC) make plan"

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

state-rm: ## Remove resource from state (usage: make state-rm RESOURCE=resource.name)
	@if [ -z "$(RESOURCE)" ]; then \
		echo "$(RED)Error: RESOURCE variable not set$(NC)"; \
		echo "Usage: make state-rm RESOURCE=resource.name"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Removing $(RESOURCE) from state...$(NC)"
	terraform state rm $(RESOURCE)
	@echo "$(GREEN)✓ Resource removed$(NC)"

state-mv: ## Move resource in state (usage: make state-mv FROM=old.name TO=new.name)
	@if [ -z "$(FROM)" ] || [ -z "$(TO)" ]; then \
		echo "$(RED)Error: FROM or TO variable not set$(NC)"; \
		echo "Usage: make state-mv FROM=old.name TO=new.name"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Moving $(FROM) to $(TO)...$(NC)"
	terraform state mv $(FROM) $(TO)
	@echo "$(GREEN)✓ Resource moved$(NC)"

taint: ## Taint a resource (usage: make taint RESOURCE=resource.name)
	@if [ -z "$(RESOURCE)" ]; then \
		echo "$(RED)Error: RESOURCE variable not set$(NC)"; \
		echo "Usage: make taint RESOURCE=resource.name"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Tainting $(RESOURCE)...$(NC)"
	terraform taint $(RESOURCE)
	@echo "$(GREEN)✓ Resource tainted$(NC)"

untaint: ## Untaint a resource (usage: make untaint RESOURCE=resource.name)
	@if [ -z "$(RESOURCE)" ]; then \
		echo "$(RED)Error: RESOURCE variable not set$(NC)"; \
		echo "Usage: make untaint RESOURCE=resource.name"; \
		exit 1; \
	fi
	@echo "$(BLUE)Untainting $(RESOURCE)...$(NC)"
	terraform untaint $(RESOURCE)
	@echo "$(GREEN)✓ Resource untainted$(NC)"

force-unlock: ## Force unlock Terraform state (usage: make force-unlock LOCK_ID=xxx)
	@if [ -z "$(LOCK_ID)" ]; then \
		echo "$(RED)Error: LOCK_ID variable not set$(NC)"; \
		echo "Usage: make force-unlock LOCK_ID=xxx"; \
		exit 1; \
	fi
	@echo "$(RED)Force unlocking state with $(LOCK_ID)...$(NC)"
	terraform force-unlock $(LOCK_ID)
	@echo "$(GREEN)✓ State unlocked$(NC)"

import: ## Import existing resource (usage: make import ADDR=resource.name ID=resource-id)
	@if [ -z "$(ADDR)" ] || [ -z "$(ID)" ]; then \
		echo "$(RED)Error: ADDR or ID variable not set$(NC)"; \
		echo "Usage: make import ADDR=resource.name ID=resource-id"; \
		exit 1; \
	fi
	@echo "$(BLUE)Importing $(ADDR) with ID $(ID)...$(NC)"
	terraform import $(ADDR) $(ID)
	@echo "$(GREEN)✓ Resource imported$(NC)"

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

docs: ## Generate Terraform documentation
	@echo "$(BLUE)Generating documentation...$(NC)"
	@terraform-docs markdown ./ > README.md 2>/dev/null || echo "terraform-docs not installed, skipping..."

version: ## Show Terraform version
	@terraform version

test: ## Run Terraform tests (if configured)
	@echo "$(BLUE)Running Terraform tests...$(NC)"
	terraform test $(TF_ARGS)

# CI/CD targets
ci: fmt validate ## Run CI checks (format + validate)

# Production targets
prod-plan: ## Production plan (with confirmation)
	@read -p "Confirming production environment. Continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		$(MAKE) plan; \
	else \
		echo "Aborted"; \
	fi

prod-apply: ## Production apply (with confirmation)
	@read -p "Confirming production environment. Continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		$(MAKE) apply; \
	else \
		echo "Aborted"; \
	fi

# Save outputs for Ansible
save-outputs: output-json ## Save outputs to file for Ansible
	@echo "$(BLUE)Saving outputs for Ansible...$(NC)"
	@terraform output -json > ../ansible/terraform-outputs.json
	@echo "$(GREEN)✓ Outputs saved to ../ansible/terraform-outputs.json$(NC)"
