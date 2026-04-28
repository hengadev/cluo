.PHONY: help \
        dev dev-down dev-logs \
        build-staging build-prod \
        push-staging push-prod \
        restart-staging restart-prod \
        deploy-staging deploy-prod \
        _check-vps

REGISTRY     := henga
API_IMAGE    := $(REGISTRY)/cluo-api
WEB_IMAGE    := $(REGISTRY)/cluo-web
MOBILE_IMAGE := $(REGISTRY)/cluo-mobile

# VPS connection — IP resolved from homelab Terraform output
HOMELAB_TF := $(HOME)/Documents/projects/homelab/terraform
VPS_IP     := $(shell cd $(HOMELAB_TF) && terraform output -raw server_ip 2>/dev/null || echo "")
SSH_KEY    := ~/.ssh/henga
SSH_USER   := deploy
VPS_SSH    := ssh -i $(SSH_KEY) $(SSH_USER)@$(VPS_IP)

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

# =============================================================================
# Local development
# =============================================================================

dev: ## Start all services locally
	docker compose up -d

dev-down: ## Stop local services
	docker compose down

dev-logs: ## Stream local service logs
	docker compose logs -f

# =============================================================================
# Build
# =============================================================================

build-staging: ## Build all images tagged :staging
	docker build -t $(API_IMAGE):staging    ./cluo_api
	docker build -t $(WEB_IMAGE):staging    ./cluo_web
	docker build -t $(MOBILE_IMAGE):staging ./cluo_mobile

build-prod: ## Build all images tagged :latest
	docker build -t $(API_IMAGE):latest    ./cluo_api
	docker build -t $(WEB_IMAGE):latest    ./cluo_web
	docker build -t $(MOBILE_IMAGE):latest ./cluo_mobile

# =============================================================================
# Push to Docker Hub (henga/*)
# =============================================================================

push-staging: ## Push :staging images to Docker Hub
	docker push $(API_IMAGE):staging
	docker push $(WEB_IMAGE):staging
	docker push $(MOBILE_IMAGE):staging

push-prod: ## Push :latest images to Docker Hub
	docker push $(API_IMAGE):latest
	docker push $(WEB_IMAGE):latest
	docker push $(MOBILE_IMAGE):latest

# =============================================================================
# Restart containers on VPS
# =============================================================================

_check-vps:
	@test -n "$(VPS_IP)" || (echo "ERROR: could not resolve VPS IP from Terraform. Run 'make init' in the homelab repo first."; exit 1)

restart-staging: _check-vps ## Pull :staging images and restart staging on VPS
	$(VPS_SSH) "cd /opt/cluo-staging && docker compose pull && docker compose up -d --remove-orphans"

restart-prod: _check-vps ## Pull :latest images and restart production on VPS
	$(VPS_SSH) "cd /opt/cluo && docker compose pull && docker compose up -d --remove-orphans"

# =============================================================================
# Full deploy (build → push → restart)
# =============================================================================

deploy-staging: build-staging push-staging restart-staging ## Build, push and restart staging end-to-end

deploy-prod: build-prod push-prod restart-prod ## Build, push and restart production end-to-end
