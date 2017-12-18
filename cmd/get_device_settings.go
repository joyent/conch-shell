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
	"sort"
)

type getDeviceSettingsArgs struct {
	cli.Helper
	Id       string `cli:"*id,serial" usage:"The id of the device"`
	KeysOnly bool   `cli:"keys-only" usage:"Only display the setting keys/names"`
}

var GetDeviceSettingsCmd = &cli.Command{
	Name: "get_device_settings",
	Desc: "Get settings for a specific device serial",
	Argv: func() interface{} { return new(getDeviceSettingsArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(ctx, &getDeviceSettingsArgs{})

		if err != nil {
			return err
		}

		argv := args.Local.(*getDeviceSettingsArgs)
		settings, err := api.GetDeviceSettings(argv.Id)
		if err != nil {
			return err
		}

		keys := make([]string, 0, len(settings))
		for k := range settings {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		if argv.KeysOnly {
			if args.Global.JSON {
				j, err := json.Marshal(keys)
				if err != nil {
					return err
				}
				fmt.Println(string(j))
				return nil
			}
			for _, k := range keys {
				fmt.Println(k)
			}
			return nil
		}

		if args.Global.JSON {
			j, err := json.Marshal(settings)
			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		}

		for _, k := range keys {
			fmt.Printf("%s : %v\n", k, settings[k])
		}
		return nil
	},
}
