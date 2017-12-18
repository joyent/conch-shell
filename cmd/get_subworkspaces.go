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
)

type getSubWorkspacesArgs struct {
	cli.Helper
	Id string `cli:"*id,uuid" usage:"ID of the workspace (required)"`
}

var GetSubWorkspacesCmd = &cli.Command{
	Name: "get_subworkspaces",
	Desc: "Get a list of subworkspaces for the given workspace ID",
	Argv: func() interface{} { return new(getSubWorkspacesArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(ctx, &getSubWorkspacesArgs{})

		if err != nil {
			return err
		}

		argv := args.Local.(*getSubWorkspacesArgs)

		workspace_id, err := uuid.FromString(argv.Id)
		if err != nil {
			return err
		}

		workspaces, err := api.GetSubWorkspaces(workspace_id)
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
