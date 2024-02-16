GO ?= go

all: help

dep: ## Update dependencies
	$(GO) mod tidy

goimports: ## Run goimports
	goimports -w .

lint: ## Run static analysis checks
	golangci-lint --timeout 5m run 

test: ## Run unittests
	$(GO) test -race -short -count=1 ./...

ready: dep goimports lint test ## Runs all checks and generations before commit

css:
	tailwindcss -i pkg/templates/static/input.css -o pkg/templates/static/output.css --watch

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: all dep goimports lint test ready help  