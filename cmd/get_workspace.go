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

type getWorkspaceArgs struct {
	cli.Helper
	Id string `cli:"*id,uuid" usage:"ID of the workspace (required)"`
}

var GetWorkspaceCmd = &cli.Command{
	Name: "get_workspace",
	Desc: "Get information about a single workspace, by UUID",
	Argv: func() interface{} { return new(getWorkspaceArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getWorkspaceArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*getWorkspaceArgs)
		workspace, err := api.GetWorkspace(argv.Id)
		if err != nil {
			return err
		}

		if args.Global.JSON == true {
			j, err := json.Marshal(workspace)

			if err != nil {
				return err
			}

			fmt.Println(string(j))
		} else {
			fmt.Printf(
				"Role: %s\nID: %s\nName: %s\nDescription: %s\n",
				workspace.Role,
				workspace.Id,
				workspace.Name,
				workspace.Description,
			)
		}
		return nil
	},
}
