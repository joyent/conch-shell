// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package workspaces

import (
	"fmt"
	"github.com/joyent/conch-shell/util"
	conch "github.com/joyent/go-conch"
	pgtime "github.com/joyent/go-conch/pg_time"
	"gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
	"regexp"
	"sort"
)

func getFailures(app *cli.Cmd) {
	var (
		full_output       = app.BoolOpt("full", false, "Instead of just presenting a datacenter summary, break results out by rack as well. Has no effect on --json")
		show_uuids        = app.BoolOpt("uuids", false, "Show UUIDs where appropriate")
		datacenter_choice = app.StringOpt("datacenter az", "", "Limit the output to a particular datacenter by UUID, partial UUID, or string name")
	)

	app.Action = func() {

		type minimalReportDevice struct {
			AssetTag          string                                   `json:"asset_tag"`
			Created           pgtime.ConchPgTime                       `json:"created, int"`
			Graduated         pgtime.ConchPgTime                       `json:"graduated"`
			HardwareProduct   uuid.UUID                                `json:"hardware_product"`
			Health            string                                   `json:"health"`
			Id                string                                   `json:"id"`
			LastSeen          pgtime.ConchPgTime                       `json:"last_seen, int"`
			Location          conch.ConchDeviceLocation                `json:"location"`
			Role              string                                   `json:"role"`
			State             string                                   `json:"state"`
			SystemUuid        uuid.UUID                                `json:"system_uuid"`
			Updated           pgtime.ConchPgTime                       `json:"updated, int"`
			Validated         pgtime.ConchPgTime                       `json:"validated, int"`
			FailedValidations map[string][]conch.ConchValidationReport `json:"failed_validations"`
		}

		type reportRack struct {
			Rack          conch.ConchRack       `json:"rack"`
			FailedDevices []minimalReportDevice `json:"failed_devices"`
		}

		type datacenterReport struct {
			Name    string                 `json:"datacenter"`
			Id      uuid.UUID              `json:"id"`
			Summary map[string]int         `json:"summary"`
			Racks   map[string]*reportRack `json:"racks"`
		}

		const (
			default_component_type = "UNKNOWN"
			default_datacenter     = "UNKNOWN"
			default_rack           = "UNKNOWN"
			default_rack_unit      = 0
		)

		full_report := make(map[string]datacenterReport)

		workspace_devices, err := util.API.GetWorkspaceDevices(
			WorkspaceUuid,
			false,
			"",
			"fail",
		)

		if err != nil {
			util.Bail(err)
		}

		for _, d := range workspace_devices {
			full_d, err := util.API.FillInDevice(d)
			if err != nil {
				util.Bail(err)
			}

			report_device := minimalReportDevice{
				full_d.AssetTag,
				full_d.Created,
				full_d.Graduated,
				full_d.HardwareProduct,
				full_d.Health,
				full_d.Id,
				full_d.LastSeen,
				full_d.Location,
				full_d.Role,
				full_d.State,
				full_d.SystemUuid,
				full_d.Updated,
				full_d.Validated,
				make(map[string][]conch.ConchValidationReport),
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
					Summary: make(map[string]int),
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
					make([]minimalReportDevice, 0),
				}
			}

			full_report[datacenter].Racks[rack].FailedDevices = append(
				full_report[datacenter].Racks[rack].FailedDevices,
				report_device,
			)

			for _, report := range full_d.Validations {
				if report.Status == conch.ConchValidationReportStatusFail {
					v_type := default_component_type
					if report.ComponentType != "" {
						v_type = report.ComponentType
					}

					report_device.FailedValidations[v_type] = append(
						report_device.FailedValidations[v_type],
						report,
					)

					if _, ok := full_report[datacenter].Summary[v_type]; ok {
						full_report[datacenter].Summary[v_type]++
					} else {
						full_report[datacenter].Summary[v_type] = 1
					}
				}
			}

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

			types := make([]string, 0)
			for k := range full_report[a].Summary {
				types = append(types, k)
			}
			sort.Strings(types)

			for _, t := range types {
				fmt.Printf("    %8s: %d\n", t, full_report[a].Summary[t])
			}

			if !*full_output {
				fmt.Println()
				continue
			}

			fmt.Println()
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

				for _, device := range rack.FailedDevices {
					if *show_uuids {
						fmt.Printf("      %s - %s:\n",
							device.Id,
							device.SystemUuid,
						)
					} else {
						fmt.Printf("      %s:\n", device.Id)
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
