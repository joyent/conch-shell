# What

`conch` is a CLI for accessing the [Conch](https://github.com/joyent/conch) API.

[![Go Report Card](https://goreportcard.com/badge/joyent/conch-shell)](https://goreportcard.com/report/joyent/conch-shell)

# Getting The App

## Binaries

Releases are available over at https://github.com/joyent/conch-shell/releases
for a handful of platforms, including macOS, Linux, and Solaris/SmartOS.

## Docker

Images are available on Docker Hub ( https://hub.docker.com/r/joyentbuildops/conch-shell/ ).
To sucessfully run the app, the app needs a persistent `.conch.json` file
mounted in `/root`. An example run line is:

```
docker run --rm -it -v /home/user/.conch.json:/root/.conch.json joyent/conch-shell:latest profile ls
```

## Joyent Employees

The latest production release at
https://github.com/joyent/conch-shell/releases/latest is certified against the
production instance.

The most recent release is *not* certified for production yet (thus the
'pre-release' tag) and probably works best against the staging instance. Grab
it at https://github.com/joyent/conch-shell/releases

# Notes

## SSL Certs

Go makes a lot of assumptions about a user's runtime environment. One assumption
is that a system holds SSL certs in a set of default directories which are
hardcoded into the go runtime by platform. If the user's runtime differs from
go's expectation, the user will receive a message like `x509: failed to load
system roots and no roots provided`.

To set a custom location for SSL certs, one can specify `SSL_CERT_DIR` or
`SSL_CERT_FILE` in the environment before running conch shell.

For instance: `SSL_CERT_FILE=/opt/certs.pem conch login`


# Copyright / License

Copyright Joyent Inc

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.

# Building

## Manual

* `go get github.com/joyent/conch-shell`. This installs the code at
  `$GOPATH/src/github.com/joyent/conch-shell`
* In the `conch-shell` checkout:
	* `make tools` - Install the necessary build tools
	* `make` - Build the application

*Always* use the Makefile to build the app, rather than `go build`. The
Makefile passes necessary build vars into the app.

## Docker

Joyent's test and release process uses Docker. If you'd like to use that
process as well, use the following Makefile targets:

* `make docker_test` - Copies the local source code into the container, builds
  the app, runs the test suite, and then runs `conch version` to verify basic
  functionality

* `make docker_release` - Checks the provided version from git and
  executes `make release`, dropping the results in the local
  `release` directory.

## Reproducible Builds

As of 1841c57, our build process no longer inserts local values that break
reproducible builds. Using docker and our `make docker_release` process, one
should be able to reproduce our builds and validate via checksum.

However, at time of writing, go itself does not support reproducible builds
when `GOPATH` changes, since it embeds that path in the binary for "debugging"
purposes. It is not possible to reproduce our builds outside the docker build
environment. 

This issue is being tracked at https://github.com/golang/go/issues/16860

