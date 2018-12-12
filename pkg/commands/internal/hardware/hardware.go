// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hardware

import (
	"errors"
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
	"io/ioutil"
	"os"
	"text/template"
)

const singleHWPTemplate = `
ID: {{ .ID }} 
  Name: {{ .Name }}
  Alias: {{ .Alias }}
  Legacy Product Name: {{ .LegacyProductName }}

  SKU:  {{ .SKU }}
  Generation Name: {{ .GenerationName }}
  Vendor: {{ .Vendor }}
  Prefix: {{ .Prefix }}

  Profile: {{ .Profile.ID }}
    Purpose: {{ .Profile.Purpose }}
    BIOS: {{ .Profile.BiosFirmware }}
    HBA Firmware: {{ .Profile.HbaFirmware }}

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
		type extendedProduct struct {
			*conch.HardwareProduct
			Vendor string `json:"vendor"`
		}
		ret, err := util.API.GetHardwareProduct(ProductUUID)
		if err != nil {
			util.Bail(err)
		}

		extRet := extendedProduct{&ret, ""}

		if util.JSON {
			util.JSONOut(extRet)
			return
		}
		t, err := template.New("hw").Parse(singleHWPTemplate)
		if err != nil {
			util.Bail(err)
		}
		if err := t.Execute(os.Stdout, extRet); err != nil {
			util.Bail(err)
		}

	}
}

func getOneSpecification(app *cli.Cmd) {
	app.Action = func() {
		ret, err := util.API.GetHardwareProduct(ProductUUID)
		if err != nil {
			util.Bail(err)
		}
		if ret.Specification == "" {
			fmt.Println("{}")
			return
		}
		fmt.Println(ret.Specification)
	}
}

func getAll(app *cli.Cmd) {
	var (
		fullOutput = app.BoolOpt("full", false, "When --ids-only is *not* used, provide additional data about the devices rather than normal truncated data.")
		idsOnly    = app.BoolOpt("ids-only", false, "Only retrieve hardware product IDs")
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
				r.HardwareVendorID.String(), // BUG(sungo) fetch the vendor name
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

func createOne(app *cli.Cmd) {
	var (
		nameOpt     = app.StringOpt("name", "", "Joyent's Name")
		aliasOpt    = app.StringOpt("alias", "", "Joyent's Name")
		vendorOpt   = app.StringOpt("vendor", "", "Vendor UUID")
		prefixOpt   = app.StringOpt("prefix", "", "Prefix")
		skuOpt      = app.StringOpt("sku", "", "SKU")
		genOpt      = app.StringOpt("generation-name generation gen", "", "Generation Name")
		legacyOpt   = app.StringOpt("legacy-name legacy", "", "Legacy Product Name")
		specOpt     = app.BoolOpt("specification spec", false, "Will provide specification as last arg")
		filePathArg = app.StringArg("FILE", "-", "Path to a JSON file to use as the specification. '-' indicates STDIN")
	)
	app.Spec = "--name --alias --vendor [OPTIONS] [FILE]"

	app.Action = func() {
		vendor, err := uuid.FromString(*vendorOpt)
		if err != nil {
			util.Bail(err)
		}

		h := conch.HardwareProduct{
			Name:              *nameOpt,
			Alias:             *aliasOpt,
			HardwareVendorID:  vendor,
			SKU:               *skuOpt,
			GenerationName:    *genOpt,
			LegacyProductName: *legacyOpt,
			Prefix:            *prefixOpt,
		}

		if *specOpt {
			var b []byte
			var err error
			if *filePathArg == "-" {
				b, err = ioutil.ReadAll(os.Stdin)
			} else {
				b, err = ioutil.ReadFile(*filePathArg)
			}
			if err != nil {
				util.Bail(err)
			}
			if len(string(b)) <= 1 {
				util.Bail(errors.New("No specification provided"))
			}
			h.Specification = string(b)
		}

		if err := util.API.SaveHardwareProduct(&h); err != nil {
			util.Bail(err)
		}

		ret, err := util.API.GetHardwareProduct(&h.ID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(ret)
			return
		}
		fmt.Println(ret.ID)
	}

}

func removeOne(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.DeleteHardwareProduct(ProductUUID); err != nil {
			util.Bail(err)
		}
	}
}

func updateOne(app *cli.Cmd) {
	var (
		nameOpt     = app.StringOpt("name", "", "Joyent's Name")
		aliasOpt    = app.StringOpt("alias", "", "Joyent's Name")
		vendorOpt   = app.StringOpt("vendor", "", "Vendor UUID")
		prefixOpt   = app.StringOpt("prefix", "", "Prefix")
		skuOpt      = app.StringOpt("sku", "", "SKU")
		genOpt      = app.StringOpt("generation-name generation gen", "", "Generation Name")
		legacyOpt   = app.StringOpt("legacy-name legacy", "", "Legacy Product Name")
		specOpt     = app.BoolOpt("specification spec", false, "Will provide specification as last arg")
		filePathArg = app.StringArg("FILE", "-", "Path to a JSON file to use as the specification. '-' indicates STDIN")
	)
	app.Spec = "[OPTIONS] [ --specification ] [FILE]"

	app.Action = func() {
		h, err := util.API.GetHardwareProduct(ProductUUID)
		if err != nil {
			util.Bail(err)
		}

		if *nameOpt != "" {
			h.Name = *nameOpt
		}

		if *aliasOpt != "" {
			h.Alias = *aliasOpt
		}

		if *vendorOpt != "" {
			vendor, err := uuid.FromString(*vendorOpt)
			if err != nil {
				util.Bail(err)
			}

			h.HardwareVendorID = vendor
		}

		if *prefixOpt != "" {
			h.Prefix = *prefixOpt
		}

		if *skuOpt != "" {
			h.SKU = *skuOpt
		}

		if *genOpt != "" {
			h.GenerationName = *genOpt
		}

		if *legacyOpt != "" {
			h.LegacyProductName = *legacyOpt
		}

		if *specOpt {
			var b []byte
			var err error
			if *filePathArg == "-" {
				b, err = ioutil.ReadAll(os.Stdin)
			} else {
				b, err = ioutil.ReadFile(*filePathArg)
			}
			if err != nil {
				util.Bail(err)
			}
			if len(string(b)) <= 1 {
				util.Bail(errors.New("No specification provided"))
			}
			h.Specification = string(b)
		}

		if err := util.API.SaveHardwareProduct(&h); err != nil {
			util.Bail(err)
		}

		ret, err := util.API.GetHardwareProduct(h.ID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(ret)
			return
		}
	}
}
