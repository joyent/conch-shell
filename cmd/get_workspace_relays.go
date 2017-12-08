// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"encoding/json"
	"fmt"
	pgtime "github.com/joyent/go-conch/pg_time"
	"github.com/mkideal/cli"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"time"
)

type getWorkspaceRelaysArgs struct {
	cli.Helper
	Id         string `cli:"*id,uuid" usage:"ID of the workspace (required)"`
	ActiveOnly bool   `cli:"active-only" usage:"Only retrieve active relays"`
	FullOutput bool   `cli:"full" usage:"When --json is used, provide full data about the devices rather than the normal truncated data"`
}

var GetWorkspaceRelaysCmd = &cli.Command{
	Name: "get_workspace_relays",
	Desc: "Get a list of relays for the given workspace ID",
	Argv: func() interface{} { return new(getWorkspaceRelaysArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getWorkspaceRelaysArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*getWorkspaceRelaysArgs)

		relays, err := api.GetWorkspaceRelays(argv.Id, argv.ActiveOnly)
		if err != nil {
			return err
		}

		if args.Global.JSON && argv.FullOutput {
			j, err := json.Marshal(relays)

			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		}

		type resultRow struct {
			Id         string             `json:"id"`
			Alias      string             `json:"asset_tag"`
			Created    pgtime.ConchPgTime `json:"created, int"`
			IpAddr     string             `json:"ipaddr"`
			SshPort    int                `json:"ssh_port"`
			Updated    pgtime.ConchPgTime `json:"updated"`
			Version    string             `json:"version"`
			NumDevices int                `json:"num_devices"`
		}

		results := make([]resultRow, 0)

		for _, r := range relays {
			num_devices := len(r.Devices)
			results = append(results, resultRow{
				r.Id,
				r.Alias,
				r.Created,
				r.IpAddr,
				r.SshPort,
				r.Updated,
				r.Version,
				num_devices,
			})
		}

		if args.Global.JSON == true {
			j, err := json.Marshal(results)

			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"ID",
			"Alias",
			"Created",
			"IP Addr",
			"SSH Port",
			"Updated",
			"Version",
			"Number of Devices",
		})

		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		for _, r := range results {
			updated := ""
			if !r.Updated.IsZero() {
				updated = r.Updated.Format(time.UnixDate)
			}

			table.Append([]string{
				r.Id,
				r.Alias,
				r.Created.Format(time.UnixDate),
				r.IpAddr,
				strconv.Itoa(r.SshPort),
				updated,
				r.Version,
				strconv.Itoa(r.NumDevices),
			})
		}

		table.Render()
		return nil
	},
}
