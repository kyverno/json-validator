KOCACHE              ?= /tmp/ko-cache

#############
# VARIABLES #
#############

GIT_SHA             := $(shell git rev-parse HEAD)
TIMESTAMP           := $(shell date '+%Y-%m-%d_%I:%M:%S%p')
GOOS                ?= $(shell go env GOOS)
GOARCH              ?= $(shell go env GOARCH)
REGISTRY            ?= ghcr.io
REPO                ?= kyverno
BACKEND_DIR         := backend
BACKEND_BIN         := $(BACKEND_DIR)/backend
LD_FLAGS            := "-s -w"
LOCAL_PLATFORM      := linux/$(GOARCH)
PLATFORMS           := linux/arm64,linux/amd64
KO_PLATFORMS        := all
KO_TAGS             := $(GIT_SHA)
IMAGE    			:= json-validator
REPO     			:= $(REGISTRY)/$(REPO)/$(IMAGE)
KO_REGISTRY         := ko.local

ifndef VERSION
APP_VERSION         := $(GIT_SHA)
else
APP_VERSION         := $(VERSION)
endif

#########
# TOOLS #
#########

TOOLS_DIR                          := $(PWD)/.tools
HELM                               := $(TOOLS_DIR)/helm
HELM_VERSION                       := v3.10.1
KO                                 := $(TOOLS_DIR)/ko
KO_VERSION                         := main #e93dbee8540f28c45ec9a2b8aec5ef8e43123966
HELM_DOCS                          := $(TOOLS_DIR)/helm-docs
HELM_DOCS_VERSION                  := v1.11.0
GCI                                := $(TOOLS_DIR)/gci
GCI_VERSION                        := v0.9.1
GOFUMPT                            := $(TOOLS_DIR)/gofumpt
GOFUMPT_VERSION                    := v0.4.0
TOOLS                              := $(HELM) $(KO) $(HELM_DOCS) $(GCI) $(GOFUMPT)

$(HELM):
	@echo Install helm... >&2
	@GOBIN=$(TOOLS_DIR) go install helm.sh/helm/v3/cmd/helm@$(HELM_VERSION)

$(KO):
	@echo Install ko... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/google/ko@$(KO_VERSION)

$(HELM_DOCS):
	@echo Install helm-docs... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/norwoodj/helm-docs/cmd/helm-docs@$(HELM_DOCS_VERSION)

$(GCI):
	@echo Install gci... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/daixiang0/gci@$(GCI_VERSION)

$(GOFUMPT):
	@echo Install gofumpt... >&2
	@GOBIN=$(TOOLS_DIR) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

.PHONY: gofumpt
gofumpt: $(GOFUMPT)
	@echo "Running gofumpt"
	@$(GOFUMPT) -w ./backend

.PHONY: fmt
fmt: gci gofumpt

.PHONY: install-tools
install-tools: $(TOOLS) ## Install tools

.PHONY: clean-tools
clean-tools: ## Remove installed tools
	@echo Clean tools... >&2
	@rm -rf $(TOOLS_DIR)

###########
# CODEGEN #
###########

.PHONY: codegen-helm-docs
codegen-helm-docs: ## Generate helm docs
	@echo Generate helm docs... >&2
	@docker run -v ${PWD}/charts:/work -w /work jnorwood/helm-docs:v1.11.0 -s file

.PHONY: verify-helm-docs
verify-helm-docs: codegen-helm-docs ## Check Helm charts are up to date
	@echo Checking helm charts are up to date... >&2
	@git --no-pager diff -- charts
	@echo 'If this test fails, it is because the git diff is non-empty after running "make codegen-helm-docs".' >&2
	@echo 'To correct this, locally run "make codegen-helm-docs", commit the changes, and re-run tests.' >&2
	@git diff --quiet --exit-code -- charts

.PHONY: build-frontend
build-frontend:
	@cd frontend && npm install && APP_VERSION=$(APP_VERSION) npm run build

.PHONY: build-backend-assets
build-backend-assets: build-frontend ## Build backend assets
	@echo Building backend assets... >&2
	@rm -rf backend/pkg/ui/dist && cp -r frontend/dist backend/pkg/ui/dist

.PHONY: build-backend
build-backend: build-backend-assets ## Build backend
	@echo Building backend... >&2
	@cd backend && go mod tidy && go build .

.PHONY: ko-build
ko-build: $(KO) build-backend-assets ## Build image (with ko)
	@echo Build image with ko... >&2
	@cd backend && LDFLAGS=$(LD_FLAGS) KOCACHE=$(KOCACHE) KO_DOCKER_REPO=$(KO_REGISTRY) \
		$(KO) build . --preserve-import-paths --tags=$(KO_TAGS) --platform=$(LOCAL_PLATFORM)

.PHONY: ko-publish
ko-publish: $(KO) ## Build and publish image (with ko)
	@echo Publishing image with ko... >&2
	@cd backend && LDFLAGS=$(LD_FLAGS) KOCACHE=$(KOCACHE) KO_DOCKER_REPO=$(REPO) \
		$(KO) build . --bare --tags=$(KO_TAGS) --platform=$(KO_PLATFORMS)
