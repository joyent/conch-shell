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
	uuid "gopkg.in/satori/go.uuid.v1"
	"sort"
	"strconv"
)

type getWorkspaceRackArgs struct {
	cli.Helper
	Id         string `cli:"*workspace" usage:"ID of the workspace (required)"`
	RackId     string `cli:"*id,rack" usage:"ID of the rack (required)"`
	SlotDetail bool   `cli:"slots" usage:"Show the devices in their slots rather than the rack summary"`
	FullOutput bool   `cli:"full" usage:"When --json is used and --ids-only is *not* used, provide full data about the devices rather than the normal truncated data"`
}

var GetWorkspaceRackCmd = &cli.Command{
	Name: "get_workspace_rack",
	Desc: "Get details about a single rack for the given workspace",
	Argv: func() interface{} { return new(getWorkspaceRackArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getWorkspaceRackArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*getWorkspaceRackArgs)

		workspace_id, err := uuid.FromString(argv.Id)
		if err != nil {
			return err
		}

		rack_id, err := uuid.FromString(argv.RackId)
		if err != nil {
			return err
		}

		rack, err := api.GetWorkspaceRack(workspace_id, rack_id)
		if err != nil {
			return err
		}

		if args.Global.JSON {
			j, err := json.MarshalIndent(rack, "", "	")
			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		}

		fmt.Printf(`
Workspace: %s
Rack ID:   %s
Name: %s
Role: %s
Datacenter: %s
`,
			workspace_id.String(),
			rack_id.String(),
			rack.Name,
			rack.Role,
			rack.Datacenter,
		)

		if argv.SlotDetail {
			fmt.Println()

			slot_nums := make([]int, 0, len(rack.Slots))
			for k := range rack.Slots {
				slot_nums = append(slot_nums, k)
			}
			sort.Sort(sort.Reverse(sort.IntSlice(slot_nums)))

			table := GetMarkdownTable()
			table.SetHeader([]string{
				"RU",
				"Occupied",
				"Name",
				"Alias",
				"Vendor",
				"Occupied By",
				"Health",
			})

			for _, ru := range slot_nums {
				slot := rack.Slots[ru]
				occupied := "X"

				occupant_id := ""
				occupant_health := ""

				if slot.Occupant.Id != "" {
					occupied = "+"
					occupant_id = slot.Occupant.Id
					occupant_health = slot.Occupant.Health
				}

				table.Append([]string{
					strconv.Itoa(ru),
					occupied,
					slot.Name,
					slot.Alias,
					slot.Vendor,
					occupant_id,
					occupant_health,
				})

			}
			table.Render()
		}

		return nil
	},
}
