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
	"strconv"
)

type getWorkspaceRacksArgs struct {
	cli.Helper
	Id         string `cli:"*id,uuid" usage:"ID of the workspace (required)"`
	FullOutput bool   `cli:"full" usage:"When --json is used and --ids-only is *not* used, provide full data about the devices rather than the normal truncated data"`
}

var GetWorkspaceRacksCmd = &cli.Command{
	Name: "get_workspace_racks",
	Desc: "Get a list of racks for the given workspace ID",
	Argv: func() interface{} { return new(getWorkspaceRacksArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getWorkspaceRacksArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*getWorkspaceRacksArgs)

		racks, err := api.GetWorkspaceRacks(argv.Id)
		if err != nil {
			return err
		}

		if args.Global.JSON {
			j, err := json.Marshal(racks)
			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		}

		table := GetMarkdownTable()
		table.SetHeader([]string{
			"ID",
			"Name",
			"Role",
			"Unit",
			"Size",
		})

		for _, r := range racks {
			table.Append([]string{
				fmt.Sprintf("%s", r.Id),
				r.Name,
				r.Role,
				strconv.Itoa(r.Unit),
				strconv.Itoa(r.Size),
			})
		}

		table.Render()

		return nil
	},
}
