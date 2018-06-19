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
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	"strconv"
)

func roleGetAll(app *cli.Cmd) {
	app.Action = func() {
		rs, err := util.API.GetGlobalRackRoles()
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
			"Rack Size",
		})

		for _, r := range rs {
			table.Append([]string{
				r.ID.String(),
				r.Name,
				strconv.Itoa(r.RackSize),
			})
		}

		table.Render()
	}

}

func roleGet(app *cli.Cmd) {
	app.Action = func() {
		r, err := util.API.GetGlobalRackRole(GRoleUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Name: %s
Rack Size: %d

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.Name,
			r.RackSize,
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}

}

func roleCreate(app *cli.Cmd) {
	var (
		nameOpt     = app.StringOpt("name n", "", "Name of the role")
		rackSizeOpt = app.IntOpt("rack-size", 0, "Rack Size")
	)
	app.Spec = "--name --rack-size [OPTIONS]"

	app.Action = func() {
		r := conch.GlobalRackRole{
			Name:     *nameOpt,
			RackSize: *rackSizeOpt,
		}

		if err := util.API.SaveGlobalRackRole(&r); err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Name: %s
Rack Size: %d

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.Name,
			r.RackSize,
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}
}

func roleUpdate(app *cli.Cmd) {
	var (
		nameOpt     = app.StringOpt("name n", "", "Name of the role")
		rackSizeOpt = app.IntOpt("rack-size", 0, "Rack Size")
	)

	app.Action = func() {
		r, err := util.API.GetGlobalRackRole(GRoleUUID)
		if err != nil {
			util.Bail(err)
		}

		if *nameOpt != "" {
			r.Name = *nameOpt
		}

		if *rackSizeOpt != 0 {
			r.RackSize = *rackSizeOpt
		}

		if err := util.API.SaveGlobalRackRole(r); err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Name: %s
Rack Size: %d

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.Name,
			r.RackSize,
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}
}

func roleDelete(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.DeleteGlobalRackRole(GRoleUUID); err != nil {
			util.Bail(err)
		}
	}
}
