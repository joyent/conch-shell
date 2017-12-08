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

type getDeviceLocationArgs struct {
	cli.Helper
	Id string `cli:"*id,serial" usage:"The id of the device"`
}

var GetDeviceLocationCmd = &cli.Command{
	Name: "get_device_location",
	Desc: "Get the location for a specific device serial",
	Argv: func() interface{} { return new(getDeviceLocationArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getDeviceLocationArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*getDeviceLocationArgs)
		location, err := api.GetDeviceLocation(argv.Id)
		if err != nil {
			return err
		}

		if args.Global.JSON {
			j, err := json.Marshal(location)
			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		}

		fmt.Printf(`
Location for device %s:
  Datacenter:
    Id:   %s
    Name: %s
  Rack:
    Id:   %s
    Name: %s
    Role: %s
    Unit: %d
`,
			argv.Id,
			location.Datacenter.Id,
			location.Datacenter.Name,
			location.Rack.Id,
			location.Rack.Name,
			location.Rack.Role,
			location.Rack.Unit,
		)
		return nil
	},
}
