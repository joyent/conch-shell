// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package global contains commands that operate on structures in the global
// domain, rather than a workspace. API "global admin" access level is required
// for these commands.
package global

import (
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/joyent/conch-shell/pkg/util"
)

// GdcUUID is the UUID of the global datacenter provided by the user
var GdcUUID uuid.UUID

// GRoomUUID is the UUID of the global datacenter room provided by the user
var GRoomUUID uuid.UUID

// GRackUUID is the UUID of the global rack provided by the user
var GRackUUID uuid.UUID

// GRoleUUID is the UUID of the global rack role provided by the user
var GRoleUUID uuid.UUID

// GLayoutUUID is the UUID of the global rack layout provided by the user
var GLayoutUUID uuid.UUID

// Init loads up the commands
func Init(app *cli.Cli) {
	app.Command(
		"global system",
		"Execute commands against objects without concern for workspaces. System admin access is required.",
		func(cmd *cli.Cmd) {
			cmd.Before = util.BuildAPIAndVerifyLogin

			cmd.Command(
				"datacenters dcs",
				"Operate on all datacenters",
				func(dcs *cli.Cmd) {
					dcs.Command(
						"get",
						"Get all datacenters",
						dcGetAll,
					)

					dcs.Command(
						"create",
						"Create a datacenter",
						dcCreate,
					)
				},
			)

			cmd.Command(
				"datacenter dc",
				"Operate on individual datacenters",
				func(dc *cli.Cmd) {
					var gdcIDStr = dc.StringArg("ID", "", "The UUID of the datacenter")

					dc.Spec = "ID"
					dc.Before = func() {
						id, err := util.MagicDatacenterID(*gdcIDStr)
						if err != nil {
							util.Bail(err)
						}
						GdcUUID = id
					}

					dc.Command(
						"get",
						"Get a datacenter",
						dcGet,
					)

					dc.Command(
						"delete rm",
						"Delete a datacenter",
						dcDelete,
					)

					dc.Command(
						"update",
						"Update a datacenter",
						dcUpdate,
					)

					dc.Command(
						"rooms",
						"Get all rooms assigned to a datacenter",
						dcGetAllRooms,
					)

					dc.Command(
						"layout-tree",
						"Get a tree of the datacenter, its rooms, racks, and layouts",
						dcAllTheThingsTree,
					)
				},
			)
			/////////////////////////////////
			cmd.Command(
				"rooms rs",
				"Operate on all rooms",
				func(rs *cli.Cmd) {
					rs.Command(
						"get",
						"Get all rooms",
						roomGetAll,
					)

					rs.Command(
						"create",
						"Create a room",
						roomCreate,
					)
				},
			)

			cmd.Command(
				"room r",
				"Operate on individual rooms",
				func(r *cli.Cmd) {
					var roomIDStr = r.StringArg("ID", "", "The UUID of the room")

					r.Spec = "ID"
					r.Before = func() {
						id, err := util.MagicGlobalRoomID(*roomIDStr)
						if err != nil {
							util.Bail(err)
						}
						GRoomUUID = id
					}

					r.Command(
						"get",
						"Get a room",
						roomGet,
					)

					r.Command(
						"delete rm",
						"Delete a room",
						roomDelete,
					)

					r.Command(
						"update",
						"Update a room",
						roomUpdate,
					)

					r.Command(
						"racks",
						"Get all racks assigned to the room",
						roomGetAllRacks,
					)

				},
			)

			/////////////////////////////////
			cmd.Command(
				"racks rks",
				"Operate on all racks",
				func(rs *cli.Cmd) {
					rs.Command(
						"get",
						"Get all racks",
						rackGetAll,
					)

					rs.Command(
						"create",
						"Create a rack",
						rackCreate,
					)
				},
			)

			cmd.Command(
				"rack rk",
				"Operate on individual racks",
				func(r *cli.Cmd) {
					var rackIDStr = r.StringArg("ID", "", "The UUID of the rack")

					r.Spec = "ID"
					r.Before = func() {
						id, err := util.MagicGlobalRackID(*rackIDStr)
						if err != nil {
							util.Bail(err)
						}
						GRackUUID = id
					}

					r.Command(
						"get",
						"Get a rack",
						rackGet,
					)

					r.Command(
						"delete rm",
						"Delete a rack",
						rackDelete,
					)

					r.Command(
						"update",
						"Update a rack",
						rackUpdate,
					)

					r.Command(
						"layout",
						"Commands for dealing with the rack's layout",
						func(l *cli.Cmd) {
							l.Command(
								"get",
								"Get the rack's layout",
								rackLayout,
							)

							l.Command(
								"import",
								"Import a layout for this rack",
								rackImportLayout,
							)

							l.Command(
								"export",
								"Export the layout for this rack",
								rackExportLayout,
							)
						},
					)
				},
			)

			/////////////////////////////////
			cmd.Command(
				"roles ros",
				"Operate on all roles",
				func(rs *cli.Cmd) {
					rs.Command(
						"get",
						"Get all roles",
						roleGetAll,
					)

					rs.Command(
						"create",
						"Create a role",
						roleCreate,
					)
				},
			)

			cmd.Command(
				"role ro",
				"Operate on individual roles",
				func(r *cli.Cmd) {
					var roleIDStr = r.StringArg("ID", "", "The UUID of the role")

					r.Spec = "ID"
					r.Before = func() {
						id, err := util.MagicGlobalRackRoleID(*roleIDStr)
						if err != nil {
							util.Bail(err)
						}
						GRoleUUID = id
					}

					r.Command(
						"get",
						"Get a role",
						roleGet,
					)

					r.Command(
						"delete rm",
						"Delete a role",
						roleDelete,
					)

					r.Command(
						"update",
						"Update a role",
						roleUpdate,
					)
				},
			)

			/////////////////////////////////
			cmd.Command(
				"layouts ls",
				"Operate on all rack layouts",
				func(rs *cli.Cmd) {

					rs.Command(
						"create",
						"Create a layout",
						layoutCreate,
					)
				},
			)

			cmd.Command(
				"layout l",
				"Operate on individual layout entries",
				func(r *cli.Cmd) {
					var layoutIDStr = r.StringArg("ID", "", "The UUID of the layout")

					r.Spec = "ID"
					r.Before = func() {
						id, err := util.MagicGlobalRackLayoutSlotID(*layoutIDStr)
						if err != nil {
							util.Bail(err)
						}
						GLayoutUUID = id
					}

					r.Command(
						"get",
						"Get a layout",
						layoutGet,
					)

					r.Command(
						"delete rm",
						"Delete a layout",
						layoutDelete,
					)

					r.Command(
						"update",
						"Update a layout",
						layoutUpdate,
					)
				},
			)
		},
	)
}
