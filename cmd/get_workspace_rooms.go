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
	"github.com/olekukonko/tablewriter"
	"os"
)

type getWorkspaceRoomsArgs struct {
	cli.Helper
	Id string `cli:"*id,uuid" usage:"ID of the workspace (required)"`
}

var GetWorkspaceRoomsCmd = &cli.Command{
	Name: "get_workspace_rooms",
	Desc: "Get a list of rooms for the given workspace ID",
	Argv: func() interface{} { return new(getWorkspaceRoomsArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getWorkspaceRoomsArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*getWorkspaceRoomsArgs)

		rooms, err := api.GetWorkspaceRooms(argv.Id)
		if err != nil {
			return err
		}

		if args.Global.JSON == true {

			j, err := json.Marshal(rooms)

			if err != nil {
				return err
			}

			fmt.Println(string(j))

		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "AZ", "Alias", "Vendor Name"})

			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")

			for _, r := range rooms {
				table.Append([]string{r.Id, r.Az, r.Alias, r.VendorName})
			}

			table.Render()
		}
		return nil
	},
}
