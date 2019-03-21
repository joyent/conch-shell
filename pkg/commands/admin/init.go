// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package admin contains administrative commands for the conch api
package admin

import (
	"net/mail"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
)

// UserEmail contains the email address of the user being operated on in
// the 'user' sub tree
var UserEmail string

// Init loads up the commands
func Init(app *cli.Cli) {
	app.Command(
		"admin",
		"Commands for various server-side administrative tasks",
		func(cmd *cli.Cmd) {
			cmd.Before = util.BuildAPIAndVerifyLogin

			cmd.Command(
				"users",
				"List all users",
				listAllUsers,
			)

			cmd.Command(
				"user",
				"Administrative commands for operating on a user",
				func(cmd *cli.Cmd) {

					var userIDStr = cmd.StringArg(
						"USER",
						"",
						"The email address of the user",
					)

					cmd.Spec = "USER"

					cmd.Before = func() {
						address, err := mail.ParseAddress(*userIDStr)
						if err != nil {
							util.Bail(err)
						}
						UserEmail = address.Address
					}

					cmd.Command(
						"get",
						"Get the basic info about a user",
						getUser,
					)

					cmd.Command(
						"revoke",
						"Revoke the auth tokens for a given user",
						revokeTokens,
					)

					cmd.Command(
						"delete rm",
						"Delete a user from conch. This *cannot* be undone",
						deleteUser,
					)

					cmd.Command(
						"create",
						"Create a new user. Does *not* assign them to a workspace",
						createUser,
					)

					cmd.Command(
						"reset",
						"Reset the password for the user",
						resetUserPassword,
					)

					cmd.Command(
						"update",
						"Update properties of the user",
						updateUser,
					)

					cmd.Command(
						"promote",
						"Promote the user to system admin",
						promoteUser,
					)

					cmd.Command(
						"demote",
						"Demote the user to a regular user",
						demoteUser,
					)

				},
			)
		},
	)
}
