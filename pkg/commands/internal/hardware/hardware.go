// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hardware

import (
	"fmt"
	"github.com/joyent/conch-shell/pkg/util"
	"gopkg.in/jawher/mow.cli.v1"
	"os"
	"text/template"
)

const singleHWPTemplate = `
ID: {{ .ID }} 
  Name: {{ .Name }}
  Alias: {{ .Alias }}
  Vendor: {{ .Vendor }}
  Prefix: {{ .Prefix }}

  Profile: {{ .Profile.ID }}
    Purpose: {{ .Profile.Purpose }}
    BIOS: {{ .Profile.BiosFirmware }}

    CPU Count: {{ .Profile.NumCPU }}
    CPU Type:  {{ .Profile.CPUType }}

    NIC Count: {{ .Profile.NumNics }}
    PSU Total: {{ .Profile.TotalPSU }}{{ if .Profile.NumUSB }}
    USB Count: {{ .Profile.NumUSB }}{{ end }}

    DIMM Count: {{ .Profile.NumDimms }}
    RAM Total:  {{ .Profile.TotalRAM }} GB
    {{ if .Profile.SASNum }} 
    SAS Count: {{ .Profile.SASNum }}
    SAS Size:  {{ .Profile.SizeSAS }}
    SAS Slots: {{ .Profile.SlotsSAS }}
    {{ end }}{{ if .Profile.NumSATA }}
    SATA Count: {{ .Profile.NumSATA }}
    SATA Size:  {{ .Profile.SizeSATA }}
    SATA Slots: {{ .Profile.SlotsSATA }}
    {{end}}{{ if .Profile.NumSSD }}
    SSD Count: {{ .Profile.NumSSD }}
    SSD Size:  {{ .Profile.SizeSSD }}
    SSD Slots: {{ .Profile.SlotsSSD }}
    {{ end }}{{ if ne .Profile.Zpool.ID.String "00000000-0000-0000-0000-000000000000"}}
    Zpool: {{ .Profile.Zpool.ID }}
      Cache:     {{ .Profile.Zpool.Cache }}
      Log:       {{ .Profile.Zpool.Log }}
      Disks Per: {{ .Profile.Zpool.DisksPer }}
      Spare:     {{ .Profile.Zpool.Spare }}
      VDEV N:    {{ .Profile.Zpool.VdevN }}
      VDEV T:    {{ .Profile.Zpool.VdevT }}
    {{ end }}
`

func getOne(app *cli.Cmd) {
	app.Action = func() {
		ret, err := util.API.GetHardwareProduct(ProductUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(ret)
			return
		}
		t, err := template.New("hw").Parse(singleHWPTemplate)
		if err != nil {
			util.Bail(err)
		}
		if err := t.Execute(os.Stdout, ret); err != nil {
			util.Bail(err)
		}

	}
}

func getAll(app *cli.Cmd) {
	var (
		fullOutput = app.BoolOpt("full", false, "When --ids-only is *not* used, provide additional data about the devices rather than normal truncated data. Note: this slows things down immensely")
		idsOnly    = app.BoolOpt("ids-only", false, "Only retrieve device IDs")
	)

	app.Action = func() {
		ret, err := util.API.GetHardwareProducts()
		if err != nil {
			util.Bail(err)
		}

		if *idsOnly {
			ids := make([]string, 0)
			for _, r := range ret {
				ids = append(ids, r.ID.String())
			}
			if util.JSON {
				util.JSONOut(ids)
			} else {
				for _, id := range ids {
					fmt.Println(id)
				}
			}
			return
		}

		if *fullOutput {
			if util.JSON {
				util.JSONOut(ret)
				return
			}
			t, err := template.New("hw").Parse(singleHWPTemplate)
			if err != nil {
				util.Bail(err)
			}
			for _, r := range ret {
				if err := t.Execute(os.Stdout, r); err != nil {
					util.Bail(err)
				}
			}
			return
		}
		type retRow struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Alias   string `json:"alias"`
			Prefix  string `json:"prefix"`
			Vendor  string `json:"vendor"`
			Purpose string `json:"purpose"`
		}
		rows := make([]retRow, 0)
		for _, r := range ret {
			rows = append(rows, retRow{
				r.ID.String(),
				r.Name,
				r.Alias,
				r.Prefix,
				r.Vendor,
				r.Profile.Purpose,
			})
		}

		if util.JSON {
			util.JSONOut(rows)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{"ID", "Name", "Alias", "Prefix", "Vendor", "Purpose"})

		for _, r := range rows {
			table.Append([]string{r.ID, r.Name, r.Alias, r.Prefix, r.Vendor, r.Purpose})
		}

		table.Render()
	}
}
