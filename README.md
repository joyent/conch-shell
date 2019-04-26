# What

`conch` is a CLI for accessing the [Conch](https://github.com/joyent/conch) API.

# Documentation

Documentation, including the build process, can be found
[here](https://joyent.github.io/conch-shell)

# Notes

## Joyent Employees

The latest production release at
https://github.com/joyent/conch-shell/releases/latest is certified against the
production instance.

The most recent release is *not* certified for production yet (thus the
'pre-release' tag) and probably works best against the staging instance. Grab
it at https://github.com/joyent/conch-shell/releases

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

