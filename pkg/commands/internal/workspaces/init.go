// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package workspaces contains commands for dealing with objects tied to a
// workspace
package workspaces

import (
	"github.com/joyent/conch-shell/pkg/util"
	"gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// WorkspaceUUID is the UUID of the workspace we're looking at, as gathered by
// the parent command
var WorkspaceUUID uuid.UUID

// RelayID is the ID of the relay we're looking at , as gathered by the parent
// command
var RelayID string

// Init loads up the commands dealing with workspaces
func Init(app *cli.Cli) {
	app.Command(
		"workspaces wss",
		"Commands for dealing with all workspaces",
		getAll,
	)
	app.Command(
		"workspace ws",
		"Commands for dealing with a single workspace",
		func(cmd *cli.Cmd) {

			var workspaceIDStr = cmd.StringArg("ID", "", "The UUID or string name of the workspace")

			cmd.Spec = "ID"

			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()

				// It's a little weird to not use := below. The problem is that
				// WorkspaceUuid is a global. If we use :=, because go can be a
				// bit weird about scoping, we get a proper err but also a
				// locally scoped version of WorkspaceUuid. If we declare err
				// separately and use =, it all works out.
				var err error
				WorkspaceUUID, err = util.MagicWorkspaceID(*workspaceIDStr)
				if err != nil {
					util.Bail(err)
				}
			}

			cmd.Command(
				"get",
				"Get details of a single workspace",
				getOne,
			)

			cmd.Command(
				"users",
				"Get a list of users for a single workspace",
				getUsers,
			)

			cmd.Command(
				"devices",
				"Get a list of devices for a single workspace",
				getDevices,
			)

			cmd.Command(
				"racks",
				"Get a list of racks for a single workspace",
				getRacks,
			)

			cmd.Command(
				"rack",
				"Get details about a single rack in a workspace",
				getRack,
			)

			cmd.Command(
				"relays",
				"Get a list of relays for a single workspace",
				getRelays,
			)

			cmd.Command(
				"rooms",
				"Get a list of rooms for a single workspace",
				getRooms,
			)

			cmd.Command(
				"subs subworkspaces ws",
				"Get a list of subworkspaces for a single workspace",
				getSubs,
			)

			cmd.Command(
				"health",
				"Get a summary of the health for a single workspace",
				getHealth,
			)
			cmd.Command(
				"failures",
				"Get failure data for a single workspace",
				getFailures,
			)

			cmd.Command(
				"invite",
				"Invite a user to this workspace",
				inviteUser,
			)

			cmd.Command(
				"relay",
				"Commands for a single relay in a workspace",
				func(cmd *cli.Cmd) {
					var relayIDStr = cmd.StringArg("ID", "", "The relay ID")
					cmd.Before = func() {
						RelayID = *relayIDStr
					}

					cmd.Spec = "ID"
					cmd.Command(
						"devices",
						"Get a list of devices for a given relay",
						getRelayDevices,
					)
				},
			)
		},
	)
}
