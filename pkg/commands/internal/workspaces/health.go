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
		fullOutput       = app.BoolOpt("full", false, "Instead of just presenting a datacenter summary, break results out by rack as well. Has no effect on --json")
		showUUIDs        = app.BoolOpt("uuids", false, "Show UUIDs where appropriate")
		platformName     = app.BoolOpt("platform-name", false, "Use the platform name (like 'Joyent-Foo-Platform-XXXX') instead of the common name (like 'Mantis Shrimp MkIII')")
		datacenterChoice = app.StringOpt("datacenter az", "", "Limit the output to a particular datacenter by UUID, partial UUID, or string name")
	)

	app.Action = func() {
		type reportRack struct {
			Rack    conch.Rack                `json:"rack"`
			Summary map[string]map[string]int `json:"summary"`
		}

		type datacenterReport struct {
			Name    string                    `json:"datacenter"`
			ID      uuid.UUID                 `json:"id"`
			Summary map[string]map[string]int `json:"summary"`
			Racks   map[string]*reportRack    `json:"racks"`
		}

		const (
			defaultHardwareType = "UNKNOWN"
			defaultDatacenter   = "UNKNOWN"
			defaultRack         = "UNKNOWN"
		)

		fullReport := make(map[string]datacenterReport)

		workspaceDevices, err := util.API.GetWorkspaceDevices(
			WorkspaceUUID,
			true,
			"",
			"",
		)

		if err != nil {
			util.Bail(err)
		}

		for _, d := range workspaceDevices {
			fullDevice, err := util.API.FillInDevice(d)
			if err != nil {
				util.Bail(err)
			}

			datacenter := defaultDatacenter
			datacenterUUID := uuid.UUID{}
			if fullDevice.Location.Datacenter.Name != "" {
				datacenter = fullDevice.Location.Datacenter.Name
				datacenterUUID = fullDevice.Location.Datacenter.ID

			}

			if *datacenterChoice != "" {
				re := regexp.MustCompile(fmt.Sprintf("^%s-", *datacenterChoice))
				if (datacenterUUID.String() != *datacenterChoice) &&
					(datacenter != *datacenterChoice) &&
					!re.MatchString(*datacenterChoice) {
					continue
				}
			}

			if _, ok := fullReport[datacenter]; !ok {
				fullReport[datacenter] = datacenterReport{
					Name:    datacenter,
					ID:      datacenterUUID,
					Summary: make(map[string]map[string]int),
					Racks:   make(map[string]*reportRack),
				}
			}

			rack := defaultRack
			if fullDevice.Location.Rack.Name != "" {
				rack = fullDevice.Location.Rack.Name
			}
			if _, ok := fullReport[datacenter].Racks[rack]; !ok {
				fullReport[datacenter].Racks[rack] = &reportRack{
					fullDevice.Location.Rack,
					make(map[string]map[string]int),
				}
			}

			hwtype := defaultHardwareType
			if *platformName {
				if fullDevice.Location.TargetHardwareProduct.Name != "" {
					hwtype = fullDevice.Location.TargetHardwareProduct.Name
				}
			} else {
				if fullDevice.Location.TargetHardwareProduct.Alias != "" {
					hwtype = fullDevice.Location.TargetHardwareProduct.Alias
				}
			}
			if _, ok := fullReport[datacenter].Summary[hwtype]; !ok {
				fullReport[datacenter].Summary[hwtype] = make(map[string]int)
			}
			if _, ok := fullReport[datacenter].Racks[rack].Summary[hwtype]; !ok {
				fullReport[datacenter].Racks[rack].Summary[hwtype] = make(map[string]int)
			}

			fullReport[datacenter].Summary[hwtype][fullDevice.Health]++
			fullReport[datacenter].Racks[rack].Summary[hwtype][fullDevice.Health]++

		}

		if util.JSON {
			util.JSONOut(fullReport)
			return
		}

		az := make([]string, 0)
		for k := range fullReport {
			az = append(az, k)
		}
		sort.Strings(az)

		for _, a := range az {
			if *showUUIDs {
				fmt.Printf("%s - %s\n", a, fullReport[a].ID)
			} else {
				fmt.Println(a)
			}
			fmt.Println("  Summary:")

			hwtypes := make([]string, 0)
			for k := range fullReport[a].Summary {
				hwtypes = append(hwtypes, k)
			}
			sort.Strings(hwtypes)

			for _, h := range hwtypes {
				fmt.Printf("    %s:\n", h)

				types := make([]string, 0)
				for k := range fullReport[a].Summary[h] {
					types = append(types, k)
				}
				sort.Strings(types)
				for _, t := range types {
					fmt.Printf("      %8s: %d\n", t, fullReport[a].Summary[h][t])
				}
				fmt.Println()
			}

			if !*fullOutput {
				fmt.Println()
				continue
			}

			fmt.Println("  Racks:")

			rackNames := make([]string, 0)
			for k := range fullReport[a].Racks {
				rackNames = append(rackNames, k)
			}
			sort.Strings(rackNames)

			for _, rackName := range rackNames {
				rack := fullReport[a].Racks[rackName]
				if *showUUIDs {
					fmt.Printf("    %s - %s:\n", rackName, rack.Rack.ID)
				} else {
					fmt.Printf("    %s:\n", rackName)
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
