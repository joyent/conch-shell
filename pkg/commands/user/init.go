// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package user contains command pertaining to the active Conch user
package user

import (
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
)

var (
	// SettingName is where the 'setting' parent command stores the name of the
	// setting we're dealing with
	SettingName string
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
				"profile",
				"View your Conch profile",
				getProfile,
			)

			cmd.Command(
				"settings",
				"Get the settings for the current user",
				getSettings,
			)

			cmd.Command(
				"setting",
				"Commands for dealing with a single setting for the current user",
				func(cmd *cli.Cmd) {

					var settingNameArg = cmd.StringArg("NAME", "", "The string name of a setting")

					cmd.Spec = "NAME"

					cmd.Before = func() {
						util.BuildAPIAndVerifyLogin()
						SettingName = *settingNameArg
					}

					cmd.Command(
						"get",
						"Get a setting for the current user",
						getSetting,
					)

					cmd.Command(
						"set",
						"Set a setting for the current user",
						setSetting,
					)

					cmd.Command(
						"delete",
						"Delete a setting for the current user",
						deleteSetting,
					)
				},
			)

			// The biggest use case for disabling these functions is security,
			// particularly when it comes to edge automation. It's probably a
			// bad idea for some automation on a random server to be able to
			// create and remove tokens.
			if !util.DisableApiTokenCRUD() {
				cmd.Command(
					"tokens",
					"List API tokens",
					listTokens,
				)

				cmd.Command(
					"token",
					"Operate on a single token",
					func(cmd *cli.Cmd) {
						cmd.Command(
							"remove del rm",
							"Remove an API token",
							removeToken,
						)

						cmd.Command(
							"create",
							"Create an API token",
							createToken,
						)

						cmd.Command(
							"get",
							"See information about a single API token",
							getToken,
						)
					},
				)
			}
		},
	)
}
