// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package updater contains commands related to updating conch
package updater

import (
	"gopkg.in/jawher/mow.cli.v1"
)

// Init loads up the commands dealing with updating
func Init(app *cli.Cli) {
	app.Command(
		"updater",
		"Commands around self-updating",
		func(cmd *cli.Cmd) {
			cmd.Command(
				"status",
				"Verify that we have the most recent revision",
				status,
			)

			cmd.Command(
				"changelog",
				"Display the latest changelog",
				changelog,
			)

			cmd.Command(
				"self-update",
				"Update the running application to the latest release",
				selfUpdate,
			)
		},
	)

}
