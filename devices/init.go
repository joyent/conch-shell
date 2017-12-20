// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package devices

import (
	"github.com/joyent/conch-shell/util"
	"gopkg.in/jawher/mow.cli.v1"
)

var DeviceSerial string

func Init(app *cli.Cli) {
	app.Command(
		"device d",
		"Commands for dealing with a single device",
		func(cmd *cli.Cmd) {

			var device_serial_str = cmd.StringArg("ID", "", "The serial of the device")

			cmd.Spec = "ID"

			cmd.Before = func() {
				util.BuildApiAndVerifyLogin()

				DeviceSerial = *device_serial_str
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

		},
	)
}
