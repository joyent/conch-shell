package workspaces

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
	"net/mail"
)

func addUser(app *cli.Cmd) {
	var (
		emailArg = app.StringArg("EMAIL", "", "The email address of the user to be added")
		roleArg  = app.StringOpt("role", "ro", "The role for the new user. Acceptable values are 'ro', 'rw', and 'admin'")
	)

	app.Spec = "EMAIL [OPTIONS]"
	app.LongDesc = `In Days Gone By, one could create a new user while adding them to a workspace.
Those days are behind us. New users must now be created via the 'admin user' interface.`

	app.Action = func() {
		var role string

		address, err := mail.ParseAddress(*emailArg)
		if err != nil {
			util.Bail(err)
		}

		email := address.Address

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

		err = util.API.AddUserToWorkspace(
			WorkspaceUUID,
			email,
			role,
		)

		if err != nil {
			util.Bail(err)
		}

		if !util.JSON {
			fmt.Println("User " + email + "has been added to workspace " + WorkspaceUUID.String() + " and they have been informed via email")
		}
	}
}

func removeUser(app *cli.Cmd) {
	var (
		emailArg = app.StringArg("EMAIL", "", "The email address of the user to be removed")
	)

	app.Spec = "EMAIL [OPTIONS]"

	app.Action = func() {
		address, err := mail.ParseAddress(*emailArg)
		if err != nil {
			util.Bail(err)
		}
		email := address.Address

		users, err := util.API.GetWorkspaceUsers(WorkspaceUUID)
		if err != nil {
			util.Bail(err)
		}

		for _, user := range users {
			if user.Email == *emailArg {
				if !uuid.Equal(user.RoleVia, uuid.UUID{}) {
					if !uuid.Equal(user.RoleVia, WorkspaceUUID) {
						ws, err := util.API.GetWorkspace(user.RoleVia)
						if err != nil {
							util.Bail(err)
						}
						util.Bail(fmt.Errorf("cannot continue. User %s has access to this workspace via the parent workspace %s", *emailArg, ws.Name))
						return
					}
				}
			}
		}

		err = util.API.RemoveUserFromWorkspace(WorkspaceUUID, email)
		if err != nil {
			util.Bail(err)
		}
		if !util.JSON {
			fmt.Println("The user has been removed from the workspace. They will not be notified")
		}
	}

}
