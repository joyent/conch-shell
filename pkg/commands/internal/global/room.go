// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package global contains commands that operate on structures in the global
// domain, rather than a workspace. API "global admin" access level is required
// for these commands.
package global

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
)

func roomGetAll(app *cli.Cmd) {
	app.Action = func() {
		rs, err := util.API.GetGlobalRooms()
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

func roomGet(app *cli.Cmd) {
	app.Action = func() {
		r, err := util.API.GetGlobalRoom(GRoomUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Datacenter ID: %s
AZ: %s
Alias: %s
Vendor Name: %s

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.DatacenterID.String(),
			r.AZ,
			r.Alias,
			r.VendorName,
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}

}

func roomCreate(app *cli.Cmd) {
	var (
		dcIDOpt       = app.StringOpt("datacenter-id dc", "", "UUID of the datacenter")
		azOpt         = app.StringOpt("az", "", "AZ Name")
		aliasOpt      = app.StringOpt("alias", "", "Room Alias")
		vendorNameOpt = app.StringOpt("vendor-name vn", "", "Vendor Name")
	)
	app.Spec = "--datacenter-id --az --alias [OPTIONS]"

	app.Action = func() {
		dcID, err := uuid.FromString(*dcIDOpt)
		if err != nil {
			util.Bail(err)
		}
		r := conch.GlobalRoom{
			DatacenterID: dcID,
			AZ:           *azOpt,
			Alias:        *aliasOpt,
			VendorName:   *vendorNameOpt,
		}

		if err := util.API.SaveGlobalRoom(&r); err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Datacenter ID: %s
AZ: %s
Alias: %s
Vendor Name: %s

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.DatacenterID.String(),
			r.AZ,
			r.Alias,
			r.VendorName,
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}
}

func roomUpdate(app *cli.Cmd) {
	var (
		dcIDOpt       = app.StringOpt("datacenter-id dc", "", "UUID of the datacenter")
		azOpt         = app.StringOpt("az", "", "AZ Name")
		aliasOpt      = app.StringOpt("alias", "", "Room Alias")
		vendorNameOpt = app.StringOpt("vendor-name vn", "", "Vendor Name")
	)

	app.Action = func() {
		r, err := util.API.GetGlobalRoom(GRoomUUID)
		if err != nil {
			util.Bail(err)
		}

		if *dcIDOpt != "" {
			dcID, err := uuid.FromString(*dcIDOpt)
			if err != nil {
				util.Bail(err)
			}
			r.DatacenterID = dcID
		}

		if *azOpt != "" {
			r.AZ = *azOpt
		}

		if *aliasOpt != "" {
			r.Alias = *aliasOpt
		}

		if *vendorNameOpt != "" {
			r.VendorName = *vendorNameOpt
		}

		if err := util.API.SaveGlobalRoom(&r); err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Datacenter ID: %s
AZ: %s
Alias: %s
Vendor Name: %s

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.DatacenterID.String(),
			r.AZ,
			r.Alias,
			r.VendorName,
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}
}

func roomDelete(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.DeleteGlobalRoom(GRoomUUID); err != nil {
			util.Bail(err)
		}
	}
}

func roomGetAllRacks(app *cli.Cmd) {
	app.Action = func() {
		r, err := util.API.GetGlobalRoom(GRoomUUID)
		if err != nil {
			util.Bail(err)
		}

		rs, err := util.API.GetGlobalRoomRacks(r)
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
			"Name",
			"Role",
		})

		for _, r := range rs {
			role, err := util.API.GetGlobalRackRole(r.RoleID)
			if err != nil {
				util.Bail(err)
			}

			table.Append([]string{
				r.ID.String(),
				r.Name,
				fmt.Sprintf("%s (%s)", role.Name, r.RoleID.String()),
			})
		}

		table.Render()
	}

}
