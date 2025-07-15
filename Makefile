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
GOBIN ?= $(BIN_DIR)

export GOBIN
export PATH := $(BIN_DIR):$(PATH)

.PHONY: tools
tools: tools-test ## Install static checkers & other binaries
	@echo "🚚 Downloading tools.."
	@mkdir -p $(GOBIN)
	@ \
	command -v golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest & \
	command -v goreleaser > /dev/null || go install github.com/goreleaser/goreleaser/v2@latest & \
	wait

tools-test: ## Install tools for testing
	@echo "🚚 Downloading tools for testing.."
	@mkdir -p $(GOBIN)
	@ \
	command -v gocover-cobertura > /dev/null || go install github.com/boumenot/gocover-cobertura@latest & \
	wait

.PHONY: lint
lint: tools ## Lint the source code
	@echo "🧹 Cleaning go.mod.."
	@go mod tidy
	@echo "🧹 Formatting files.."
	@go fmt ./...
	@echo "🧹 Vetting go.mod.."
	@go vet ./...
	@echo "🧹 GoCI Lint.."
	@$(BIN_DIR)/golangci-lint fmt ./...
	@$(BIN_DIR)/golangci-lint run ./...
	@echo "🧹 Check GoReleaser.."
	@$(BIN_DIR)/goreleaser check

.PHONY: lint-ci
lint-ci: ## Lint the source code in CI mode
	@echo "🧹 Cleaning go.mod.."
	@go mod tidy
	@echo "🧹 Formatting files.."
	@go fmt ./...
	@echo "🧹 Vetting go.mod.."
	@go vet ./...

.PHONY: run
run: ## Run Frens
	@go run -ldflags $(LDFLAGS_COMMON) main.go -- $(filter-out $@,$(MAKECMDGOALS))

.PHONY: build
build: ## Build Frens
	@echo "🔨 Building binary.."
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

.PHONY: test-ci
test-ci: tools-test test ## Run tests in the CI mode
	@gocover-cobertura < coverage.txt > coverage.xml

copyright: ## Apply copyrights to all files
	@echo "🧹 Applying license headers"
	@docker run  --rm -v $(CURDIR):/github/workspace ghcr.io/apache/skywalking-eyes/license-eye:4021a396bf07b2136f97c12708476418a8157d72 -v info -c .licenserc.yaml header fix

license: copyright

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
