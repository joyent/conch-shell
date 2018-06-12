// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package global

import (
	"fmt"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
	conch "github.com/joyent/go-conch"
)

func dcGetAll(app *cli.Cmd) {
	app.Action = func() {
		d, err := util.API.GetGlobalDatacenters()
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(d)
			return
		}
		table := util.GetMarkdownTable()
		table.SetHeader([]string{
			"ID",
			"Region",
			"Vendor",
			"Vendor Name",
			"Location",
		})

		for _, dc := range d {
			table.Append([]string{
				dc.ID.String(),
				dc.Region,
				dc.Vendor,
				dc.VendorName,
				dc.Location,
			})
		}

		table.Render()
	}
}

func dcGet(app *cli.Cmd) {
	app.Action = func() {
		d, err := util.API.GetGlobalDatacenter(GdcUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(d)
			return
		}

		fmt.Printf(`
ID: %s
Region: %s
Vendor: %s
Vendor Name: %s
Location: %s

Created: %s
Updated: %s

`,
			d.ID.String(),
			d.Region,
			d.Vendor,
			d.VendorName,
			d.Location,
			util.TimeStr(d.Created),
			util.TimeStr(d.Updated),
		)
	}
}

func dcCreate(app *cli.Cmd) {
	var (
		regionOpt     = app.StringOpt("region", "", "Region identifier")
		vendorOpt     = app.StringOpt("vendor", "", "Vendor")
		vendorNameOpt = app.StringOpt("vendor-name", "", "Vendor Name")
		locationOpt   = app.StringOpt("location", "", "Location")
	)
	app.Spec = "--region --vendor --location [OPTIONS]"

	app.Action = func() {
		d := conch.GlobalDatacenter{
			Region:     *regionOpt,
			Vendor:     *vendorOpt,
			VendorName: *vendorNameOpt,
			Location:   *locationOpt,
		}

		err := util.API.SaveGlobalDatacenter(&d)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(d)
			return
		}

		fmt.Printf(`
ID: %s
Region: %s
Vendor: %s
Vendor Name: %s
Location: %s

Created: %s
Updated: %s

`,
			d.ID.String(),
			d.Region,
			d.Vendor,
			d.VendorName,
			d.Location,
			util.TimeStr(d.Created),
			util.TimeStr(d.Updated),
		)

	}
}

func dcDelete(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.DeleteGlobalDatacenter(GdcUUID); err != nil {
			util.Bail(err)
		}
	}
}

func dcUpdate(app *cli.Cmd) {
	var (
		regionOpt     = app.StringOpt("region", "", "Region identifier")
		vendorOpt     = app.StringOpt("vendor", "", "Vendor")
		vendorNameOpt = app.StringOpt("vendor-name", "", "Vendor Name")
		locationOpt   = app.StringOpt("location", "", "Location")
	)

	app.Action = func() {
		d, err := util.API.GetGlobalDatacenter(GdcUUID)
		if err != nil {
			util.Bail(err)
		}

		if *regionOpt != "" {
			d.Region = *regionOpt
		}
		if *vendorOpt != "" {
			d.Vendor = *vendorOpt
		}
		if *vendorNameOpt != "" {
			d.VendorName = *vendorNameOpt
		}
		if *locationOpt != "" {
			d.Location = *locationOpt
		}

		if err := util.API.SaveGlobalDatacenter(&d); err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(d)
			return
		}

		fmt.Printf(`
ID: %s
Region: %s
Vendor: %s
Vendor Name: %s
Location: %s

Created: %s
Updated: %s

`,
			d.ID.String(),
			d.Region,
			d.Vendor,
			d.VendorName,
			d.Location,
			util.TimeStr(d.Created),
			util.TimeStr(d.Updated),
		)

	}

}
