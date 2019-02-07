// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package devices

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
)

const extendedDeviceTemplate = `
Serial: {{ .ID }}
Hostname: {{.Hostname }}
Asset Tag: {{ .AssetTag }}
Health: {{ .Health }}
System UUID: {{ .SystemUUID }}{{ if .IsTritonSetup }}
Set up for Triton: {{ .TritonSetup.Local }}
  - UUID: {{ .TritonUUID }}{{- end }}{{ if .IsGraduated }}
Graduated: {{ .Graduated.Local }}{{- end }}{{ if .IsValidated }}
Validated: {{ .Validated.Local }}{{- end }}

IPMI: {{ .IPMI }}{{ if .LatestReportIsInvalid }}

** LATEST REPORT IS INVALID **{{- end }}

Location:
  Datacenter: {{ .Location.Datacenter.Name }}
    Vendor:   {{ .Location.Datacenter.VendorName }}

  Rack: {{ .Location.Rack.Name }} 
    - Role: {{ .Location.Rack.Role }}
    - RU:   {{ .Location.Rack.Unit }} of {{ .Location.Rack.Size }}
	- ID:   {{ .Location.Rack.ID }}

Created:      {{ .Created.Local }}
Last Seen:    {{ .LastSeen.Local }}
Last Updated: {{ .Updated.Local }}

{{ if .IsTritonSetup }}
Triton Setup: {{ .TritonSetup.Local }}
Triton UUID:  {{ .TritonUUID }}
{{ end -}}

Hardware:
  SKU: {{ .SKU }}
  Name: {{ .HardwareName }}
{{ if len .Nics }}
Network Interfaces:{{ range .Nics }}
  - {{ .IfaceName }}
    MAC: {{ .MAC }}
    Vendor: {{ .IfaceVendor }}
    Type: {{ .IfaceType }}{{if len .PeerSwitch }}

    Peer: {{ .PeerSwitch }}
      Port: {{ .PeerPort }}
      MAC: {{ .PeerMac }}{{ end }}
{{ end }}{{- end }}
{{ if len .Disks }}
Disks:{{range $name, $slots := .Enclosures}}
  Enclosure: {{ $name }}{{ range $slots }}
    Slot: {{ .Slot }}
        SN:     {{ .SerialNumber }}
        HBA:    {{ .HBA }}
        Type:   {{ .DriveType }}
        Vendor: {{ .Vendor }}
        Model:  {{ .Model }}
        Transport: {{ .Transport }}
        Size:   {{ .Size }}
        Health: {{ .Health }}
        Firmware: {{ .Firmware }}
{{ end }}{{ end }}{{ end }}
{{ if len .Validations }}
Validations:{{ range .Validations }}
  - {{ .Name }}{{ range .Validations }}{{ if .Passed }}
    - pass: {{ .Name }}{{ else }}
	- FAIL: {{ .Name }}
      Results:{{ range .Results }}
        - {{ .Message }}
          Category: {{ .Category }}{{- if len .ComponentID }}
          ComponentID: {{ .ComponentID }}{{ end }}
          Status: {{ .Status }}
{{ end }}{{ end }}{{ end }}{{ end }}{{ end }}
`

func getOne(app *cli.Cmd) {
	var extended = app.BoolOpt("extended", false, "Only affects JSON output. Alters the device structure to provide better access to disk data and provides access to the most recent validation results")
	app.Action = func() {
		if util.JSON {
			if *extended {
				d, err := util.API.GetExtendedDevice(DeviceSerial)
				if err != nil {
					util.Bail(err)
				}
				util.JSONOut(d)
				return
			}

			d, err := util.API.GetDevice(DeviceSerial)
			if err != nil {
				util.Bail(err)
			}
			util.JSONOut(d)
			return
		}

		ed, err := util.API.GetExtendedDevice(DeviceSerial)
		if err != nil {
			util.Bail(err)
		}

		t, err := template.New("extended_device").Parse(extendedDeviceTemplate)
		if err != nil {
			util.Bail(err)
		}

		if err := t.Execute(os.Stdout, ed); err != nil {
			util.Bail(err)
		}
	}
}

func getLocation(app *cli.Cmd) {
	app.Action = func() {
		location, err := util.API.GetDeviceLocation(DeviceSerial)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(location)
			return
		}

		fmt.Printf(`Location for device %s:
  Datacenter:
    Id:   %s
    Name: %s
  Rack:
    Id:   %s
    Name: %s
    Role: %s
    Unit: %d
`,
			DeviceSerial,
			location.Datacenter.ID,
			location.Datacenter.Name,
			location.Rack.ID,
			location.Rack.Name,
			location.Rack.Role,
			location.Rack.Unit,
		)
	}
}

func getSettings(app *cli.Cmd) {
	var keysOnly = app.BoolOpt("keys-only", false, "Only display the setting keys/names")
	app.Action = func() {
		settings, err := util.API.GetDeviceSettings(DeviceSerial)
		if err != nil {
			util.Bail(err)
		}

		keys := make([]string, 0, len(settings))
		for k := range settings {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		if *keysOnly {
			if util.JSON {
				util.JSONOut(keys)
				return
			}

			for _, k := range keys {
				fmt.Println(k)
			}
			return
		}

		if util.JSON {
			util.JSONOut(settings)
			return
		}

		for _, k := range keys {
			fmt.Printf("%s : %v\n", k, settings[k])
		}
	}
}

func getIPMI(app *cli.Cmd) {
	app.Action = func() {

		ipmi, err := util.API.GetDeviceIPMI(DeviceSerial)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(map[string]string{"ipmi": ipmi})
		} else {
			fmt.Println(ipmi)
		}
	}
}

func getSetting(app *cli.Cmd) {
	app.Action = func() {

		setting, err := util.API.GetDeviceSetting(DeviceSerial, DeviceSettingName)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(map[string]string{DeviceSettingName: setting})
		} else {
			fmt.Println(setting)
		}
	}
}

func setSetting(app *cli.Cmd) {
	var settingValueArg = app.StringArg("VALUE", "", "Value of the setting")
	app.Spec = "VALUE"
	app.Action = func() {
		err := util.API.SetDeviceSetting(
			DeviceSerial,
			DeviceSettingName,
			*settingValueArg,
		)
		if err != nil {
			util.Bail(err)
		}
	}
}

func deleteSetting(app *cli.Cmd) {
	app.Action = func() {
		err := util.API.DeleteDeviceSetting(
			DeviceSerial,
			DeviceSettingName,
		)
		if err != nil {
			util.Bail(err)
		}
	}
}

func graduate(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.GraduateDevice(DeviceSerial); err != nil {
			util.Bail(err)
		}
	}
}

func tritonReboot(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.DeviceTritonReboot(DeviceSerial); err != nil {
			util.Bail(err)
		}
	}
}

func setTritonUUID(app *cli.Cmd) {
	var (
		tritonUUID = app.StringArg("UUID", "", "The Triton UUID")
	)
	app.Spec = "UUID"

	app.Action = func() {
		u, err := uuid.FromString(*tritonUUID)
		if err != nil {
			util.Bail(err)
		}

		if err := util.API.SetDeviceTritonUUID(DeviceSerial, u); err != nil {
			util.Bail(err)
		}
	}
}

func markTritonSetup(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.MarkDeviceTritonSetup(DeviceSerial); err != nil {
			util.Bail(err)
		}
	}
}

func setAssetTag(app *cli.Cmd) {
	var (
		assetTagArg = app.StringArg("TAG", "", "The asset tag")
	)
	app.Spec = "TAG"
	app.Action = func() {
		if err := util.API.SetDeviceAssetTag(DeviceSerial, *assetTagArg); err != nil {
			util.Bail(err)
		}

	}
}

func getAssetTag(app *cli.Cmd) {
	app.Action = func() {
		d, err := util.API.GetDevice(DeviceSerial)
		if err != nil {
			util.Bail(err)
		}
		fmt.Println(d.AssetTag)
	}
}

func getReport(app *cli.Cmd) {
	app.Action = func() {
		util.JSON = true
		d, err := util.API.GetDevice(DeviceSerial)
		if err != nil {
			util.Bail(err)
		}
		j, err := json.MarshalIndent(d.LatestReport, "", "  ")
		if err != nil {
			util.Bail(err)
		}
		fmt.Println(string(j))
	}
}

func getTags(app *cli.Cmd) {
	var keysOnly = app.BoolOpt("keys-only", false, "Only display the tag keys/names")
	app.Action = func() {
		settings, err := util.API.GetDeviceTags(DeviceSerial)
		if err != nil {
			util.Bail(err)
		}

		keys := make([]string, 0, len(settings))
		for k := range settings {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		if *keysOnly {
			if util.JSON {
				util.JSONOut(keys)
				return
			}

			for _, k := range keys {
				fmt.Println(k)
			}
			return
		}

		if util.JSON {
			util.JSONOut(settings)
			return
		}

		for _, k := range keys {
			fmt.Printf("%s : %v\n", k, settings[k])
		}
	}
}

func getTag(app *cli.Cmd) {
	app.Action = func() {

		setting, err := util.API.GetDeviceTag(DeviceSerial, DeviceTagName)
		if err != nil {
			util.Bail(err)
		}

		tag := strings.TrimPrefix(DeviceTagName, "tag.")

		if util.JSON {
			util.JSONOut(map[string]string{tag: setting})
		} else {
			fmt.Println(setting)
		}
	}
}

func setTag(app *cli.Cmd) {
	var settingValueArg = app.StringArg("VALUE", "", "Value of the tag")
	app.Spec = "VALUE"
	app.Action = func() {
		err := util.API.SetDeviceTag(
			DeviceSerial,
			DeviceTagName,
			*settingValueArg,
		)
		if err != nil {
			util.Bail(err)
		}
	}
}

func deleteTag(app *cli.Cmd) {
	app.Action = func() {
		err := util.API.DeleteDeviceTag(
			DeviceSerial,
			DeviceTagName,
		)
		if err != nil {
			util.Bail(err)
		}
	}
}
