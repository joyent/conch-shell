.PHONY: sane clean deps update_deps docs_server

CONCH_VERSION="0.0.0"
CONCH_BUILD_TIME=`date +%s`
CONCH_GIT_REV=`git describe --always --abbrev --dirty`


conch: main.go config/** util/** workspaces/** user/** devices/** reports/** templates/**.go
	go build -ldflags="-X main.Version=${CONCH_VERSION} -X main.BuildTime=${CONCH_BUILD_TIME} -X main.GitRev=${CONCH_GIT_REV}" -v -o conch

clean: 
	rm -f conch

sane:
	gofmt -w -s config/*.go util/*.go workspaces/*.go devices/*.go user/*.go reports/*.go templates/*.go

deps:
	glide install

vendor: glide.lock
	glide update

update_deps:
	glide update

docs_server:
	godoc -http=:6060 -v -goroot ./
