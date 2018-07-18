//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Bowery/prompt"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/config"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
)

func newProfile(app *cli.Cmd) {
	var (
		nameOpt      = app.StringOpt("name", "", "Profile name. Must be unique")
		userOpt      = app.StringOpt("user", "", "API User name")
		apiOpt       = app.StringOpt("api url", "", "API URL")
		passwordOpt  = app.StringOpt("password pass", "", "API Password")
		workspaceOpt = app.StringOpt("workspace ws", "", "Default workspace")
		overwriteOpt = app.BoolOpt("overwrite force", false, "Overwrite any profile with a matching name")
	)

	app.Action = func() {
		p := &config.ConchProfile{}

		password := *passwordOpt

		if *nameOpt == "" {
			s, err := prompt.Basic("Profile Name:", true)
			if err != nil {
				util.Bail(err)
			}

			p.Name = s
		} else {
			p.Name = *nameOpt
		}

		if !*overwriteOpt {
			if _, ok := util.Config.Profiles[p.Name]; ok {
				util.Bail(
					fmt.Errorf(
						"A profile already exists with name '%s'",
						p.Name,
					),
				)
			}
		}

		if *userOpt == "" {
			s, err := prompt.Basic("User Name:", true)
			if err != nil {
				util.Bail(err)
			}
			p.User = s

		} else {
			p.User = *userOpt
		}

		if password == "" {
			s, err := prompt.Password("Password:")
			if err != nil {
				util.Bail(err)
			}

			password = s
		}

		if *apiOpt == "" {
			s, err := prompt.BasicDefault("API URL:", "https://conch.joyent.us")
			if err != nil {
				util.Bail(err)
			}
			p.BaseURL = s
		} else {
			p.BaseURL = *apiOpt
		}

		api := &conch.Conch{
			BaseURL: p.BaseURL,
		}
		if util.UserAgent != "" {
			api.UA = util.UserAgent
		}

		err := api.Login(p.User, password)
		if err != nil {
			util.Bail(err)
		}

		p.Session = api.Session
		p.JWToken = api.JWToken

		util.API = api

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
		util.WriteConfig()
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

		util.WriteConfig()
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
			"Workspace ID",
			"Workspace Name",
			"API URL",
			"Expires",
		})

		for _, prof := range util.Config.Profiles {
			active := ""
			if prof.Active {
				active = "*"
			}
			workspaceUUID := ""
			workspaceName := ""
			if !uuid.Equal(prof.WorkspaceUUID, uuid.UUID{}) {
				workspaceUUID = prof.WorkspaceUUID.String()
				if len(prof.WorkspaceName) > 0 {
					workspaceName = prof.WorkspaceName
				} else {
					// BUG(sungo): This is a transition point since workspace
					// names only get set during profile creation.
					workspaceName = "**UNKNOWN**"
				}
			}
			expires := "Unknown"
			if prof.Expires > 0 {
				expires = util.TimeStr(time.Unix(int64(prof.Expires), 0))
			}

			table.Append([]string{
				active,
				prof.Name,
				prof.User,
				workspaceUUID,
				workspaceName,
				prof.BaseURL,
				expires,
			})
		}
		table.Render()
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

		util.WriteConfig()
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
				fmt.Errorf("Profile '%s' does not exist", *profileArg),
			)
		}

		util.WriteConfig()
	}
}

func refreshJWT(app *cli.Cmd) {
	app.Action = func() {
		util.BuildAPI()
		if util.ActiveProfile == nil {
			util.Bail(errors.New("No active profile. Please use 'conch profile' to create or set an active profile"))
		}

		if err := util.API.VerifyLogin(0, true); err != nil {
			util.Bail(err)
		}
		util.ActiveProfile.Session = util.API.Session
		util.ActiveProfile.JWToken = util.API.JWToken
		util.ActiveProfile.Expires = util.API.Expires
		util.WriteConfig()

		expires := time.Unix(int64(util.API.Expires), 0)
		if util.JSON {
			util.JSONOut(struct {
				Expires time.Time `json:"expires"`
			}{expires})

			return
		}

		fmt.Printf(
			"Auth for profile '%s' now expires at %s\n",
			util.ActiveProfile.Name,
			util.TimeStr(expires),
		)
	}
}

func revokeJWT(app *cli.Cmd) {
	var forceOpt = app.BoolOpt("force", false, "Perform destructive actions")
	app.Spec = "--force"

	app.Action = func() {
		if *forceOpt {
			util.BuildAPIAndVerifyLogin()
			if err := util.API.RevokeOwnTokens(); err != nil {
				util.Bail(err)
			}

			if !util.JSON {
				fmt.Println("Tokens revoked.")
			}
		}
	}
}

func relogin(app *cli.Cmd) {
	var (
		passwordOpt = app.StringOpt("password pass", "", "API Password")
	)

	app.Action = func() {
		util.BuildAPI()

		password := *passwordOpt

		if password == "" {
			s, err := prompt.Password("Password:")
			if err != nil {
				util.Bail(err)
			}

			password = s
		}

		err := util.API.Login(util.ActiveProfile.User, password)
		if err != nil {
			util.Bail(err)
		}

		util.ActiveProfile.Session = util.API.Session
		util.ActiveProfile.JWToken = util.API.JWToken

		util.WriteConfig()
	}
}
