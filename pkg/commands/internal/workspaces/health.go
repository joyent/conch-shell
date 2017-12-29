// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package workspaces

import (
	"fmt"
	"github.com/joyent/conch-shell/pkg/util"
	conch "github.com/joyent/go-conch"
	"gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
	"regexp"
	"sort"
)

func getHealth(app *cli.Cmd) {
	var (
		full_output       = app.BoolOpt("full", false, "Instead of just presenting a datacenter summary, break results out by rack as well. Has no effect on --json")
		show_uuids        = app.BoolOpt("uuids", false, "Show UUIDs where appropriate")
		platform_name     = app.BoolOpt("platform-name", false, "Use the platform name (like 'Joyent-Foo-Platform-XXXX') instead of the common name (like 'Mantis Shrimp MkIII')")
		datacenter_choice = app.StringOpt("datacenter az", "", "Limit the output to a particular datacenter by UUID, partial UUID, or string name")
	)

	app.Action = func() {
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
		)

		full_report := make(map[string]datacenterReport)

		workspace_devices, err := util.API.GetWorkspaceDevices(
			WorkspaceUuid,
			true,
			"",
			"",
		)

		if err != nil {
			util.Bail(err)
		}

		for _, d := range workspace_devices {
			full_d, err := util.API.FillInDevice(d)
			if err != nil {
				util.Bail(err)
			}

			datacenter := default_datacenter
			datacenter_uuid := uuid.UUID{}
			if full_d.Location.Datacenter.Name != "" {
				datacenter = full_d.Location.Datacenter.Name
				datacenter_uuid = full_d.Location.Datacenter.Id

			}

			if *datacenter_choice != "" {
				re := regexp.MustCompile(fmt.Sprintf("^%s-", *datacenter_choice))
				if (datacenter_uuid.String() != *datacenter_choice) &&
					(datacenter != *datacenter_choice) &&
					!re.MatchString(*datacenter_choice) {
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
			if *platform_name {
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

		if util.JSON {
			util.JsonOut(full_report)
			return
		}

		az := make([]string, 0)
		for k := range full_report {
			az = append(az, k)
		}
		sort.Strings(az)

		for _, a := range az {
			if *show_uuids {
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

			if !*full_output {
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
				if *show_uuids {
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
	}
}
