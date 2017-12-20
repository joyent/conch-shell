// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package user

import (
	"fmt"
	"github.com/joyent/conch-shell/util"
	conch "github.com/joyent/go-conch"
	"gopkg.in/jawher/mow.cli.v1"
	"strings"
)

func getSettings(app *cli.Cmd) {
	app.Before = util.BuildApiAndVerifyLogin
	app.Action = func() {
		settings, err := util.API.GetUserSettings()
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JsonOut(settings)
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
	app.Before = util.BuildApiAndVerifyLogin

	var setting_id_str = app.StringArg("ID", "", "Setting name")
	app.Spec = "ID"

	app.Action = func() {
		setting, err := util.API.GetUserSetting(*setting_id_str)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JsonOut(setting)
		} else {
			fmt.Println(setting)
		}
	}
}

// Login() is exported so that it can be used as a first level command as well
// as a nested one

// BUG(sungo): prompt for data if args are empty
func Login(app *cli.Cmd) {
	var (
		api_str      = app.StringOpt("api", "https://conch.joyent.us", "The url of the API server")
		user_str     = app.StringOpt("user u", "", "The user name to log in with")
		password_str = app.StringOpt("password p", "", "The user's password")
	)

	app.Action = func() {
		api := &conch.Conch{
			BaseUrl: strings.TrimRight(*api_str, "/"),
			User:    *user_str,
		}

		if err := api.Login(*password_str); err != nil {
			util.Bail(err)
		}

		if api.Session == "" {
			util.Bail(conch.ConchNoSessionData)
		}

		util.Config.Api = api.BaseUrl
		util.Config.User = api.User
		util.Config.Session = api.Session

		if err := util.Config.SerializeToFile(util.Config.Path); err == nil {

			fmt.Printf("Success. Config written to %s\n", util.Config.Path)
		}

	}
}
