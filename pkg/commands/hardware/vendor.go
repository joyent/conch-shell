// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hardware

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
)

func displayHardwareVendor(v conch.HardwareVendor) {
	fmt.Printf(`ID: %s
Name: %s
Created: %s
Updated: %s
`,
		v.ID.String(),
		v.Name,
		util.TimeStr(v.Created),
		util.TimeStr(v.Updated),
	)
}

func getOneVendor(app *cli.Cmd) {
	app.Action = func() {
		ret, err := util.API.GetHardwareVendor(HardwareVendorName)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(ret)
			return
		}

		if err != nil {
			util.Bail(err)
		}
		displayHardwareVendor(ret)
	}
}

func getAllVendors(app *cli.Cmd) {
	app.Action = func() {
		ret, err := util.API.GetHardwareVendors()
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(ret)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{
			"ID",
			"Name",
			"Created",
			"Updated",
		})
		for _, v := range ret {
			table.Append([]string{
				v.ID.String(),
				v.Name,
				util.TimeStr(v.Created),
				util.TimeStr(v.Updated),
			})
		}
		table.Render()
	}
}

func createOneVendor(app *cli.Cmd) {
	app.Action = func() {
		v := conch.HardwareVendor{Name: HardwareVendorName}

		err := util.API.SaveHardwareVendor(&v)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(v)
			return
		}
		displayHardwareVendor(v)
	}

}

func deleteOneVendor(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.DeleteHardwareVendor(HardwareVendorName); err != nil {
			util.Bail(err)
		}
	}

}
