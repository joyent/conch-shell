.PHONY: clean deps update_deps docs_server all

CONCH_VERSION="0.0.0"
CONCH_BUILD_TIME=`date +%s`
CONCH_GIT_REV=`git describe --always --abbrev --dirty`


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
