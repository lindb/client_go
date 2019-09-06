.PHONY: test

pre-test: build ## go generate mock file.
	if [ ! -e ./bin/golangci-lint ]; then \
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s $(GOLANGCI_LINT_VERSION); \
	fi
	GO111MODULE=on ./bin/golangci-lint run


test:  pre-test ## Run test cases. (Args: GOLANGCI_LINT_VERSION=latest)
	GO111MODULE=on go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

build: ## build trigger downloading depenciencies
	GO111MODULE=on go build -o ./bin/example github.com/lindb/client_go/example
