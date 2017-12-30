// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package mbo

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/joyent/conch-shell/pkg/util"
	conch "github.com/joyent/go-conch"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"gopkg.in/montanaflynn/stats.v0"
	uuid "gopkg.in/satori/go.uuid.v1"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"time"
)

func DurationFormatCsv(t time.Duration) (pretty string) {
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

func PrettyComponentType(ugly string, category string) (pretty string) {
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

type ComponentFailReport struct {
	DeviceId string                      `json:"device_id"`
	Created  time.Time                   `json:"created"`
	Result   conch.ConchValidationReport `json:"validation_result"`
}

type ComponentFail struct {
	FirstFail ComponentFailReport `json:"first_fail"`
	FirstPass ComponentFailReport `json:"first_pass"`
}

type MantaDevice map[string]ComponentFail

type TypeReportDevice struct {
	DeviceId        string
	FailureType     string
	ComponentName   string
	RemediationTime time.Duration
	FirstFail       ComponentFailReport
	FirstPass       ComponentFailReport
}

type TypeReport struct {
	All     []float64
	Mean    time.Duration
	Median  time.Duration
	Count   int64
	Devices []TypeReportDevice
}

func (data *TypeReport) Calc() {
	mean, _ := stats.Mean(data.All)
	median, _ := stats.Median(data.All)

	data.Mean = time.Duration(mean)
	data.Median = time.Duration(median)
}

type DatacenterReport struct {
	Name string
	Id   uuid.UUID

	TimesByType          map[string]*TypeReport
	TimesBySubType       map[string]map[string]*TypeReport
	TimesByVendorAndType map[string]map[string]*TypeReport
}

type MantaReport struct {
	Raw           map[string]MantaDevice
	Processed     map[string]DatacenterReport
	BeenProcessed bool
}

func (manta_report *MantaReport) NewFromFile(path string) (err error) {
	var manta_report_raw map[string]MantaDevice
	report_path, err := homedir.Expand(path)
	if err != nil {
		return err
	}

	manta_report_json, err := ioutil.ReadFile(report_path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(manta_report_json, &manta_report_raw); err != nil {
		return err
	}

	manta_report.Raw = manta_report_raw
	manta_report.BeenProcessed = false
	manta_report.Processed = make(map[string]DatacenterReport)
	return nil
}

func (manta_report *MantaReport) NewFromUrl(url string) (err error) {
	var manta_report_raw map[string]MantaDevice
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bodyBytes, &manta_report_raw); err != nil {
		return err
	}

	manta_report.Raw = manta_report_raw
	manta_report.BeenProcessed = false
	manta_report.Processed = make(map[string]DatacenterReport)

	return nil
}

func (manta_report *MantaReport) Process(datacenter_choice string, remediation_min int) {
	null_uuid := uuid.UUID{}
	peer_re := regexp.MustCompile("_peer$")

	hardware_products := make(map[uuid.UUID]conch.ConchHardwareProduct)

	if util.Pretty {
		fmt.Println("Fetching hardware products...")
		util.Spin.Start()
	}

	prods, err := util.API.GetHardwareProducts()
	if util.Pretty {
		util.Spin.Stop()
	}

	if err != nil {
		util.Bail(err)
	}

	for _, prod := range prods {
		hardware_products[prod.Id] = prod
	}

	report := make(map[string]DatacenterReport)

	var p *mpb.Progress
	var bar *mpb.Bar
	if util.Pretty {
		p = mpb.New()
		bar = p.AddBar(int64(len(manta_report.Raw)),
			mpb.AppendDecorators(
				decor.Percentage(3, decor.DSyncSpace),
			),
		)
	}
	defer func() { p.Stop() }()
	if util.Pretty {
		fmt.Println("Processing manta report records....")
	}

	for serial, failures := range manta_report.Raw {
		if util.Pretty {
			bar.Increment()
		}
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

		if datacenter_choice != "" {
			re := regexp.MustCompile(fmt.Sprintf("^%s-", datacenter_choice))
			if (datacenter_uuid.String() != datacenter_choice) &&
				(datacenter != datacenter_choice) &&
				!re.MatchString(datacenter_choice) {
				continue
			}
		}

		vendor := "UNKNOWN"
		if _, ok := hardware_products[device.HardwareProduct]; ok {
			vendor = hardware_products[device.HardwareProduct].Vendor
		}

		times_by_type := make(map[string]*TypeReport)
		times_by_subtype := make(map[string]map[string]*TypeReport)
		times_by_vendor := make(map[string]map[string]*TypeReport)

		zero_duration, err := time.ParseDuration("0s")
		if _, ok := report[datacenter]; !ok {
			report[datacenter] = DatacenterReport{
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
			times_by_vendor[vendor] = make(map[string]*TypeReport)
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
			if remediation_time.Seconds() < float64(remediation_min) {
				continue
			}

			full_failure := TypeReportDevice{
				serial,
				failure_type,
				component_name,
				remediation_time,
				failure.FirstFail,
				failure.FirstPass,
			}

			if _, ok := times_by_type[failure_type]; !ok {
				times_by_type[failure_type] = &TypeReport{
					make([]float64, 0),
					zero_duration,
					zero_duration,
					0,
					make([]TypeReportDevice, 0),
				}
			}
			times_by_type[failure_type].All = append(
				times_by_type[failure_type].All,
				float64(remediation_time),
			)
			times_by_type[failure_type].Count++
			times_by_type[failure_type].Devices = append(
				times_by_type[failure_type].Devices,
				full_failure,
			)

			if _, ok := times_by_vendor[vendor][failure_type]; !ok {
				times_by_vendor[vendor][failure_type] = &TypeReport{
					make([]float64, 0),
					zero_duration,
					zero_duration,
					0,
					make([]TypeReportDevice, 0),
				}
			}
			times_by_vendor[vendor][failure_type].All = append(
				times_by_vendor[vendor][failure_type].All,
				float64(remediation_time),
			)
			times_by_vendor[vendor][failure_type].Count++
			times_by_vendor[vendor][failure_type].Devices = append(
				times_by_vendor[vendor][failure_type].Devices,
				full_failure,
			)

			if _, ok := times_by_subtype[failure_type]; !ok {
				times_by_subtype[failure_type] = make(map[string]*TypeReport)
			}

			if _, ok := times_by_subtype[failure_type][component_name]; !ok {
				times_by_subtype[failure_type][component_name] = &TypeReport{
					make([]float64, 0),
					zero_duration,
					zero_duration,
					0,
					make([]TypeReportDevice, 0),
				}
			}

			times_by_subtype[failure_type][component_name].All = append(
				times_by_subtype[failure_type][component_name].All,
				float64(remediation_time),
			)
			times_by_subtype[failure_type][component_name].Count++
			times_by_subtype[failure_type][component_name].Devices = append(
				times_by_subtype[failure_type][component_name].Devices,
				full_failure,
			)

		}

	}
	if util.Pretty {
		fmt.Println("Complete...")
		p.Stop()
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

	manta_report.Processed = report
	manta_report.BeenProcessed = true
}

func (manta_report *MantaReport) AsText(full_output bool, include_vendors bool, include_components bool) (output string) {
	if !manta_report.BeenProcessed {
		return ""
	}

	report := manta_report.Processed

	var output_buff bytes.Buffer
	az_names := make([]string, 0)
	for name := range report {
		az_names = append(az_names, name)
	}
	sort.Strings(az_names)

	for _, name := range az_names {
		az := report[name]
		output_buff.WriteString(fmt.Sprintf("%s:\n", az.Name))

		if full_output || include_vendors {
			output_buff.WriteString(fmt.Sprintln("  By Vendor:"))
			vendors := make([]string, 0)
			for v := range az.TimesByVendorAndType {
				vendors = append(vendors, v)
			}
			sort.Strings(vendors)

			for _, vendor := range vendors {
				output_buff.WriteString(fmt.Sprintf("    %s:\n", vendor))

				vendor_data := az.TimesByVendorAndType[vendor]

				time_types := make([]string, 0)
				for t := range vendor_data {
					time_types = append(time_types, t)
				}
				sort.Strings(time_types)

				for _, time_type := range time_types {
					data := vendor_data[time_type]
					output_buff.WriteString(fmt.Sprintf(
						"      %s: (%d)\n",
						time_type,
						data.Count,
					))
					output_buff.WriteString(fmt.Sprintf(
						"        Mean   : %s\n",
						data.Mean,
					))
					output_buff.WriteString(fmt.Sprintf(
						"        Median : %s\n",
						data.Median,
					))
				}
				output_buff.WriteString(fmt.Sprintln())
			}
		}

		time_types := make([]string, 0)
		for t := range az.TimesByType {
			time_types = append(time_types, t)
		}
		sort.Strings(time_types)

		output_buff.WriteString(fmt.Sprintln("  By Component Type:"))

		for _, time_type := range time_types {
			data := az.TimesByType[time_type]

			output_buff.WriteString(fmt.Sprintln())
			output_buff.WriteString(fmt.Sprintf(
				"    %s: (%d)\n",
				time_type,
				data.Count,
			))
			output_buff.WriteString(fmt.Sprintf(
				"      Mean   : %s\n",
				data.Mean,
			))
			output_buff.WriteString(fmt.Sprintf(
				"      Median : %s\n",
				data.Median,
			))

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

			if full_output || include_components {
				output_buff.WriteString(fmt.Sprintln())
				output_buff.WriteString(fmt.Sprintln("      By Component:"))

				sub_types := make([]string, 0)
				for t := range az.TimesBySubType[time_type] {
					sub_types = append(sub_types, t)
				}
				sort.Strings(sub_types)

				for _, sub_type := range sub_types {
					sub_data := az.TimesBySubType[time_type][sub_type]
					pretty_sub_type := PrettyComponentType(
						sub_type,
						time_type,
					)
					output_buff.WriteString(fmt.Sprintf(
						"        %s: (%d)\n",
						pretty_sub_type,
						sub_data.Count,
					))
					output_buff.WriteString(fmt.Sprintf(
						"          Mean   : %s\n",
						sub_data.Mean,
					))
					output_buff.WriteString(fmt.Sprintf(
						"          Median : %s\n",
						sub_data.Median,
					))
				}
			}
		}
		output_buff.WriteString(fmt.Sprintln())
	}
	return output_buff.String()

}

func (manta_report *MantaReport) AsCsv() (data string) {

	csv_vendor := make([][]string, 0)
	csv_vendor = append(csv_vendor, []string{
		"Datacenter",
		"Vendor",
		"Type",
		"Failure Count",
		"Mean",
		"Median",
	})

	csv_component := make([][]string, 0)
	csv_component = append(csv_component, []string{
		"Datacenter",
		"Type",
		"Component",
		"Failure Count",
		"Mean",
		"Median",
	})

	for name, az := range manta_report.Processed {
		for vendor, vendor_data := range az.TimesByVendorAndType {
			for time_type, data := range vendor_data {
				csv_vendor = append(csv_vendor, []string{
					name,
					vendor,
					time_type,
					strconv.FormatInt(data.Count, 10),
					DurationFormatCsv(data.Mean),
					DurationFormatCsv(data.Median),
				})
			}
		}
		for time_type, data := range az.TimesBySubType {
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

			for sub_type, sub_data := range data {
				pretty_sub_type := PrettyComponentType(
					sub_type,
					time_type,
				)

				csv_component = append(csv_component, []string{
					name,
					time_type,
					pretty_sub_type,
					strconv.FormatInt(sub_data.Count, 10),
					DurationFormatCsv(sub_data.Mean),
					DurationFormatCsv(sub_data.Median),
				})
			}
		}
	}

	var output_buff bytes.Buffer
	w := csv.NewWriter(&output_buff)
	w.WriteAll(csv_vendor)
	output_buff.WriteString("\n")
	w.WriteAll(csv_component)

	return output_buff.String()
}
