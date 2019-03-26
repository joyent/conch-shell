FROM golang:1.12.1-alpine AS build
ENV CGO_ENABLED 0

RUN apk add --no-cache --update make git perl-utils dep shadow

ARG CACHE_BUSTER="wat"

ENV PATH "/go/bin:${PATH}"

RUN go get honnef.co/go/tools/cmd/staticcheck

RUN mkdir -p /go/src/github.com/joyent/conch-shell/
WORKDIR /go/src/github.com/joyent/conch-shell/

ARG VCS_REF="master"
ARG VERSION="v0.0.0-dirty"
LABEL org.label-schema.vcs-ref $VCS_REF
LABEL org.label-schema.version $VERSION 

COPY . /go/src/github.com/joyent/conch-shell/

RUN make

FROM scratch
COPY --from=build /go/src/github.com/joyent/conch-shell/bin/conch /bin/conch
COPY --from=build /etc/ssl /etc/ssl

ENTRYPOINT [ "/bin/conch", "--no-version-check" ]
CMD ["version"]
