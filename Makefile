.PHONY: help \
        dev dev-down dev-logs \
        build-desktop-linux-local install-desktop-linux \
        generate-signing-key \
        build-staging build-staging-api build-staging-web build-staging-mobile \
        build-prod    build-prod-api    build-prod-web    build-prod-mobile \
        push-staging  push-staging-api  push-staging-web  push-staging-mobile \
        push-prod     push-prod-api     push-prod-web     push-prod-mobile \
        release-staging release release-desktop release-desktop-linux \
        release-desktop-staging release-desktop-linux-staging \
        restart-staging restart-staging-api restart-staging-web restart-staging-mobile \
        restart-prod    restart-prod-api    restart-prod-web    restart-prod-mobile \
        deploy-staging deploy-staging-api deploy-staging-web deploy-staging-mobile \
        deploy-prod    deploy-prod-api    deploy-prod-web    deploy-prod-mobile \
        _check-vps

# Environment: staging or production (default). Controls S3 prefix for manifest + binaries.
ENV ?= production
# Dry-run mode: validates version bump, build, and signing without uploading to S3.
DRY_RUN ?=

DESKTOP_S3_BUCKET := cluo-assets-prod

ifeq ($(ENV),staging)
DESKTOP_S3_PREFIX := staging/desktop
else
DESKTOP_S3_PREFIX := desktop
endif

DESKTOP_MANIFEST_URL := https://$(DESKTOP_S3_BUCKET).s3.eu-central-1.amazonaws.com/$(DESKTOP_S3_PREFIX)/manifest.json
RELEASE_NOTES        ?=

# Signing keys for manifest verification
# Private key file (hex-encoded Ed25519). Generate with: make generate-signing-key
SIGNING_KEY_FILE     := $(HOME)/.config/cluo/signing-private.key
PUBLIC_KEY_FILE      := $(HOME)/.config/cluo/signing-public.key
# Public key is read from file for ldflags injection
PUBLIC_KEY           := $(shell cat $(PUBLIC_KEY_FILE) 2>/dev/null || echo "")
SIGN_MANIFEST        := cd cluo_desktop && go run ./cmd/sign-manifest

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
# Signing key management
# =============================================================================

generate-signing-key: ## Generate Ed25519 key pair for manifest signing
	@mkdir -p $(dir $(SIGNING_KEY_FILE))
	$(SIGN_MANIFEST) genkey $(SIGNING_KEY_FILE) $(PUBLIC_KEY_FILE)
	@echo ""
	@echo "==> Add this to your build ldflags:"
	@echo "    -X cluo_desktop/updater.PublicKey=$$(cat $(PUBLIC_KEY_FILE))"

# =============================================================================
# Desktop release (Windows binary → S3 → signed manifest.json)
# Prerequisites: mingw-w64, wails CLI, AWS CLI, jq, signing key
# Usage: make release-desktop VERSION=1.2.0
#        make release-desktop VERSION=1.2.0 RELEASE_NOTES="Fix X, add Y"
#        make release-desktop VERSION=1.2.0 ENV=staging
#        make release-desktop VERSION=1.2.0 DRY_RUN=1
# =============================================================================

release-desktop: ## Build and release cluo_desktop for Windows — VERSION=x.y.z [ENV=production|staging] [DRY_RUN=1]
	@test -n "$(VERSION)" || (echo "ERROR: VERSION is required. Usage: make release-desktop VERSION=1.0.0"; exit 1)
	@test -f $(SIGNING_KEY_FILE) || (echo "ERROR: Signing key not found. Run 'make generate-signing-key' first."; exit 1)
	@echo "==> [$(ENV)] Step 1/5: Bumping version in wails.json..."
	cd cluo_desktop && jq --arg v "$(VERSION)" '.version = $$v' wails.json > wails.json.tmp && mv wails.json.tmp wails.json
	@echo "==> [$(ENV)] Step 2/5: Building cluo_desktop v$(VERSION) for Windows (mingw cross-compile)..."
	cd cluo_desktop && CC=x86_64-w64-mingw32-gcc wails build -platform windows/amd64 \
		-ldflags "-X cluo_desktop/updater.Version=$(VERSION) -X cluo_desktop/updater.ManifestURL=$(DESKTOP_MANIFEST_URL) -X cluo_desktop/updater.PublicKey=$(PUBLIC_KEY)"
	@test -f cluo_desktop/build/bin/cluo_desktop.exe || (echo "ERROR: wails build did not produce binary"; exit 1)
	@set -e; \
	BINARY=cluo_desktop/build/bin/cluo_desktop.exe; \
	CHECKSUM=$$(sha256sum $$BINARY | awk '{print "sha256:"$$1}'); \
	DOWNLOAD_URL=https://$(DESKTOP_S3_BUCKET).s3.eu-central-1.amazonaws.com/$(DESKTOP_S3_PREFIX)/v$(VERSION)/cluo_desktop_windows_amd64.exe; \
	echo "==> [$(ENV)] Step 3/5: Checksum: $$CHECKSUM"; \
	printf '{\n  "version": "%s",\n  "release_notes": "%s",\n  "downloads": {\n    "windows_amd64": "%s"\n  },\n  "checksums": {\n    "windows_amd64": "%s"\n  }\n}\n' \
		"$(VERSION)" "$(RELEASE_NOTES)" "$$DOWNLOAD_URL" "$$CHECKSUM" \
		> /tmp/cluo_desktop_manifest.json; \
	($(SIGN_MANIFEST) sign /tmp/cluo_desktop_manifest.json $(SIGNING_KEY_FILE)); \
	if [ "$(DRY_RUN)" = "1" ]; then \
		echo "[DRY RUN] Step 4/5: Skipping binary upload to s3://$(DESKTOP_S3_BUCKET)/$(DESKTOP_S3_PREFIX)/v$(VERSION)/cluo_desktop_windows_amd64.exe"; \
		echo "[DRY RUN] Step 5/5: Skipping manifest upload to s3://$(DESKTOP_S3_BUCKET)/$(DESKTOP_S3_PREFIX)/manifest.json"; \
		echo "[DRY RUN] Release v$(VERSION) validated (no uploads). Manifest URL: $(DESKTOP_MANIFEST_URL)"; \
	else \
		echo "==> [$(ENV)] Step 4/5: Uploading binary to S3..."; \
		aws s3 cp $$BINARY \
			s3://$(DESKTOP_S3_BUCKET)/$(DESKTOP_S3_PREFIX)/v$(VERSION)/cluo_desktop_windows_amd64.exe \
			--content-type application/octet-stream; \
		echo "==> [$(ENV)] Step 5/5: Uploading signed manifest to S3..."; \
		aws s3 cp /tmp/cluo_desktop_manifest.json \
			s3://$(DESKTOP_S3_BUCKET)/$(DESKTOP_S3_PREFIX)/manifest.json \
			--content-type application/json \
			--cache-control "no-cache,no-store,must-revalidate"; \
		echo "==> [$(ENV)] Released v$(VERSION). Manifest: $(DESKTOP_MANIFEST_URL)"; \
	fi

# =============================================================================
# Desktop release — Linux (native build, no cross-compiler required)
# Merges linux_amd64 into the existing manifest, preserving other platforms.
# Prerequisites: wails CLI, AWS CLI, jq, signing key
# Usage: make release-desktop-linux VERSION=1.2.0
#        make release-desktop-linux VERSION=1.2.0 ENV=staging
#        make release-desktop-linux VERSION=1.2.0 DRY_RUN=1
# =============================================================================

release-desktop-linux: ## Build and release cluo_desktop for Linux — VERSION=x.y.z [ENV=production|staging] [DRY_RUN=1]
	@test -n "$(VERSION)" || (echo "ERROR: VERSION is required. Usage: make release-desktop-linux VERSION=1.0.0"; exit 1)
	@test -f $(SIGNING_KEY_FILE) || (echo "ERROR: Signing key not found. Run 'make generate-signing-key' first."; exit 1)
	@echo "==> [$(ENV)] Step 1/5: Bumping version in wails.json..."
	cd cluo_desktop && jq --arg v "$(VERSION)" '.version = $$v' wails.json > wails.json.tmp && mv wails.json.tmp wails.json
	@echo "==> [$(ENV)] Step 2/5: Building cluo_desktop v$(VERSION) for Linux (native)..."
	cd cluo_desktop && wails build -platform linux/amd64 -tags webkit2_41 \
		-ldflags "-X cluo_desktop/updater.Version=$(VERSION) -X cluo_desktop/updater.ManifestURL=$(DESKTOP_MANIFEST_URL) -X cluo_desktop/updater.PublicKey=$(PUBLIC_KEY)"
	@test -f cluo_desktop/build/bin/cluo_desktop || (echo "ERROR: wails build did not produce binary"; exit 1)
	@set -e; \
	BINARY=cluo_desktop/build/bin/cluo_desktop; \
	CHECKSUM=$$(sha256sum $$BINARY | awk '{print "sha256:"$$1}'); \
	DOWNLOAD_URL=https://$(DESKTOP_S3_BUCKET).s3.eu-central-1.amazonaws.com/$(DESKTOP_S3_PREFIX)/v$(VERSION)/cluo_desktop_linux_amd64; \
	echo "==> [$(ENV)] Step 3/5: Checksum: $$CHECKSUM"; \
	if [ "$(DRY_RUN)" = "1" ]; then \
		printf '{"downloads":{},"checksums":{}}' > /tmp/cluo_current_manifest.json; \
	else \
		aws s3 cp s3://$(DESKTOP_S3_BUCKET)/$(DESKTOP_S3_PREFIX)/manifest.json /tmp/cluo_current_manifest.json 2>/dev/null \
			|| printf '{"downloads":{},"checksums":{}}' > /tmp/cluo_current_manifest.json; \
	fi; \
	jq --arg ver "$(VERSION)" --arg notes "$(RELEASE_NOTES)" \
		--arg url "$$DOWNLOAD_URL" --arg cs "$$CHECKSUM" \
		'.version = $$ver | .release_notes = $$notes | .downloads.linux_amd64 = $$url | .checksums.linux_amd64 = $$cs | del(.signature)' \
		/tmp/cluo_current_manifest.json > /tmp/cluo_desktop_manifest.json; \
	($(SIGN_MANIFEST) sign /tmp/cluo_desktop_manifest.json $(SIGNING_KEY_FILE)); \
	if [ "$(DRY_RUN)" = "1" ]; then \
		echo "[DRY RUN] Step 4/5: Skipping binary upload to s3://$(DESKTOP_S3_BUCKET)/$(DESKTOP_S3_PREFIX)/v$(VERSION)/cluo_desktop_linux_amd64"; \
		echo "[DRY RUN] Step 5/5: Skipping manifest upload to s3://$(DESKTOP_S3_BUCKET)/$(DESKTOP_S3_PREFIX)/manifest.json"; \
		echo "[DRY RUN] Release v$(VERSION) validated (no uploads). Manifest URL: $(DESKTOP_MANIFEST_URL)"; \
	else \
		echo "==> [$(ENV)] Step 4/5: Uploading binary to S3..."; \
		aws s3 cp $$BINARY \
			s3://$(DESKTOP_S3_BUCKET)/$(DESKTOP_S3_PREFIX)/v$(VERSION)/cluo_desktop_linux_amd64 \
			--content-type application/octet-stream; \
		echo "==> [$(ENV)] Step 5/5: Uploading signed manifest to S3..."; \
		aws s3 cp /tmp/cluo_desktop_manifest.json \
			s3://$(DESKTOP_S3_BUCKET)/$(DESKTOP_S3_PREFIX)/manifest.json \
			--content-type application/json \
			--cache-control "no-cache,no-store,must-revalidate"; \
		echo "==> [$(ENV)] Released v$(VERSION) for linux/amd64. Manifest: $(DESKTOP_MANIFEST_URL)"; \
	fi

# =============================================================================
# Desktop release — staging convenience targets
# These delegate to release-desktop / release-desktop-linux with ENV=staging.
# =============================================================================

release-desktop-staging: ## Build and release cluo_desktop for Windows (staging) — VERSION=x.y.z [DRY_RUN=1]
	$(MAKE) release-desktop VERSION="$(VERSION)" ENV=staging RELEASE_NOTES="$(RELEASE_NOTES)" DRY_RUN="$(DRY_RUN)"

release-desktop-linux-staging: ## Build and release cluo_desktop for Linux (staging) — VERSION=x.y.z [DRY_RUN=1]
	$(MAKE) release-desktop-linux VERSION="$(VERSION)" ENV=staging RELEASE_NOTES="$(RELEASE_NOTES)" DRY_RUN="$(DRY_RUN)"

# =============================================================================
# Local development
# =============================================================================

dev: ## Start all services locally
	docker compose up -d

dev-down: ## Stop local services
	docker compose down

dev-logs: ## Stream local service logs
	docker compose logs -f

build-desktop-linux-local: ## Build cluo_desktop for Linux locally (no S3 upload) — VERSION=x.y.z optional
	@echo "==> Building cluo_desktop for Linux (local, no release)..."
	cd cluo_desktop && wails build -platform linux/amd64 -tags webkit2_41 \
		-ldflags "-X cluo_desktop/updater.Version=$(or $(VERSION),0.0.0-local) -X cluo_desktop/updater.ManifestURL=$(DESKTOP_MANIFEST_URL) -X cluo_desktop/updater.PublicKey=$(PUBLIC_KEY)"
	@echo "==> Binary: cluo_desktop/build/bin/cluo_desktop"

install-desktop-linux: ## Install cluo_desktop for current user (binary + icon + .desktop entry for Wofi/launchers)
	@bash cluo_desktop/build/linux/install.sh

# =============================================================================
# Build
# =============================================================================

build-staging: ## Build all images tagged :staging
	docker build -t $(API_IMAGE):staging    ./cluo_api
	docker build -t $(WEB_IMAGE):staging    ./cluo_web
	docker build --build-arg PUBLIC_APP_ENV=staging -t $(MOBILE_IMAGE):staging ./cluo_mobile

build-staging-api: ## Build cluo-api image tagged :staging
	docker build -t $(API_IMAGE):staging ./cluo_api

build-staging-web: ## Build cluo-web image tagged :staging
	docker build -t $(WEB_IMAGE):staging ./cluo_web

build-staging-mobile: ## Build cluo-mobile image tagged :staging
	docker build --build-arg PUBLIC_APP_ENV=staging -t $(MOBILE_IMAGE):staging ./cluo_mobile

build-prod: ## Build all images tagged :latest
	docker build -t $(API_IMAGE):latest    ./cluo_api
	docker build -t $(WEB_IMAGE):latest    ./cluo_web
	docker build --build-arg PUBLIC_APP_ENV=production -t $(MOBILE_IMAGE):latest ./cluo_mobile

build-prod-api: ## Build cluo-api image tagged :latest
	docker build -t $(API_IMAGE):latest ./cluo_api

build-prod-web: ## Build cluo-web image tagged :latest
	docker build -t $(WEB_IMAGE):latest ./cluo_web

build-prod-mobile: ## Build cluo-mobile image tagged :latest
	docker build --build-arg PUBLIC_APP_ENV=production -t $(MOBILE_IMAGE):latest ./cluo_mobile

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

release: ## Build, sign, and publish everything — VERSION=x.y.z [ENV=production|staging] [DRY_RUN=1]
	@test -n "$(VERSION)" || (echo "ERROR: VERSION is required. Usage: make release VERSION=1.0.0"; exit 1)
	@echo "==> Atomic release for v$(VERSION) [$(ENV)]"
	@echo "==> Step 1/4: Building Docker images..."
	$(MAKE) build-prod
	@echo "==> Step 2/4: Pushing Docker images..."
	$(MAKE) push-prod
	@echo "==> Step 3/4: Building desktop binary and publishing manifest..."
	$(MAKE) release-desktop VERSION=$(VERSION) RELEASE_NOTES="$(RELEASE_NOTES)" ENV=$(ENV) DRY_RUN=$(DRY_RUN)
	@echo "==> Step 4/4: Complete!"
	@echo ""
	@echo "==> Release v$(VERSION) [$(ENV)] published:"
	@echo "    Docker images: $(API_IMAGE):latest, $(WEB_IMAGE):latest, $(MOBILE_IMAGE):latest"
	@echo "    Desktop manifest: $(DESKTOP_MANIFEST_URL)"
	@echo ""
	@echo "==> Don't forget: make deploy-prod to restart server containers"

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
