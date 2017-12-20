// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package reports

import (
	"github.com/joyent/conch-shell/util"
	"gopkg.in/jawher/mow.cli.v1"
)

func Init(app *cli.Cli) {
	app.Command(
		"reports rep",
		"Various read-only reports",
		func(cmd *cli.Cmd) {
			cmd.Before = util.BuildApiAndVerifyLogin

			cmd.Command(
				"mbo_hardware_failures mbo mhf",
				"MBO hardware failure report",
				mboHardwareFailures,
			)

			cmd.Command(
				"mbo_graphs",
				"Sets up a local webserver that provides various graphs about the mbo data",
				mboHardwareFailureGraphListener,
			)
		},
	)
}
