GO           ?= go
GOFMT        ?= $(GO)fmt
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
GOOPTS       ?=
GOHOSTOS     ?= $(shell $(GO) env GOHOSTOS)
GOHOSTARCH   ?= $(shell $(GO) env GOHOSTARCH)
GOMODULENAME   ?= $(shell $(GO) list -m)
GO_VERSION        ?= $(shell $(GO) version)
GO_VERSION_NUMBER ?= $(word 3, $(GO_VERSION))

GOLANGCI_LINT :=
GOLANGCI_LINT_OPTS ?=
GOLANGCI_LINT_VERSION ?= v1.45.2

PROTOC       ?= protoc
PROTOC_OPTS ?= --go_opt=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types
PROTOC_OPTS := $(PROTOC_OPTS) -I$(shell go env GOMODCACHE)/github.com/gogo/protobuf@v1.3.2/protobuf/
PROTOC_OPTS := $(PROTOC_OPTS) -I./api

# golangci-lint only supports linux, darwin and windows platforms on i386/amd64.
# windows isn't included here because of the path separator being different.
ifeq ($(GOHOSTOS),$(filter $(GOHOSTOS),linux darwin))
	ifeq ($(GOHOSTARCH),$(filter $(GOHOSTARCH),amd64 i386))
		# If we're in CI and there is an Actions file, that means the linter
		# is being run in Actions, so we don't need to run it here.
		ifeq (,$(CIRCLE_JOB))
			GOLANGCI_LINT := $(FIRST_GOPATH)/bin/golangci-lint
		else ifeq (,$(wildcard .github/workflows/golangci-lint.yml))
			GOLANGCI_LINT := $(FIRST_GOPATH)/bin/golangci-lint
		endif
	endif
endif

pkgs          = ./...


proto:
	for protofile in `find -name "*.proto"`; \
	do \
		$(PROTOC) --go_out=module=${GOMODULENAME}:. ${PROTOC_OPTS}  $${protofile}; \
	done
#	$(PROTOC) --go_out=module=${GOMODULENAME}:. ./api/types/capacity.proto
#	protoc -I$(shell go env GOMODCACHE)/github.com/gogo/protobuf@v1.3.2/protobuf/ \
#		-I./config \
#		-I./pkg/utils/capacity/ \
#		-I./pkg/utils/fs/ \
#		--go_opt=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types
#		--go_out=./config \
#		./config/config.proto

idas:
	go build -ldflags="-s -w" -o dist/idas ./cmd/idas

.PHONY: common-lint
common-lint: $(GOLANGCI_LINT)
ifdef GOLANGCI_LINT
	@echo ">> running golangci-lint"
	$(GOLANGCI_LINT) run $(pkgs)
endif

ifdef GOLANGCI_LINT
$(GOLANGCI_LINT):
	mkdir -p $(FIRST_GOPATH)/bin
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_LINT_VERSION)/install.sh \
		| sed -e '/install -d/d' \
		| sh -s -- -b $(FIRST_GOPATH)/bin $(GOLANGCI_LINT_VERSION)
endif
