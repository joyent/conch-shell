package cmd

import (
	"encoding/json"
	"fmt"
	conch "github.com/joyent/go-conch"
	"github.com/mkideal/cli"
	"github.com/olekukonko/tablewriter"
	"os"
	"time"
)

type getWorkspaceDevicesArgs struct {
	cli.Helper
	Id         string `cli:"*id,uuid" usage:"ID of the workspace (required)"`
	IdsOnly    bool   `cli:"ids-only" usage:"Only retrieve device IDs"`
	Graduated  string `cli:"graduated" usage:"Filter by the 'graduated' field"`
	Health     string `cli:"health" usage:"Filter by the 'health' field using the string provided"`
	FullOutput bool   `cli:"full" usage:"When --json is used and --ids-only is *not* used, provide full data about the devices rather than the normal truncated data"`
}

var GetWorkspaceDevicesCmd = &cli.Command{
	Name: "get_workspace_devices",
	Desc: "Get a list of devices for the given workspace ID",
	Argv: func() interface{} { return new(getWorkspaceDevicesArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getWorkspaceDevicesArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*getWorkspaceDevicesArgs)

		devices, err := api.GetWorkspaceDevices(argv.Id, argv.IdsOnly, argv.Graduated, argv.Health)
		if err != nil {
			return err
		}

		if args.Global.JSON && argv.FullOutput {
			j, err := json.Marshal(devices)

			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		}

		type resultRow struct {
			Id       string `json:"id"`
			AssetTag string `json:"asset_tag"`
			Created  string `json:"created, int"`
			LastSeen string `json:"last_seen"`
			Health   string `json:"health"`
			Flags    string `json:"flags"`
		}

		type resultItem struct {
			Device conch.ConchDevice
			Flags  string
		}

		results := make([]resultItem, 0)

		for _, d := range devices {

			flags := ""

			if !d.Deactivated.IsZero() {
				flags += "X"
			}

			if !d.Validated.IsZero() {
				flags += "v"
			}

			if !d.Graduated.IsZero() {
				flags += "g"
			}
			results = append(results, resultItem{d, flags})
		}

		if args.Global.JSON == true {
			if argv.IdsOnly {
				out := make([]string, 0)
				for _, v := range devices {
					out = append(out, v.Id)
				}
				j, err := json.Marshal(out)

				if err != nil {
					return err
				}
				fmt.Println(string(j))
				return nil
			}

			out := make([]resultRow, 0)
			for _, v := range results {
				d := v.Device
				flags := v.Flags

				last_seen := ""
				if !d.LastSeen.IsZero() {
					last_seen = fmt.Sprintf("%d", d.LastSeen.UTC().Unix())
				}

				out = append(out, resultRow{
					d.Id,
					d.AssetTag,
					fmt.Sprintf("%d", d.Created.UTC().Unix()),
					last_seen,
					d.Health,
					flags,
				})

			}
			j, err := json.Marshal(out)

			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		}

		if argv.IdsOnly {
			for _, v := range devices {
				fmt.Println(v.Id)
			}
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"ID",
			"Asset Tag",
			"Created",
			"Last Seen",
			"Health",
			"Flags",
		})

		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		for _, r := range results {
			d := r.Device
			flags := r.Flags

			last_seen := ""
			if !d.LastSeen.IsZero() {
				last_seen = d.LastSeen.Format(time.UnixDate)
			}

			table.Append([]string{
				d.Id,
				d.AssetTag,
				d.Created.Format(time.UnixDate),
				last_seen,
				d.Health,
				flags,
			})
		}

		table.Render()
		return nil
	},
}
