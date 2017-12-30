// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package user contains command pertaining to the active Conch user
package user

import (
	"gopkg.in/jawher/mow.cli.v1"
)

// Init loads up the user commands
func Init(app *cli.Cli) {
	app.Command(
		"user u",
		"Commands for dealing with the current user",
		func(cmd *cli.Cmd) {
			// Because login happens in here, we can't VerifyLogin blindly.
			// Everyone has to do that on their own if they way.

			cmd.Command(
				"settings",
				"Get the settings for the current user",
				getSettings,
			)

			cmd.Command(
				"setting",
				"Get a setting for the current user",
				getSetting,
			)

			cmd.Command(
				"login",
				"Log in",
				Login,
			)

		},
	)
}
