package workspaces

import (
	"github.com/joyent/conch-shell/pkg/util"
	"gopkg.in/jawher/mow.cli.v1"
)

func inviteUser(app *cli.Cmd) {
	var (
		emailArg = app.StringArg("EMAIL", "", "The email address of the user to be invited")
		roleArg  = app.StringOpt("role", "Read-only", "The role for the new user")
	)

	app.Spec = "EMAIL [OPTIONS]"

	app.Action = func() {
		err := util.API.InviteUser(
			WorkspaceUUID,
			*emailArg,
			*roleArg,
		)

		if err != nil {
			util.Bail(err)
		}
	}
}
