// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"fmt"
	conch "github.com/joyent/go-conch"
	"github.com/mkideal/cli"
	"gopkg.in/montanaflynn/stats.v0"
	uuid "gopkg.in/satori/go.uuid.v1"
	"regexp"
	"sort"
	"time"
)

type mboComponentFailReport struct {
	DeviceId string                      `json:"device_id"`
	Created  time.Time                   `json:"created"`
	Result   conch.ConchValidationReport `json:"validation_result"`
}

type mboComponentFail struct {
	FirstFail mboComponentFailReport `json:"first_fail"`
	FirstPass mboComponentFailReport `json:"first_pass"`
}

type mboMantaDevice map[string]mboComponentFail
type mboMantaReport map[string]mboMantaDevice

type mboHardwareFailureArgs struct {
	cli.Helper
	MantaReport mboMantaReport `cli:"manta-report" usage:"The Manta job output file" parser:"jsonfile"`
	Datacenter  string         `cli:"datacenter" usage:"Limit the output to a particular datacenter"`
	Components  bool           `cli:"include-components" usage:"Breakout failures by component name, as well as type"`
	Vendors     bool           `cli:"include-vendors" usage:"Include vendor data"`
	Full        bool           `cli:"full" usage:"Include all data. --include-components and --include-vendors are ignored"`
	CSV         bool           `cli:"csv" usage:"Output report as CSV."`
}

type mboTypeReport struct {
	All    []float64
	Mean   time.Duration
	Median time.Duration
	Count  int64
}

func mboPrettyComponentType(ugly string, category string) (pretty string) {
	switch ugly {
	case "bios_firmware_version":
		pretty = "BIOS Firmware Revision"
	case "product_name":
		if category == "BIOS" {
			pretty = "Firmware Programming Issue"
		} else {
			pretty = "Product Name"
		}
	case "sas_hdd_num":
		pretty = "Number of SAS HDDs"
	case "sas_ssd_num":
		pretty = "Number of SAS SSDs"
	case "usb_hdd_num":
		pretty = "Number of USB HDDs"
	case "links_up":
		pretty = "Number of Active Links"
	case "nics_num":
		pretty = "Number of Network Interfaces"
	case "num_peer_switch_ports":
		pretty = "Number of Peer Switch Ports"
	case "num_switch_peers":
		pretty = "Number of Switch Peers"
	case "switch_peer":
		pretty = "Switch Peer"
	case "dimm_count":
		pretty = "DIMM Count"
	case "ram_total":
		pretty = "Total RAM Size"
	default:
		pretty = category
	}
	return pretty
}

func mboCalcTimes(data *mboTypeReport) {
	mean, _ := stats.Mean(data.All)
	median, _ := stats.Median(data.All)

	data.Mean = time.Duration(mean)
	data.Median = time.Duration(median)
}

var MboHardwareFailureCmd = &cli.Command{
	Name: "mbo_hardware_failure",
	Desc: "Report that shows info about hardware failures, as per LP-42570027",
	Argv: func() interface{} { return new(mboHardwareFailureArgs) },
	Fn: func(ctx *cli.Context) error {

		type datacenterReport struct {
			Name string
			Id   uuid.UUID

			TimesByType          map[string]*mboTypeReport
			TimesBySubType       map[string]map[string]*mboTypeReport
			TimesByVendorAndType map[string]map[string]*mboTypeReport
		}

		/*****************/

		args, _, api, err := GetStarted(&mboHardwareFailureArgs{}, ctx)

		if err != nil {
			return err
		}

		null_uuid := uuid.UUID{}
		peer_re := regexp.MustCompile("_peer$")

		report := make(map[string]datacenterReport)

		argv := args.Local.(*mboHardwareFailureArgs)
		for serial, failures := range argv.MantaReport {
			device, err := api.GetDevice(serial)
			if err != nil {
				continue
			}

			if uuid.Equal(device.HardwareProduct, null_uuid) {
				continue
			}

			hardware_product, err := api.GetHardwareProduct(device.HardwareProduct)
			if err != nil {
				continue
			}
			vendor := hardware_product.Vendor
			if vendor == "" {
				vendor = "UNKNOWN"
			}

			datacenter := "UNKNOWN"
			datacenter_uuid := uuid.UUID{}

			if device.Location.Datacenter.Name != "" {
				datacenter = device.Location.Datacenter.Name
				datacenter_uuid = device.Location.Datacenter.Id
			}

			if argv.Datacenter != "" {
				if datacenter != argv.Datacenter {
					continue
				}
			}

			times_by_type := make(map[string]*mboTypeReport)
			times_by_subtype := make(map[string]map[string]*mboTypeReport)
			times_by_vendor := make(map[string]map[string]*mboTypeReport)

			zero_duration, err := time.ParseDuration("0s")
			if _, ok := report[datacenter]; !ok {
				report[datacenter] = datacenterReport{
					datacenter,
					datacenter_uuid,
					times_by_type,
					times_by_subtype,
					times_by_vendor,
				}
			} else {
				times_by_type = report[datacenter].TimesByType
				times_by_subtype = report[datacenter].TimesBySubType
				times_by_vendor = report[datacenter].TimesByVendorAndType
			}

			if _, ok := times_by_vendor[vendor]; !ok {
				times_by_vendor[vendor] = make(map[string]*mboTypeReport)
			}

			for _, failure := range failures {
				failure_type := failure.FirstPass.Result.ComponentType
				if (failure_type == "") || (failure_type == "Undetermined") {
					failure_type = "UNKNOWN"
				}

				component_name := failure.FirstPass.Result.ComponentName
				if (component_name == "") || (component_name == "Undetermined") {
					component_name = "UNKNOWN"
				}

				if peer_re.MatchString(component_name) {
					component_name = "switch_peer"
				}

				t_fail := failure.FirstFail.Created
				if t_fail.IsZero() {
					continue
				}

				t_pass := failure.FirstPass.Created
				if t_pass.IsZero() {
					continue
				}

				if _, ok := times_by_type[failure_type]; !ok {
					times_by_type[failure_type] = &mboTypeReport{
						make([]float64, 0),
						zero_duration,
						zero_duration,
						0,
					}
				}
				times_by_type[failure_type].All = append(
					times_by_type[failure_type].All,
					float64(t_pass.Sub(t_fail)),
				)
				times_by_type[failure_type].Count++

				if _, ok := times_by_vendor[vendor][failure_type]; !ok {
					times_by_vendor[vendor][failure_type] = &mboTypeReport{
						make([]float64, 0),
						zero_duration,
						zero_duration,
						0,
					}
				}
				times_by_vendor[vendor][failure_type].All = append(
					times_by_vendor[vendor][failure_type].All,
					float64(t_pass.Sub(t_fail)),
				)
				times_by_vendor[vendor][failure_type].Count++

				if _, ok := times_by_subtype[failure_type]; !ok {
					times_by_subtype[failure_type] = make(map[string]*mboTypeReport)
				}

				if _, ok := times_by_subtype[failure_type][component_name]; !ok {
					times_by_subtype[failure_type][component_name] = &mboTypeReport{
						make([]float64, 0),
						zero_duration,
						zero_duration,
						0,
					}
				}

				times_by_subtype[failure_type][component_name].All = append(
					times_by_subtype[failure_type][component_name].All,
					float64(t_pass.Sub(t_fail)),
				)
				times_by_subtype[failure_type][component_name].Count++

			}

		}

		for _, az := range report {
			for _, time_data := range az.TimesByType {
				mboCalcTimes(time_data)
			}

			for _, type_data := range az.TimesBySubType {
				for _, sub_type := range type_data {
					mboCalcTimes(sub_type)
				}
			}

			for _, vendor_data := range az.TimesByVendorAndType {
				for _, type_data := range vendor_data {
					mboCalcTimes(type_data)
				}
			}

		}

		az_names := make([]string, 0)
		for name := range report {
			az_names = append(az_names, name)
		}
		sort.Strings(az_names)

		for _, name := range az_names {
			az := report[name]
			if !argv.CSV {
				fmt.Printf("%s:\n", az.Name)
			}

			if argv.Full || argv.Vendors {
				if !argv.CSV {
					fmt.Println("  By Vendor:")
				}
				vendors := make([]string, 0)
				for v := range az.TimesByVendorAndType {
					vendors = append(vendors, v)
				}
				sort.Strings(vendors)

				for _, vendor := range vendors {
					if !argv.CSV {
						fmt.Printf("    %s:\n", vendor)
					}

					vendor_data := az.TimesByVendorAndType[vendor]

					time_types := make([]string, 0)
					for t := range vendor_data {
						time_types = append(time_types, t)
					}
					sort.Strings(time_types)

					for _, time_type := range time_types {
						data := vendor_data[time_type]

						if !argv.CSV {
							fmt.Printf("      %s: (%d)\n", time_type, data.Count)
							fmt.Printf("        Mean   : %s\n", data.Mean)
							fmt.Printf("        Median : %s\n", data.Median)
						}
					}
					if !argv.CSV {
						fmt.Println()
					}
				}
			}

			time_types := make([]string, 0)
			for t := range az.TimesByType {
				time_types = append(time_types, t)
			}
			sort.Strings(time_types)

			if !argv.CSV {
				fmt.Println("  By Component Type:")
			}

			for _, time_type := range time_types {
				data := az.TimesByType[time_type]

				if !argv.CSV {
					fmt.Println()
					fmt.Printf("    %s: (%d)\n", time_type, data.Count)
					fmt.Printf("      Mean   : %s\n", data.Mean)
					fmt.Printf("      Median : %s\n", data.Median)
				}

				switch time_type {
				case "SAS_SSD":
					continue
				case "SATA_SSD":
					continue
				case "SAS_HDD":
					continue
				case "CPU":
					continue
				}

				if argv.Full || argv.Components {
					if !argv.CSV {
						fmt.Println()
						fmt.Printf("      By Component:\n")
					}
					sub_types := make([]string, 0)
					for t := range az.TimesBySubType[time_type] {
						sub_types = append(sub_types, t)
					}
					sort.Strings(sub_types)

					for _, sub_type := range sub_types {
						sub_data := az.TimesBySubType[time_type][sub_type]
						pretty_sub_type := mboPrettyComponentType(
							sub_type,
							time_type,
						)

						if !argv.CSV {
							fmt.Printf(
								"        %s: (%d)\n",
								pretty_sub_type,
								sub_data.Count,
							)
							fmt.Printf(
								"          Mean   : %s\n",
								sub_data.Mean,
							)
							fmt.Printf(
								"          Median : %s\n",
								sub_data.Median,
							)
						}
					}
				}
			}

			if !argv.CSV {
				fmt.Println()
			}
		}

		return nil
	},
}
