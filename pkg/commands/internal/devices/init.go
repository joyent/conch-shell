// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package devices contains commands pertaining to individual devices
package devices

import (
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// DeviceSerial represents the serial of the device, gathered from the parent
// command
var DeviceSerial string

// DeviceServiceUUID represents the UUID of the device service being used in
// the relevant command tree
var DeviceServiceUUID uuid.UUID

// DeviceRoleUUID represents the UUID of the device role being used in the
// relevant command tree
var DeviceRoleUUID uuid.UUID

// DeviceSettingName is the name of the device setting being active upon
var DeviceSettingName string

// Init loads up all the device related commands
func Init(app *cli.Cli) {
	app.Command(
		"device-services dss",
		"Commands for dealing with device services",
		func(cmd *cli.Cmd) {
			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
			}

			cmd.Command(
				"get",
				"Get a list all available device services",
				getAllDeviceServices,
			)

			cmd.Command(
				"create",
				"Create a new device service",
				createDeviceService,
			)
		},
	)

	app.Command(
		"device-service ds",
		"Commands for dealing with a single device service",
		func(cmd *cli.Cmd) {
			var deviceServiceIDArg = cmd.StringArg(
				"ID",
				"",
				"The ID or name of the device service",
			)
			cmd.Spec = "ID"

			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
				var err error
				DeviceServiceUUID, err = util.MagicDeviceServiceID(*deviceServiceIDArg)
				if err != nil {
					util.Bail(err)
				}
			}

			cmd.Command(
				"get",
				"Get info about a single device service",
				getOneDeviceService,
			)

			cmd.Command(
				"delete rm",
				"Delete a single device service",
				deleteDeviceService,
			)

			cmd.Command(
				"modify update",
				"Update a device service",
				modifyDeviceService,
			)

		},
	)

	app.Command(
		"device-roles drs",
		"Commands for dealing with device roles",
		func(cmd *cli.Cmd) {
			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
			}

			cmd.Command(
				"get",
				"Get a list of all available device roles",
				getAllDeviceRoles,
			)

			cmd.Command(
				"create",
				"Create a new device role",
				createDeviceRole,
			)
		},
	)
	app.Command(
		"device-role dr",
		"Commands for dealing with a single device role",
		func(cmd *cli.Cmd) {
			var deviceRoleIDArg = cmd.StringArg(
				"ID",
				"",
				"The ID of the device role",
			)
			cmd.Spec = "ID"

			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()

				var err error
				DeviceRoleUUID, err = util.MagicDeviceRoleID(*deviceRoleIDArg)
				if err != nil {
					util.Bail(err)
				}
			}

			cmd.Command(
				"get",
				"Get info about a single device role",
				getOneDeviceRole,
			)

			cmd.Command(
				"delete rm",
				"Delete a single device role",
				deleteDeviceRole,
			)

			cmd.Command(
				"modify update",
				"Update a device role",
				modifyDeviceRole,
			)

			cmd.Command(
				"add-service as",
				"Add a device service to a role",
				addServiceToDeviceRole,
			)

			cmd.Command(
				"remove-service rms",
				"Remove a device service from a role",
				removeServiceFromDeviceRole,
			)

		},
	)
	app.Command(
		"device d",
		"Commands for dealing with a single device",
		func(cmd *cli.Cmd) {

			var deviceSerialStr = cmd.StringArg("ID", "", "The serial of the device")

			cmd.Spec = "ID"

			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()

				DeviceSerial = *deviceSerialStr
			}

			cmd.Command(
				"get",
				"Get the details of a single device",
				getOne,
			)

			cmd.Command(
				"location",
				"Get the location of a single device",
				getLocation,
			)

			cmd.Command(
				"settings",
				"Get the settings for a single device",
				getSettings,
			)

			cmd.Command(
				"setting",
				"Get the value of a single setting for a single device",
				func(cmd *cli.Cmd) {
					var settingNameArg = cmd.StringArg(
						"NAME",
						"",
						"The name of the setting",
					)
					cmd.Spec = "NAME"

					cmd.Before = func() {
						DeviceSettingName = *settingNameArg
					}

					cmd.Command(
						"get",
						"Get a particular device setting",
						getSetting,
					)

					cmd.Command(
						"set",
						"Set a particular device setting",
						setSetting,
					)

					cmd.Command(
						"delete rm",
						"Delete a particular device setting",
						deleteSetting,
					)
				},
			)

			cmd.Command(
				"graduate",
				"Mark a device as 'graduated'. WARNING: This is a one-way operation that cannot be undone",
				graduate,
			)

			cmd.Command(
				"asset_tag",
				"Subcommands that deal with asset tags",
				func(cmd *cli.Cmd) {

					cmd.Command(
						"get",
						"get a device's asset tag",
						getAssetTag,
					)

					cmd.Command(
						"set",
						"Set a device's asset tag",
						setAssetTag,
					)
				},
			)

			cmd.Command(
				"report",
				"Get the latest recorded device report as JSON",
				getReport,
			)

			cmd.Command(
				"triton",
				"Subcommands that deal with various Triton related settings",
				func(cmd *cli.Cmd) {
					cmd.Command(
						"reboot",
						"Mark a device as rebooted into Triton. WARNING: This is a one-way operation that cannot be undone",
						tritonReboot,
					)

					cmd.Command(
						"uuid",
						"Set the Triton UUID. WARNING: This is a one-way operation that cannot be undone",
						setTritonUUID,
					)

					cmd.Command(
						"setup",
						"Mark the device as having been setup in Triton. WARNING: This is a one-way operation that cannot be undone",
						markTritonSetup,
					)
				},
			)
		},
	)
}
