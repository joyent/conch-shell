// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package profile contains commands pertaining to login profiles
package profile

import (
	"github.com/jawher/mow.cli"
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
				"refresh",
				"Refresh the auth token for the active profile",
				refreshJWT,
			)

			cmd.Command(
				"relogin",
				"Log in again, preserving all other profile data",
				relogin,
			)

			cmd.Command(
				"change-password",
				"Change the password associated with this profile",
				changePassword,
			)

			cmd.Command(
				"revoke-tokens",
				"Revoke all auth tokens. User must log in again after this.",
				revokeJWT,
			)

			cmd.Command(
				"set",
				"Change profile settings",
				func(cmd *cli.Cmd) {
					cmd.Command(
						"workspace ws",
						"Set the workspace (by name or ID) for the active profile",
						setWorkspace,
					)

					cmd.Command(
						"active",
						"Change which profile is active",
						setActive,
					)

					cmd.Command(
						"version-check vc",
						"Enable/disable version checking",
						func(cmd *cli.Cmd) {
							cmd.Command(
								"enable",
								"Enable version checking",
								enableVersionCheck,
							)

							cmd.Command(
								"disable",
								"Disable version checking",
								disableVersionCheck,
							)
						},
					)
				},
			)

		},
	)
}
