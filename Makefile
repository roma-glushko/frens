BIN_DIR=$(PWD)/tmp/bin
VERSION_PACKAGE := nasmijewelry.com/shop/internal/version
VERSION ?= $(shell git describe --long --always --abbrev=15)
COMMIT ?= $(shell git describe --dirty --long --always --abbrev=15)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= "latest"

LDFLAGS_COMMON := "-X $(VERSION_PACKAGE).GitCommit=$(COMMIT) -X $(VERSION_PACKAGE).Version=$(VERSION) -X $(VERSION_PACKAGE).BuildDate=$(BUILD_DATE)"

.PHONY: help
help:
	@echo "🛠️ Frens :: Dev Commands\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: install-bin
install-bin: ## Install static checkers & other binaries
	@echo "🚚 Downloading binaries.."
	@GOBIN=$(BIN_DIR) go install mvdan.cc/gofumpt@latest
	@GOBIN=$(BIN_DIR) go install github.com/air-verse/air@latest
	@GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@GOBIN=$(BIN_DIR) go install github.com/g4s8/envdoc@latest
	@GOBIN=$(BIN_DIR) go install github.com/denis-tingaikin/go-header/cmd/go-header@latest
	@GOBIN=$(BIN_DIR) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: lint
lint: install-bin ## Lint the source code
	@echo "🧹 Cleaning go.mod.."
	@go mod tidy
	@echo "🧹 Formatting files.."
	@go fmt ./...
	@$(BIN_DIR)/gofumpt -l -w .
	@echo "🧹 Vetting go.mod.."
	@go vet ./...
	@echo "🧹 GoCI Lint.."
	@$(BIN_DIR)/golangci-lint run ./...

.PHONY: run
run: ## Run Frens
	@go run -ldflags $(LDFLAGS_COMMON) main.go -- $(filter-out $@,$(MAKECMDGOALS))

.PHONY: build
build: ## Build Frens
	@echo "🔨Building Frens binary.."
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@go build -ldflags $(LDFLAGS_COMMON) -o ./dist/frens;

.PHONY: test
test: ## Run tests
	@go test -v -count=1 -race -shuffle=on -coverprofile=coverage.txt ./...
