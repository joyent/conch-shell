// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package admin

import (
	"fmt"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
)

func revokeTokens(app *cli.Cmd) {
	var (
		userOpt  = app.StringOpt("user", "", "UUID or email address of user")
		forceOpt = app.BoolOpt("force", false, "Perform destructive actions")
	)
	app.Spec = "--user --force"

	app.Action = func() {
		if *forceOpt {
			if *userOpt == "" {
				return
			}
			if err := util.API.RevokeUserTokens(*userOpt); err != nil {
				util.Bail(err)
			}

			if !util.JSON {
				fmt.Printf("Tokens revoked for %s.\n", *userOpt)
			}
		}
	}
}
