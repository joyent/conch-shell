//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package profile

import (
	"encoding/json"
	"fmt"
	"github.com/blang/semver"
	"github.com/joyent/conch-shell/pkg/config"
	"github.com/joyent/conch-shell/pkg/util"
	conch "github.com/joyent/go-conch"
	"github.com/tcnksm/go-input"
	"gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
	"strings"
)

func newProfile(app *cli.Cmd) {
	var (
		nameOpt       = app.StringOpt("name", "", "Profile name. Must be unique")
		userOpt       = app.StringOpt("user", "", "API User name")
		apiOpt        = app.StringOpt("api url", "", "API URL")
		apiVersionOpt = app.StringOpt("version api_version", "1.0.0", "API Version")
		passwordOpt   = app.StringOpt("password pass", "", "API Password")
		workspaceOpt  = app.StringOpt("workspace ws", "", "Default workspace")
		overwriteOpt  = app.BoolOpt("overwrite force", false, "Overwrite any profile with a matching name")
	)

	app.Action = func() {
		p := &config.ConchProfile{}

		password := *passwordOpt

		ui := input.DefaultUI()

		if *nameOpt == "" {
			s, err := ui.Ask(
				"Profile Name",
				&input.Options{
					Loop:      true,
					Required:  true,
					HideOrder: true,
				},
			)
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
			s, err := ui.Ask(
				"User Name",
				&input.Options{
					Loop:      true,
					Required:  true,
					HideOrder: true,
				},
			)
			if err != nil {
				util.Bail(err)
			}
			p.User = s

		} else {
			p.User = *userOpt
		}

		if password == "" {
			s, err := ui.Ask(
				"Password",
				&input.Options{
					Loop:      true,
					Required:  true,
					HideOrder: true,
					Mask:      true,
				},
			)
			if err != nil {
				util.Bail(err)
			}

			password = s
		}

		if *apiOpt == "" {
			s, err := ui.Ask(
				"API URL",
				&input.Options{
					Default:   "https://conch.joyent.us",
					Loop:      true,
					Required:  true,
					HideOrder: true,
				},
			)
			if err != nil {
				util.Bail(err)
			}

			p.BaseURL = s
		} else {
			p.BaseURL = *apiOpt
		}

		_, err := semver.Make(*apiVersionOpt)
		if err != nil {
			util.Bail(err)
		}

		api := &conch.Conch{
			BaseURL:    p.BaseURL,
			APIVersion: *apiVersionOpt,
		}
		if util.UserAgent != "" {
			api.UA = util.UserAgent
		}

		err = api.Login(p.User, password)
		if err != nil {
			util.Bail(err)
		}

		p.Session = api.Session

		util.API = api

		if *workspaceOpt == "" {
			p.WorkspaceUUID = uuid.UUID{}
		} else {
			p.WorkspaceUUID, err = util.MagicWorkspaceID(*workspaceOpt)
			if err != nil {
				util.Bail(err)
			}
		}

		if len(util.Config.Profiles) == 0 {
			p.Active = true
		}

		util.Config.Profiles[p.Name] = p
		if err := util.Config.SerializeToFile(util.Config.Path); err != nil {
			util.Bail(err)
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
		if err := util.Config.SerializeToFile(util.Config.Path); err != nil {
			util.Bail(err)
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
			"Workspace ID",
			"API URL",
			"API Version",
		})

		for _, prof := range util.Config.Profiles {
			active := ""
			if prof.Active {
				active = "*"
			}
			workspace := ""
			if !uuid.Equal(prof.WorkspaceUUID, uuid.UUID{}) {
				workspace = prof.WorkspaceUUID.String()
			}
			table.Append([]string{
				active,
				prof.Name,
				prof.User,
				workspace,
				prof.BaseURL,
				prof.APIVersion,
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
		util.ActiveProfile.WorkspaceUUID = workspaceUUID
		if err := util.Config.SerializeToFile(util.Config.Path); err != nil {
			util.Bail(err)
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
				fmt.Errorf("Profile '%s' does not exist", *profileArg),
			)
		}
		if err := util.Config.SerializeToFile(util.Config.Path); err != nil {
			util.Bail(err)
		}
	}
}

func setAPIVersion(app *cli.Cmd) {
	var (
		versionArg = app.StringArg("VERSION", "", "SemVer version string")
	)
	app.Spec = "VERSION"

	app.Action = func() {

		apiVer := strings.TrimLeft(*versionArg, "v")
		_, err := semver.Make(apiVer)
		if err != nil {
			util.Bail(err)
		}

		util.ActiveProfile.APIVersion = apiVer
		if err := util.Config.SerializeToFile(util.Config.Path); err != nil {
			util.Bail(err)
		}
	}

}
