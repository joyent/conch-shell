// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package profile contains commands pertaining to login profiles
package profile

import (
	"gopkg.in/jawher/mow.cli.v1"
)

// Init loads up the profile commands
func Init(app *cli.Cli) {
	app.Command(
		"profile prof",
		"Commands for creating and adjusting login profiles",
		func(cmd *cli.Cmd) {
			// Because login happens in here, we can't VerifyLogin blindly.
			// Everyone has to do that on their own if they want.

			cmd.Command(
				"new create add",
				"Create a new login profile",
				newProfile,
			)

			cmd.Command(
				"delete del rm",
				"Delete a profile",
				deleteProfile,
			)
			cmd.Command(
				"list ls",
				"List all known profiles",
				listProfiles,
			)

			cmd.Command(
				"set",
				"Change profile settings",
				func(cmd *cli.Cmd) {
					cmd.Command(
						"workspace ws",
						"Set the workspace ID for the active profile",
						setWorkspace,
					)

					cmd.Command(
						"active",
						"Change which profile is active",
						setActive,
					)
				},
			)

		},
	)
}
