// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package hardware contains commands related to hardware products, profiles,
// and the like
package hardware

import (
	"errors"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// ProductUUID is the UUID of the hardware product we're looking at, as
// gathered by the parent command
var ProductUUID uuid.UUID

// Init loads up the hardware commands
func Init(app *cli.Cli) {
	app.Command(
		"hardware h",
		"Commands for dealing with hardware",
		func(cmd *cli.Cmd) {
			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
			}

			cmd.Command(
				"products ps",
				"Get a list of hardware products",
				getAll,
			)

			cmd.Command(
				"product p",
				"Deal with a single hardware product",
				func(cmd *cli.Cmd) {
					var productIDStr = cmd.StringArg("ID", "", "The UUID, name, or alias of a hardware product")
					cmd.Spec = "ID"

					cmd.Before = func() {
						ProductUUID, _ = util.MagicProductID(*productIDStr)
						if uuid.Equal(ProductUUID, uuid.UUID{}) {
							util.Bail(errors.New("Could not resolve the hardware product ID, name, or alias"))
						}
					}

					cmd.Command(
						"get",
						"Get a single hardware product",
						getOne,
					)
				},
			)
		},
	)
}
