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
	"github.com/joyent/go-conch/pg_time"
	"github.com/mkideal/cli"
	uuid "gopkg.in/satori/go.uuid.v1"
	"sort"
)

type reportFailureArgs struct {
	cli.Helper
	Id         string `cli:"*workspace_id,workspace_uuid,workspace" usage:"ID of the workspace (required)"`
	Breakout   bool   `cli:"breakout" usage:"Instead of just presenting a datacenter summary, breakout results by rack as well (Ignored in the presence of --json)"`
	Uuids      bool   `cli:"uuids" usage:"Show UUIDs where appropriate"`
	Datacenter string `cli:"datacenter" usage:"Limit the output to a particular datacenter UUID"`
}

var ReportFailureCmd = &cli.Command{
	Name: "report_failure",
	Desc: "Report that shows info about hardware failures",
	Argv: func() interface{} { return new(reportFailureArgs) },
	Fn: func(ctx *cli.Context) error {

		type minimalReportDevice struct {
			AssetTag          string                                   `json:"asset_tag"`
			Created           pg_time.ConchPgTime                      `json:"created, int"`
			Graduated         pg_time.ConchPgTime                      `json:"graduated"`
			HardwareProduct   uuid.UUID                                `json:"hardware_product"`
			Health            string                                   `json:"health"`
			Id                string                                   `json:"id"`
			LastSeen          pg_time.ConchPgTime                      `json:"last_seen, int"`
			Location          conch.ConchDeviceLocation                `json:"location"`
			Role              string                                   `json:"role"`
			State             string                                   `json:"state"`
			SystemUuid        uuid.UUID                                `json:"system_uuid"`
			Updated           pg_time.ConchPgTime                      `json:"updated, int"`
			Validated         pg_time.ConchPgTime                      `json:"validated, int"`
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
		/*****************/

		args, _, api, err := GetStarted(&reportFailureArgs{}, ctx)

		if err != nil {
			return err
		}

		argv := args.Local.(*reportFailureArgs)

		var workspace_devices []conch.ConchDevice

		workspace_devices, err = api.GetWorkspaceDevices(
			argv.Id,
			false,
			"",
			"fail",
		)

		if err != nil {
			return err
		}

		for _, d := range workspace_devices {
			full_d, err := api.FillInDevice(d)
			if err != nil {
				return err
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
			if argv.Datacenter != "" {
				if datacenter_uuid.String() != argv.Datacenter {
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

			types := make([]string, 0)
			for k := range full_report[a].Summary {
				types = append(types, k)
			}
			sort.Strings(types)

			for _, t := range types {
				fmt.Printf("    %8s: %d\n", t, full_report[a].Summary[t])
			}

			if !argv.Breakout {
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
				if argv.Uuids {
					fmt.Printf("    %s - %s:\n", rack_name, rack.Rack.Id)
				} else {
					fmt.Printf("    %s:\n", rack_name)
				}

				for _, device := range rack.FailedDevices {
					if argv.Uuids {
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

		return nil
	},
}
