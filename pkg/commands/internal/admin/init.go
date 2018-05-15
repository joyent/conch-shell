// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package admin contains administrative commands for the conch api
package admin

import (
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
)

// Init loads up the commands
func Init(app *cli.Cli) {
	app.Command(
		"admin",
		"Commands for various server-side administrative tasks",
		func(cmd *cli.Cmd) {
			cmd.Before = util.BuildAPIAndVerifyLogin

			cmd.Command(
				"revoke-tokens",
				"Revoke the auth tokens for a given user",
				revokeTokens,
			)
		},
	)
}
