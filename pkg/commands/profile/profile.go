//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package profile

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/joyent/conch-shell/pkg/config"
	"github.com/joyent/conch-shell/pkg/util"
)

func newProfile(app *cli.Cmd) {
	var (
		nameOpt      = app.StringOpt("name", "", "Profile name. Must be unique")
		overwriteOpt = app.BoolOpt("overwrite force", false, "Overwrite any profile with a matching name")
		workspaceOpt = app.StringOpt("workspace ws", "", "Default workspace")

		tokenOpt = app.StringOpt("token", "", "The API token, generally provided by the web UI")

		envOpt = app.StringOpt("environment env", "production", "Specify the environment: production, staging, development (provide URL in the --url parameter)")
		urlOpt = app.StringOpt("url", "", "If the environment is 'development', this defines the API URL. Ignored otherwise")
	)

	app.Spec = "--name --token [OPTIONS]"

	app.Action = func() {
		var err error
		var p *config.ConchProfile
		if _, ok := util.Config.Profiles[*nameOpt]; ok {
			if !*overwriteOpt {
				util.Bail(
					fmt.Errorf(
						"a profile already exists with name '%s'",
						*nameOpt,
					),
				)
			}

			p = util.Config.Profiles[*nameOpt]
		} else {
			p = &config.ConchProfile{}
			p.Name = *nameOpt
		}

		switch *envOpt {
		case "production":
			p.BaseURL = config.ProductionURL
		case "staging":
			p.BaseURL = config.StagingURL
		case "development":
			if *urlOpt == "" {
				util.Bail(errors.New("please provide --url"))
			}
			p.BaseURL = *urlOpt
		default:
			p.BaseURL = config.ProductionURL
		}

		/***/

		util.API = &conch.Conch{
			BaseURL: p.BaseURL,
		}

		/***/

		p.Token = config.Token(*tokenOpt)
		util.API.Token = *tokenOpt

		if ok, err := util.API.VerifyToken(); !ok {
			util.Bail(err)
		}

		if *workspaceOpt == "" {
			p.WorkspaceUUID = uuid.UUID{}
		} else {
			p.WorkspaceUUID, err = util.MagicWorkspaceID(*workspaceOpt)
			if err != nil {
				util.Bail(err)
			}

			ws, err := util.API.GetWorkspace(p.WorkspaceUUID)
			if err != nil {
				util.Bail(err)
			}

			p.WorkspaceName = ws.Name
		}

		if len(util.Config.Profiles) == 0 {
			p.Active = true
		}

		util.Config.Profiles[p.Name] = p
		util.WriteConfigForce()

		if !util.JSON {
			fmt.Printf("Done. Config written to %s\n", util.Config.Path)
		}

	}
}

func deleteProfile(app *cli.Cmd) {
	var (
		nameArg = app.StringArg("NAME", "", "Name of the profile to delete")
	)

	app.Spec = "NAME"

	app.Action = func() {
		delete(util.Config.Profiles, *nameArg)
		switch len(util.Config.Profiles) {
		case 0:
			fmt.Println("WARNING: No profiles remain")
		case 1:
			for _, prof := range util.Config.Profiles {
				fmt.Printf("Only one profile remains. Setting profile '%s' to active.\n", prof.Name)
				prof.Active = true
				break
			}
		}

		util.WriteConfigForce()
		if !util.JSON {
			fmt.Printf("Done. Config written to %s\n", util.Config.Path)
		}

	}

}

func listProfiles(app *cli.Cmd) {
	app.Action = func() {
		table := util.GetMarkdownTable()

		if util.JSON {
			j, err := json.Marshal(util.Config.Profiles)
			if err != nil {
				util.Bail(err)
			}
			fmt.Println(string(j))
			return
		}

		table.SetHeader([]string{
			"Active",
			"Name",
			"User",
			"Workspace Name",
			"API URL",
		})

		for _, prof := range util.Config.Profiles {
			active := ""
			if prof.Active {
				if util.IgnoreConfig {
					active = "*?"
				} else {
					active = "*"
				}
			}
			workspaceName := ""
			if !uuid.Equal(prof.WorkspaceUUID, uuid.UUID{}) {
				if len(prof.WorkspaceName) > 0 {
					workspaceName = prof.WorkspaceName
				}
			}

			table.Append([]string{
				active,
				prof.Name,
				prof.User,
				workspaceName,
				prof.BaseURL,
			})
		}
		table.Render()
		if util.IgnoreConfig {
			fmt.Println("\n? The active profile has been overridden by the use of a token")
		}
	}
}

func setWorkspace(app *cli.Cmd) {
	var (
		workspaceArg = app.StringArg("ID", "", "Workspace name or ID")
	)
	app.Spec = "ID"

	app.Before = func() {
		util.BuildAPIAndVerifyLogin()
	}

	app.Action = func() {
		if util.ActiveProfile == nil {
			util.Bail(errors.New("there is no active profile. Please use 'profile set active' to mark a profile as active"))
		}

		workspaceUUID, err := util.MagicWorkspaceID(*workspaceArg)
		if err != nil {
			util.Bail(err)
		}

		ws, err := util.API.GetWorkspace(workspaceUUID)
		if err != nil {
			util.Bail(err)
		}

		util.ActiveProfile.WorkspaceUUID = ws.ID
		util.ActiveProfile.WorkspaceName = ws.Name

		util.WriteConfigForce()
		if !util.JSON {
			fmt.Printf("Done. Config written to %s\n", util.Config.Path)
		}

	}
}

func setActive(app *cli.Cmd) {
	var (
		profileArg = app.StringArg("NAME", "", "Profile name")
	)
	app.Spec = "NAME"

	app.Action = func() {
		if _, ok := util.Config.Profiles[*profileArg]; ok {
			for _, prof := range util.Config.Profiles {
				if prof.Name == *profileArg {
					prof.Active = true
				} else {
					prof.Active = false
				}
			}
		} else {
			util.Bail(
				fmt.Errorf("profile '%s' does not exist", *profileArg),
			)
		}

		util.WriteConfigForce()
		if !util.JSON {
			fmt.Printf("Done. Config written to %s\n", util.Config.Path)
		}

	}
}

func revokeJWT(app *cli.Cmd) {
	var (
		forceOpt   = app.BoolOpt("force", false, "Perform destructive actions")
		revokeAuth = app.BoolOpt("auth-only", false, "Revoke auth tokens, not API tokens. This will force you to log in again on the website")
		tokenAuth  = app.BoolOpt("tokens-only", false, "Revoke all API tokens. This will likely break all your automations and your ability to continue using the shell so use this carefully")
		allAuth    = app.BoolOpt("all", false, "The nuclear option. Revoke all auth *and* API tokens, forcing you to login again *and* to generate new API tokens for automation processes, including the shell. Use this very carefully")
	)
	app.Spec = "--force (--auth-only | --tokens-only | --all)"

	app.Action = func() {
		if !*forceOpt {
			return
		}
		util.BuildAPI()

		if *allAuth {
			if err := util.API.RevokeMyTokensAndLogins(); err != nil {
				util.Bail(err)
			}

			if !util.JSON {
				fmt.Println("Login and API tokens revoked")
			}
			return
		}

		if *revokeAuth {
			if err := util.API.RevokeMyLogins(); err != nil {
				util.Bail(err)
			}

			if !util.JSON {
				fmt.Println("Login tokens revoked")
			}
			return
		}
		if *tokenAuth {
			if err := util.API.RevokeMyTokens(); err != nil {
				util.Bail(err)
			}

			if !util.JSON {
				fmt.Println("API tokens revoked")
			}
			return
		}
	}
}

func setToken(cmd *cli.Cmd) {
	var tokenArg = cmd.StringArg("TOKEN", "", "An API token")
	cmd.Spec = "TOKEN"

	cmd.Action = func() {
		if util.ActiveProfile == nil {
			util.Bail(errors.New("there is no active profile. Please use 'profile set active' to mark a profile as active"))
		}

		util.ActiveProfile.Token = config.Token(*tokenArg)
		util.Token = *tokenArg

		util.WriteConfigForce()
	}

}
