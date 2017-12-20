// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package reports

// BUG(sungo): use getHardwareProducts and just look up by uuid rather than fetching it for every device

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/joyent/conch-shell/util"
	conch "github.com/joyent/go-conch"
	homedir "github.com/mitchellh/go-homedir"
	"gopkg.in/jawher/mow.cli.v1"
	"gopkg.in/montanaflynn/stats.v0"
	uuid "gopkg.in/satori/go.uuid.v1"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

func mboDurationFormatCsv(t time.Duration) (pretty string) {
	seconds := int64(t.Seconds()) % 60
	minutes := int64(t.Minutes()) % 60
	hours := int64(t.Hours()) % 24

	days := int64(t/(24*time.Hour)) % 365 % 7
	weeks := int64(t/(24*time.Hour)) / 7 % 52

	// To make this work as a duration in Excel and Google Sheets, the duration
	// string must be HH:MM:SS so we need to add things back in.
	// I'm also ignoring years here on purpose.
	hours = hours + (days * 24) + (weeks * 7 * 24)

	return fmt.Sprintf(
		"%s:%s:%s",
		strconv.FormatInt(hours, 10),
		strconv.FormatInt(minutes, 10),
		strconv.FormatInt(seconds, 10),
	)
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

type mboTypeReport struct {
	All    []float64
	Mean   time.Duration
	Median time.Duration
	Count  int64
}

func (data *mboTypeReport) Calc() {
	mean, _ := stats.Mean(data.All)
	median, _ := stats.Median(data.All)

	data.Mean = time.Duration(mean)
	data.Median = time.Duration(median)
}

func mboHardwareFailures(app *cli.Cmd) {
	var (
		full_output        = app.BoolOpt("full", false, "Instead of just presenting a datacenter summary, break results out by rack as well. Has no effect on --json")
		datacenter_choice  = app.StringOpt("datacenter az", "", "Limit the output to a particular datacenter by UUID, partial UUID, or string name")
		csv_output         = app.BoolOpt("csv", false, "Output report as CSV. Assumes --full and overrides --json")
		include_vendors    = app.BoolOpt("include-vendors", false, "Include vendor data")
		include_components = app.BoolOpt("include-components", false, "Break out failures by components")
		manta_report_path  = app.StringOpt("manta-report", "", "Path to Manta job output file")
		remediation_min = app.IntOpt("remediation-minimum", 90, "For a failure to be considered, its remediation time must be greater than or equal to this number")
	)

	app.Spec = "--manta-report [--full] [--csv] [--include-vendors] [--include-components] [--datacenter] [--remediation-minimum]"

	app.Action = func() {

		type datacenterReport struct {
			Name string
			Id   uuid.UUID

			TimesByType          map[string]*mboTypeReport
			TimesBySubType       map[string]map[string]*mboTypeReport
			TimesByVendorAndType map[string]map[string]*mboTypeReport
		}

		/*********************************/

		if *csv_output {
			*full_output = true
		}

		report_path, err := homedir.Expand(*manta_report_path)
		if err != nil {
			util.Bail(err)
		}

		manta_report_json, err := ioutil.ReadFile(report_path)
		if err != nil {
			util.Bail(err)
		}

		var manta_report mboMantaReport
		if err := json.Unmarshal(manta_report_json, &manta_report); err != nil {
			util.Bail(err)
		}

		csv_vendor := make([][]string, 0)
		csv_vendor = append(csv_vendor, []string{
			"Datacenter",
			"Vendor",
			"Type",
			"Count",
			"Mean",
			"Median",
		})

		csv_component := make([][]string, 0)
		csv_component = append(csv_component, []string{
			"Datacenter",
			"Type",
			"Component",
			"Count",
			"Mean",
			"Median",
		})

		null_uuid := uuid.UUID{}
		peer_re := regexp.MustCompile("_peer$")

		report := make(map[string]datacenterReport)

		for serial, failures := range manta_report {
			device, err := util.API.GetDevice(serial)
			if err != nil {
				continue
			}

			if uuid.Equal(device.HardwareProduct, null_uuid) {
				continue
			}

			datacenter := "UNKNOWN"
			datacenter_uuid := uuid.UUID{}

			if device.Location.Datacenter.Name != "" {
				datacenter = device.Location.Datacenter.Name
				datacenter_uuid = device.Location.Datacenter.Id
			}

			if *datacenter_choice != "" {
				re := regexp.MustCompile(fmt.Sprintf("^%s-", *datacenter_choice))
				if (datacenter_uuid.String() != *datacenter_choice) &&
					(datacenter != *datacenter_choice) &&
					!re.MatchString(*datacenter_choice) {
					continue
				}
			}

			hardware_product, err := util.API.GetHardwareProduct(device.HardwareProduct)
			if err != nil {
				continue
			}
			vendor := hardware_product.Vendor
			if vendor == "" {
				vendor = "UNKNOWN"
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

				remediation_time := t_pass.Sub(t_fail)
				if remediation_time.Seconds() < float64(*remediation_min) {
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
					float64(remediation_time),
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
					float64(remediation_time),
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
					float64(remediation_time),
				)
				times_by_subtype[failure_type][component_name].Count++

			}

		}

		for _, az := range report {
			for _, time_data := range az.TimesByType {
				time_data.Calc()
			}

			for _, type_data := range az.TimesBySubType {
				for _, sub_type := range type_data {
					sub_type.Calc()
				}
			}

			for _, vendor_data := range az.TimesByVendorAndType {
				for _, type_data := range vendor_data {
					type_data.Calc()
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
			if !*csv_output {
				fmt.Printf("%s:\n", az.Name)
			}

			if *full_output || *include_vendors {
				if !*csv_output {
					fmt.Println("  By Vendor:")
				}
				vendors := make([]string, 0)
				for v := range az.TimesByVendorAndType {
					vendors = append(vendors, v)
				}
				sort.Strings(vendors)

				for _, vendor := range vendors {
					if !*csv_output {
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

						csv_vendor = append(csv_vendor, []string{
							name,
							vendor,
							time_type,
							strconv.FormatInt(data.Count, 10),
							mboDurationFormatCsv(data.Mean),
							mboDurationFormatCsv(data.Median),
						})

						if !*csv_output {
							fmt.Printf("      %s: (%d)\n", time_type, data.Count)
							fmt.Printf("        Mean   : %s\n", data.Mean)
							fmt.Printf("        Median : %s\n", data.Median)
						}
					}
					if !*csv_output {
						fmt.Println()
					}
				}
			}

			time_types := make([]string, 0)
			for t := range az.TimesByType {
				time_types = append(time_types, t)
			}
			sort.Strings(time_types)

			if !*csv_output {
				fmt.Println("  By Component Type:")
			}

			for _, time_type := range time_types {
				data := az.TimesByType[time_type]

				if !*csv_output {
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

				if *full_output || *include_components {
					if !*csv_output {
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

						csv_component = append(csv_component, []string{
							name,
							time_type,
							pretty_sub_type,
							strconv.FormatInt(sub_data.Count, 10),
							mboDurationFormatCsv(sub_data.Mean),
							mboDurationFormatCsv(sub_data.Median),
						})

						if !*csv_output {
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

			if !*csv_output {
				fmt.Println()
			}
		}

		if *csv_output {
			w := csv.NewWriter(os.Stdout)
			w.WriteAll(csv_vendor)
			fmt.Println()
			w.WriteAll(csv_component)
		}
	}
}
