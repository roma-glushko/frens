.PHONY: help
help:
	@echo "🛠️ Dev Commands\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

VERSION_PACKAGE := github.com/roma-glushko/frens/internal/version
VERSION ?= $(shell git describe --long --always --abbrev=15)
COMMIT ?= $(shell git describe --dirty --long --always --abbrev=15)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= "latest"

LDFLAGS_COMMON := "-X $(VERSION_PACKAGE).GitCommit=$(COMMIT) -X $(VERSION_PACKAGE).Version=$(VERSION) -X $(VERSION_PACKAGE).BuildDate=$(BUILD_DATE)"

BIN_DIR=$(PWD)/tmp/bin

.PHONY: install-tools
install-tools: ## Install static checkers & other binaries
	@echo "🚚 Downloading tools.."
	@GOBIN=$(BIN_DIR) go install mvdan.cc/gofumpt@latest
	@GOBIN=$(BIN_DIR) go install github.com/air-verse/air@latest
	@GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@GOBIN=$(BIN_DIR) go install github.com/g4s8/envdoc@latest
	@GOBIN=$(BIN_DIR) go install github.com/denis-tingaikin/go-header/cmd/go-header@latest
	@GOBIN=$(BIN_DIR) go install github.com/goreleaser/goreleaser/v2@latest

.PHONY: lint
lint: install-tools ## Lint the source code
	@echo "🧹 Cleaning go.mod.."
	@go mod tidy
	@echo "🧹 Formatting files.."
	@go fmt ./...
	@$(BIN_DIR)/gofumpt -l -w .
	@echo "🧹 Vetting go.mod.."
	@go vet ./...
	@echo "🧹 GoCI Lint.."
	@$(BIN_DIR)/golangci-lint run ./...
	@echo "🧹Check GoReleaser.."
	@$(BIN_DIR)/goreleaser check

.PHONY: run
run: ## Run Frens
	@go run -ldflags $(LDFLAGS_COMMON) main.go -- $(filter-out $@,$(MAKECMDGOALS))

.PHONY: build
build: ## Build Frens
	@echo "🔨Building binary.."
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@go build -ldflags $(LDFLAGS_COMMON) -o ./dist/frens;

.PHONY: gen
gen: ## Generate Go code
	@go generate ./...

.PHONY: gen-check
gen-check: gen ## Check if Go code needs to be generated
	@git diff --exit-code

.PHONY: test
test: ## Run tests
	@go test -v -count=1 -race -shuffle=on -coverprofile=coverage.txt ./...

VENDOR ?= roma-glushko
PROJECT ?= frens
SOURCE ?= https://github.com/roma-glushko/frens
LICENSE ?= Apache-2.0
DESCRIPTION ?= "A friendship management application for introverts and not only. Build relationships with people that lasts."
REPOSITORY ?= roma-glushko/frens

.PHONY: image
image: ## Build docker image
	@echo "🛠️ Building image.."
	@echo "- Version: $(VERSION)"
	@echo "- Commit: $(COMMIT)"
	@echo "- Build Date: $(BUILD_DATE)"
	@docker build . -t $(REPOSITORY):$(VERSION) -f Dockerfile \
	--build-arg VERSION="$(VERSION)" \
	--build-arg COMMIT="$(COMMIT)" \
	--build-arg BUILD_DATE="$(BUILD_DATE)" \
	--label=org.opencontainers.image.vendor="$(VENDOR)" \
	--label=org.opencontainers.image.title="$(PROJECT)" \
	--label=org.opencontainers.image.revision="$(COMMIT)" \
	--label=org.opencontainers.image.version="$(VERSION)" \
	--label=org.opencontainers.image.created="$(BUILD_DATE)" \
	--label=org.opencontainers.image.source="$(SOURCE)" \
	--label=org.opencontainers.image.licenses="$(LICENSE)" \
	--label=org.opencontainers.image.description=$(DESCRIPTION)

.PHONY: image-lint
image-lint: ## Lint Dockerfile
	@hadolint Dockerfile
