// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package devices

import (
	"fmt"
	"github.com/joyent/conch-shell/pkg/util"
	"github.com/joyent/go-conch"
	"gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
	"strconv"
)

func getAllDeviceRoles(app *cli.Cmd) {
	app.Action = func() {
		roles, err := util.API.GetDeviceRoles()
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(roles)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{
			"Id",
			"Description",
			"Hardware Product ID",
			"Created",
			"Updated",
			"Services Count",
		})

		for _, r := range roles {
			table.Append([]string{
				r.ID.String(),
				r.Description,
				r.HardwareProductID.String(),
				r.Created.String(),
				r.Updated.String(),
				strconv.Itoa(len(r.Services)),
			})
		}
		table.Render()
	}
}

func getOneDeviceRole(app *cli.Cmd) {
	app.Action = func() {
		r, err := util.API.GetDeviceRole(DeviceRoleUUID)
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(r)
			return
		}
		fmt.Println("ID: " + r.ID.String())
		fmt.Println("Created: " + r.Created.String())
		fmt.Println("Updated: " + r.Updated.String())

		fmt.Println("\nDescription: " + r.Description)
		fmt.Println("Hardware Product ID: " + r.HardwareProductID.String())
		fmt.Println("Services:")
		for _, sID := range r.Services {
			s, err := util.API.GetDeviceService(sID)
			if err != nil {
				fmt.Printf("  - %s (Error occured in lookup: %s)\n", sID, err)
			} else {
				fmt.Printf("  - %s [ %s ]\n", s.ID, s.Name)
			}
		}
	}
}

func createDeviceRole(app *cli.Cmd) {
	var descrArg = app.StringOpt("description d", "", "Description of the role")
	var hwProductArg = app.StringOpt("hardware-product h", "", "Hardware Product UUID")

	app.Spec = "--description --hardware-product"

	app.Action = func() {
		id, err := uuid.FromString(*hwProductArg)
		if err != nil {
			util.Bail(err)
		}

		role := &conch.DeviceRole{
			Description:       *descrArg,
			HardwareProductID: id,
		}
		if err := util.API.SaveDeviceRole(role); err != nil {
			util.Bail(err)
		}
		fmt.Println("New Device Role ID: " + role.ID.String())
	}
}

func deleteDeviceRole(app *cli.Cmd) {
	app.Action = func() {
		err := util.API.DeleteDeviceRole(DeviceRoleUUID)
		if err != nil {
			util.Bail(err)
		}
	}
}

func modifyDeviceRole(app *cli.Cmd) {
	var descrArg = app.StringOpt("description d", "", "Description of the role")
	var hwProductArg = app.StringOpt("hardware-product h", "", "Hardware Product UUID")

	app.Action = func() {
		r, err := util.API.GetDeviceRole(DeviceRoleUUID)
		if err != nil {
			util.Bail(err)
		}

		if *descrArg != "" {
			r.Description = *descrArg
		}

		if *hwProductArg != "" {
			id, err := uuid.FromString(*hwProductArg)
			if err != nil {
				util.Bail(err)
			}
			r.HardwareProductID = id
		}

		err = util.API.SaveDeviceRole(&r)

		if err != nil {
			util.Bail(err)
		}
	}
}

func addServiceToDeviceRole(app *cli.Cmd) {
	var serviceIDArg = app.StringArg("SERVICE", "", "UUID of the service")

	app.Spec = "SERVICE"

	app.Action = func() {
		r, err := util.API.GetDeviceRole(DeviceRoleUUID)
		if err != nil {
			util.Bail(err)
		}

		id, err := uuid.FromString(*serviceIDArg)
		if err != nil {
			util.Bail(err)
		}

		s, err := util.API.GetDeviceService(id)
		if err != nil {
			util.Bail(err)
		}

		err = util.API.AddServiceToDeviceRole(r, s)
		if err != nil {
			util.Bail(err)
		}

	}

}
func removeServiceFromDeviceRole(app *cli.Cmd) {
	var serviceIDArg = app.StringArg("SERVICE", "", "UUID of the service")

	app.Spec = "SERVICE"

	app.Action = func() {
		r, err := util.API.GetDeviceRole(DeviceRoleUUID)
		if err != nil {
			util.Bail(err)
		}

		id, err := uuid.FromString(*serviceIDArg)
		if err != nil {
			util.Bail(err)
		}

		s, err := util.API.GetDeviceService(id)
		if err != nil {
			util.Bail(err)
		}

		err = util.API.RemoveServiceFromDeviceRole(r, s)
		if err != nil {
			util.Bail(err)
		}

	}

}
