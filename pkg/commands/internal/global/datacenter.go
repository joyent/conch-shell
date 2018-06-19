// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package global

import (
	"fmt"

	gotree "github.com/DiSiqueira/GoTree"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	"sort"
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

func dcGetAllRooms(app *cli.Cmd) {
	app.Action = func() {
		d, err := util.API.GetGlobalDatacenter(GdcUUID)
		if err != nil {
			util.Bail(err)
		}

		rs, err := util.API.GetGlobalDatacenterRooms(d)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(rs)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{
			"ID",
			"Datacenter ID",
			"AZ",
			"Alias",
			"Vendor Name",
		})

		for _, r := range rs {
			table.Append([]string{
				r.ID.String(),
				r.DatacenterID.String(),
				r.AZ,
				r.Alias,
				r.VendorName,
			})
		}
		table.Render()
	}
}

func dcAllTheThingsTree(app *cli.Cmd) {
	app.Action = func() {
		hwProds := make(map[string]conch.HardwareProduct)

		d, err := util.API.GetGlobalDatacenter(GdcUUID)
		if err != nil {
			util.Bail(err)
		}

		tree := gotree.GTStructure{}
		tree.Name = fmt.Sprintf("DC: %s (%s)", d.Region, d.ID)

		rs, err := util.API.GetGlobalDatacenterRooms(d)
		if err != nil {
			util.Bail(err)
		}

		for _, room := range rs {
			roomTree := gotree.GTStructure{}
			roomTree.Name = fmt.Sprintf("Room: %s (%s)", room.AZ, room.ID)

			racks, err := util.API.GetGlobalRoomRacks(room)
			if err != nil {
				util.Bail(err)
			}

			for _, rack := range racks {
				rackTree := gotree.GTStructure{}
				rackTree.Name = fmt.Sprintf("Rack: %s (%s)", rack.Name, rack.ID)

				ls, err := util.API.GetGlobalRackLayout(rack)
				if err != nil {
					util.Bail(err)
				}

				sort.Sort(byRUStart(ls))
				for _, layout := range ls {
					var hw conch.HardwareProduct

					hw, ok := hwProds[layout.ProductID.String()]
					if !ok {
						hw, err = util.API.GetHardwareProduct(layout.ProductID)
						if err != nil {
							util.Bail(err)
						}
						hwProds[layout.ProductID.String()] = hw
					}

					layoutTree := gotree.GTStructure{}
					layoutTree.Name = fmt.Sprintf(
						"RU: %d | Product: %s",
						layout.RUStart,
						hw.Name,
					)
					rackTree.Items = append(rackTree.Items, layoutTree)
				}

				roomTree.Items = append(roomTree.Items, rackTree)
			}

			tree.Items = append(tree.Items, roomTree)
		}

		gotree.PrintTree(tree)
	}
}
