// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package workspaces

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/pgtime"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
)

func getFailures(app *cli.Cmd) {
	var (
		fullOutput       = app.BoolOpt("full", false, "Instead of just presenting a datacenter summary, break results out by rack as well. Has no effect on --json")
		showUUIDs        = app.BoolOpt("uuids", false, "Show UUIDs where appropriate")
		datacenterChoice = app.StringOpt("datacenter az", "", "Limit the output to a particular datacenter by UUID, partial UUID, or string name")
	)

	app.Action = func() {

		type minimalReportDevice struct {
			AssetTag          string                              `json:"asset_tag"`
			Created           pgtime.PgTime                       `json:"created"`
			Graduated         pgtime.PgTime                       `json:"graduated"`
			HardwareProduct   uuid.UUID                           `json:"hardware_product"`
			Health            string                              `json:"health"`
			ID                string                              `json:"id"`
			LastSeen          pgtime.PgTime                       `json:"last_seen"`
			Location          conch.DeviceLocation                `json:"location"`
			Role              uuid.UUID                           `json:"role"`
			State             string                              `json:"state"`
			SystemUUID        uuid.UUID                           `json:"system_uuid"`
			Updated           pgtime.PgTime                       `json:"updated"`
			Validated         pgtime.PgTime                       `json:"validated"`
			FailedValidations map[string][]conch.ValidationReport `json:"failed_validations"`
		}

		type reportRack struct {
			Rack          conch.Rack            `json:"rack"`
			FailedDevices []minimalReportDevice `json:"failed_devices"`
		}

		type datacenterReport struct {
			Name    string                 `json:"datacenter"`
			ID      uuid.UUID              `json:"id"`
			Summary map[string]int         `json:"summary"`
			Racks   map[string]*reportRack `json:"racks"`
		}

		const (
			defaultComponentType = "UNKNOWN"
			defaultDatacenter    = "UNKNOWN"
			defaultRack          = "UNKNOWN"
		)

		fullReport := make(map[string]datacenterReport)

		workspaceDevices, err := util.API.GetWorkspaceDevices(
			WorkspaceUUID,
			false,
			"",
			"fail",
		)

		if err != nil {
			util.Bail(err)
		}

		for _, d := range workspaceDevices {
			fullDevice, err := util.API.FillInDevice(d)
			if err != nil {
				util.Bail(err)
			}

			reportDevice := minimalReportDevice{
				fullDevice.AssetTag,
				fullDevice.Created,
				fullDevice.Graduated,
				fullDevice.HardwareProduct,
				fullDevice.Health,
				fullDevice.ID,
				fullDevice.LastSeen,
				fullDevice.Location,
				fullDevice.Role,
				fullDevice.State,
				fullDevice.SystemUUID,
				fullDevice.Updated,
				fullDevice.Validated,
				make(map[string][]conch.ValidationReport),
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
					Summary: make(map[string]int),
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
					make([]minimalReportDevice, 0),
				}
			}

			fullReport[datacenter].Racks[rack].FailedDevices = append(
				fullReport[datacenter].Racks[rack].FailedDevices,
				reportDevice,
			)

			for _, report := range fullDevice.Validations {
				if report.Status == conch.ValidationReportStatusFail {
					vType := defaultComponentType
					if report.ComponentType != "" {
						vType = report.ComponentType
					}

					reportDevice.FailedValidations[vType] = append(
						reportDevice.FailedValidations[vType],
						report,
					)

					if _, ok := fullReport[datacenter].Summary[vType]; ok {
						fullReport[datacenter].Summary[vType]++
					} else {
						fullReport[datacenter].Summary[vType] = 1
					}
				}
			}

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

			types := make([]string, 0)
			for k := range fullReport[a].Summary {
				types = append(types, k)
			}
			sort.Strings(types)

			for _, t := range types {
				fmt.Printf("    %8s: %d\n", t, fullReport[a].Summary[t])
			}

			if !*fullOutput {
				fmt.Println()
				continue
			}

			fmt.Println()
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

				for _, device := range rack.FailedDevices {
					if *showUUIDs {
						fmt.Printf("      %s - %s:\n",
							device.ID,
							device.SystemUUID,
						)
					} else {
						fmt.Printf("      %s:\n", device.ID)
					}
					for _, t := range types {
						if _, ok := device.FailedValidations[t]; !ok {
							continue
						}
						fmt.Printf("        %s:\n", t)
						for _, validation := range device.FailedValidations[t] {
							fmt.Printf(
								"          %s : %s\n",
								validation.ComponentName,
								validation.Log,
							)
						}
						fmt.Println()
					}
				}
				fmt.Println()
			}

			fmt.Println()
		}
	}
}
