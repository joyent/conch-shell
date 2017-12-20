// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package devices

import (
	"fmt"
	"github.com/joyent/conch-shell/util"
	"github.com/joyent/go-conch"
	"gopkg.in/jawher/mow.cli.v1"
	"sort"
)

func getOne(app *cli.Cmd) {
	var full_output = app.BoolOpt("full", false, "When global --json is used, provide full data about the devices rather than normal truncated data")
	app.Action = func() {
		device, err := util.API.GetDevice(DeviceSerial)
		if err != nil {
			util.Bail(err)
		}
		devices := make([]conch.ConchDevice, 0)
		devices = append(devices, device)

		util.DisplayDevices(devices, *full_output)
	}
}

func getLocation(app *cli.Cmd) {
	app.Action = func() {
		location, err := util.API.GetDeviceLocation(DeviceSerial)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JsonOut(location)
			return
		}

		fmt.Printf(`Location for device %s:
  Datacenter:
    Id:   %s
    Name: %s
  Rack:
    Id:   %s
    Name: %s
    Role: %s
    Unit: %d
`,
			DeviceSerial,
			location.Datacenter.Id,
			location.Datacenter.Name,
			location.Rack.Id,
			location.Rack.Name,
			location.Rack.Role,
			location.Rack.Unit,
		)
		return
	}
}

func getSettings(app *cli.Cmd) {
	var keys_only = app.BoolOpt("keys-only", false, "Only display the setting keys/names")
	app.Action = func() {
		settings, err := util.API.GetDeviceSettings(DeviceSerial)
		if err != nil {
			util.Bail(err)
		}

		keys := make([]string, 0, len(settings))
		for k := range settings {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		if *keys_only {
			if util.JSON {
				util.JsonOut(keys)
				return
			}

			for _, k := range keys {
				fmt.Println(k)
			}
			return
		}

		if util.JSON {
			util.JsonOut(settings)
			return
		}

		for _, k := range keys {
			fmt.Printf("%s : %v\n", k, settings[k])
		}
		return
	}
}

func getSetting(app *cli.Cmd) {
	var setting_str = app.StringArg("SETTING", "", "The name of the setting to retrieve")
	app.Spec = "SETTING"

	app.Action = func() {

		setting, err := util.API.GetDeviceSetting(DeviceSerial, *setting_str)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JsonOut(map[string]string{*setting_str: setting})
		} else {
			fmt.Println(setting)
		}
	}
}
