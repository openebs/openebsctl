# Copyright 2020 The OpenEBS Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
GOBIN := $(or $(shell go env GOBIN 2>/dev/null), $(shell go env GOPATH 2>/dev/null)/bin)

PACKAGES = $(shell go list ./... | grep -v 'vendor\|pkg/generated')

#name openebsctl to be a kubectl-plugin
OPENEBSCTL=kubectl-openebs

.PHONY: all
all: license-check deps verify-deps openebsctl


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

.PHONY: golint
golint:
	@echo "+ Installing golint"
	@GO111MODULE=on go get -u golang.org/x/lint/golint;
	@echo "--> Running golint"
	@golint -set_exit_status $(PACKAGES)
	@echo "Golint successful !"
	@echo "--------------------------------"

.PHONY: license-check
license-check:
	@echo "--> Checking license header..."
	@licRes=$$(for file in $$(find . -type f -regex '.*\.sh\|.*\.go\|.*Docker.*\|.*\Makefile*' ! -path './vendor/*') ; do \
               awk 'NR<=3' $$file | grep -Eq "(Copyright|generated|GENERATED)" || echo $$file; \
       done); \
       if [ -n "$${licRes}" ]; then \
               echo "license header checking failed:"; echo "$${licRes}"; \
               exit 1; \
       fi
	@echo "--> Done checking license."
