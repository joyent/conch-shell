.PHONY: clean deps update_deps docs_server release all

CONCH_VERSION="0.0.0"
CONCH_BUILD_TIME=`date +%s`
CONCH_GIT_REV=`git describe --always --abbrev --dirty`

BUILD = go build -ldflags="-X main.Version=${CONCH_VERSION} -X main.BuildTime=${CONCH_BUILD_TIME} -X main.GitRev=${CONCH_GIT_REV}" 

all: bin/conch bin/conch-mbo

bin/conch: pkg/**/*.go cmd/conch/*.go
	go build -ldflags="-X main.Version=${CONCH_VERSION} -X main.BuildTime=${CONCH_BUILD_TIME} -X main.GitRev=${CONCH_GIT_REV}" -v -o bin/conch cmd/conch/conch.go

bin/conch-mbo: pkg/**/*.go cmd/conch-mbo/*.go
	go build -ldflags="-X main.Version=${CONCH_VERSION} -X main.BuildTime=${CONCH_BUILD_TIME} -X main.GitRev=${CONCH_GIT_REV}" -v -o bin/conch-mbo cmd/conch-mbo/conch-mbo.go

clean: 
	rm -f bin/conch bin/conch-mbo

deps:
	glide install

vendor: glide.lock
	glide update

update_deps:
	glide update

docs_server:
	godoc -http=:6060 -v -goroot ./


release:
	@mkdir -p release/${CONCH_VERSION}
	@mkdir -p release/${CONCH_VERSION}/darwin-amd64
	GOOS=darwin GOARCH=amd64 ${BUILD} -o release/${CONCH_VERSION}/darwin-amd64/conch cmd/conch/conch.go
	@mkdir -p release/${CONCH_VERSION}/linux-amd64
	GOOS=linux GOARCH=amd64 ${BUILD} -o release/${CONCH_VERSION}/linux-amd64/conch cmd/conch/conch.go
	@mkdir -p release/${CONCH_VERSION}/linux-arm
	GOOS=linux GOARCH=arm ${BUILD} -o release/${CONCH_VERSION}/linux-arm/conch cmd/conch/conch.go
	@mkdir -p release/${CONCH_VERSION}/solaris-amd64
	GOOS=solaris GOARCH=amd64 ${BUILD} -o release/${CONCH_VERSION}/solaris-amd64/conch cmd/conch/conch.go
