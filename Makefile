.PHONY: help test deps

# Ref: https://gist.github.com/prwhite/8168133
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

GOMOCK_VERSION = "v1.5.0"

gomock: ## go generate mock file.
	go install "github.com/golang/mock/mockgen@$(GOMOCK_VERSION)"
	go list ./... |grep -v '/gomock' | xargs go generate -v

header: ## check and add license header.
	sh addlicense.sh

lint: ## run lint
ifeq (, $(shell which golangci-lint))
	# binary will be $(go env GOPATH)/bin/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.51.2
else
	echo "Found golangci-lint"
endif
	golangci-lint run ./...

test: header lint ## Run test cases.
	go install "github.com/rakyll/gotest@v0.0.6"
	gotest -v -race -coverprofile=coverage.out -covermode=atomic ./...

deps:  ## Update vendor.
	go mod verify
	go mod tidy -v
