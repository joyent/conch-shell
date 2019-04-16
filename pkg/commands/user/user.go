// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"text/template"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
)

const userProfile = `
Name: {{.Name}}
Email: {{.Email}}

Created: {{.Created.Local}}
Last Login: {{.LastLogin.Local}}
{{if len .Workspaces}}
Workspaces:{{ range .Workspaces }}
  Name: {{.Name}}
  Role: {{.Role}}
  Description: {{.Description}}
{{end}}{{end}}

`

func getProfile(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin
	app.Action = func() {
		profile, err := util.API.GetUserProfile()
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(profile)
			return
		}

		sort.Sort(profile.Workspaces)

		t, err := template.New("profile").Parse(userProfile)
		if err != nil {
			util.Bail(err)
		}

		if err := t.Execute(os.Stdout, profile); err != nil {
			util.Bail(err)
		}

	}
}

func getSettings(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin
	app.Action = func() {
		settings, err := util.API.GetUserSettings()
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(settings)
		} else {
			if len(settings) > 0 {
				for k, v := range settings {
					fmt.Printf("%s: %v\n", k, v)
				}
			}
		}
	}
}

func getSetting(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin

	app.Action = func() {
		setting, err := util.API.GetUserSetting(SettingName)
		if err != nil {
			util.Bail(err)
		}

		var value interface{}

		v, ok := setting.(map[string]interface{})
		if ok {
			value = v[SettingName]
		} else {
			value = setting
		}

		if util.JSON {
			util.JSONOut(value)
		} else {
			fmt.Println(value)
		}
	}
}

func setSetting(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin

	var settingValueArg = app.StringArg("VALUE", "", "Setting value as JSON string")
	app.Spec = "VALUE"

	app.Action = func() {
		var userData interface{}
		err := json.Unmarshal([]byte(*settingValueArg), &userData)

		if err != nil {
			// If the value doesn't parse properly as JSON, we assume it's
			// literal. This catches the single-value case where we want
			// { "foo": "bar" } by just letting the user pass in a name of
			// "foo" and a value of "bar"

			// The perhaps surprising side effect is that crappy JSON will
			// enter the database as a string.
			userData = *settingValueArg
		}

		data := make(map[string]interface{})
		data[SettingName] = userData

		err = util.API.SetUserSetting(SettingName, data)
		if err != nil {
			util.Bail(err)
		}
	}
}

func deleteSetting(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin

	app.Action = func() {
		err := util.API.DeleteUserSetting(SettingName)
		if err != nil {
			util.Bail(err)
		}
	}
}

func listTokens(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin

	app.Action = func() {
		tokens, err := util.API.GetMyApiTokens()
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(tokens)
			return
		}

		sort.Sort(tokens)

		table := util.GetMarkdownTable()
		table.SetHeader([]string{"Name", "Created", "Last Used"})

		for _, t := range tokens {
			timeStr := ""
			if !t.LastUsed.IsZero() {
				timeStr = util.TimeStr(t.LastUsed)
			}

			table.Append([]string{
				t.Name,
				util.TimeStr(t.Created),
				timeStr,
			})
		}

		table.Render()
	}
}

func createToken(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin

	var nameArg = app.StringArg("NAME", "", "Name for the token")
	app.Spec = "NAME"

	app.Action = func() {
		token, err := util.API.CreateMyToken(*nameArg)
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(token)
			return
		}

		fmt.Println("***")
		fmt.Println("*** This is the *only* time the token string will be shown.")
		fmt.Println("***")
		fmt.Println("*** Please make sure to record it as the string cannot be retrieved later ")
		fmt.Println("***")
		fmt.Println()
		fmt.Printf("Name: %s\n", token.Name)
		fmt.Printf("Token: %s    <--- Write this down\n", token.Token)
		fmt.Println()
	}
}

func removeToken(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin

	var nameArg = app.StringArg("NAME", "", "Name for the token")
	app.Spec = "NAME"

	app.Action = func() {
		err := util.API.DeleteMyToken(*nameArg)
		if err != nil {
			util.Bail(err)
		}
	}
}

func getToken(cmd *cli.Cmd) {
	cmd.Before = util.BuildAPIAndVerifyLogin

	var nameArg = cmd.StringArg("NAME", "", "Name for the token")
	cmd.Spec = "NAME"

	cmd.Action = func() {
		token, err := util.API.GetMyToken(*nameArg)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(token)
			return
		}

		lastUsed := "[ Never Used ]"
		if !token.LastUsed.IsZero() {
			lastUsed = util.TimeStr(token.LastUsed)
		}

		fmt.Printf(`
Name: %s
Created: %s
Last Used: %s
`,
			token.Name,
			util.TimeStr(token.Created),
			lastUsed,
		)
	}
}
