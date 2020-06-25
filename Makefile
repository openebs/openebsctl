# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
GOBIN := $(or $(shell go env GOBIN 2>/dev/null), $(shell go env GOPATH 2>/dev/null)/bin)

PACKAGES = $(shell go list ./... | grep -v 'vendor\|pkg/generated')

#name mayactl to be a kubectl-plugin
MAYACTL=kubectl-mayactl

.PHONY: all
all: deps verify-deps mayactl


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

.PHONY: mayactl
mayactl:
	@echo "----------------------------"
	@echo "--> mayactl                    "
	@echo "----------------------------"
	@PNAME=MAYACTL CTLNAME=${MAYACTL} sh -c "'$(PWD)/build.sh'"
	@echo "--> Removing old directory..."
	@sudo rm -rf /usr/local/bin/${MAYACTL}
	@echo "----------------------------"
	@echo "copying new mayactl"
	@echo "----------------------------"
	@sudo mkdir -p  /usr/local/bin/
	@sudo cp -a "$(PWD)/bin/MAYACTL/${MAYACTL}"  /usr/local/bin/${MAYACTL}
	@echo "=> copied to /usr/local/bin"
