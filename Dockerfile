FROM golang:1.11.1-alpine
RUN apk add --update make git perl-utils

ENV CGO_ENABLED 0

ARG CACHE_BUSTER="wat"

RUN go get github.com/golang/dep/cmd/dep && \
	go get github.com/alecthomas/gometalinter && \
	gometalinter --install

RUN mkdir -p /go/src/github.com/joyent/conch-shell/
WORKDIR /go/src/github.com/joyent/conch-shell/

ARG VCS_REF="master"
ARG VERSION="v0.0.0-dirty"

COPY . /go/src/github.com/joyent/conch-shell/
RUN make default

ENTRYPOINT [ "bin/conch" ]
