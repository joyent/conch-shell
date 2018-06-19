// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package devices

import (
	"fmt"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
)

func getAllDeviceServices(app *cli.Cmd) {
	app.Action = func() {
		services, err := util.API.GetDeviceServices()
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(services)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{"Id", "Name", "Created", "Updated"})

		for _, s := range services {
			table.Append([]string{
				s.ID.String(),
				s.Name,
				s.Created.String(),
				s.Updated.String(),
			})
		}
		table.Render()
	}
}

func getOneDeviceService(app *cli.Cmd) {
	app.Action = func() {
		s, err := util.API.GetDeviceService(DeviceServiceUUID)
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(s)
			return
		}
		fmt.Println("ID: " + s.ID.String())
		fmt.Println("Name: " + s.Name)
		fmt.Println("Created: " + s.Created.String())
		fmt.Println("Updated: " + s.Updated.String())
	}
}

func createDeviceService(app *cli.Cmd) {
	var nameArg = app.StringOpt("name n", "", "Name of the service")

	app.Spec = "--name"

	app.Action = func() {
		service := &conch.DeviceService{Name: *nameArg}
		if err := util.API.SaveDeviceService(service); err != nil {
			util.Bail(err)
		}
		fmt.Println("New Device Service ID: " + service.ID.String())
	}
}

func deleteDeviceService(app *cli.Cmd) {
	app.Action = func() {
		err := util.API.DeleteDeviceService(DeviceServiceUUID)
		if err != nil {
			util.Bail(err)
		}
	}
}

func modifyDeviceService(app *cli.Cmd) {
	var nameArg = app.StringOpt("name n", "", "Name of the service")
	app.Spec = "--name"
	app.Action = func() {
		s, err := util.API.GetDeviceService(DeviceServiceUUID)
		if err != nil {
			util.Bail(err)
		}

		s.Name = *nameArg
		err = util.API.SaveDeviceService(&s)

		if err != nil {
			util.Bail(err)
		}
	}
}
