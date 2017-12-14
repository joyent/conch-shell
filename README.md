# What

`conch` is a CLI for accessing the [Conch](https://github.com/joyent/conch) API.

# Copyright / License

Copyright 2017, Joyent Inc

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.

# Getting The App

A dev snapshot is available at http://us-east.manta.joyent.com/sungo/public/conch-shell 
and will be updated whenever @sungo decides that's a good idea.

Official releases are available at https://github.com/joyent/conch-shell/releases

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

* Install [Glide](https://glide.sh/)
* `cd $GOPATH/src/github.com/joyent/conch-shell`
* Run `make`
* Run `./conch`

# Notes

*Always* use the Makefile to build the app. The Makefile passes necessary build
vars into the app. 

Before committing or sending a PR, run `make sane`. This will run `go vet` and
`gofmt`.

