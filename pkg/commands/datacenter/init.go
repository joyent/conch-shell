// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package datacenter

import (
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/joyent/conch-shell/pkg/util"
)

// GdcUUID is the UUID of the global datacenter provided by the user
var GdcUUID uuid.UUID

// Init loads up the commands
func Init(app *cli.Cli) {

	app.Command(
		"datacenters dcs",
		"Operate on all datacenters",
		func(cmd *cli.Cmd) {
			cmd.Before = util.BuildAPIAndVerifyLogin
			cmd.Command(
				"get",
				"Get all datacenters",
				dcGetAll,
			)

			cmd.Command(
				"create",
				"Create a datacenter",
				dcCreate,
			)
		},
	)

	app.Command(
		"datacenter dc",
		"Operate on individual datacenters",
		func(cmd *cli.Cmd) {
			var gdcIDStr = cmd.StringArg("ID", "", "The UUID of the datacenter")

			cmd.Spec = "ID"
			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
				id, err := util.MagicDatacenterID(*gdcIDStr)
				if err != nil {
					util.Bail(err)
				}
				GdcUUID = id
			}

			cmd.Command(
				"get",
				"Get a datacenter",
				dcGet,
			)

			cmd.Command(
				"delete rm",
				"Delete a datacenter",
				dcDelete,
			)

			cmd.Command(
				"update",
				"Update a datacenter",
				dcUpdate,
			)

			cmd.Command(
				"rooms",
				"Get all rooms assigned to a datacenter",
				dcGetAllRooms,
			)

			cmd.Command(
				"layout-tree",
				"Get a tree of the datacenter, its rooms, racks, and layouts",
				dcAllTheThingsTree,
			)
		},
	)
}
