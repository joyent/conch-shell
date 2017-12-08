package cmd

import (
	//	pgtime "github.com/joyent/go-conch/pg_time"
	conch "github.com/joyent/go-conch"
	"github.com/mkideal/cli"
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

		return DisplayDevices(relay.Devices, args.Global.JSON, argv.FullOutput)
	},
}
