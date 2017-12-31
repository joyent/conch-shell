// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package mbo provides processing and formatting for data coming out of the
// MBO hardware failure Manta job
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

// DurationFormatCsv formats time.Duration into a string that Excel-ish
// products recognize as a duration
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

// PrettyComponentType takes a raw component type string like 'sas_hdd_num' and
// provides a human friendly name like 'Number of SAS HDDs'
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

// ComponentFailReport represents a validation result, munged for the report
type ComponentFailReport struct {
	DeviceID string                 `json:"device_id"`
	Created  time.Time              `json:"created"`
	Result   conch.ValidationReport `json:"validation_result"`
}

// ComponentFail reports the first time a validation failed and the first time
// it succeeded
type ComponentFail struct {
	FirstFail ComponentFailReport `json:"first_fail"`
	FirstPass ComponentFailReport `json:"first_pass"`
}

// MantaDevice is a map of ComponentFails, keyed by component type
type MantaDevice map[string]ComponentFail

// TypeReportDevice is our munged version of the Component reports, gathering
// device information, the failure data, and remediation time based on the
// timestamps of the FailReports
type TypeReportDevice struct {
	DeviceID        string
	FailureType     string
	ComponentName   string
	RemediationTime time.Duration
	FirstFail       ComponentFailReport
	FirstPass       ComponentFailReport
}

// TypeReport is a summary of failures for a component type
type TypeReport struct {
	All     []float64
	Mean    time.Duration
	Median  time.Duration
	Count   int64
	Devices []TypeReportDevice
}

// Calc calculates mean and median remediation times for a Type Report
func (data *TypeReport) Calc() {
	mean, _ := stats.Mean(data.All)
	median, _ := stats.Median(data.All)

	data.Mean = time.Duration(mean)
	data.Median = time.Duration(median)
}

// DatacenterReport gathers TypeReports for a datacenter, bundling them by
// type, subtype, and vendor
type DatacenterReport struct {
	Name string
	ID   uuid.UUID

	TimesByType          map[string]*TypeReport
	TimesBySubType       map[string]map[string]*TypeReport
	TimesByVendorAndType map[string]map[string]*TypeReport
}

// MantaReport represents both the raw (from json) and processed-for-our-usage
// versions of the Manta job output
type MantaReport struct {
	Raw           map[string]MantaDevice
	Processed     map[string]DatacenterReport
	BeenProcessed bool
}

// NewFromFile loads and parses the JSON output of the Manta job from a file on
// disk
func (r *MantaReport) NewFromFile(path string) (err error) {
	var rawReport map[string]MantaDevice
	reportPath, err := homedir.Expand(path)
	if err != nil {
		return err
	}

	j, err := ioutil.ReadFile(reportPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(j, &rawReport); err != nil {
		return err
	}

	r.Raw = rawReport
	r.BeenProcessed = false
	r.Processed = make(map[string]DatacenterReport)
	return nil
}

// NewFromURL loads and parses the JSON output of the Manta job from an HTTP/S
// URL
func (r *MantaReport) NewFromURL(url string) (err error) {
	var rawReport map[string]MantaDevice
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bodyBytes, &rawReport); err != nil {
		return err
	}

	r.Raw = rawReport
	r.BeenProcessed = false
	r.Processed = make(map[string]DatacenterReport)

	return nil
}

// Process is the heavy lifter, turning the raw JSON into all the various
// summaries and aggregates
func (r *MantaReport) Process(datacenterChoice string, remediationMin int) {
	nullUUID := uuid.UUID{}
	peerRe := regexp.MustCompile("_peer$")

	hardwareProducts := make(map[uuid.UUID]conch.HardwareProduct)

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
		hardwareProducts[prod.ID] = prod
	}

	report := make(map[string]DatacenterReport)

	var p *mpb.Progress
	var bar *mpb.Bar
	if util.Pretty {
		p = mpb.New()
		bar = p.AddBar(int64(len(r.Raw)),
			mpb.AppendDecorators(
				decor.Percentage(3, decor.DSyncSpace),
			),
		)
		defer func() { p.Stop() }()
		fmt.Println("Processing manta report records....")
	}

	for serial, failures := range r.Raw {
		if util.Pretty {
			bar.Increment()
		}
		device, err := util.API.GetDevice(serial)
		if err != nil {
			continue
		}

		if uuid.Equal(device.HardwareProduct, nullUUID) {
			continue
		}

		datacenter := "UNKNOWN"
		datacenterUUID := uuid.UUID{}

		if device.Location.Datacenter.Name != "" {
			datacenter = device.Location.Datacenter.Name
			datacenterUUID = device.Location.Datacenter.ID
		}

		if datacenterChoice != "" {
			re := regexp.MustCompile(fmt.Sprintf("^%s-", datacenterChoice))
			if (datacenterUUID.String() != datacenterChoice) &&
				(datacenter != datacenterChoice) &&
				!re.MatchString(datacenterChoice) {
				continue
			}
		}

		vendor := "UNKNOWN"
		if _, ok := hardwareProducts[device.HardwareProduct]; ok {
			vendor = hardwareProducts[device.HardwareProduct].Vendor
		}

		timesType := make(map[string]*TypeReport)
		timesSubType := make(map[string]map[string]*TypeReport)
		timesVendor := make(map[string]map[string]*TypeReport)

		zeroDuration, _ := time.ParseDuration("0s")
		if _, ok := report[datacenter]; !ok {
			report[datacenter] = DatacenterReport{
				datacenter,
				datacenterUUID,
				timesType,
				timesSubType,
				timesVendor,
			}
		} else {
			timesType = report[datacenter].TimesByType
			timesSubType = report[datacenter].TimesBySubType
			timesVendor = report[datacenter].TimesByVendorAndType
		}

		if _, ok := timesVendor[vendor]; !ok {
			timesVendor[vendor] = make(map[string]*TypeReport)
		}

		for _, failure := range failures {
			failureType := failure.FirstPass.Result.ComponentType

			if (failureType == "") || (failureType == "Undetermined") {
				failureType = "UNKNOWN"
			}

			componentName := failure.FirstPass.Result.ComponentName
			if (componentName == "") || (componentName == "Undetermined") {
				componentName = "UNKNOWN"
			}

			if peerRe.MatchString(componentName) {
				componentName = "switch_peer"
			}

			tFail := failure.FirstFail.Created
			if tFail.IsZero() {
				continue
			}

			tPass := failure.FirstPass.Created
			if tPass.IsZero() {
				continue
			}

			remediationTime := tPass.Sub(tFail)
			if remediationTime.Seconds() < float64(remediationMin) {
				continue
			}

			fullFailure := TypeReportDevice{
				serial,
				failureType,
				componentName,
				remediationTime,
				failure.FirstFail,
				failure.FirstPass,
			}

			if _, ok := timesType[failureType]; !ok {
				timesType[failureType] = &TypeReport{
					make([]float64, 0),
					zeroDuration,
					zeroDuration,
					0,
					make([]TypeReportDevice, 0),
				}
			}
			timesType[failureType].All = append(
				timesType[failureType].All,
				float64(remediationTime),
			)
			timesType[failureType].Count++
			timesType[failureType].Devices = append(
				timesType[failureType].Devices,
				fullFailure,
			)

			if _, ok := timesVendor[vendor][failureType]; !ok {
				timesVendor[vendor][failureType] = &TypeReport{
					make([]float64, 0),
					zeroDuration,
					zeroDuration,
					0,
					make([]TypeReportDevice, 0),
				}
			}
			timesVendor[vendor][failureType].All = append(
				timesVendor[vendor][failureType].All,
				float64(remediationTime),
			)
			timesVendor[vendor][failureType].Count++
			timesVendor[vendor][failureType].Devices = append(
				timesVendor[vendor][failureType].Devices,
				fullFailure,
			)

			if _, ok := timesSubType[failureType]; !ok {
				timesSubType[failureType] = make(map[string]*TypeReport)
			}

			if _, ok := timesSubType[failureType][componentName]; !ok {
				timesSubType[failureType][componentName] = &TypeReport{
					make([]float64, 0),
					zeroDuration,
					zeroDuration,
					0,
					make([]TypeReportDevice, 0),
				}
			}

			timesSubType[failureType][componentName].All = append(
				timesSubType[failureType][componentName].All,
				float64(remediationTime),
			)
			timesSubType[failureType][componentName].Count++
			timesSubType[failureType][componentName].Devices = append(
				timesSubType[failureType][componentName].Devices,
				fullFailure,
			)

		}

	}
	if util.Pretty {
		fmt.Println("Complete...")
		p.Stop()
	}

	for _, az := range report {
		for _, t := range az.TimesByType {
			t.Calc()
		}

		for _, t := range az.TimesBySubType {
			for _, s := range t {
				s.Calc()
			}
		}

		for _, v := range az.TimesByVendorAndType {
			for _, t := range v {
				t.Calc()
			}
		}

	}

	r.Processed = report
	r.BeenProcessed = true
}

// AsText turns the processed data into a text report. If Process has not been
// called, an empty string is returned
func (r *MantaReport) AsText(fullOutput bool, includeVendors bool, includeComponents bool) (output string) {
	if !r.BeenProcessed {
		return ""
	}

	report := r.Processed

	var outputBuf bytes.Buffer
	azNames := make([]string, 0)
	for name := range report {
		azNames = append(azNames, name)
	}
	sort.Strings(azNames)

	for _, name := range azNames {
		az := report[name]
		outputBuf.WriteString(fmt.Sprintf("%s:\n", az.Name))

		if fullOutput || includeVendors {
			outputBuf.WriteString(fmt.Sprintln("  By Vendor:"))
			vendors := make([]string, 0)
			for v := range az.TimesByVendorAndType {
				vendors = append(vendors, v)
			}
			sort.Strings(vendors)

			for _, vendor := range vendors {
				outputBuf.WriteString(fmt.Sprintf("    %s:\n", vendor))

				vendorData := az.TimesByVendorAndType[vendor]

				timeTypes := make([]string, 0)
				for t := range vendorData {
					timeTypes = append(timeTypes, t)
				}
				sort.Strings(timeTypes)

				for _, timeType := range timeTypes {
					data := vendorData[timeType]
					outputBuf.WriteString(fmt.Sprintf(
						"      %s: (%d)\n",
						timeType,
						data.Count,
					))
					outputBuf.WriteString(fmt.Sprintf(
						"        Mean   : %s\n",
						data.Mean,
					))
					outputBuf.WriteString(fmt.Sprintf(
						"        Median : %s\n",
						data.Median,
					))
				}
				outputBuf.WriteString(fmt.Sprintln())
			}
		}

		timeTypes := make([]string, 0)
		for t := range az.TimesByType {
			timeTypes = append(timeTypes, t)
		}
		sort.Strings(timeTypes)

		outputBuf.WriteString(fmt.Sprintln("  By Component Type:"))

		for _, timeType := range timeTypes {
			data := az.TimesByType[timeType]

			outputBuf.WriteString(fmt.Sprintln())
			outputBuf.WriteString(fmt.Sprintf(
				"    %s: (%d)\n",
				timeType,
				data.Count,
			))
			outputBuf.WriteString(fmt.Sprintf(
				"      Mean   : %s\n",
				data.Mean,
			))
			outputBuf.WriteString(fmt.Sprintf(
				"      Median : %s\n",
				data.Median,
			))

			switch timeType {
			case "SAS_SSD":
				continue
			case "SATA_SSD":
				continue
			case "SAS_HDD":
				continue
			case "CPU":
				continue
			}

			if fullOutput || includeComponents {
				outputBuf.WriteString(fmt.Sprintln())
				outputBuf.WriteString(fmt.Sprintln("      By Component:"))

				subTypes := make([]string, 0)
				for t := range az.TimesBySubType[timeType] {
					subTypes = append(subTypes, t)
				}
				sort.Strings(subTypes)

				for _, subType := range subTypes {
					subData := az.TimesBySubType[timeType][subType]
					prettySubType := PrettyComponentType(
						subType,
						timeType,
					)
					outputBuf.WriteString(fmt.Sprintf(
						"        %s: (%d)\n",
						prettySubType,
						subData.Count,
					))
					outputBuf.WriteString(fmt.Sprintf(
						"          Mean   : %s\n",
						subData.Mean,
					))
					outputBuf.WriteString(fmt.Sprintf(
						"          Median : %s\n",
						subData.Median,
					))
				}
			}
		}
		outputBuf.WriteString(fmt.Sprintln())
	}
	return outputBuf.String()

}

// AsCsv takes the processed report data and returns a csv
func (r *MantaReport) AsCsv() (data string) {

	csvVendor := make([][]string, 0)
	csvVendor = append(csvVendor, []string{
		"Datacenter",
		"Vendor",
		"Type",
		"Failure Count",
		"Mean",
		"Median",
	})

	csvComponent := make([][]string, 0)
	csvComponent = append(csvComponent, []string{
		"Datacenter",
		"Type",
		"Component",
		"Failure Count",
		"Mean",
		"Median",
	})

	for name, az := range r.Processed {
		for vendor, vendorData := range az.TimesByVendorAndType {
			for timeType, data := range vendorData {
				csvVendor = append(csvVendor, []string{
					name,
					vendor,
					timeType,
					strconv.FormatInt(data.Count, 10),
					DurationFormatCsv(data.Mean),
					DurationFormatCsv(data.Median),
				})
			}
		}
		for timeType, data := range az.TimesBySubType {
			switch timeType {
			case "SAS_SSD":
				continue
			case "SATA_SSD":
				continue
			case "SAS_HDD":
				continue
			case "CPU":
				continue
			}

			for subType, subData := range data {
				prettySubType := PrettyComponentType(
					subType,
					timeType,
				)

				csvComponent = append(csvComponent, []string{
					name,
					timeType,
					prettySubType,
					strconv.FormatInt(subData.Count, 10),
					DurationFormatCsv(subData.Mean),
					DurationFormatCsv(subData.Median),
				})
			}
		}
	}

	var outputBuf bytes.Buffer
	w := csv.NewWriter(&outputBuf)
	w.WriteAll(csvVendor)
	outputBuf.WriteString("\n")
	w.WriteAll(csvComponent)

	return outputBuf.String()
}
