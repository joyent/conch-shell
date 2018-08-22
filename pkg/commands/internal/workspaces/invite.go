package workspaces

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
)

func inviteUser(app *cli.Cmd) {
	var (
		emailArg = app.StringArg("EMAIL", "", "The email address of the user to be invited")
		roleArg  = app.StringOpt("role", "ro", "The role for the new user. Acceptable values are 'ro', 'rw', and 'admin'")
	)

	app.Spec = "EMAIL [OPTIONS]"

	app.Action = func() {
		var role string
		switch *roleArg {
		case "ro", "rw", "admin":
			role = *roleArg
		case "Read-only":
			if !util.JSON {
				fmt.Println("Role 'Read-only' is now called 'ro'. Will substitute the new value.")
			}
			role = "ro"
		case "Administrator":
			if !util.JSON {
				fmt.Println("Role 'Adminstrator' is now called 'admin'. Will substitute the new value.")
			}
			role = "admin"

		case "Integrator", "DC Operations", "Integrator Manager":
			if !util.JSON {
				fmt.Printf("Role '%s' is now called 'rw'. Will substitute the new value.\n", *roleArg)
			}
			role = "rw"

		default:
			if !util.JSON {
				fmt.Println("Unknown role name. Falling back to 'ro'")
			}
			role = "ro"
		}

		err := util.API.InviteUser(
			WorkspaceUUID,
			*emailArg,
			role,
		)

		if err != nil {
			util.Bail(err)
		}
	}
}
