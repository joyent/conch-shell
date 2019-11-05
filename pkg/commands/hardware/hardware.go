// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hardware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"text/template"

	cli "github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/joyent/conch-shell/pkg/util"
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
      Raid LUN Count {{ .Profile.RaidLunNum }}

    DIMM Count: {{ .Profile.NumDimms }}
    RAM Total:  {{ .Profile.TotalRAM }} GB

    Drives:
    {{ if .Profile.SasHddNum }}
      SAS HDD:
        Count: {{ .Profile.SasHddNum }}
        Size:  {{ .Profile.SasHddSize }}
        Slots: {{ .Profile.SasHddSlots }}
    {{ end }}{{ if .Profile.SataHddNum }}
      SATA HDD:
        Count: {{ .Profile.SataHddNum }}
        Size:  {{ .Profile.SataHddSize }}
        Slots: {{ .Profile.SataHddSlots }}
    {{ end }}{{ if .Profile.SataSsdNum }}
      SATA SSD:
        Count: {{ .Profile.SataSsdNum }}
        Size:  {{ .Profile.SataSsdSize }}
        Slots: {{ .Profile.SataSsdSlots }}
    {{ end }}{{ if .Profile.NvmeSsdNum }}
      NVME SSD:
        Count: {{ .Profile.NvmeSsdNum }}
        Size:  {{ .Profile.NvmeSsdSize }}
        Slots: {{ .Profile.NvmeSsdSlots }}
    {{ end }}
`

func generateDumpableProduct(p conch.HardwareProduct) interface{} {
	// This is all to unhide empty fields marked 'omitempty' in the main
	// struct
	profile := p.Profile
	dumpProfile := struct {
		*conch.HardwareProfile
		HbaFirmware  string `json:"hba_firmware"`
		NvmeSsdNum   int    `json:"nvme_ssd_num"`
		NvmeSsdSize  int    `json:"nvme_ssd_size"`
		NvmeSsdSlots string `json:"nvme_ssd_slots"`
		RaidLunNum   int    `json:"raid_lun_num"`
		SataHddNum   int    `json:"sata_hdd_num"`
		SataHddSize  int    `json:"sata_hdd_size"`
		SataHddSlots string `json:"sata_hdd_slots"`
		SataSsdNum   int    `json:"sata_ssd_num"`
		SataSsdSize  int    `json:"sata_ssd_size"`
		SataSsdSlots string `json:"sata_ssd_slots"`
		TotalPSU     int    `json:"psu_total"`

		ID bool `json:"id,omitempty"`
	}{
		HardwareProfile: &profile,
		HbaFirmware:     profile.HbaFirmware,
		NvmeSsdNum:      profile.NvmeSsdNum,
		NvmeSsdSize:     profile.NvmeSsdSize,
		NvmeSsdSlots:    profile.NvmeSsdSlots,
		RaidLunNum:      profile.RaidLunNum,
		SataHddNum:      profile.SataHddNum,
		SataHddSize:     profile.SataHddSize,
		SataHddSlots:    profile.SataHddSlots,
		SataSsdNum:      profile.SataSsdNum,
		SataSsdSize:     profile.SataSsdSize,
		SataSsdSlots:    profile.SataSsdSlots,
		TotalPSU:        profile.TotalPSU,
	}

	dumpStruct := struct {
		*conch.HardwareProduct
		Prefix            string      `json:"prefix"`
		GenerationName    string      `json:"generation_name"`
		LegacyProductName string      `json:"legacy_product_name"`
		SKU               string      `json:"sku"`
		Profile           interface{} `json:"hardware_product_profile"`

		Specification bool `json:"specification,omitempty"`
		Created       bool `json:"created,omitempty"`
		Updated       bool `json:"updated,omitempty"`
		ID            bool `json:"id,omitempty"`
	}{
		HardwareProduct:   &p,
		Prefix:            p.Prefix,
		GenerationName:    p.GenerationName,
		LegacyProductName: p.LegacyProductName,
		SKU:               p.SKU,
		Profile:           dumpProfile,
	}

	return dumpStruct
}

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

		if util.JSON {
			util.JSONOut(ret)
			return
		}
		var vendor_name string

		if !uuid.Equal(ret.HardwareVendorID, uuid.UUID{}) {
			vendor, err := util.API.GetHardwareVendorByID(ret.HardwareVendorID)
			if err != nil {
				util.Bail(err)
			}
			vendor_name = vendor.Name
		}

		extRet := extendedProduct{&ret, vendor_name}

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
		h, err := util.API.GetHardwareProduct(ProductUUID)
		if err != nil {
			util.Bail(err)
		}

		var specification string

		if h.Specification == nil {
			specification = "{}"
		} else if reflect.TypeOf(h.Specification).String() == "string" {
			specification = h.Specification.(string)
		} else {
			j, err := json.Marshal(h.Specification)
			if err != nil {
				util.Bail(err)
			}

			specification = string(j)
		}

		fmt.Println(specification)
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
			SKU     string `json:"sku"`
			Name    string `json:"name"`
			Alias   string `json:"alias"`
			Prefix  string `json:"prefix"`
			Vendor  string `json:"vendor"`
			Purpose string `json:"purpose"`
		}
		rows := make([]retRow, 0)
		for _, r := range ret {
			var vendor_name string

			if !uuid.Equal(r.HardwareVendorID, uuid.UUID{}) {
				vendor, err := util.API.GetHardwareVendorByID(r.HardwareVendorID)
				if err != nil {
					util.Bail(err)
				}
				vendor_name = vendor.Name
			}

			rows = append(rows, retRow{
				r.ID.String(),
				r.SKU,
				r.Name,
				r.Alias,
				r.Prefix,
				vendor_name,
				r.Profile.Purpose,
			})
		}

		if util.JSON {
			util.JSONOut(rows)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{"ID", "SKU", "Name", "Alias", "Prefix", "Vendor", "Purpose"})

		for _, r := range rows {
			table.Append([]string{r.ID, r.SKU, r.Name, r.Alias, r.Prefix, r.Vendor, r.Purpose})
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
				util.Bail(errors.New("no specification provided"))
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
				util.Bail(errors.New("no specification provided"))
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

func dumpTemplate(cmd *cli.Cmd) {
	// My kingdom for comments in JSON

	cmd.LongDesc = `This is a JSON template for a hardware product, including its hardware profile.
It is used in creating a new hardware product and profile. If you need to update a product and profile, see "conch hardware product :id import".

The specification field of the hardware product is explicitly not supported here since dedicated commands exist to deal with those.

JSON does not allow comments so the following is an example with the required fields marked clearly.

{
     "name": "",  # REQUIRED
     "alias": "", # REQUIRED
     "hardware_vendor_id": "00000000-0000-0000-0000-000000000000", # REQUIRED - see "conch hardware vendors"
     "prefix": "",
     "generation_name": "",
     "legacy_product_name": "",
     "sku": "",
     "hardware_product_profile": {
          "bios_firmware": "", # REQUIRED
          "cpu_type": "",      # REQUIRED
          "cpu_num": 0,        # REQUIRED
          "dimms_num": 0,      # REQUIRED
          "nics_num": 0,       # REQUIRED
          "usb_num": 0,        # REQUIRED
          "purpose": "",       # REQUIRED
          "ram_total": 0,      # REQUIRED
          "rack_unit": 0,      # REQUIRED
          "sas_hdd_num": 0,
          "sas_hdd_size": 0,
          "sas_hdd_slots": "",
          "hba_firmware": "",
          "nvme_ssd_num": 0,
          "nvme_ssd_size": 0,
          "nvme_ssd_slots": "",
          "raid_lun_num": 0,
          "sata_hdd_num": 0,
          "sata_hdd_size": 0,
          "sata_hdd_slots": "",
          "sata_ssd_num": 0,
          "sata_ssd_size": 0,
          "sata_ssd_slots": "",
          "psu_total": 0
     }
}

`
	cmd.Action = func() {
		util.JSONOutIndent(generateDumpableProduct(conch.HardwareProduct{}))
	}
}

func verifyHardwareProduct(p conch.HardwareProduct) {
	if p.Name == "" {
		util.Bail(errors.New("'name' field is required"))
	}

	if p.Alias == "" {
		util.Bail(errors.New("'alias' field is required"))
	}

	if uuid.Equal(p.HardwareVendorID, uuid.UUID{}) {
		util.Bail(errors.New("'hardware_vendor_id' field is required"))
	}

}

func importNewProductJson(app *cli.Cmd) {
	var (
		filePathArg = app.StringArg("FILE", "-", "Path to a JSON file to use as the data source. '-' indicates STDIN")
	)
	app.Spec = "FILE"
	app.Action = func() {
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
			util.Bail(errors.New("no data provided"))
		}

		p := conch.HardwareProduct{}

		if err := json.Unmarshal(b, &p); err != nil {
			util.Bail(err)
		}

		verifyHardwareProduct(p)

		if err := util.API.SaveHardwareProduct(&p); err != nil {
			util.Bail(err)
		}

		ret, err := util.API.GetHardwareProduct(p.ID)
		if err != nil {
			util.Bail(err)
		}

		util.JSONOutIndent(ret)

	}
}

func importChangedProductJson(app *cli.Cmd) {
	var (
		filePathArg = app.StringArg("FILE", "-", "Path to a JSON file to use as the data source. '-' indicates STDIN")
	)
	app.Spec = "FILE"
	app.Action = func() {

		p, err := util.API.GetHardwareProduct(ProductUUID)
		if err != nil {
			util.Bail(err)
		}

		var b []byte
		if *filePathArg == "-" {
			b, err = ioutil.ReadAll(os.Stdin)
		} else {
			b, err = ioutil.ReadFile(*filePathArg)
		}
		if err != nil {
			util.Bail(err)
		}
		if len(string(b)) <= 1 {
			util.Bail(errors.New("no data provided"))
		}

		id := p.ID
		if err := json.Unmarshal(b, &p); err != nil {
			util.Bail(err)
		}
		p.ID = id

		verifyHardwareProduct(p)

		if err := util.API.SaveHardwareProduct(&p); err != nil {
			util.Bail(err)
		}

		ret, err := util.API.GetHardwareProduct(p.ID)
		if err != nil {
			util.Bail(err)
		}

		util.JSONOutIndent(generateDumpableProduct(ret))

	}
}
func exportProductJson(cmd *cli.Cmd) {
	cmd.Action = func() {
		p, err := util.API.GetHardwareProduct(ProductUUID)
		if err != nil {
			util.Bail(err)
		}
		util.JSONOutIndent(generateDumpableProduct(p))
	}
}
