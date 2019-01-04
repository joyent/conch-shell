// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import (
	"encoding/json"
	"fmt"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
)

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
