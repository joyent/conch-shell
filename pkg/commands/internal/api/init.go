// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package api contains commands that allow direct API access
package api

import (
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
)

// Init loads up the commands dealing with direct api access
func Init(app *cli.Cli) {
	app.Command(
		"api",
		"Execute raw API commands",
		func(cmd *cli.Cmd) {
			cmd.Before = util.BuildAPIAndVerifyLogin
			cmd.Command(
				"get",
				"Perform an HTTP get against a provided URL",
				get,
			)
			cmd.Command(
				"delete",
				"Perform an HTTP DELETE against a provided URL",
				deleteAPI,
			)

		},
	)

}
