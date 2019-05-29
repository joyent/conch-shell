// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package rack

import (
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/joyent/conch-shell/pkg/util"
)

// GRackUUID is the UUID of the global rack provided by the user
var GRackUUID uuid.UUID

// Init loads up the commands
func Init(app *cli.Cli) {
	app.Command(
		"racks rks",
		"Operate on all racks",
		func(cmd *cli.Cmd) {
			cmd.Before = util.BuildAPIAndVerifyLogin
			cmd.Command(
				"get",
				"Get all racks",
				rackGetAll,
			)

			cmd.Command(
				"create",
				"Create a rack",
				rackCreate,
			)
		},
	)

	app.Command(
		"rack rk",
		"Operate on individual racks",
		func(r *cli.Cmd) {
			var rackIDStr = r.StringArg("ID", "", "The UUID of the rack")

			r.Spec = "ID"
			r.Before = func() {
				util.BuildAPIAndVerifyLogin()
				id, err := util.MagicRackID(*rackIDStr)
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
				"phase",
				"Get the rack's phase",
				rackPhaseGet,
			)

			r.Command(
				"set",
				"Change various settings on a rack",
				func(cmd *cli.Cmd) {
					cmd.Command(
						"phase",
						"Set a rack's phase",
						rackPhaseSet,
					)
				},
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

			r.Command(
				"assign",
				"Assign devices to slots in this rack using JSON artifacts",
				rackAssign,
			)

			r.Command(
				"assignments",
				"Dump a JSON extract of the devices assigned to this rack's slots. Intended for use with 'assign'",
				rackAssignments,
			)

		},
	)

}
