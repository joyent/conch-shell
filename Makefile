.PHONY: clean deps update_deps docs_server release all changelog

CONCH_VERSION="0.2.1"
CONCH_BUILD_TIME=`date +%s`
CONCH_GIT_REV=`git describe --always --abbrev --dirty`

UNAME_S=$(shell uname -s)

BUILD_ARGS = -ldflags="-X main.Version=${CONCH_VERSION} -X main.BuildTime=${CONCH_BUILD_TIME} -X main.GitRev=${CONCH_GIT_REV}" 

BUILD = go build ${BUILD_ARGS} 

BINARIES = bin/conch

ifneq ($(UNAME_S), SunOS)
BINARIES += bin/conch-mbo
endif

all: ${BINARIES}

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

release: vendor
	GOOS=darwin GOARCH=amd64 ${BUILD} -o release/conch-mbo-darwin-amd64 cmd/conch-mbo/conch-mbo.go
	GOOS=darwin GOARCH=amd64 ${BUILD} -o release/conch-darwin-amd64 cmd/conch/conch.go
	GOOS=linux GOARCH=amd64 ${BUILD} -o release/conch-linux-amd64 cmd/conch/conch.go
	GOOS=linux GOARCH=arm ${BUILD} -o release/conch-linux-arm cmd/conch/conch.go
	GOOS=solaris GOARCH=amd64 ${BUILD} -o release/conch-solaris-amd64 cmd/conch/conch.go
	GOOS=freebsd GOARCH=amd64 ${BUILD} -o release/conch-freebsd-amd64 cmd/conch/conch.go

# gem install github_changelog_generator
changelog:
	github_changelog_generator -u joyent -p conch-shell #--due-tag ${CONCH_VERSION}
