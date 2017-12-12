// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	config "github.com/joyent/conch-shell/config"
	conch "github.com/joyent/go-conch"
	pgtime "github.com/joyent/go-conch/pg_time"
	"github.com/mkideal/cli"
	"github.com/olekukonko/tablewriter"
	"os"
	"time"
)

var (
	ConchConfigurationError = errors.New("No configuration data found or parse error")
	ConchNoApiSessionData   = errors.New("No API session data found")
)

type CliArgs struct {
	Local  interface{}
	Global *GlobalArgs
}

type MinimalDevice struct {
	Id       string             `json:"id"`
	AssetTag string             `json:"asset_tag"`
	Created  pgtime.ConchPgTime `json:"created,int"`
	LastSeen pgtime.ConchPgTime `json:"last_seen,int"`
	Health   string             `json:"health"`
	Flags    string             `json:"flags"`
}

// GetStarted handles the initial logic of parsing arguments, loading the JSON
// config file and verifying that the login credentials are still valid.
// Pretty much every command should start by using this function.
//
// Pro-tip: To cast args.Local to your command's arguments struct, use:
//   argv := args.Local.(*myLocalArgs)
func GetStarted(argv interface{}, ctx *cli.Context) (args *CliArgs, cfg *config.ConchConfig, api *conch.Conch, err error) {
	globals := &GlobalArgs{}
	if err := ctx.GetArgvList(argv, globals); err != nil {
		return nil, nil, nil, err
	}

	args = &CliArgs{
		Local:  argv,
		Global: globals,
	}

	cfg, err = config.NewFromJsonFile(globals.ConfigPath)
	if err != nil {
		return nil, nil, nil, ConchConfigurationError
	}

	if cfg.Session == "" {
		return nil, nil, nil, ConchNoApiSessionData
	}

	api = &conch.Conch{
		BaseUrl: cfg.Api,
		User:    cfg.User,
		Session: cfg.Session,
	}

	if err = api.VerifyLogin(); err != nil {
		return nil, nil, nil, err
	}

	return args, cfg, api, nil
}

// GenerateDeviceFlags() is an abstraction to make sure that the 'flags' field
// for ConchDevices remains uniform
func GenerateDeviceFlags(d conch.ConchDevice) (flags string) {
	flags = ""

	if !d.Deactivated.IsZero() {
		flags += "X"
	}

	if !d.Validated.IsZero() {
		flags += "v"
	}

	if !d.Graduated.IsZero() {
		flags += "g"
	}
	return flags
}

// TableizeMinimalDevices() is an abstraction to make sure that tables of
// ConchDevices-turned-MinimalDevices are uniform
func TableizeMinimalDevices(devices []MinimalDevice, table *tablewriter.Table) *tablewriter.Table {
	table.SetHeader([]string{
		"ID",
		"Asset Tag",
		"Created",
		"Last Seen",
		"Health",
		"Flags",
	})

	for _, d := range devices {
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
			d.Flags,
		})
	}

	return table
}

// DisplayDevices() is an abstraction to make sure that the output of
// ConchDevices is uniform, be it tables, json, or full json
func DisplayDevices(devices []conch.ConchDevice, json_output bool, full_output bool) (err error) {
	minimals := make([]MinimalDevice, 0)
	for _, d := range devices {
		minimals = append(minimals, MinimalDevice{
			d.Id,
			d.AssetTag,
			d.Created,
			d.LastSeen,
			d.Health,
			GenerateDeviceFlags(d),
		})
	}

	if json_output {
		var j []byte
		if full_output {
			j, err = json.Marshal(devices)
		} else {
			j, err = json.Marshal(minimals)
		}
		if err != nil {
			return err
		}
		fmt.Println(string(j))
		return nil
	}

	TableizeMinimalDevices(minimals, GetMarkdownTable()).Render()

	return nil
}

func GetMarkdownTable() (table *tablewriter.Table) {
	table = tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	return table
}
