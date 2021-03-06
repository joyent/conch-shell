// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package workspaces contains commands for dealing with objects tied to a
// workspace
package workspaces

import (
	"errors"
	"fmt"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/joyent/conch-shell/pkg/util"
)

// WorkspaceUUID is the UUID of the workspace we're looking at, as gathered by
// the parent command
var WorkspaceUUID uuid.UUID

// RelayID is the ID of the relay we're looking at , as gathered by the parent
// command
var RelayID string

// RackUUID is the UUID of the rack we're working with, as gathered by the
// parent command
var RackUUID uuid.UUID

// Init loads up the commands dealing with workspaces
func Init(app *cli.Cli) {
	app.Command(
		"workspaces wss",
		"Get a list of all workspaces",
		getAll,
	)
	app.Command(
		"workspace ws",
		"Commands for dealing with a single workspace",
		func(cmd *cli.Cmd) {

			var workspaceIDStr = cmd.StringArg("ID", "", "The UUID or string name of the workspace")

			cmd.Spec = "[ID]"

			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
				var newUUID uuid.UUID
				if len(*workspaceIDStr) > 0 {
					newUUID, _ = util.MagicWorkspaceID(*workspaceIDStr)
					if uuid.Equal(newUUID, uuid.UUID{}) {
						util.Bail(fmt.Errorf("workspace %s does not exist or you do not have permission to access it", *workspaceIDStr))
					}
					WorkspaceUUID = newUUID
					return
				}
				if uuid.Equal(util.ActiveProfile.WorkspaceUUID, uuid.UUID{}) {
					util.Bail(errors.New("no workspace was found in the active profile"))
				}
				WorkspaceUUID = util.ActiveProfile.WorkspaceUUID
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
				"Subcommands that deal with an individual rack",
				func(cmd *cli.Cmd) {
					var rackIDStr = cmd.StringArg("ID", "", "The rack ID")

					cmd.Before = func() {
						var err error
						RackUUID, err = util.MagicWorkspaceRackID(WorkspaceUUID, *rackIDStr)
						if err != nil {
							util.Bail(err)
						}
					}

					cmd.Spec = "ID"

					cmd.Command(
						"assign",
						"Assign devices to slots in this rack using JSON artifacts",
						assignRack,
					)

					cmd.Command(
						"assignments",
						"Dump a JSON extract of the devices assigned to this rack's slots. Intended for use with 'assign'",
						assignmentsRack,
					)

					cmd.Command(
						"get",
						"Get details about a single rack in a workspace",
						getRack,
					)

					cmd.Command(
						"add",
						"Add a single rack to a workspace",
						addRack,
					)

					cmd.Command(
						"remove delete rm",
						"Remove a single rack from a workspace",
						deleteRack,
					)

				},
			)

			cmd.Command(
				"relays",
				"Get a list of relays for a single workspace",
				getRelays,
			)

			cmd.Command(
				"subs subworkspaces ws",
				"Get a list of subworkspaces for a single workspace",
				getSubs,
			)

			cmd.Command(
				"add-user add invite",
				"Add an existing user to this workspace",
				addUser,
			)

			cmd.Command(
				"remove-user",
				"Remove a user from this workspace",
				removeUser,
			)

			cmd.Command(
				"create",
				"Create various items inside the given workspace",
				func(cmd *cli.Cmd) {
					cmd.Command(
						"subworkspace sub",
						"Create a subworkspace",
						createSubWorkspace,
					)
				},
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
