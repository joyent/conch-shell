// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package devices contains commands pertaining to individual devices
package devices

import (
	"github.com/joyent/conch-shell/pkg/util"
	"gopkg.in/jawher/mow.cli.v1"
)

// DeviceSerial represents the serial of the device, gathered from the parent
// command
var DeviceSerial string

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
				"settings",
				"Get the seettings for a single device",
				getSettings,
			)
			cmd.Command(
				"setting",
				"Get the value of a single setting for a single device",
				getSetting,
			)

			cmd.Command(
				"graduate",
				"Mark a device as 'graduated'. WARNING: This is a one-way operation that cannot be undone",
				graduate,
			)

			cmd.Command(
				"triton_reboot",
				"Mark a device as rebooted into Triton. WARNING: This is a one-way operation that cannot be undone",
				tritonReboot,
			)

			cmd.Command(
				"triton_uuid",
				"Set the Triton UUID",
				setTritonUUID,
			)

			cmd.Command(
				"triton_setup",
				"Mark the device as having been setup in Triton. WARNING: This is a one-way operation that cannot be undone",
				markTritonSetup,
			)
		},
	)
}
