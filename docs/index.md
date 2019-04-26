# Introduction

The Conch ecosystem is designed to make the deployment of new server hardware
easier, specifically targetting equipment to be used in the Joyent
SmartDatacenter product line. Conch has two major backend systems. First, edge
software boots new hardware, upgrades firmware, performs burn-in testing, and
gathers the general state of the hardware. (This software is currently closed
source.) Second, this edge data is fed into the open source [Conch API
server](https://github.com/joyent/conch) where the data is processed, validated,
stored, and reported upon.

The Conch Shell is a CLI designed to interact with the Conch API, allowing
operators to prepare datacenter builds, monitor those builds, and perform
validation actions. While Conch has a [web
UI](https://github.com/joyent/conch-ui), the shell's intention is to allow
operators to do their work entirely in the CLI environment. Further, the conch
shell specifically targets scriptability, allowing operators to automate all
aspects of the Conch API.

```
$ conch wss

| ROLE  |                  ID                  | NAME        | DESCRIPTION      |
|-------|--------------------------------------|-------------|------------------|
| admin | 67409695-3ee3-4d45-8366-8213090b158a | GLOBAL      | Global workspace |
| admin | 8c86ff85-1c0a-413c-a544-88e46d65a370 | conch-dev   | Conch Dev        |
| admin | 6e5aacae-fc39-4fca-8245-b62c7b231ca1 | us-east-1   |                  |
| admin | 8267bc59-307d-428b-a336-e381fe512aa3 | us-west-1   |                  |
```

# User Documentation

* [How To Login](auth)
  * [Deeper Dive on API Tokens, including commands](tokens)

*The main focus for the Conch Shell is Joyent's production environments. While
the shell should operate against any instance of the Conch API, some bits of the
documentation are Joyent specific.*

*Further, this site and all its documentation is targetted towards the `master`
branch. If you need documentation for a specific release, see the `docs`
directory in the git tag of your choice.*

# Obtaining The App

## Binaries

Joyent provides pre-compiled binaries for various platforms via [Github
Releases](https://github.com/joyent/conch-shell/releases). These binaries
are built in Joyent's CI infrastructure using Docker to ensure a consistent
build environment. If you'd like us to consider adding support for your favorite
platform, file a Github issue.

For the SmartOS/Illumos users in the crowd, the solaris binary works fine in
those environments.

Support is offered via [Github
Issues](https://github.com/joyent/conch-shell/issues) for the current release
only. 


## Building The App

### Requirements

* [Go](https://golang.org/) - At the time of writing, you need Go 1.12 or higher.
  See the Dockerfile for the exact version that we're using.
* GNU Make - Our Makefile is a bit fancy and only works under GNU make

### The Build

* `make tools` will install the linters and infrastructure bits
* `make` will install dependencies, run the tests, and build the binaries. The
  new binaries will be deposited in `bin/`

Always use the Makefile. The build process, via the Makefile, adds important
information to the code base that is required for its successful operation.

Several configuration options are available only during the build process and
can be specified in the environment. See the first few lines of the Makefile.

### Building A Release

* You'll need [docker](https://www.docker.com/).
* `make release` will execute our build process via Docker, depositing the
  results in `release/`

# Copyright / License

Copyright Joyent Inc

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, you can obtain one at <http://mozilla.org/MPL/2.0/>

