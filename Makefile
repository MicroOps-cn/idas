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

DONT_FIND := -name vendor -prune -o -name .git -prune -o -name .cache -prune -o -name .pkg -prune -o

PROTOC       ?= protoc
#PROTOC_OPTS ?= -I ./vendor/github.com/gogo/protobuf:./api:./:vendor
#PROTOC_OPTS := $(PROTOC_OPTS) --gogoslick_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor,plugins=grpc,paths=source_relative:.
# Protobuf files
PROTO_DEFS := $(shell find . $(DONT_FIND) -type f -name '*.proto' -print)
PROTO_GOS := $(patsubst %.proto,%.pb.go,$(PROTO_DEFS))

PROTOC_OPTS ?= --gogo_opt=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor
PROTOC_OPTS := $(PROTOC_OPTS) -I$(shell $(GO) list -f "{{ .Dir }}" -m github.com/gogo/protobuf)/protobuf/
PROTOC_OPTS := $(PROTOC_OPTS) -I$(shell $(GO) list -f "{{ .Dir }}" -m github.com/gogo/protobuf)/
PROTOC_OPTS := $(PROTOC_OPTS) -I./api
PROTOC_OPTS := $(PROTOC_OPTS) --gogo_out=module=${GOMODULENAME}:.

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

clean-protos:
	rm -rf $(PROTO_GOS)
protos: clean-protos $(PROTO_GOS)

%.pb.go:
	@# The store-gateway RPC is based on Thanos which uses relative references to other protos, so we need
	@# to configure all such relative paths. `gogo/protobuf` is used by it.
	case "$@" in	\
		vendor*)			\
			;;					\
		*)						\
			$(PROTOC) $(PROTOC_OPTS) ./$(patsubst %.pb.go,%.proto,$@); \
			;;					\
		esac

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
