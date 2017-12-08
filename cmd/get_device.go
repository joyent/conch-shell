// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	conch "github.com/joyent/go-conch"
	"github.com/mkideal/cli"
)

type getDeviceArgs struct {
	cli.Helper
	Id         string `cli:"*id,serial" usage:"The id of the device"`
	FullOutput bool   `cli:"full" usage:"When --json is used, provide full data about the device rather than the normal truncated data"`
}

var GetDeviceCmd = &cli.Command{
	Name: "get_device",
	Desc: "Get data about a specific device serial",
	Argv: func() interface{} { return new(getDeviceArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getDeviceArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*getDeviceArgs)
		device, err := api.GetDevice(argv.Id)
		if err != nil {
			return err
		}

		devices := make([]conch.ConchDevice, 0)
		devices = append(devices, device)

		return DisplayDevices(devices, args.Global.JSON, argv.FullOutput)
	},
}
