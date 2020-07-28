# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
GOBIN := $(or $(shell go env GOBIN 2>/dev/null), $(shell go env GOPATH 2>/dev/null)/bin)

PACKAGES = $(shell go list ./... | grep -v 'vendor\|pkg/generated')

#name openebsctl to be a kubectl-plugin
OPENEBSCTL=kubectl-openebs

.PHONY: all
all: deps verify-deps openebsctl


# deps ensures fresh go.mod and go.sum.
.PHONY: deps
deps:
	@echo "--> Tidying up submodules"
	@go mod tidy
	@echo "--> Veryfying submodules"
	@go mod verify

# to verify the deps are in sync
.PHONY: verify-deps
verify-deps: deps
	@if !(git diff --quiet HEAD -- go.sum go.mod); then \
		echo "go module files are out of date, please commit the changes to go.mod and go.sum"; exit 1; \
	fi

.PHONY: format
format:
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

#.PHONY: test
#test:
#	go test ./...

.PHONY: openebsctl
openebsctl:
	@echo "----------------------------"
	@echo "--> openebsctl                    "
	@echo "----------------------------"
	@PNAME=OPENEBSCTL CTLNAME=${OPENEBSCTL} sh -c "'$(PWD)/build.sh'"
	@echo "--> Removing old directory..."
	@sudo rm -rf /usr/local/bin/${OPENEBSCTL}
	@echo "----------------------------"
	@echo "copying new openebsctl"
	@echo "----------------------------"
	@sudo mkdir -p  /usr/local/bin/
	@sudo cp -a "$(PWD)/bin/OPENEBSCTL/${OPENEBSCTL}"  /usr/local/bin/${OPENEBSCTL}
	@echo "=> copied to /usr/local/bin"
