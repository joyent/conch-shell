.PHONY: clean deps update_deps docs_server release all changelog

CONCH_VERSION="0.1.2"
CONCH_BUILD_TIME=`date +%s`
CONCH_GIT_REV=`git describe --always --abbrev --dirty`

BUILD_ARGS = -ldflags="-X main.Version=${CONCH_VERSION} -X main.BuildTime=${CONCH_BUILD_TIME} -X main.GitRev=${CONCH_GIT_REV}" 

RELEASE_TARGET = release/darwin-amd64/conch release/linux-amd64/conch release/linux-arm/conch
UNAME_S = $(shell uname -s)

ifeq ($(UNAME_S),SunOS)
RELEASE_TARGET += release/solaris-amd64/conch
endif


BUILD = go build ${BUILD_ARGS} 
all: bin/conch bin/conch-mbo

bin/conch: pkg/**/*.go cmd/conch/*.go vendor
	${BUILD} -o bin/conch cmd/conch/conch.go

bin/conch-mbo: pkg/**/*.go cmd/conch-mbo/*.go vendor
	${BUILD} -o bin/conch-mbo cmd/conch-mbo/conch-mbo.go

clean: 
	rm -f bin/conch bin/conch-mbo

vendor: deps

deps:
	dep ensure -v -vendor-only

update_deps:
	dep ensure -update -v

docs_server:
	godoc -http=:6060 -v -goroot ./

release: release/${CONCH_VERSION} ${RELEASE_TARGET}

release/${CONCH_VERSION}:
	@mkdir -p release/${CONCH_VERSION}

release/darwin-amd64/conch:
	@mkdir -p release/${CONCH_VERSION}/darwin-amd64
	GOOS=darwin GOARCH=amd64 ${BUILD} -o release/${CONCH_VERSION}/darwin-amd64/conch cmd/conch/conch.go

release/linux-amd64/conch:
	@mkdir -p release/${CONCH_VERSION}/linux-amd64
	GOOS=linux GOARCH=amd64 ${BUILD} -o release/${CONCH_VERSION}/linux-amd64/conch cmd/conch/conch.go

release/linux-arm/conch:
	@mkdir -p release/${CONCH_VERSION}/linux-arm
	GOOS=linux GOARCH=arm ${BUILD} -o release/${CONCH_VERSION}/linux-arm/conch cmd/conch/conch.go

release/solaris-amd64/conch:
	@mkdir -p release/${CONCH_VERSION}/solaris-amd64
	GOOS=solaris GOARCH=amd64 ${BUILD} -o release/${CONCH_VERSION}/solaris-amd64/conch cmd/conch/conch.go

# gem install github_changelog_generator
changelog:
	github_changelog_generator -u joyent -p conch-shell #--due-tag ${CONCH_VERSION}
