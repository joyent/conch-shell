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

type getWorkspacesArgs struct {
	cli.Helper
}

var GetWorkspacesCmd = &cli.Command{
	Name: "get_workspaces",
	Desc: "Get a list of workspaces and their IDs",
	Argv: func() interface{} { return new(getWorkspacesArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getWorkspacesArgs{}, ctx)

		if err != nil {
			return err
		}

		//argv := args.Local.(*getWorkspacesArgs)

		workspaces, err := api.GetWorkspaces()
		if err != nil {
			return err
		}

		if args.Global.JSON == true {

			j, err := json.Marshal(workspaces)

			if err != nil {
				return err
			}

			fmt.Println(string(j))

		} else {
			table := GetMarkdownTable()
			table.SetHeader([]string{"Role", "Id", "Name", "Description"})

			for _, w := range workspaces {
				table.Append([]string{w.Role, w.Id.String(), w.Name, w.Description})
			}

			table.Render()
		}
		return nil
	},
}
