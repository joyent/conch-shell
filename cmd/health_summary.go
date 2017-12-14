// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"encoding/json"
	"fmt"
	conch "github.com/joyent/go-conch"
	"github.com/mkideal/cli"
	uuid "gopkg.in/satori/go.uuid.v1"
	"sort"
)

type healthSummaryArgs struct {
	cli.Helper
	Id           string `cli:"*workspace_id,workspace_uuid,workspace" usage:"ID of the workspace (required)"`
	Breakout     bool   `cli:"breakout" usage:"Instead of just presenting a datacenter summary, breakout results by rack as well (Ignored in the presence of --json)"`
	Uuids        bool   `cli:"uuids" usage:"Show UUIDs where appropriate"`
	PlatformName bool   `cli:"platform-name" usage:"Use the platform name (like 'Joyent-Foo-Platform-XXXX') instead of common name (like 'Mantis Shrimp MkIII')"`
	Datacenter   string `cli:"datacenter" usage:"Limit the output to a particular datacenter UUID"`
}

var HealthSummaryCmd = &cli.Command{
	Name: "health_summary",
	Desc: "Report that shows info about health summary by hardware type",
	Argv: func() interface{} { return new(healthSummaryArgs) },
	Fn: func(ctx *cli.Context) error {

		type reportRack struct {
			Rack    conch.ConchRack           `json:"rack"`
			Summary map[string]map[string]int `json:"summary"`
		}

		type datacenterReport struct {
			Name    string                    `json:"datacenter"`
			Id      uuid.UUID                 `json:"id"`
			Summary map[string]map[string]int `json:"summary"`
			Racks   map[string]*reportRack    `json:"racks"`
		}

		const (
			default_hardware_type = "UNKNOWN"
			default_datacenter    = "UNKNOWN"
			default_rack          = "UNKNOWN"
			default_rack_unit     = 0
		)

		full_report := make(map[string]datacenterReport)

		args, _, api, err := GetStarted(&healthSummaryArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*healthSummaryArgs)

		workspace_id, err := uuid.FromString(argv.Id)
		if err != nil {
			return err
		}

		workspace_devices, err := api.GetWorkspaceDevices(
			workspace_id,
			true,
			"",
			"",
		)

		if err != nil {
			return err
		}

		for _, d := range workspace_devices {
			full_d, err := api.FillInDevice(d)
			if err != nil {
				return err
			}

			datacenter := default_datacenter
			datacenter_uuid := uuid.UUID{}
			if full_d.Location.Datacenter.Name != "" {
				datacenter = full_d.Location.Datacenter.Name
				datacenter_uuid = full_d.Location.Datacenter.Id

			}

			if argv.Datacenter != "" {
				if datacenter_uuid.String() != argv.Datacenter {
					continue
				}
			}

			if _, ok := full_report[datacenter]; !ok {
				full_report[datacenter] = datacenterReport{
					Name:    datacenter,
					Id:      datacenter_uuid,
					Summary: make(map[string]map[string]int),
					Racks:   make(map[string]*reportRack),
				}
			}

			rack := default_rack
			if full_d.Location.Rack.Name != "" {
				rack = full_d.Location.Rack.Name
			}
			if _, ok := full_report[datacenter].Racks[rack]; !ok {
				full_report[datacenter].Racks[rack] = &reportRack{
					full_d.Location.Rack,
					make(map[string]map[string]int),
				}
			}

			hwtype := default_hardware_type
			if argv.PlatformName {
				if full_d.Location.TargetHardwareProduct.Name != "" {
					hwtype = full_d.Location.TargetHardwareProduct.Name
				}
			} else {
				if full_d.Location.TargetHardwareProduct.Alias != "" {
					hwtype = full_d.Location.TargetHardwareProduct.Alias
				}
			}
			if _, ok := full_report[datacenter].Summary[hwtype]; !ok {
				full_report[datacenter].Summary[hwtype] = make(map[string]int)
			}
			if _, ok := full_report[datacenter].Racks[rack].Summary[hwtype]; !ok {
				full_report[datacenter].Racks[rack].Summary[hwtype] = make(map[string]int)
			}

			full_report[datacenter].Summary[hwtype][full_d.Health]++
			full_report[datacenter].Racks[rack].Summary[hwtype][full_d.Health]++

		}

		if args.Global.JSON {
			j, err := json.Marshal(full_report)
			if err != nil {
				return err
			}
			fmt.Println(string(j))
			return nil
		}

		az := make([]string, 0)
		for k := range full_report {
			az = append(az, k)
		}
		sort.Strings(az)

		for _, a := range az {
			if argv.Uuids {
				fmt.Printf("%s - %s\n", a, full_report[a].Id)
			} else {
				fmt.Println(a)
			}
			fmt.Println("  Summary:")

			hwtypes := make([]string, 0)
			for k := range full_report[a].Summary {
				hwtypes = append(hwtypes, k)
			}
			sort.Strings(hwtypes)

			for _, h := range hwtypes {
				fmt.Printf("    %s:\n", h)

				types := make([]string, 0)
				for k := range full_report[a].Summary[h] {
					types = append(types, k)
				}
				sort.Strings(types)
				for _, t := range types {
					fmt.Printf("      %8s: %d\n", t, full_report[a].Summary[h][t])
				}
				fmt.Println()
			}

			if !argv.Breakout {
				fmt.Println()
				continue
			}

			fmt.Println("  Racks:")

			rack_names := make([]string, 0)
			for k := range full_report[a].Racks {
				rack_names = append(rack_names, k)
			}
			sort.Strings(rack_names)

			for _, rack_name := range rack_names {
				rack := full_report[a].Racks[rack_name]
				if argv.Uuids {
					fmt.Printf("    %s - %s:\n", rack_name, rack.Rack.Id)
				} else {
					fmt.Printf("    %s:\n", rack_name)
				}

				hwtypes := make([]string, 0)
				for k := range rack.Summary {
					hwtypes = append(hwtypes, k)
				}
				sort.Strings(hwtypes)

				for _, h := range hwtypes {
					fmt.Printf("      %s:\n", h)

					types := make([]string, 0)
					for k := range rack.Summary[h] {
						types = append(types, k)
					}
					sort.Strings(types)
					for _, t := range types {
						fmt.Printf("        %8s: %d\n", t, rack.Summary[h][t])
					}
					fmt.Println()
				}
				fmt.Println()
			}
		}

		return nil
	},
}
