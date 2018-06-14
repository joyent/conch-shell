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
	"strconv"
)

func layoutGetAll(app *cli.Cmd) {
	app.Action = func() {
		rs, err := util.API.GetGlobalRackLayoutSlots()
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
			"Rack ID",
			"Product ID",
			"RU Start",
		})

		for _, r := range rs {
			table.Append([]string{
				r.ID.String(),
				r.RackID.String(),
				r.ProductID.String(),
				strconv.Itoa(r.RUStart),
			})
		}

		table.Render()
	}

}

func layoutGet(app *cli.Cmd) {
	app.Action = func() {
		r, err := util.API.GetGlobalRackLayoutSlot(GLayoutUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Rack ID: %s
Product ID: %s
RU Start: %d

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.RackID.String(),
			r.ProductID.String(),
			r.RUStart,
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}
}

func layoutCreate(app *cli.Cmd) {
	var (
		rackIDOpt    = app.StringOpt("rack-id", "", "UUID of the rack")
		productIDOpt = app.StringOpt("product-id", "", "UUID of the hardware product")
		ruStartOpt   = app.IntOpt("ru-start ru", 0, "Rack unit start number")
	)

	app.Spec = "--rack-id --product-id --ru-start [OPTIONS]"

	app.Action = func() {
		rackID, err := uuid.FromString(*rackIDOpt)
		if err != nil {
			util.Bail(err)
		}
		productID, err := uuid.FromString(*productIDOpt)
		if err != nil {
			util.Bail(err)
		}

		r := conch.GlobalRackLayoutSlot{
			RackID:    rackID,
			ProductID: productID,
			RUStart:   *ruStartOpt,
		}

		if err := util.API.SaveGlobalRackLayoutSlot(&r); err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Rack ID: %s
Product ID: %s
RU Start: %d

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.RackID.String(),
			r.ProductID.String(),
			r.RUStart,
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}
}

func layoutUpdate(app *cli.Cmd) {
	var (
		rackIDOpt    = app.StringOpt("rack-id", "", "UUID of the rack")
		productIDOpt = app.StringOpt("product-id", "", "UUID of the hardware product")
		ruStartOpt   = app.IntOpt("ru-start ru", 0, "Rack unit start number")
	)

	app.Action = func() {
		r, err := util.API.GetGlobalRackLayoutSlot(GLayoutUUID)
		if err != nil {
			util.Bail(err)
		}

		if *rackIDOpt != "" {
			rackID, err := uuid.FromString(*rackIDOpt)
			if err != nil {
				util.Bail(err)
			}
			r.RackID = rackID
		}

		if *productIDOpt != "" {
			productID, err := uuid.FromString(*productIDOpt)
			if err != nil {
				util.Bail(err)
			}
			r.ProductID = productID
		}

		if *ruStartOpt != 0 {
			r.RUStart = *ruStartOpt
		}

		if err := util.API.SaveGlobalRackLayoutSlot(r); err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		fmt.Printf(`
ID: %s
Rack ID: %s
Product ID: %s
RU Start: %d

Created: %s
Updated: %s

`,
			r.ID.String(),
			r.RackID.String(),
			r.ProductID.String(),
			r.RUStart,
			util.TimeStr(r.Created),
			util.TimeStr(r.Updated),
		)
	}
}

func layoutDelete(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.DeleteGlobalRackLayoutSlot(GLayoutUUID); err != nil {
			util.Bail(err)
		}
	}
}
