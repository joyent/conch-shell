// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package global contains commands that operate on strucutres in the global
// domain, rather than a workspace. API "global admin" access level is required
// for these commands.
package global

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
	conch "github.com/joyent/go-conch"
	uuid "gopkg.in/satori/go.uuid.v1"
)

func rackGetAll(app *cli.Cmd) {
	app.Action = func() {
		rs, err := util.API.GetGlobalRacks()
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
			"Datacenter Room ID",
			"Name",
			"Role ID",
		})

		for _, r := range rs {
			table.Append([]string{
				r.ID.String(),
				r.DatacenterRoomID.String(),
				r.Name,
				r.RoleID.String(),
			})
		}

		table.Render()
	}

}

func rackGet(app *cli.Cmd) {
	app.Action = func() {
		r, err := util.API.GetGlobalRack(GRackUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Datacenter Room ID: %s
Name: %s
Role ID: %s

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.DatacenterRoomID.String(),
			r.Name,
			r.RoleID.String(),
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}

}

func rackCreate(app *cli.Cmd) {
	var (
		dcIDOpt   = app.StringOpt("datacenter-room-id dr", "", "UUID of the datacenter room")
		roleIDOpt = app.StringOpt("role-id r", "", "UUID of the rack role")
		nameOpt   = app.StringOpt("name n", "", "Name of the rack")
	)
	app.Spec = "--datacenter-room-id --role-id --name [OPTIONS]"

	app.Action = func() {
		dcID, err := uuid.FromString(*dcIDOpt)
		if err != nil {
			util.Bail(err)
		}
		roleID, err := uuid.FromString(*roleIDOpt)
		if err != nil {
			util.Bail(err)
		}

		r := conch.GlobalRack{
			DatacenterRoomID: dcID,
			RoleID:           roleID,
			Name:             *nameOpt,
		}

		if err := util.API.SaveGlobalRack(&r); err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Datacenter Room ID: %s
Name: %s
Role ID: %s

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.DatacenterRoomID.String(),
			r.Name,
			r.RoleID.String(),
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}
}

func rackUpdate(app *cli.Cmd) {
	var (
		dcIDOpt   = app.StringOpt("datacenter-room-id dr", "", "UUID of the datacenter room")
		roleIDOpt = app.StringOpt("role-id r", "", "UUID of the rack role")
		nameOpt   = app.StringOpt("name n", "", "Name of the rack")
	)

	app.Action = func() {
		r, err := util.API.GetGlobalRack(GRackUUID)
		if err != nil {
			util.Bail(err)
		}
		if *dcIDOpt != "" {
			dcID, err := uuid.FromString(*dcIDOpt)
			if err != nil {
				util.Bail(err)
			}
			r.DatacenterRoomID = dcID
		}

		if *roleIDOpt != "" {
			roleID, err := uuid.FromString(*roleIDOpt)
			if err != nil {
				util.Bail(err)
			}
			r.RoleID = roleID
		}

		if *nameOpt != "" {
			r.Name = *nameOpt
		}

		if err := util.API.SaveGlobalRack(&r); err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Datacenter Room ID: %s
Name: %s
Role ID: %s

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.DatacenterRoomID.String(),
			r.Name,
			r.RoleID.String(),
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}
}
func rackDelete(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.DeleteGlobalRack(GRackUUID); err != nil {
			util.Bail(err)
		}
	}
}
