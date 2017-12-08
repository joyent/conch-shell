// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mkideal/cli"
)

type getSettingArgs struct {
	cli.Helper
	Key string `cli:"key" usage:"The setting name"`
}

var GetSettingCmd = &cli.Command{
	Name: "get_setting",
	Desc: "Get the value of the provided setting for the current user",
	Argv: func() interface{} { return new(getSettingArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getSettingArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*getSettingArgs)
		setting, err := api.GetUserSetting(argv.Key)
		if err != nil {
			return err
		}

		if args.Global.JSON == true {
			j, err := json.Marshal(setting)

			if err != nil {
				return err
			}

			fmt.Println(string(j))
		} else {
			fmt.Println(setting)
		}
		return nil
	},
}
