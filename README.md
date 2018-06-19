# What

`conch` is a CLI for accessing the [Conch](https://github.com/joyent/conch) API.

[![Go Report Card](https://goreportcard.com/badge/joyent/conch-shell)](https://goreportcard.com/report/joyent/conch-shell) [![Travis-CI](https://travis-ci.org/joyent/conch-shell.svg?branch=master)](https://travis-ci.org/joyent/conch-shell)

# Getting The App

Releases are available over at https://github.com/joyent/conch-shell/releases
for a handful of platforms, including macOS, Linux, and Solaris/SmartOS.

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

* `go get github.com/joyent/conch-shell`. This installs the code at
  `$GOPATH/src/github.com/joyent/conch-shell`
* In the `conch-shell` checkout:
	* `make tools` - Install the necessary build tools
	* `make` - Build the application

# Notes

*Always* use the Makefile to build the app. The Makefile passes necessary build
vars into the app. 

