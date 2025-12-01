# ============================================================================
# CLUO API COMMANDS
# ============================================================================
# Supports monolith architecture with domain-based organization

.PHONY: test-all test-unit test-integration test-coverage test-race test-benchmark test-verbose test-smoke test-clean test-help
.PHONY: test-unit-client test-unit-case test-unit-document test-unit-media test-unit-user
.PHONY: test-unit-common test-unit-infrastructure test-unit-application test-unit-interface
.PHONY: test-integration-client test-integration-case test-integration-document
.PHONY: test-domain-client test-domain-case test-domain-document test-domain-media test-domain-user
.PHONY: test-infrastructure-client test-infrastructure-case test-infrastructure-document test-infrastructure-media test-infrastructure-report test-infrastructure-token
.PHONY: encx-validate encx-generate encx-clean encx-build encx-help

# ============================================================================
# ENCRYPTION COMMANDS (ENCX)
# ============================================================================

encx-help: ## Show ENCRYPTION (ENCX) help message
	@echo 'ENCX Encryption Commands:'
	@echo ''
	@echo 'Available Commands:'
	@echo '  encx-install      Install ENCX CLI globally (recommended)'
	@echo '  encx-validate      Validate ENCX struct tags before generation'
	@echo '  encx-generate      Generate ENCX encryption code for all structs'
	@echo '  encx-build         Build ENCX generator tool locally (optional)'
	@echo '  encx-clean         Clean generated ENCX files (custom helper)'
	@echo ''
	@echo 'Development Workflow:'
	@echo '  1. Design structs with encx tags'
	@echo '  2. Run make encx-validate'
	@echo '  3. Run make encx-generate'
	@echo '  4. Test your changes'
	@echo '  5. Commit both source and generated code'
	@echo ''
	@echo 'Direct ENCX Commands (without make):'
	@echo '    go run ./cmd/encx-gen validate -v .'
	@echo '    go run ./cmd/encx-gen generate -v .'

encx-install: ## Install ENCX CLI globally (recommended method)
	@echo "Installing ENCX CLI globally..."
	@go install github.com/hengadev/encx/cmd/encx-gen@latest
	@echo "ENCX CLI installed globally!"
	@echo "You can now use: encx-gen validate -v ."

encx-build: ## Build ENCX generator tool locally (optional development method)
	@echo "Building ENCX generator tool locally..."
	@mkdir -p bin
	@go build -o bin/encx-gen github.com/hengadev/encx/cmd/encx-gen
	@echo "ENCX generator built: bin/encx-gen"
	@echo "You can now use: ./bin/encx-gen validate -v ."

encx-validate: ## Validate ENCX struct tags
	@echo "Validating ENCX struct tags..."
	@if command -v encx-gen >/dev/null 2>&1; then \
		echo "Using global encx-gen CLI..."; \
		encx-gen validate -v .; \
	elif [ -f "./bin/encx-gen" ]; then \
		echo "Using local encx-gen binary..."; \
		./bin/encx-gen validate -v .; \
	else \
		echo "Using go run (no installation needed)..."; \
		go run github.com/hengadev/encx/cmd/encx-gen validate -v .; \
	fi
	@echo "ENCX validation completed!"

encx-generate: ## Generate ENCX encryption code
	@echo "Generating ENCX encryption code..."
	@echo "Running validation before generation..."
	@if command -v encx-gen >/dev/null 2>&1; then \
		echo "Using global encx-gen CLI..."; \
		encx-gen validate -v . || exit 1; \
		echo "Validation passed. Generating code..."; \
		encx-gen generate -v .; \
	elif [ -f "./bin/encx-gen" ]; then \
		echo "Using local encx-gen binary..."; \
		./bin/encx-gen validate -v . || exit 1; \
		echo "Validation passed. Generating code..."; \
		./bin/encx-gen generate -v .; \
	else \
		echo "Using go run (no installation needed)..."; \
		go run github.com/hengadev/encx/cmd/encx-gen validate -v . || exit 1; \
		echo "Validation passed. Generating code..."; \
		go run github.com/hengadev/encx/cmd/encx-gen generate -v .; \
	fi
	@echo "ENCX code generation completed!"
	@echo "Generated files (suffix: _encx.go) have been created alongside source files"

encx-clean: ## Clean generated ENCX files (custom helper)
	@echo "Cleaning generated ENCX files..."
	@find . -name "*_encx.go" -type f -delete
	@echo "Generated ENCX files cleaned!"

# ============================================================================
# TESTING COMMANDS
# ============================================================================

# Test settings
TEST_TIMEOUT?=30m
TEST_VERBOSE?=-v
TEST_RACE?=-race
TEST_COVERAGE?=-cover
TEST_COVERAGE_HTML?=-coverprofile=coverage.out -covermode=atomic

# Domain paths (monolith structure)
DOMAINS=client case document media user
COMMON_PATHS=common/errs common/middleware common/contracts common/testutils
INFRASTRUCTURE_PATHS=infrastructure/postgres infrastructure/s3 infrastructure/redis infrastructure/vault
APPLICATION_PATHS=application/client application/case application/document application/media
INTERFACE_PATHS=interface/auth interface/case interface/client interface/document interface/reports interface/token

# Integration test paths
INTEGRATION_ROOT=test/integration
INTEGRATION_DOMAINS=client case document

# Test root paths
INTERNAL_ROOT=internal
PKG_ROOT=pkg

test-help: ## Show testing help message
	@echo 'Cluo API Backend Commands:'
	@echo ''
	@echo 'Encryption Commands (ENCX):'
	@echo '  encx-install      Install ENCX CLI globally (recommended)'
	@echo '  encx-validate      Validate ENCX struct tags before generation'
	@echo '  encx-generate      Generate ENCX encryption code for all structs'
	@echo '  encx-build         Build ENCX generator tool locally (optional)'
	@echo '  encx-clean         Clean generated ENCX files (custom helper)'
	@echo ''
	@echo 'Primary Test Commands:'
	@echo '  test-all           Run all unit and integration tests'
	@echo '  test-unit          Run all unit tests across all modules'
	@echo '  test-integration   Run all integration tests'
	@echo ''
	@echo 'Domain-Specific Unit Tests:'
	@echo '  test-unit-client      Run unit tests for client domain'
	@echo '  test-unit-case        Run unit tests for case domain'
	@echo '  test-unit-document    Run unit tests for document domain'
	@echo '  test-unit-media       Run unit tests for media domain'
	@echo '  test-unit-user        Run unit tests for user domain'
	@echo ''
	@echo 'Architecture Layer Tests:'
	@echo '  test-unit-common           Run unit tests for common utilities'
	@echo '  test-unit-infrastructure   Run unit tests for infrastructure layer'
	@echo '  test-unit-application      Run unit tests for application layer'
	@echo '  test-unit-interface        Run unit tests for interface layer'
	@echo ''
	@echo 'Domain-Specific Integration Tests:'
	@echo '  test-integration-client     Run integration tests for client domain'
	@echo '  test-integration-case       Run integration tests for case domain'
	@echo '  test-integration-document   Run integration tests for document domain'
	@echo ''
	@echo 'Layer-Specific Tests:'
	@echo '  test-domain-client          Run domain tests for client'
	@echo '  test-infrastructure-client  Run infrastructure tests for client'
	@echo ''
	@echo 'Advanced Testing Commands:'
	@echo '  test-coverage      Generate HTML coverage report for all modules'
	@echo '  test-race          Run tests with race condition detection'
	@echo '  test-benchmark     Run performance benchmarks'
	@echo '  test-verbose       Run tests with detailed output'
	@echo '  test-smoke         Run quick validation (short tests only)'
	@echo ''
	@echo 'Generic Specific Test Commands:'
	@echo '  test-file          Run specific test file (usage: make test-file TEST_PATH=<path>)'
	@echo '  test-func          Run specific test function (usage: make test-func TEST_NAME=<func> TEST_PATH=<path>)'
	@echo '  test-pkg           Run tests for specific package (usage: make test-pkg PKG_PATH=<path>)'
	@echo ''
	@echo 'Utility Commands:'
	@echo '  test-clean         Clean test artifacts and coverage files'
	@echo '  test-deps          Check test dependencies'

# Primary Commands
test-all: test-unit test-integration ## Run all unit and integration tests

test-unit: ## Run all unit tests across all modules
	@echo "Running all unit tests..."
	@echo "Testing domain modules..."
	@for domain in $(DOMAINS); do \
		if [ -d "$(INTERNAL_ROOT)/domain/$$domain" ]; then \
			echo "Testing domain/$$domain..."; \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/domain/$$domain/... || exit 1; \
		fi; \
	done
	@echo "Testing common modules..."
	@for path in $(COMMON_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			echo "Testing $$path..."; \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/... || exit 1; \
		fi; \
	done
	@echo "Testing infrastructure modules..."
	@for path in $(INFRASTRUCTURE_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			echo "Testing $$path..."; \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/... || exit 1; \
		fi; \
	done
	@echo "Testing application modules..."
	@for path in $(APPLICATION_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			echo "Testing $$path..."; \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/... || exit 1; \
		fi; \
	done
	@echo "Testing interface modules..."
	@for path in $(INTERFACE_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			echo "Testing $$path..."; \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/... || exit 1; \
		fi; \
	done
	@echo "Testing pkg modules..."
	@if [ -d "$(PKG_ROOT)" ]; then \
		echo "Testing pkg..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(PKG_ROOT)/... || exit 1; \
	fi
	@echo "All unit tests completed!"

test-integration: ## Run all integration tests
	@echo "Running all integration tests..."
	@for domain in $(INTEGRATION_DOMAINS); do \
		if [ -d "$(INTEGRATION_ROOT)/$$domain" ]; then \
			echo "Testing integration for $$domain..."; \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTEGRATION_ROOT)/$$domain/... || exit 1; \
		else \
			echo "Warning: $(INTEGRATION_ROOT)/$$domain not found, skipping..."; \
		fi; \
	done
	@echo "All integration tests completed!"

# Domain-Specific Unit Tests
test-unit-client: ## Run unit tests for client domain
	@echo "Running unit tests for client domain..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/domain/client/... ./$(INTERNAL_ROOT)/application/client/... ./$(INTERNAL_ROOT)/infrastructure/postgres/client/... ./$(INTERNAL_ROOT)/interface/client/...

test-unit-case: ## Run unit tests for case domain
	@echo "Running unit tests for case domain..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/domain/case/... ./$(INTERNAL_ROOT)/application/case/... ./$(INTERNAL_ROOT)/infrastructure/postgres/case/... ./$(INTERNAL_ROOT)/interface/case/...

test-unit-document: ## Run unit tests for document domain
	@echo "Running unit tests for document domain..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/domain/document/... ./$(INTERNAL_ROOT)/application/document/... ./$(INTERNAL_ROOT)/infrastructure/postgres/document/... ./$(INTERNAL_ROOT)/interface/document/...

test-unit-media: ## Run unit tests for media domain
	@echo "Running unit tests for media domain..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/domain/media/... ./$(INTERNAL_ROOT)/infrastructure/postgres/media/... ./$(INTERNAL_ROOT)/infrastructure/s3/...

test-unit-user: ## Run unit tests for user domain
	@echo "Running unit tests for user domain..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/domain/user/...

# Architecture Layer Tests
test-unit-common: ## Run unit tests for common utilities
	@echo "Running unit tests for common modules..."
	@for path in $(COMMON_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/...; \
		fi; \
	done

test-unit-infrastructure: ## Run unit tests for infrastructure layer
	@echo "Running unit tests for infrastructure modules..."
	@for path in $(INFRASTRUCTURE_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/...; \
		fi; \
	done

test-unit-application: ## Run unit tests for application layer
	@echo "Running unit tests for application modules..."
	@for path in $(APPLICATION_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/...; \
		fi; \
	done

test-unit-interface: ## Run unit tests for interface layer
	@echo "Running unit tests for interface modules..."
	@for path in $(INTERFACE_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/...; \
		fi; \
	done

# Domain-Specific Integration Tests
test-integration-client: ## Run integration tests for client domain
	@if [ -d "$(INTEGRATION_ROOT)/client" ]; then \
		echo "Running integration tests for client domain..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTEGRATION_ROOT)/client/...; \
	else \
		echo "Warning: $(INTEGRATION_ROOT)/client not found, skipping..."; \
	fi

test-integration-case: ## Run integration tests for case domain
	@if [ -d "$(INTEGRATION_ROOT)/case" ]; then \
		echo "Running integration tests for case domain..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTEGRATION_ROOT)/case/...; \
	else \
		echo "Warning: $(INTEGRATION_ROOT)/case not found, skipping..."; \
	fi

test-integration-document: ## Run integration tests for document domain
	@if [ -d "$(INTEGRATION_ROOT)/document" ]; then \
		echo "Running integration tests for document domain..."; \
		go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTEGRATION_ROOT)/document/...; \
	else \
		echo "Warning: $(INTEGRATION_ROOT)/document not found, skipping..."; \
	fi

# Layer-Specific Tests (for debugging specific layers)
test-domain-client: ## Run domain tests for client
	@echo "Running domain tests for client..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/domain/client/...

test-infrastructure-client: ## Run infrastructure tests for client
	@echo "Running infrastructure tests for client..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/infrastructure/postgres/client/...

test-infrastructure-case: ## Run infrastructure tests for case
	@echo "Running infrastructure tests for case..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/infrastructure/postgres/case/...

test-infrastructure-document: ## Run infrastructure tests for document
	@echo "Running infrastructure tests for document..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/infrastructure/postgres/document/...

test-infrastructure-media: ## Run infrastructure tests for media
	@echo "Running infrastructure tests for media..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/infrastructure/postgres/media/...

test-infrastructure-report: ## Run infrastructure tests for report
	@echo "Running infrastructure tests for report..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/infrastructure/postgres/report/...

test-infrastructure-token: ## Run infrastructure tests for token
	@echo "Running infrastructure tests for token..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/infrastructure/postgres/token/...

# Advanced Testing Commands
test-coverage: ## Generate HTML coverage report for all modules
	@echo "Generating coverage report for all modules..."
	@rm -f coverage.out coverage.html
	@echo "mode: atomic" > coverage.out
	@echo "Collecting coverage from all modules..."
	@# Domain modules
	@for domain in $(DOMAINS); do \
		if [ -d "$(INTERNAL_ROOT)/domain/$$domain" ]; then \
			go test $(TEST_COVERAGE_HTML) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/domain/$$domain/... -o $$domain-domain-coverage.out 2>/dev/null || true; \
			if [ -f $$domain-domain-coverage.out ]; then \
				tail -n +2 $$domain-domain-coverage.out >> coverage.out; \
				rm $$domain-domain-coverage.out; \
			fi; \
		fi; \
	done
	@# Common modules
	@for path in $(COMMON_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			module_name=$$(echo $$path | tr '/' '_'); \
			go test $(TEST_COVERAGE_HTML) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/... -o $$module_name-coverage.out 2>/dev/null || true; \
			if [ -f $$module_name-coverage.out ]; then \
				tail -n +2 $$module_name-coverage.out >> coverage.out; \
				rm $$module_name-coverage.out; \
			fi; \
		fi; \
	done
	@# Infrastructure modules
	@for path in $(INFRASTRUCTURE_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			module_name=$$(echo $$path | tr '/' '_'); \
			go test $(TEST_COVERAGE_HTML) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/... -o $$module_name-coverage.out 2>/dev/null || true; \
			if [ -f $$module_name-coverage.out ]; then \
				tail -n +2 $$module_name-coverage.out >> coverage.out; \
				rm $$module_name-coverage.out; \
			fi; \
		fi; \
	done
	@# Application modules
	@for path in $(APPLICATION_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			module_name=$$(echo $$path | tr '/' '_'); \
			go test $(TEST_COVERAGE_HTML) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/... -o $$module_name-coverage.out 2>/dev/null || true; \
			if [ -f $$module_name-coverage.out ]; then \
				tail -n +2 $$name-coverage.out >> coverage.out; \
				rm $$module_name-coverage.out; \
			fi; \
		fi; \
	done
	@# Interface modules
	@for path in $(INTERFACE_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			module_name=$$(echo $$path | tr '/' '_'); \
			go test $(TEST_COVERAGE_HTML) -timeout $(TEST_TIMEOUT) ./$(INTERNAL_ROOT)/$$path/... -o $$module_name-coverage.out 2>/dev/null || true; \
			if [ -f $$module_name-coverage.out ]; then \
				tail -n +2 $$module_name-coverage.out >> coverage.out; \
				rm $$module_name-coverage.out; \
			fi; \
		fi; \
	done
	@# Pkg modules
	@if [ -d "$(PKG_ROOT)" ]; then \
		go test $(TEST_COVERAGE_HTML) -timeout $(TEST_TIMEOUT) ./$(PKG_ROOT)/... -o pkg-coverage.out 2>/dev/null || true; \
		if [ -f pkg-coverage.out ]; then \
			tail -n +2 pkg-coverage.out >> coverage.out; \
			rm pkg-coverage.out; \
		fi; \
	fi
	@if [ -s coverage.out ]; then \
		go tool cover -html=coverage.out -o coverage.html; \
		echo "Coverage report generated: coverage.html"; \
		echo "Overall coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"; \
	else \
		echo "No coverage data collected."; \
	fi

test-race: ## Run tests with race condition detection
	@echo "Running tests with race condition detection..."
	@$(MAKE) test-unit TEST_RACE="$(TEST_RACE)"
	@$(MAKE) test-integration TEST_RACE="$(TEST_RACE)"

test-benchmark: ## Run performance benchmarks
	@echo "Running performance benchmarks..."
	@for domain in $(DOMAINS); do \
		if [ -d "$(INTERNAL_ROOT)/domain/$$domain" ]; then \
			echo "Benchmarking $$domain domain..."; \
			go test -bench=. -benchmem ./$(INTERNAL_ROOT)/domain/$$domain/... -timeout $(TEST_TIMEOUT) || true; \
		fi; \
	done
	@for path in $(INFRASTRUCTURE_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			echo "Benchmarking $$path..."; \
			go test -bench=. -benchmem ./$(INTERNAL_ROOT)/$$path/... -timeout $(TEST_TIMEOUT) || true; \
		fi; \
	done

test-verbose: ## Run tests with detailed output
	@echo "Running tests with verbose output..."
	@$(MAKE) test-unit TEST_VERBOSE="-v"
	@$(MAKE) test-integration TEST_VERBOSE="-v"

test-smoke: ## Run quick validation (short tests only)
	@echo "Running smoke tests (short tests only)..."
	@for domain in $(DOMAINS); do \
		if [ -d "$(INTERNAL_ROOT)/domain/$$domain" ]; then \
			echo "Smoke testing $$domain domain..."; \
			go test -short -timeout 5m ./$(INTERNAL_ROOT)/domain/$$domain/... || true; \
		fi; \
	done
	@for path in $(COMMON_PATHS); do \
		if [ -d "$(INTERNAL_ROOT)/$$path" ]; then \
			echo "Smoke testing $$path..."; \
			go test -short -timeout 5m ./$(INTERNAL_ROOT)/$$path/... || true; \
		fi; \
	done

# Specific Test Commands
test-file: ## Run specific test file (usage: make test-file TEST_PATH=<path>)
	@if [ -z "$(TEST_PATH)" ]; then \
		echo "Error: TEST_PATH is required. Usage: make test-file TEST_PATH=<path>"; \
		exit 1; \
	fi
	@echo "Running tests for file: $(TEST_PATH)"
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) $(TEST_PATH)

test-func: ## Run specific test function (usage: make test-func TEST_NAME=<func> TEST_PATH=<path>)
	@if [ -z "$(TEST_NAME)" ] || [ -z "$(TEST_PATH)" ]; then \
		echo "Error: TEST_NAME and TEST_PATH are required. Usage: make test-func TEST_NAME=<func> TEST_PATH=<path>"; \
		exit 1; \
	fi
	@echo "Running test function: $(TEST_NAME) in $(TEST_PATH)"
	@dir_path=$$(dirname "$(TEST_PATH)"); go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) -run $(TEST_NAME) "./$$dir_path"

test-pkg: ## Run tests for specific package (usage: make test-pkg PKG_PATH=<path>)
	@if [ -z "$(PKG_PATH)" ]; then \
		echo "Error: PKG_PATH is required. Usage: make test-pkg PKG_PATH=<path>"; \
		exit 1; \
	fi
	@echo "Running tests for package: $(PKG_PATH)"
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./$(PKG_PATH)/...

# Utility Commands
test-clean: ## Clean test artifacts and coverage files
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@rm -f *-coverage.out
	@find . -name "*.test" -delete
	@find . -name "*.prof" -delete
	@echo "Test artifacts cleaned!"

test-deps: ## Check test dependencies
	@echo "Checking test dependencies..."
	@go mod download
	@go mod verify
	@echo "Test dependencies verified!"

# Development helper commands
test-quick: ## Quick test run for current directory
	@echo "Running tests in current directory..."
	@go test $(TEST_VERBOSE) -timeout $(TEST_TIMEOUT) ./...

test-watch: ## Watch for changes and run tests (requires entr)
	@if command -v entr >/dev/null 2>&1; then \
		echo "Watching for changes... (Ctrl+C to stop)"; \
		find . -name "*.go" | entr -r make test-quick; \
	else \
		echo "Error: 'entr' command not found. Install it to use test-watch."; \
		echo "On macOS: brew install entr"; \
		echo "On Ubuntu: apt-get install entr"; \
	fi
