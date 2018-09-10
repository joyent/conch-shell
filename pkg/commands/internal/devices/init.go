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
)

// DeviceSerial represents the serial of the device, gathered from the parent
// command
var DeviceSerial string

// DeviceSettingName is the name of the device setting being acted upon
var DeviceSettingName string

// DeviceTagName is the name of the device tag being acted upon
var DeviceTagName string

// Init loads up all the device related commands
func Init(app *cli.Cli) {
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
				"ipmi",
				"Get the IPMI address for a single device",
				getIPMI,
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

			// TAGS
			cmd.Command(
				"tags",
				"Get the tags for a single device",
				getTags,
			)

			cmd.Command(
				"tag",
				"Get the value of a single tag for a single device",
				func(cmd *cli.Cmd) {
					var tagNameArg = cmd.StringArg(
						"NAME",
						"",
						"The name of the tag",
					)
					cmd.Spec = "NAME"

					cmd.Before = func() {
						DeviceTagName = *tagNameArg
					}

					cmd.Command(
						"get",
						"Get a particular device tag",
						getTag,
					)

					cmd.Command(
						"set",
						"Set a particular device tag",
						setTag,
					)

					cmd.Command(
						"delete rm",
						"Delete a particular device tag",
						deleteTag,
					)
				},
			)
		},
	)
}
