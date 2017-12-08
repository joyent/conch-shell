package cmd

import (
	"encoding/json"
	"fmt"
	//	pgtime "github.com/joyent/go-conch/pg_time"
	conch "github.com/joyent/go-conch"
	"github.com/mkideal/cli"
	"github.com/olekukonko/tablewriter"
	"os"
	//	"strconv"
)

type getRelayDevicesArgs struct {
	cli.Helper
	WorkspaceId string `cli:"*workspace-id,workspace-uuid,workspace,ws" usage:"ID of the workspace (required)"`
	Id          string `cli:"*id,relay-id,relay" usage:"ID of the relay (required)"`
	FullOutput  bool   `cli:"full" usage:"When --json is used, provide full data about the devices rather than the normal truncated data"`
}

var GetRelayDevicesCmd = &cli.Command{
	Name: "get_relay_devices",
	Desc: "Get a list of relays for the given workspace ID",
	Argv: func() interface{} { return new(getRelayDevicesArgs) },
	Fn: func(ctx *cli.Context) error {
		args, _, api, err := GetStarted(&getRelayDevicesArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*getRelayDevicesArgs)

		relays, err := api.GetWorkspaceRelays(argv.WorkspaceId, false)
		if err != nil {
			return err
		}

		var relay conch.ConchRelay
		found_relay := false
		for _, r := range relays {
			if r.Id == argv.Id {
				relay = r
				found_relay = true
			}
		}
		if found_relay == false {
			return conch.ConchDataNotFound
		}
		minimals := make([]MinimalDevice, 0)
		for _, d := range relay.Devices {
			minimals = append(minimals, MinimalDevice{
				d.Id,
				d.AssetTag,
				d.Created,
				d.LastSeen,
				d.Health,
				GenerateDeviceFlags(d),
			})
		}

		if args.Global.JSON {
			var j []byte
			if argv.FullOutput {
				j, err = json.Marshal(relay.Devices)
			} else {
				j, err = json.Marshal(minimals)
			}
			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		}

		TableizeMinimalDevices(minimals, tablewriter.NewWriter(os.Stdout)).Render()

		return nil
	},
}
