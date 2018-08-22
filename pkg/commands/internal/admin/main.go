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
		forceOpt = app.BoolOpt("force", false, "Perform destructive actions")
	)
	app.Spec = "--force"

	app.Action = func() {
		if !*forceOpt {
			return
		}

		if err := util.API.RevokeUserTokens(UserEmail); err != nil {
			util.Bail(err)
		}

		if !util.JSON {
			fmt.Printf("Tokens revoked for %s.\n", UserEmail)
		}
	}
}

func deleteUser(app *cli.Cmd) {
	var (
		forceOpt       = app.BoolOpt("force", false, "Perform destructive actions")
		clearTokensOpt = app.BoolOpt("clear-tokens", false, "Purge the user's API tokens")
	)
	app.Spec = "--force [OPTIONS]"

	app.Action = func() {
		if !*forceOpt {
			return
		}

		if err := util.API.DeleteUser(UserEmail, *clearTokensOpt); err != nil {
			util.Bail(err)
		}

		if !util.JSON {
			fmt.Println("User " + UserEmail + " deleted.")
		}
	}
}

func createUser(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.CreateUser(UserEmail, "", ""); err != nil {
			util.Bail(err)
		}
		if !util.JSON {
			fmt.Println("User " + UserEmail + " created. An email has been sent containing their new password")
		}
	}
}

func resetUserPassword(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.ResetUserPassword(UserEmail); err != nil {
			util.Bail(err)
		}
		if !util.JSON {
			fmt.Println("The password for " + UserEmail + " has been reset. An email has been sent containing their new password")
		}
	}
}
