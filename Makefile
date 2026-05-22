.PHONY: help \
        dev dev-down dev-logs \
        build-staging build-staging-api build-staging-web build-staging-mobile \
        build-prod    build-prod-api    build-prod-web    build-prod-mobile \
        push-staging  push-staging-api  push-staging-web  push-staging-mobile \
        push-prod     push-prod-api     push-prod-web     push-prod-mobile \
        release-staging release \
        release-desktop \
        restart-staging restart-staging-api restart-staging-web restart-staging-mobile \
        restart-prod    restart-prod-api    restart-prod-web    restart-prod-mobile \
        deploy-staging deploy-staging-api deploy-staging-web deploy-staging-mobile \
        deploy-prod    deploy-prod-api    deploy-prod-web    deploy-prod-mobile \
        _check-vps

DESKTOP_MANIFEST_URL := https://cluo-assets-production.s3.eu-central-1.amazonaws.com/desktop/manifest.json
DESKTOP_S3_BUCKET    := cluo-assets-production
RELEASE_NOTES        ?=

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
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'

# =============================================================================
# Desktop release (Windows binary → S3 → manifest.json)
# Prerequisites: mingw-w64 (sudo apt install mingw-w64), wails CLI, AWS CLI
# Usage: make release-desktop VERSION=1.2.0
#        make release-desktop VERSION=1.2.0 RELEASE_NOTES="Fix X, add Y"
# =============================================================================

release-desktop: ## Build and release cluo_desktop for Windows — VERSION=x.y.z required
	@test -n "$(VERSION)" || (echo "ERROR: VERSION is required. Usage: make release-desktop VERSION=1.0.0"; exit 1)
	@echo "==> Building cluo_desktop v$(VERSION) for Windows (mingw cross-compile)..."
	cd cluo_desktop && CC=x86_64-w64-mingw32-gcc wails build -platform windows/amd64 \
		-ldflags "-X cluo_desktop/updater.Version=$(VERSION) -X cluo_desktop/updater.ManifestURL=$(DESKTOP_MANIFEST_URL)"
	@set -e; \
	BINARY=cluo_desktop/build/bin/cluo_desktop.exe; \
	CHECKSUM=$$(sha256sum $$BINARY | awk '{print "sha256:"$$1}'); \
	DOWNLOAD_URL=https://$(DESKTOP_S3_BUCKET).s3.eu-central-1.amazonaws.com/desktop/v$(VERSION)/cluo_desktop_windows_amd64.exe; \
	echo "==> Checksum: $$CHECKSUM"; \
	echo "==> Uploading binary to S3..."; \
	aws s3 cp $$BINARY \
		s3://$(DESKTOP_S3_BUCKET)/desktop/v$(VERSION)/cluo_desktop_windows_amd64.exe \
		--content-type application/octet-stream; \
	echo "==> Updating manifest.json..."; \
	printf '{\n  "version": "%s",\n  "release_notes": "%s",\n  "downloads": {\n    "windows_amd64": "%s"\n  },\n  "checksums": {\n    "windows_amd64": "%s"\n  }\n}\n' \
		"$(VERSION)" "$(RELEASE_NOTES)" "$$DOWNLOAD_URL" "$$CHECKSUM" \
		> /tmp/cluo_desktop_manifest.json; \
	aws s3 cp /tmp/cluo_desktop_manifest.json \
		s3://$(DESKTOP_S3_BUCKET)/desktop/manifest.json \
		--content-type application/json \
		--cache-control "no-cache,no-store,must-revalidate"; \
	echo "==> Released v$(VERSION). Manifest: $(DESKTOP_MANIFEST_URL)"

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

build-staging-api: ## Build cluo-api image tagged :staging
	docker build -t $(API_IMAGE):staging ./cluo_api

build-staging-web: ## Build cluo-web image tagged :staging
	docker build -t $(WEB_IMAGE):staging ./cluo_web

build-staging-mobile: ## Build cluo-mobile image tagged :staging
	docker build -t $(MOBILE_IMAGE):staging ./cluo_mobile

build-prod: ## Build all images tagged :latest
	docker build -t $(API_IMAGE):latest    ./cluo_api
	docker build -t $(WEB_IMAGE):latest    ./cluo_web
	docker build -t $(MOBILE_IMAGE):latest ./cluo_mobile

build-prod-api: ## Build cluo-api image tagged :latest
	docker build -t $(API_IMAGE):latest ./cluo_api

build-prod-web: ## Build cluo-web image tagged :latest
	docker build -t $(WEB_IMAGE):latest ./cluo_web

build-prod-mobile: ## Build cluo-mobile image tagged :latest
	docker build -t $(MOBILE_IMAGE):latest ./cluo_mobile

# =============================================================================
# Push to Docker Hub (henga/*)
# =============================================================================

push-staging: ## Push all :staging images to Docker Hub
	docker push $(API_IMAGE):staging
	docker push $(WEB_IMAGE):staging
	docker push $(MOBILE_IMAGE):staging

push-staging-api: ## Push cluo-api :staging to Docker Hub
	docker push $(API_IMAGE):staging

push-staging-web: ## Push cluo-web :staging to Docker Hub
	docker push $(WEB_IMAGE):staging

push-staging-mobile: ## Push cluo-mobile :staging to Docker Hub
	docker push $(MOBILE_IMAGE):staging

push-prod: ## Push all :latest images to Docker Hub
	docker push $(API_IMAGE):latest
	docker push $(WEB_IMAGE):latest
	docker push $(MOBILE_IMAGE):latest

push-prod-api: ## Push cluo-api :latest to Docker Hub
	docker push $(API_IMAGE):latest

push-prod-web: ## Push cluo-web :latest to Docker Hub
	docker push $(WEB_IMAGE):latest

push-prod-mobile: ## Push cluo-mobile :latest to Docker Hub
	docker push $(MOBILE_IMAGE):latest

# =============================================================================
# Release (build + push to Docker Hub)
# =============================================================================

release-staging: build-staging push-staging ## Build and push all :staging images to Docker Hub

release: build-prod push-prod ## Build and push all :latest images to Docker Hub

# =============================================================================
# Restart containers on VPS
# =============================================================================

_check-vps:
	@test -n "$(VPS_IP)" || (echo "ERROR: could not resolve VPS IP from Terraform. Run 'make init' in the homelab repo first."; exit 1)

restart-staging: _check-vps ## Pull all :staging images and restart staging on VPS
	$(VPS_SSH) "cd /opt/cluo-staging && docker compose pull && docker compose up -d --remove-orphans"

restart-staging-api: _check-vps ## Pull and restart cluo-staging-api on VPS
	$(VPS_SSH) "cd /opt/cluo-staging && docker compose pull cluo-staging-api && docker compose up -d --no-deps cluo-staging-api"

restart-staging-web: _check-vps ## Pull and restart cluo-staging-web on VPS
	$(VPS_SSH) "cd /opt/cluo-staging && docker compose pull cluo-staging-web && docker compose up -d --no-deps cluo-staging-web"

restart-staging-mobile: _check-vps ## Pull and restart cluo-staging-mobile on VPS
	$(VPS_SSH) "cd /opt/cluo-staging && docker compose pull cluo-staging-mobile && docker compose up -d --no-deps cluo-staging-mobile"

restart-prod: _check-vps ## Pull all :latest images and restart production on VPS
	$(VPS_SSH) "cd /opt/cluo && docker compose pull && docker compose up -d --remove-orphans"

restart-prod-api: _check-vps ## Pull and restart cluo-prod-api on VPS
	$(VPS_SSH) "cd /opt/cluo && docker compose pull cluo-prod-api && docker compose up -d --no-deps cluo-prod-api"

restart-prod-web: _check-vps ## Pull and restart cluo-prod-web on VPS
	$(VPS_SSH) "cd /opt/cluo && docker compose pull cluo-prod-web && docker compose up -d --no-deps cluo-prod-web"

restart-prod-mobile: _check-vps ## Pull and restart cluo-prod-mobile on VPS
	$(VPS_SSH) "cd /opt/cluo && docker compose pull cluo-prod-mobile && docker compose up -d --no-deps cluo-prod-mobile"

# =============================================================================
# Full deploy (build → push → restart)
# =============================================================================

deploy-staging: build-staging push-staging restart-staging ## Build, push and restart all staging end-to-end

deploy-staging-api: build-staging-api push-staging-api restart-staging-api ## Build, push and restart cluo-api staging

deploy-staging-web: build-staging-web push-staging-web restart-staging-web ## Build, push and restart cluo-web staging

deploy-staging-mobile: build-staging-mobile push-staging-mobile restart-staging-mobile ## Build, push and restart cluo-mobile staging

deploy-prod: build-prod push-prod restart-prod ## Build, push and restart all production end-to-end

deploy-prod-api: build-prod-api push-prod-api restart-prod-api ## Build, push and restart cluo-api production

deploy-prod-web: build-prod-web push-prod-web restart-prod-web ## Build, push and restart cluo-web production

deploy-prod-mobile: build-prod-mobile push-prod-mobile restart-prod-mobile ## Build, push and restart cluo-mobile production
