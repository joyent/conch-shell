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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

func displayOneGlobalRack(r conch.GlobalRack) {

	fmt.Printf(`
ID: %s
Datacenter Room ID: %s
Name: %s
Role ID: %s
Serial Number: %s
Asset Tag: %s

Created: %s
Updated: %s

`,
		r.ID.String(),
		r.DatacenterRoomID.String(),
		r.Name,
		r.RoleID.String(),
		r.SerialNumber,
		r.AssetTag,
		util.TimeStr(r.Created),
		util.TimeStr(r.Updated),
	)
}

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
			"Serial Number",
			"Asset Tag",
		})

		for _, r := range rs {
			table.Append([]string{
				r.ID.String(),
				r.DatacenterRoomID.String(),
				r.Name,
				r.RoleID.String(),
				r.SerialNumber,
				r.AssetTag,
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

		displayOneGlobalRack(r)
	}
}

func rackCreate(app *cli.Cmd) {
	var (
		dcIDOpt     = app.StringOpt("datacenter-room-id dr", "", "UUID of the datacenter room")
		roleIDOpt   = app.StringOpt("role-id r", "", "UUID of the rack role")
		nameOpt     = app.StringOpt("name n", "", "Name of the rack")
		snOpt       = app.StringOpt("serial-number sn", "", "Serial number")
		assetTagOpt = app.StringOpt("asset-tag a", "", "Asset tag")
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
			SerialNumber:     *snOpt,
			AssetTag:         *assetTagOpt,
		}

		if err := util.API.SaveGlobalRack(&r); err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}

		displayOneGlobalRack(r)
	}
}

func rackUpdate(app *cli.Cmd) {
	var (
		dcIDOpt     = app.StringOpt("datacenter-room-id dr", "", "UUID of the datacenter room")
		roleIDOpt   = app.StringOpt("role-id r", "", "UUID of the rack role")
		nameOpt     = app.StringOpt("name n", "", "Name of the rack")
		snOpt       = app.StringOpt("serial-number sn", "", "Serial number")
		assetTagOpt = app.StringOpt("asset-tag a", "", "Asset tag")
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

		if *snOpt != "" {
			r.SerialNumber = *snOpt
		}

		if *assetTagOpt != "" {
			r.AssetTag = *assetTagOpt
		}

		if err := util.API.SaveGlobalRack(&r); err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(r)
			return
		}
		displayOneGlobalRack(r)
	}
}
func rackDelete(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.DeleteGlobalRack(GRackUUID); err != nil {
			util.Bail(err)
		}
	}
}

type byRUStart []conch.GlobalRackLayoutSlot

func (b byRUStart) Len() int {
	return len(b)
}
func (b byRUStart) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byRUStart) Less(i, j int) bool {
	return b[i].RUStart > b[j].RUStart
}

func rackLayout(app *cli.Cmd) {
	app.Action = func() {
		r, err := util.API.GetGlobalRack(GRackUUID)
		if err != nil {
			util.Bail(err)
		}

		rs, err := util.API.GetGlobalRackLayout(r)
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
			"Product",
			"RU Start",
		})

		sort.Sort(byRUStart(rs))

		for _, r := range rs {
			prod, err := util.API.GetHardwareProduct(r.ProductID)
			if err != nil {
				util.Bail(err)
			}
			table.Append([]string{
				r.ID.String(),
				prod.Name,
				strconv.Itoa(r.RUStart),
			})
		}

		table.Render()
	}
}

type importLayoutSlot struct {
	RUStart      int       `json:"ru_start"`
	ProductID    uuid.UUID `json:"product_id,omitempty"`
	ProductName  string    `json:"product_name,omitempty"`
	ProductAlias string    `json:"product_alias,omitempty"`
}

type importLayout []importLayoutSlot

func rackExportLayout(cmd *cli.Cmd) {
	cmd.Action = func() {
		rack, err := util.API.GetGlobalRack(GRackUUID)
		if err != nil {
			util.Bail(err)
		}

		// Get the current state of the world
		existingLayout, err := util.API.GetGlobalRackLayout(rack)
		if err != nil {
			util.Bail(err)
		}

		var output importLayout
		if len(existingLayout) == 0 {
			output = append(output, importLayoutSlot{
				RUStart:      0,
				ProductID:    uuid.UUID{},
				ProductName:  "Product Name",
				ProductAlias: "Product Alias",
			})
		} else {

			for _, l := range existingLayout {
				hw, err := util.API.GetHardwareProduct(l.ProductID)
				if err != nil {
					util.Bail(err)
				}
				output = append(output, importLayoutSlot{
					RUStart:      l.RUStart,
					ProductID:    hw.ID,
					ProductName:  hw.Name,
					ProductAlias: hw.Alias,
				})
			}
		}

		util.JSONOutIndent(output)
	}
}

func rackImportLayout(cmd *cli.Cmd) {
	var (
		filePathArg  = cmd.StringArg("FILE", "-", "Path to a JSON file that defines the layout. '-' indicates STDIN")
		overwriteOpt = cmd.BoolOpt("overwrite", false, "If the rack has an existing layout, *overwrite* it. This is a destructive action")
	)

	cmd.Spec = "[OPTIONS] [FILE]"
	cmd.Action = func() {
		util.JSON = true
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

		var importedLayout importLayout

		if err := json.Unmarshal(b, &importedLayout); err != nil {
			util.Bail(err)
		}

		rack, err := util.API.GetGlobalRack(GRackUUID)
		if err != nil {
			util.Bail(err)
		}

		// Get the current state of the world
		existingLayout, err := util.API.GetGlobalRackLayout(rack)
		if err != nil {
			util.Bail(err)
		}

		if len(existingLayout) > 0 {
			if !*overwriteOpt {
				util.Bail(errors.New("rack already has a layout. Use --overwrite to overwrite"))
			}
		}

		// We need to support the use of product names and aliases in the
		// import so they're readable by humans. We lack a way of doing API
		// lookups on these properties so we pull them all down and create maps
		// on our own.
		productsL, err := util.API.GetHardwareProducts()
		if err != nil {
			util.Bail(err)
		}

		productsAlias := make(map[string]conch.HardwareProduct)
		productsName := make(map[string]conch.HardwareProduct)
		productsID := make(map[string]conch.HardwareProduct)

		for _, p := range productsL {
			productsAlias[p.Alias] = p
			productsName[p.Name] = p
			productsID[p.ID.String()] = p
		}

		var finalLayout []conch.GlobalRackLayoutSlot

		for _, l := range importedLayout {
			if uuid.Equal(l.ProductID, uuid.UUID{}) {
				if l.ProductName != "" {
					p, ok := productsName[l.ProductName]
					if ok {
						l.ProductID = p.ID
					}
				} else if l.ProductAlias != "" {
					p, ok := productsAlias[l.ProductAlias]
					if ok {
						l.ProductID = p.ID
					}
				} else {
					util.Bail(fmt.Errorf(
						"ru_start %d entry does not have a product id, name, or alias",
						l.RUStart,
					))
				}

				if uuid.Equal(l.ProductID, uuid.UUID{}) {
					util.Bail(fmt.Errorf(
						"ru_start %d entry does not have a product id, name, or alias",
						l.RUStart,
					))
				}
			} else {
				_, ok := productsID[l.ProductID.String()]
				if !ok {
					util.Bail(errors.New("Product ID " + l.ProductID.String() + " is unknown"))
				}
			}
			s := conch.GlobalRackLayoutSlot{
				RackID:    GRackUUID,
				ProductID: l.ProductID,
				RUStart:   l.RUStart,
			}

			finalLayout = append(finalLayout, s)
		}

		// If the rack has a layout, and the user asked us to, nuke the
		// existing layout
		if *overwriteOpt {
			for _, s := range existingLayout {
				err := util.API.DeleteGlobalRackLayoutSlot(s.ID)
				if err != nil {
					util.Bail(err)
				}
			}
		}

		// Yes, technically, doing this in two loops is unnecessary. It gets us
		// a bit of sanity though. First loop verifies the data and the second
		// loop acts on it. That way, if the first loop runs into problems, we
		// haven't changed any data yet.
		for _, s := range finalLayout {
			err := util.API.SaveGlobalRackLayoutSlot(&s)
			if err != nil {
				util.Bail(err)
			}
		}

	}
}
