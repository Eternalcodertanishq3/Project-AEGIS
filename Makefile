# ============================================================================
# Project AEGIS — Build Automation
# ============================================================================
#
# Usage:
#   make frontend     Build the React frontend
#   make backend      Build the Go binary for the current OS/arch
#   make build        Build frontend + backend
#   make run          Build and run locally
#   make build-all    Cross-compile for all supported targets
#   make test         Run the Go test suite
#   make clean        Remove all build artifacts
#
# ============================================================================

# Project metadata
BINARY_NAME  := aegis
VERSION      ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME   := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GO_MODULE    := $(shell head -1 go.mod 2>/dev/null | awk '{print $$2}')

# Directories
FRONTEND_DIR := frontend
BACKEND_DIR  := backend
DIST_DIR     := dist
CMD_DIR      := $(BACKEND_DIR)/cmd/aegis

# Go build flags
LDFLAGS := -ldflags "-s -w \
	-X main.Version=$(VERSION) \
	-X main.BuildTime=$(BUILD_TIME)"

# Build targets (GOOS/GOARCH pairs)
TARGETS := \
	windows/amd64 \
	linux/amd64 \
	linux/arm64 \
	darwin/arm64

# ============================================================================
# Primary Targets
# ============================================================================

.PHONY: all
all: build ## Build frontend and backend (default)

.PHONY: build
build: frontend backend ## Build frontend, then backend

.PHONY: run
run: build ## Build and run AEGIS locally
	./$(DIST_DIR)/$(BINARY_NAME)

# ============================================================================
# Frontend
# ============================================================================

.PHONY: frontend
frontend: ## Build the React frontend
	@echo "━━━ Building frontend ━━━"
	cd $(FRONTEND_DIR) && npm ci --silent && npm run build
	@echo "✓ Frontend built → $(FRONTEND_DIR)/dist/"

.PHONY: frontend-dev
frontend-dev: ## Start frontend dev server (hot-reload)
	cd $(FRONTEND_DIR) && npm run dev

# ============================================================================
# Backend
# ============================================================================

.PHONY: backend
backend: ## Build the Go binary for the current OS/arch
	@echo "━━━ Building backend ($(shell go env GOOS)/$(shell go env GOARCH)) ━━━"
	@mkdir -p $(DIST_DIR)
	go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)$(shell go env GOEXE) ./$(CMD_DIR)
	@echo "✓ Backend built → $(DIST_DIR)/$(BINARY_NAME)$(shell go env GOEXE)"

# ============================================================================
# Cross-Compilation
# ============================================================================

.PHONY: build-all
build-all: frontend ## Cross-compile for all supported targets
	@echo "━━━ Cross-compiling for all targets ━━━"
	@mkdir -p $(DIST_DIR)
	@$(foreach target,$(TARGETS),\
		$(eval GOOS := $(word 1,$(subst /, ,$(target)))) \
		$(eval GOARCH := $(word 2,$(subst /, ,$(target)))) \
		$(eval EXT := $(if $(filter windows,$(GOOS)),.exe,)) \
		echo "  Building $(GOOS)/$(GOARCH)..." && \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) \
			-o $(DIST_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(EXT) ./$(CMD_DIR) && \
	) true
	@echo ""
	@echo "✓ All targets built:"
	@ls -lh $(DIST_DIR)/$(BINARY_NAME)-* 2>/dev/null || dir $(DIST_DIR)\$(BINARY_NAME)-* 2>nul
	@echo ""

# ============================================================================
# Testing
# ============================================================================

.PHONY: test
test: ## Run Go test suite
	@echo "━━━ Running tests ━━━"
	go test -v -race -count=1 ./...

.PHONY: test-short
test-short: ## Run Go tests (short mode, skip integration)
	go test -short -count=1 ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@mkdir -p $(DIST_DIR)
	go test -coverprofile=$(DIST_DIR)/coverage.out ./...
	go tool cover -html=$(DIST_DIR)/coverage.out -o $(DIST_DIR)/coverage.html
	@echo "✓ Coverage report → $(DIST_DIR)/coverage.html"

# ============================================================================
# Code Quality
# ============================================================================

.PHONY: lint
lint: ## Run linters
	go vet ./...
	@echo "✓ go vet passed"

.PHONY: fmt
fmt: ## Format Go source files
	gofmt -s -w .
	@echo "✓ Formatted"

# ============================================================================
# Utilities
# ============================================================================

.PHONY: clean
clean: ## Remove all build artifacts
	@echo "━━━ Cleaning ━━━"
	rm -rf $(DIST_DIR)
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules/.cache
	@echo "✓ Clean"

.PHONY: deps
deps: ## Download all dependencies
	go mod download
	cd $(FRONTEND_DIR) && npm ci --silent
	@echo "✓ Dependencies installed"

.PHONY: dev
dev: ## Run backend in development mode (assumes frontend already built)
	go run ./$(CMD_DIR)

.PHONY: version
version: ## Print version info
	@echo "Version:    $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Module:  $(GO_MODULE)"

# ============================================================================
# Help
# ============================================================================

.PHONY: help
help: ## Show this help message
	@echo "Project AEGIS — Build Targets"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'
	@echo ""
