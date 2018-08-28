// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package devices

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"text/template"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
)

const singleDeviceTemplate = `
Serial  {{ .D.ID }}
Health: {{ .D.Health }}
State:	{{ .D.State }}

Location:
  Datacenter: {{ .D.Location.Datacenter.Name }}
  Rack: {{ .D.Location.Rack.Name }} - RU {{ .D.Location.Rack.Unit }}

{{ if .D.AssetTag }}
Asset Tag: {{ .D.AssetTag }}
{{ end -}}

Created:   {{ .D.Created.Local.Format .DF }}
Last Seen: {{ .D.LastSeen.Local.Format .DF }}

Graduated: {{ .D.Graduated.Local.Format .DF }}
{{ if .TritonSetup }}
Triton Setup: {{ .D.TritonSetup.Local.Format .DF }}
Triton UUID:  {{ .D.TritonUUID }}
{{ end -}}
{{ with $r := .D.LatestReport }}
Hardware:
  {{ if $.HardwareMatches }}Model: {{ $.D.Location.TargetHardwareProduct.Alias }} // {{ $.D.Location.TargetHardwareProduct.Name }}{{ else }}Target Model: {{ $.D.Location.TargetHardwareProduct.Alias }} // {{ $.D.Location.TargetHardwareProduct.Name }}{{ end }}

  BIOS Version: {{ $r.BiosVersion }}
  CPU: {{ $r.Processor.count }}x {{ $r.Processor.type}}
  RAM: {{ $r.Memory.total }}GB ({{ $r.Memory.count }} sticks)

  NICS:{{range $k, $n := $r.Interfaces }}
    - {{ $k }}{{ if $n.state }} - {{ $n.state }}{{ end }}{{ if $n.ipaddr }}
      IP:   {{ $n.ipaddr }}{{ end }}
      MAC:  {{ $n.mac }}{{ if $n.mtu }}
      MTU:  {{ $n.mtu }}{{ end }}
      Type: {{ $n.vendor }} {{ $n.product }}{{ if $n.peer_text }}
      Peer: {{ $n.peer_text }}{{ end }}
  {{ end }}
  Disks:{{ range $k, $d := $r.Disks }}
    - {{ $k }}{{ if $d.health }} - {{ $d.health }}{{ end }}
      Device: {{ $d.device }}
      Model:  {{ $d.model }}{{if $d.guid }}
      GUID:   {{ $d.guid }}{{ end }}{{ if $d.drive_type }}
      Type:   {{ $d.drive_type }}{{ end }}
      Vendor: {{ $d.vendor }}{{ if $d.firmware }}
      Firmware:  {{ $d.firmware }}{{ end }}{{ if $d.enclosure }}
      Enclosure: {{ $d.enclosure }}{{ end }}{{ if $d.slot }}
      Slot:   {{ $d.slot }}{{ end }}
{{ end }}{{ end }}
`

const deviceValidationsTemplate = `
  Validations:{{ if eq (len .) 0 }} NONE{{ else }}{{ range . }}
    - {{ .ComponentType }} 
      Component Name: {{ .ComponentName }}
      Status: {{ if eq .Status 1 }}OK{{ else }}FAIL{{ end }}
      Log: {{ .Log }}
{{ end }}{{ end }}
`

func getOne(app *cli.Cmd) {
	var (
		fullOutput      = app.BoolOpt("full", false, "Provide full data about the devices rather than normal truncated data")
		showValidations = app.BoolOpt("validations", false, "When --full is used without --json, display a list of current validation reports")
		showFailed      = app.BoolOpt("failed-validations", false, "When --full is used with --json, display a list of current validation reports, only those that failed. (Overrides --validations)")
	)

	app.Action = func() {
		device, err := util.API.GetDevice(DeviceSerial)
		if err != nil {
			util.Bail(err)
		}

		if *fullOutput && !util.JSON {
			t, err := template.New("device").Parse(singleDeviceTemplate)
			if err != nil {
				util.Bail(err)
			}
			hm := false
			if uuid.Equal(
				device.HardwareProduct,
				device.Location.TargetHardwareProduct.ID,
			) {
				hm = true
			}

			data := struct {
				D               conch.Device
				TritonSetup     bool
				HardwareMatches bool
				DF              string
			}{
				device,
				!device.TritonSetup.IsZero(),
				hm,
				util.DateFormat,
			}

			if err := t.Execute(os.Stdout, data); err != nil {
				util.Bail(err)
			}

			if *showValidations || *showFailed {
				t, err := template.New("validations").Parse(deviceValidationsTemplate)
				if err != nil {
					util.Bail(err)
				}

				validations := make([]conch.ValidationReport, 0)

				if *showFailed {
					for _, v := range device.Validations {
						if v.Status == 0 {
							validations = append(validations, v)
						}
					}

				} else {
					validations = device.Validations
				}

				if err := t.Execute(os.Stdout, validations); err != nil {
					util.Bail(err)
				}

			}
			return
		}

		devices := make([]conch.Device, 0)
		devices = append(devices, device)

		_ = util.DisplayDevices(devices, *fullOutput)
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
		if err := util.API.SetDeviceSetting(DeviceSerial, DeviceSettingName, *settingValueArg); err != nil {
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
		if d.LatestReport.SerialNumber == "" {
			util.Bail(errors.New("Device has not yet reported"))
		}
		j, err := json.MarshalIndent(d.LatestReport, "", "  ")
		if err != nil {
			util.Bail(err)
		}
		fmt.Println(string(j))
	}
}
