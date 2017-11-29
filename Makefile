.PHONY: sane clean deps update_deps docs_server

CONCH_VERSION="0.0.1"
CONCH_BUILD_TIME=`date +%s`
CONCH_GIT_REV=`git describe --always --abbrev --dirty`


conch:
	go build -ldflags="-X main.Version=${CONCH_VERSION} -X main.BuildTime=${CONCH_BUILD_TIME} -X main.GitRev=${CONCH_GIT_REV}" -v -o conch

clean: 
	rm -f conch

sane:
	go tool vet -all -v config/*.go cmd/*.go *.go
	gofmt -w -s config/*.go cmd/*.go *.go

deps:
	glide install

update_deps:
	glide update

docs_server:
	godoc -http=:6060 -v -goroot ./
