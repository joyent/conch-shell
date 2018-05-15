// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package relay contains commands that pertain to relay devices
package relay

import (
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
)

// RelaySerial represents the serial of the relay, gathered from the parent
// command
var RelaySerial string

// Init loads up all the device related commands
func Init(app *cli.Cli) {
	app.Command(
		"relay r",
		"Commands for dealing with a single relay",
		func(cmd *cli.Cmd) {

			var relaySerialStr = cmd.StringArg("ID", "", "The serial of the relay")

			cmd.Spec = "ID"

			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()

				RelaySerial = *relaySerialStr
			}

			cmd.Command(
				"register",
				"Register the relay",
				register,
			)
		},
	)
	app.Command(
		"relays rs",
		"See a list of all known relays",
		getAllRelays,
	)
}
