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

type getDeviceSettingArgs struct {
	cli.Helper
	Id  string `cli:"*id,serial" usage:"The id of the device"`
	Key string `cli:"*key,name" usage:"The setting name/key"`
}

var GetDeviceSettingCmd = &cli.Command{
	Name: "get_device_setting",
	Desc: "Get a single setting for a specific device serial",
	Argv: func() interface{} { return new(getDeviceSettingArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getDeviceSettingArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*getDeviceSettingArgs)
		setting, err := api.GetDeviceSetting(argv.Id, argv.Key)
		if err != nil {
			return err
		}

		if args.Global.JSON {
			j, err := json.Marshal(map[string]string{argv.Key: setting})
			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		}

		fmt.Println(setting)
		return nil
	},
}
