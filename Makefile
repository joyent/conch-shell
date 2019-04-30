VERSION ?= $(shell git describe --tags --abbrev=0 | sed 's/^v//')
DISABLE_API_VERSION_CHECK ?= 0
DISABLE_API_TOKEN_CRUD ?= 0
DISABLE_ADMIN_FUNCTIONS ?= 0

# Pass in a different value please. Please?
TOKEN_OBFUSCATION_KEY ?= "eig0Ahcoi4phepoow2Wee8ahfoe3een4shebahz0Uhu8O"

build: vendor clean test all ## Test and build binaries for local architecture into bin/

.PHONY: docker_test
docker_test: ## run a test build in docker
	docker/test.bash

.PHONY: docker_release
docker_release: ## Build all release binaries and checksums in docker
	docker/release.bash

.PHONY: clean
clean: ## Remove build products from bin/ and release/
	rm -rf bin
	rm -rf release

vendor: ## Install dependencies
	dep ensure -v -vendor-only

.PHONY: deps
deps: ## Update dependencies to latest version
	dep ensure -v

.PHONY: test
test: ## Ensure that code matchs best practices and run tests
	staticcheck ./...
	go test -v ./pkg/conch ./pkg/util ./pkg/config ./pkg/conch/uuid

.PHONY: tools
tools: ## Download and install all dev/code tools
	@echo "==> Installing dev tools"
	go get -u github.com/golang/dep/cmd/dep
	go get -u honnef.co/go/tools/cmd/staticcheck


################################
# Dynamic Fanciness            #
# aka Why GNU make Is Required #
################################

PLATFORMS  := darwin-amd64 linux-amd64 solaris-amd64 freebsd-amd64 openbsd-amd64 linux-arm
BINARIES   := conch conch-minimal tester corpus
RELEASE_BINARIES := conch

BINS       := $(foreach bin,$(BINARIES),bin/$(bin)) 
RELEASES   := $(foreach bin,$(RELEASE_BINARIES),release/$(bin))

GIT_REV    := $(shell git describe --always --abbrev --dirty --long)
FLAGS_PATH := github.com/joyent/conch-shell/pkg/util
LD_FLAGS   := -ldflags="-X github.com/joyent/conch-shell/pkg/config.ObfuscationKey=${TOKEN_OBFUSCATION_KEY} -X $(FLAGS_PATH).Version=$(VERSION) -X $(FLAGS_PATH).GitRev=$(GIT_REV) -X $(FLAGS_PATH).FlagsDisableApiVersionCheck=$(DISABLE_API_VERSION_CHECK) -X $(FLAGS_PATH).FlagsDisableApiTokenCRUD=$(DISABLE_API_TOKEN_CRUD) -X $(FLAGS_PATH).FlagsNoAdmin=$(DISABLE_ADMIN_FUNCTIONS)"
BUILD      := CGO_ENABLED=0 go build $(LD_FLAGS) 

####

all: $(BINS) ## Build all binaries

.PHONY: release
release: vendor test $(RELEASES) ## Build release binaries with checksums

bin/%:
	@mkdir -p bin
	@echo "> Building bin/$(subst bin/,,$@)"
	@$(BUILD) -o bin/$(subst bin/,,$@) cmd/$(subst bin/,,$@)/*.go

os   = $(firstword $(subst -, ,$1))
arch = $(lastword $(subst -, ,$1))

define release_me
	$(eval BIN:=$(subst release/,,$@))
	$(eval GOOS:=$(call os, $(platform)))
	$(eval GOARCH:=$(call arch, $(platform)))
	$(eval RPATH:=release/$(BIN)-$(GOOS)-$(GOARCH))

	@echo "> Building $(RPATH)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(BUILD) -o $(RPATH) cmd/$(BIN)/*.go
	shasum -a 256 $(RPATH) > $(RPATH).sha256
endef


release/%:
	@mkdir -p release
	$(foreach platform,$(PLATFORMS),$(call release_me))


.PHONY: help
help: ## Display this help message
	@echo "GNU make(1) targets:"
	@grep -E '^[a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'


