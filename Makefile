CONCH_VERSION=`git describe --tags --abbrev=0 | sed 's/^v//'`
CONCH_BUILD_TIME=`date +%s`
CONCH_GIT_REV=`git describe --always --abbrev --dirty --long`

BUILD_HOST=`hostname -s`
BUILD_WHO="${USER}@${BUILD_HOST}"

FLAGS_PATH=github.com/joyent/conch-shell/pkg/util
BUILD_ARGS = -ldflags="-X ${FLAGS_PATH}.BuildHost=${BUILD_WHO} -X ${FLAGS_PATH}.Version=${CONCH_VERSION} -X ${FLAGS_PATH}.BuildTime=${CONCH_BUILD_TIME} -X ${FLAGS_PATH}.GitRev=${CONCH_GIT_REV}"

BUILD = go build ${BUILD_ARGS} 

default: clean vendor check bin/conch bin/tester ## By default, run 'clean', 'check', 'bin/conch', 'bin/tester'

first-run: tools ## Install all the dependencies needed to build and test

.PHONY: docker_test
docker_test: ## run a test build in docker
	docker/test.bash

.PHONY: docker_release
docker_release: ## Build all release binaries and checksums in docker
	docker/release.bash

.PHONY: bin/conch
bin/conch: pkg/**/*.go cmd/conch/*.go vendor ## Build bin/conch
	@echo "==> building bin/conch"
	${BUILD} -o bin/conch cmd/conch/conch.go


clean: ## Remove build products from bin/
	@echo "==> Removing build products from bin/"
	rm -f bin/conch bin/tester

.PHONY: vendor
vendor: ## Install dependencies
	dep ensure -v -vendor-only

update_deps: ## Update dependencies
	dep ensure -update -v

.PHONY: release
release: vendor check ## Build binaries for all supported platforms
	@echo "==> Building for all the platforms"
	GOOS=darwin GOARCH=amd64 ${BUILD} -o release/conch-darwin-amd64 cmd/conch/conch.go
	GOOS=linux GOARCH=amd64 ${BUILD} -o release/conch-linux-amd64 cmd/conch/conch.go
	GOOS=linux GOARCH=arm ${BUILD} -o release/conch-linux-arm cmd/conch/conch.go
	GOOS=solaris GOARCH=amd64 ${BUILD} -o release/conch-solaris-amd64 cmd/conch/conch.go
	GOOS=freebsd GOARCH=amd64 ${BUILD} -o release/conch-freebsd-amd64 cmd/conch/conch.go
	GOOS=openbsd GOARCH=amd64 ${BUILD} -o release/conch-openbsd-amd64 cmd/conch/conch.go

.PHONY: checksums
checksums: ## Build checksums for all release binaries
	@echo "==> Building checksums"
	@rm -f release/*.sha256
	@cd release && find . -type f -iname conch-\* -print0 | xargs -0 -n 1 -I {} sh -c 'shasum -a 256 {} > "{}.sha256"'

.PHONY: staticcheck
staticcheck: ## Run staticcheck
	staticcheck ./...

.PHONY: check
check: staticcheck ## Ensure that code matchs best practices and run tests
	@echo "==> Tests for pkg/conch"
	@cd pkg/conch && go test -v 

.PHONY: fasttest
fasttest: staticcheck ## Run staticheck, also go test with -failfast
	@echo "==> Tests for pkg/conch"
	@cd pkg/conch && go test -failfast -v 

tools: ## Download and install all dev/code tools
	@echo "==> Installing dev tools"
	go get -u github.com/golang/dep/cmd/dep
	go get -u honnef.co/go/tools/cmd/staticcheck


.PHONY: help
help: ## Display this help message
	@echo "GNU make(1) targets:"
	@grep -E '^[a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'


bin/tester: internal/pkg/cmd/tester/*.go cmd/tester/*.go vendor ## Build bin/tester
	@echo "==> building bin/tester"
	go build ${BUILD_ARGS} -o bin/tester cmd/tester/main.go


.PHONY: tester_release
tester_release: vendor check  ## Build binaries for conch-api-tester
	@echo "==> Building conch-api-tester for all the platforms"
	GOOS=darwin GOARCH=amd64 ${BUILD} -o release/conch-api-tester-darwin-amd64 cmd/tester/main.go
	GOOS=linux GOARCH=amd64 ${BUILD} -o release/conch-api-tester-linux-amd64 cmd/tester/main.go
	GOOS=solaris GOARCH=amd64 ${BUILD} -o release/conch-api-tester-solaris-amd64 cmd/tester/main.go


