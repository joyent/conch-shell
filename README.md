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

# Building

* `go get github.com/joyent/conch-shell`. This installs the code at
  `$GOPATH/src/github.com/joyent/conch-shell`
* In the `conch-shell` checkout:
	* `make tools` - Install the necessary build tools
	* `make` - Build the application

## Variations

### Build against go-conch master

* Open `Gopkg.toml`. 
* Find the constraint block for go-conch
* Remove the line that begins with `version`
* Add a new line containing `branch = "master"`
* Run `dep ensure`
* Run `make`

### Build against a local checkout of go-conch

* Checkout out go-conch to `$GOPATH/src/github.com/joyent/conch-shell`
* Run `go get gopkg.in/saturi/go.uuid.v1`
* In the `conch-shell` checkout:
	* Install dependencies with `dep ensure`
	* Remove the directory `vendor/github.com/joyent`
	* Remove the directory `vendor/gopkg.in/satori`
* In the `go-conch` checkout:
	* If it exists, remove `vendor/gopkg.in/satori` 
* In the `conch-shell` checkout, run `make`
* Do *NOT* run `dep` in either checkout as this will re-install the conflicting
  dependencies


# Notes

*Always* use the Makefile to build the app. The Makefile passes necessary build
vars into the app. 

