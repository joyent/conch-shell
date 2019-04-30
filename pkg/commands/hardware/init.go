// Copyright Joyent, Inc.
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
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/joyent/conch-shell/pkg/util"
)

// ProductUUID is the UUID of the hardware product we're looking at, as
// gathered by the parent command
var ProductUUID uuid.UUID

// HardwareVendorName ...
var HardwareVendorName string

// Init loads up the hardware commands
func Init(app *cli.Cli) {
	app.Command(
		"hardware h",
		"Commands for dealing with hardware products",
		func(cmd *cli.Cmd) {
			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
			}

			cmd.Command(
				"products ps",
				"Deal with hardware products",
				func(cmd *cli.Cmd) {

					cmd.Command(
						"get",
						"Get a list of all hardware products",
						getAll,
					)

					cmd.Command(
						"create",
						"Create a hardware product",
						createOne,
					)

					cmd.Command(
						"template",
						"Dumping a JSON template for a hardware product. Used in creating a new product and profile",
						dumpTemplate,
					)

					cmd.Command(
						"import",
						"Import a JSON file that defines a new hardware product",
						importNewProductJson,
					)
				},
			)

			cmd.Command(
				"product p",
				"Deal with a single hardware product",
				func(cmd *cli.Cmd) {
					var productIDStr = cmd.StringArg("ID", "", "The UUID, name, or SKU of a hardware product")
					cmd.Spec = "ID"

					cmd.Before = func() {
						ProductUUID, _ = util.MagicProductID(*productIDStr)
						if uuid.Equal(ProductUUID, uuid.UUID{}) {
							util.Bail(errors.New("could not resolve the hardware product ID, name, or SKU"))
						}
					}

					cmd.Command(
						"get",
						"Get a single hardware product",
						getOne,
					)

					cmd.Command(
						"get_specification",
						"Get the hardware specification json blob",
						getOneSpecification,
					)

					cmd.Command(
						"delete rm",
						"Delete a hardware product",
						removeOne,
					)

					cmd.Command(
						"update up",
						"Update a hardware product",
						updateOne,
					)

					cmd.Command(
						"export",
						"Dump the JSON representation of a hardware product and profile. Intended for use with 'import'",
						exportProductJson,
					)

					cmd.Command(
						"import",
						"Update an existing hardware product and profile using a JSON file",
						importChangedProductJson,
					)

				},
			)

			cmd.Command(
				"vendors vs",
				"Get a list of all hardware vendors",
				getAllVendors,
			)

			cmd.Command(
				"vendor v",
				"Deal with a hardware vendor",
				func(cmd *cli.Cmd) {
					var vendorNameStr = cmd.StringArg("NAME", "", "The name of the hardware vendor")
					cmd.Spec = "NAME"

					cmd.Before = func() {
						HardwareVendorName = *vendorNameStr
					}

					cmd.Command(
						"get",
						"Get a single vendor",
						getOneVendor,
					)

					cmd.Command(
						"create make mk",
						"Create a single vendor",
						createOneVendor,
					)

					cmd.Command(
						"delete rm ",
						"Delete a single vendor",
						deleteOneVendor,
					)
				},
			)
		},
	)
}
