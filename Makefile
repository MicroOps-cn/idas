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

DONT_FIND := -name vendor -prune -o -name .git -prune -o -name .cache -prune -o -name .pkg -prune -o -name test -prune -o

PROTOC       ?= protoc
#PROTOC_OPTS ?= -I ./vendor/github.com/gogo/protobuf:./api:./:vendor
#PROTOC_OPTS := $(PROTOC_OPTS) --gogoslick_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor,plugins=grpc,paths=source_relative:.
# Protobuf files
PROTO_DEFS := $(shell find . $(DONT_FIND) -type f -name '*.proto' -print)
PROTO_GOS := $(shell find . $(DONT_FIND) -type f -name '*.pb.go' -print)
PROTOC_OPTS ?=
GOGO_OPT := $(GOGO_OPT)Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor,
GOGO_OPT := $(GOGO_OPT)Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,
GOGO_OPT := $(GOGO_OPT)Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,
GOGO_OPT := $(GOGO_OPT)Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,
GOGO_OPT := $(GOGO_OPT)Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,
GOGO_OPT := $(GOGO_OPT)Mgoogle/api/http.proto=github.com/gogo/googleapis/google/api,
GOGO_OPT := $(GOGO_OPT)Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,
PROTOC_OPTS := $(PROTOC_OPTS) --gogo_opt=$(GOGO_OPT)
PROTOC_OPTS := $(PROTOC_OPTS) -I$(shell $(GO) list -f "{{ .Dir }}" -m github.com/gogo/protobuf)/protobuf/
PROTOC_OPTS := $(PROTOC_OPTS) -I$(shell $(GO) list -f "{{ .Dir }}" -m github.com/gogo/protobuf)/
PROTOC_OPTS := $(PROTOC_OPTS) -I/home/sunlinyao/go/pkg/mod/github.com/gogo/googleapis@v1.4.0/
PROTOC_OPTS := $(PROTOC_OPTS) -I./api
PROTOC_OPTS := $(PROTOC_OPTS) --gogo_out=plugins=grpc,module=${GOMODULENAME}:gogo_out
PROTOC_OPTS := $(PROTOC_OPTS) --grpc-gateway_out=${GOGO_OPT}:gogo_out

BASE_PATH =

IMAGE_SUFFIX :=

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

GitCommit   = $(shell git rev-parse --short HEAD)
BuildDate   = $(shell date +%Y-%m-%dT%H:%M:%S%Z)
GoVersion   = $(shell go version|awk '{print $$3}')
Platform    = $(shell go version|awk '{print $$4}')
Version     ?= $(shell cat version)
LDFlags     := -w -s -X 'lampao/pkg/utils/version.GitCommit=$(GitCommit)'
LDFlags     += -X 'lampao/pkg/utils/version.BuildDate=$(BuildDate)'
LDFlags     += -X 'lampao/pkg/utils/version.GoVersion=$(GoVersion)'
LDFlags     += -X 'lampao/pkg/utils/version.Platform=$(Platform)'
LDFlags     += -X 'lampao/pkg/utils/version.Version=$(Version).$(GitCommit)'

info:
	@echo "Version: $(Version)"
	@echo "Git commit: $(GitCommit)"
	@echo "Build date: $(BuildDate)"
	@echo "Platform: $(Platform)"
	@echo "Go version: $(GoVersion)"

.PHONY: all
all: common-check_license protos common-lint test idas

.PHONY: clean-protos
clean-protos:
	@if [ -n "$(PROTO_GOS)" ];then \
    	echo ">> rm -f $(PROTO_GOS);" \
		rm -f $(PROTO_GOS); \
	fi \

.PHONY: protos
protos: clean-protos
	@echo ">> generate golang code of protobuf"
	PROTOC="$(PROTOC)" PROTOC_OPTS="$(PROTOC_OPTS)" PROTO_DEFS="$(PROTO_DEFS)" GOMODULENAME="$(GOMODULENAME)" ./scripts/make_proto.sh
#	grep -HoP '(?<=option go_package = ")[^;]+' $(PROTO_DEFS)|awk -F: '{pkgs[$$2]=1;protos[$$1]=$$2}END{for(pkg in pkgs){for(proto in protos){if(pkg==protos[proto]){printf("%s ",proto)}};print("")}}'|while read line; do \
#  		$(PROTOC) $(PROTOC_OPTS) $$line || echo "$(PROTOC) $(PROTOC_OPTS) $$line" && exit 1 ;\
#	done
#	find . $(DONT_FIND) -type f -name '*.pb.go' -print|while read line; do \
#  		sed -i ':label;N;s/\nvar E_\S\+ = gogoproto.E_\S\+\n//;b label' $$line; \
#  		sed -i '/gogoproto "github.com\/gogo\/protobuf\/gogoproto/d' $$line; \
#  	done;


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

.PHONY: common-check_license
common-check_license:
	@echo ">> checking license header"
	@licRes=$$(for file in $$(find . -type f -iname '*.go' ! -path './pkg/transport/*'  ! -path './test/*' ! -path './vendor/*') ; do \
               awk 'NR<=4' $$file | grep -Eq "(Copyright|generated|GENERATED)" || echo $$file; \
       done); \
       if [ -n "$${licRes}" ]; then \
               echo "license header checking failed:"; echo "$${licRes}"; \
               exit 1; \
       fi

.PHONY: idas
idas:
	CGO_ENABLED=0 go build -ldflags="-s -w $(LDFlags)" -o dist/idas ./cmd/idas

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

.PHONY: test
test:
	go test -tags make_test -cover -race -count=1 ./...

.PHONY: openapi
openapi:
	go run cmd/openapi/main.go -o public/config/idas.json
	cd public && yarn openapi
	go run scripts/sync_to_public.go

.PHONY: ui
ui:
	rm -rf public/src/.umi-production/
	cd public && yarn install && yarn run build --basePath='$(BASE_PATH)/admin/' --apiPath='$(BASE_PATH)/'
	rm -rf pkg/transport/static && cp -r public/dist pkg/transport/static

