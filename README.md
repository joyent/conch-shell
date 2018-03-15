# What

`conch` is a CLI for accessing the [Conch](https://github.com/joyent/conch) API.

[![Go Report Card](https://goreportcard.com/badge/joyent/conch-shell)](https://goreportcard.com/report/joyent/conch-shell) [![Travis-CI](https://travis-ci.org/joyent/conch-shell.svg?branch=master)](https://travis-ci.org/joyent/conch-shell)

# Getting The App

Releases are available over at https://github.com/joyent/conch-shell/releases
for a handful of platforms, including macOS, Linux, and Solaris/SmartOS.

# Copyright / License

Copyright Joyent Inc

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.

# Setup and Build

## Setting up Go

* Install [Go](https://golang.org/)
* If you're using the [standard go workspace
  layout](https://golang.org/doc/code.html#Workspaces) (and you really should
  be), make sure that `$GOPATH/bin` is in `$PATH`
  * `export GOPATH=$(go env GOPATH); export PATH="$GOPATH/bin:$PATH"`
  * Typically, `$GOPATH` is `~/go` but it doesn't have to be. Mine is
    `~/src/go`. The important part is the layout underneath `$GOPATH`
* Add the following incantation to `~/.gitconfig` which will cause `go get` to
  use ssh to access github rather than https. That's necessary for private repos
  like this one.

```
[url "git@github.com:"]
	insteadOf = https://github.com/
```

## Check out the code

* Run `go get github.com/joyent/conch-shell`
* The code will end up in `$GOPATH/src/github.com/joyent/conch-shell`

## Building

* Install [dep](https://golang.github.io/dep/docs/installation.html)
* `cd $GOPATH/src/github.com/joyent/conch-shell`
* Run `make`
* Run `./bin/conch`

# Notes

*Always* use the Makefile to build the app. The Makefile passes necessary build
vars into the app. 

