# CLUO Ansible Makefile
# Usage: make <target>

.PHONY: help init validate syntax-check lint run deploy check dry-run diff install-galaxy

# Default target
.DEFAULT_GOAL := help

# Variables
INVENTORY ?= inventory.yml
PLAYBOOK ?= site.yml
VAULT_FILE ?= group_vars/vault.yml
ASK_VAULT ?= true
TAGS ?=
SKIP_TAGS ?=
LIMIT ?=
VERBOSE ?=
extra_vars ?=

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Build ansible-playbook command
ANSIBLE_PLAYBOOK = ansible-playbook $(PLAYBOOK) -i $(INVENTORY)
ifdef TAGS
	ANSIBLE_PLAYBOOK += --tags $(TAGS)
endif
ifdef SKIP_TAGS
	ANSIBLE_PLAYBOOK += --skip-tags $(SKIP_TAGS)
endif
ifdef LIMIT
	ANSIBLE_PLAYBOOK += --limit $(LIMIT)
endif
ifdef VERBOSE
	ANSIBLE_PLAYBOOK += $(VERBOSE)
endif
ifdef extra_vars
	ANSIBLE_PLAYBOOK += $(extra_vars)
endif

help: ## Show this help message
	@echo "$(BLUE)CLUO Ansible Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Usage:$(NC) make [target]"
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make deploy              # Full deployment"
	@echo "  make deploy TAGS=docker  # Deploy only Docker role"
	@echo "  make run LIMIT=host1     # Run on specific host"
	@echo "  make check               # Dry run (check mode)"

init: ## Initialize Ansible environment
	@echo "$(BLUE)Initializing Ansible...$(NC)"
	@mkdir -p group_vars host_vars
	@if [ ! -f $(INVENTORY) ]; then \
		echo "$(YELLOW)Creating inventory from example...$(NC)"; \
		cp inventory.yml.example $(INVENTORY); \
	fi
	@if [ ! -f group_vars/all.yml ]; then \
		echo "$(YELLOW)Creating group_vars/all.yml from example...$(NC)"; \
		cp group_vars/all.yml.example group_vars/all.yml; \
	fi
	@echo "$(GREEN)✓ Init complete$(NC)"
	@echo "$(YELLOW)Edit $(INVENTORY) and group_vars/all.yml to configure$(NC)"

validate: ## Validate Ansible configuration
	@echo "$(BLUE)Validating configuration...$(NC)"
	@if [ ! -f $(INVENTORY) ]; then \
		echo "$(RED)Error: $(INVENTORY) not found$(NC)"; \
		echo "Run 'make init' first"; \
		exit 1; \
	fi
	@ansible-playbook $(PLAYBOOK) -i $(INVENTORY) --syntax-check
	@echo "$(GREEN)✓ Configuration is valid$(NC)"

syntax-check: validate ## Alias for validate

lint: ## Run Ansible lint
	@echo "$(BLUE)Running Ansible lint...$(NC)"
	@which ansible-lint > /dev/null || (echo "$(RED)ansible-lint not installed$(NC)" && exit 1)
	@ansible-lint $(PLAYBOOK) || true
	@echo "$(GREEN)✓ Lint complete$(NC)"

install-galaxy: ## Install Ansible collections from requirements (if requirements.yml exists)
	@echo "$(BLUE)Installing Ansible collections...$(NC)"
	@if [ -f requirements.yml ]; then \
		ansible-galaxy collection install -r requirements.yml; \
	else \
		echo "$(YELLOW)No requirements.yml found, skipping...$(NC)"; \
	fi

ping: ## Test connectivity to all hosts
	@echo "$(BLUE)Pinging hosts...$(NC)"
	@ansible all -i $(INVENTORY) -m ping

setup: ## Run Ansible setup module (gather facts)
	@echo "$(BLUE)Gathering facts...$(NC)"
	@ansible all -i $(INVENTORY) -m setup

check: validate ## Dry run (check mode)
	@echo "$(BLUE)Running in check mode...$(NC)"
	@$(ANSIBLE_PLAYBOOK) --check --diff

dry-run: check ## Alias for check

diff: validate ## Run with diff output
	@echo "$(BLUE)Running with diff output...$(NC)"
	@$(ANSIBLE_PLAYBOOK) --diff

run: validate ## Run playbook (default)
	@echo "$(BLUE)Running playbook...$(NC)"
	@$(ANSIBLE_PLAYBOOK)

deploy: validate ## Deploy application (full run)
	@echo "$(BLUE)Deploying CLUO...$(NC)"
	@$(ANSIBLE_PLAYBOOK)
	@echo "$(GREEN)✓ Deployment complete$(NC)"

# Quick deploy tags
deploy-security: validate ## Deploy security roles only
	@echo "$(BLUE)Deploying security roles...$(NC)"
	@$(ANSIBLE_PLAYBOOK) --tags security,hardening,ssh,firewall,fail2ban
	@echo "$(GREEN)✓ Security deployment complete$(NC)"

deploy-docker: validate ## Deploy Docker only
	@echo "$(BLUE)Deploying Docker...$(NC)"
	@$(ANSIBLE_PLAYBOOK) --tags docker
	@echo "$(GREEN)✓ Docker deployment complete$(NC)"

deploy-app: validate ## Deploy application only
	@echo "$(BLUE)Deploying application...$(NC)"
	@$(ANSIBLE_PLAYBOOK) --tags app,deploy,user
	@echo "$(GREEN)✓ Application deployment complete$(NC)"

deploy-monitoring: validate ## Deploy monitoring only
	@echo "$(BLUE)Deploying monitoring...$(NC)"
	@$(ANSIBLE_PLAYBOOK) --tags monitoring
	@echo "$(GREEN)✓ Monitoring deployment complete$(NC)"

deploy-backup: validate ## Deploy backup setup only
	@echo "$(BLUE)Setting up backups...$(NC)"
	@$(ANSIBLE_PLAYBOOK) --tags backup
	@echo "$(GREEN)✓ Backup setup complete$(NC)"

# Run with verbose output
verbose: validate ## Run with verbose output
	@echo "$(BLUE)Running with verbose output...$(NC)"
	@$(ANSIBLE_PLAYBOOK) -v

debug: validate ## Run with debug output (very verbose)
	@echo "$(BLUE)Running with debug output...$(NC)"
	@$(ANSIBLE_PLAYBOOK) -vvv

# Vault operations
vault-init: ## Initialize vault (create encrypted file)
	@echo "$(BLUE)Initializing vault...$(NC)"
	@if [ -f $(VAULT_FILE) ]; then \
		echo "$(YELLOW)Vault file already exists$(NC)"; \
		echo "Edit with: make vault-edit"; \
	else \
		cp group_vars/vault.yml.example $(VAULT_FILE); \
		ansible-vault encrypt $(VAULT_FILE); \
		echo "$(GREEN)✓ Vault initialized$(NC)"; \
	fi

vault-edit: ## Edit vault file
	@echo "$(BLUE)Editing vault...$(NC)"
	@if [ ! -f $(VAULT_FILE) ]; then \
		echo "$(RED)Error: $(VAULT_FILE) not found$(NC)"; \
		echo "Run 'make vault-init' first"; \
		exit 1; \
	fi
	@ansible-vault edit $(VAULT_FILE)

vault-view: ## View vault file
	@echo "$(BLUE)Viewing vault...$(NC)"
	@if [ ! -f $(VAULT_FILE) ]; then \
		echo "$(RED)Error: $(VAULT_FILE) not found$(NC)"; \
		exit 1; \
	fi
	@ansible-vault view $(VAULT_FILE)

vault-rekey: ## Rekey vault (change password)
	@echo "$(BLUE)Rekeying vault...$(NC)"
	@ansible-vault rekey $(VAULT_FILE)

vault-decrypt: ## Decrypt vault file (use with caution)
	@echo "$(RED)Decrypting vault file...$(NC)"
	@read -p "This will decrypt the vault file. Continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		ansible-vault decrypt $(VAULT_FILE); \
		echo "$(GREEN)✓ Vault decrypted$(NC)"; \
	else \
		echo "Aborted"; \
	fi

vault-encrypt: ## Encrypt vault file
	@echo "$(BLUE)Encrypting vault file...$(NC)"
	@ansible-vault encrypt $(VAULT_FILE)
	@echo "$(GREEN)✓ Vault encrypted$(NC)"

# Inventory operations
inventory-show: ## Show inventory
	@echo "$(BLUE)Inventory:$(NC)"
	@ansible-inventory -i $(INVENTORY) --list

inventory-graph: ## Show inventory graph
	@echo "$(BLUE)Inventory graph:$(NC)"
	@ansible-inventory -i $(INVENTORY) --graph

hosts: ## List all hosts
	@echo "$(BLUE)All hosts:$(NC)"
	@ansible all -i $(INVENTORY) --list-hosts

# Operations on remote hosts
status: ## Check service status on hosts
	@echo "$(BLUE)Checking service status...$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "cd /opt/cluo && docker compose ps" \
		--become

logs: ## View application logs
	@echo "$(BLUE)Fetching logs...$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "cd /opt/cluo && docker compose logs --tail=50" \
		--become

logs-follow: ## Follow application logs
	@echo "$(BLUE)Following logs...$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "cd /opt/cluo && docker compose logs -f" \
		--become

restart: ## Restart application services
	@echo "$(BLUE)Restarting services...$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "cd /opt/cluo && docker compose restart" \
		--become
	@echo "$(GREEN)✓ Services restarted$(NC)"

stop: ## Stop application services
	@echo "$(BLUE)Stopping services...$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "cd /opt/cluo && docker compose down" \
		--become
	@echo "$(GREEN)✓ Services stopped$(NC)"

start: ## Start application services
	@echo "$(BLUE)Starting services...$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "cd /opt/cluo && docker compose up -d" \
		--become
	@echo "$(GREEN)✓ Services started$(NC)"

update: ## Pull latest application code
	@echo "$(BLUE)Updating application...$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "cd /opt/cluo && git pull" \
		--become
	@echo "$(GREEN)✓ Code updated. Run 'make restart' to restart services$(NC)"

# Backup operations
backup: ## Trigger backup on remote host
	@echo "$(BLUE)Triggering backup...$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "cd /opt/cluo && ./backup.sh" \
		--become --become-user=cluo
	@echo "$(GREEN)✓ Backup triggered$(NC)"

backup-list: ## List backups on remote host
	@echo "$(BLUE)Listing backups...$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "ls -lh /opt/cluo/backups/" \
		--become

# System operations
update-system: ## Run system updates on hosts
	@echo "$(BLUE)Updating system packages...$(NC)"
	@ansible webservers -i $(INVENTORY) -m apt \
		-a "upgrade=dist update_cache=yes" \
		--become
	@echo "$(GREEN)✓ System updated$(NC)"

reboot: ## Reboot hosts (with confirmation)
	@echo "$(RED)Rebooting hosts...$(NC)"
	@read -p "Reboot all hosts? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		ansible webservers -i $(INVENTORY) -m reboot --become; \
		echo "$(GREEN)✓ Reboot initiated$(NC)"; \
	else \
		echo "Aborted"; \
	fi

uptime: ## Check host uptime
	@echo "$(BLUE)Host uptime:$(NC)"
	@ansible all -i $(INVENTORY) -m shell -a "uptime"

disk: ## Check disk usage
	@echo "$(BLUE)Disk usage:$(NC)"
	@ansible all -i $(INVENTORY) -m shell -a "df -h"

memory: ## Check memory usage
	@echo "$(BLUE)Memory usage:$(NC)"
	@ansible all -i $(INVENTORY) -m shell -a "free -h"

# Docker operations
docker-ps: ## List Docker containers
	@echo "$(BLUE)Docker containers:$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'" \
		--become

docker-stats: ## Show Docker stats
	@echo "$(BLUE)Docker stats:$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "docker stats --no-stream" \
		--become

docker-images: ## List Docker images
	@echo "$(BLUE)Docker images:$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "docker images --format 'table {{.Repository}}\t{{.Tag}}\t{{.Size}}'" \
		--become

docker-prune: ## Prune unused Docker resources
	@echo "$(BLUE)Pruning Docker resources...$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "docker system prune -f" \
		--become
	@echo "$(GREEN)✓ Prune complete$(NC)"

# Security operations
security-status: ## Check security status (fail2ban, ufw)
	@echo "$(BLUE)Security status:$(NC)"
	@echo ""
	@echo "$(YELLOW)UFW Status:$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "ufw status" --become
	@echo ""
	@echo "$(YELLOW)Fail2ban Status:$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "fail2ban-client status" --become

fail2ban-status: ## Show detailed Fail2ban status
	@echo "$(BLUE)Fail2ban status:$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "fail2ban-client status" --become
	@echo ""
	@echo "$(YELLOW)SSH Jail:$(NC)"
	@ansible webservers -i $(INVENTORY) -m shell \
		-a "fail2ban-client status sshd" --become

# Monitoring
metrics: ## Show basic metrics
	@echo "$(BLUE)System metrics:$(NC)"
	@ansible all -i $(INVENTORY) -m shell \
		-a "echo 'CPU: '; top -bn1 | grep 'Cpu(s)' | awk '{print $$2}' | cut -d'%' -f1; echo 'Memory: '; free | grep Mem | awk '{printf \"%.1f%%\n\", $$3/$$2 * 100.0}'; echo 'Disk: '; df -h / | tail -1 | awk '{print $$5}'"

# CI/CD targets
ci: validate ## Run CI checks
	@echo "$(BLUE)Running CI checks...$(NC)"
	@$(MAKE) validate
	@$(MAKE) lint
	@echo "$(GREEN)✓ CI checks passed$(NC)"

# Production targets
prod-deploy: validate ## Production deployment (with confirmation)
	@read -p "Confirming production environment. Continue? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		$(MAKE) deploy; \
	else \
		echo "Aborted"; \
	fi

# Documentation
docs: ## Generate Ansible documentation
	@echo "$(BLUE)Generating documentation...$(NC)"
	@ansible-doc --type role --list > /dev/null

# Version info
version: ## Show Ansible version
	@echo "$(BLUE)Ansible version:$(NC)"
	@ansible --version
	@echo ""
	@echo "$(BLUE)Python version:$(NC)"
	@python --version
	@echo ""
	@echo "$(BLUE)Installed collections:$(NC)"
	@ansible-galaxy collection list

clean: ## Clean temporary files
	@echo "$(YELLOW)Cleaning temporary files...$(NC)"
	@rm -f *.retry
	@rm -rf __pycache__
	@rm -rf .cache
	@find . -name "*.pyc" -delete
	@find . -name "*.pyo" -delete
	@echo "$(GREEN)✓ Clean complete$(NC)"
