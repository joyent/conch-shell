// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package workspaces

import (
	"fmt"
	"github.com/joyent/conch-shell/pkg/util"
	"github.com/joyent/go-conch"
	"github.com/joyent/go-conch/pgtime"
	"gopkg.in/jawher/mow.cli.v1"
	"sort"
	"strconv"
)

func getAll(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin

	app.Action = func() {
		workspaces, err := util.API.GetWorkspaces()
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(workspaces)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{"Role", "Id", "Name", "Description"})

		for _, w := range workspaces {
			table.Append([]string{w.Role, w.ID.String(), w.Name, w.Description})
		}

		table.Render()
	}
}

func getOne(app *cli.Cmd) {
	app.Action = func() {
		workspace, err := util.API.GetWorkspace(WorkspaceUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(workspace)
			return
		}

		fmt.Printf(
			"Role: %s\nID: %s\nName: %s\nDescription: %s\n",
			workspace.Role,
			workspace.ID.String(),
			workspace.Name,
			workspace.Description,
		)
	}
}

func getUsers(app *cli.Cmd) {
	app.Action = func() {
		users, err := util.API.GetWorkspaceUsers(WorkspaceUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(users)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{"Name", "Email", "Role"})

		for _, u := range users {
			table.Append([]string{u.Name, u.Email, u.Role})
		}

		table.Render()
	}
}

func getDevices(app *cli.Cmd) {

	var (
		fullOutput = app.BoolOpt("full", false, "When --ids-only is *not* used, provide additional data about the devices rather than normal truncated data. Note: this slows things down immensely")
		idsOnly    = app.BoolOpt("ids-only", false, "Only retrieve device IDs")
		graduated  = app.StringOpt("graduated", "", "Filter by the 'graduated' field")
		health     = app.StringOpt("health", "", "Filter by the 'health' field")
	)

	app.Action = func() {
		devices, err := util.API.GetWorkspaceDevices(
			WorkspaceUUID,
			*idsOnly,
			*graduated,
			*health,
		)
		if err != nil {
			util.Bail(err)
		}

		if *idsOnly {
			ids := make([]string, 0)
			if util.JSON {
				for _, d := range devices {
					ids = append(ids, d.ID)
				}
				util.JSONOut(ids)
				return
			}
			for _, d := range devices {
				fmt.Println(d.ID)
			}
			return
		}

		if *fullOutput {
			filledIn := make([]conch.Device, 0)
			for _, d := range devices {
				fullDevice, err := util.API.FillInDevice(d)
				if err != nil {
					util.Bail(err)
				}
				filledIn = append(filledIn, fullDevice)
			}
			devices = filledIn
		}

		if err := util.DisplayDevices(devices, *fullOutput); err != nil {
			util.Bail(err)
		}
	}
}

func getRacks(app *cli.Cmd) {
	app.Action = func() {
		racks, err := util.API.GetWorkspaceRacks(WorkspaceUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(racks)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{
			"ID",
			"Name",
			"Role",
			"Unit",
			"Size",
		})

		for _, r := range racks {
			table.Append([]string{
				r.ID.String(),
				r.Name,
				r.Role,
				strconv.Itoa(r.Unit),
				strconv.Itoa(r.Size),
			})
		}

		table.Render()
	}
}

func getRack(app *cli.Cmd) {
	var (
		slotDetail = app.BoolOpt("slots", false, "Show details about each rack slot")
	)

	app.Action = func() {
		rack, err := util.API.GetWorkspaceRack(WorkspaceUUID, RackUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(rack)
			return
		}

		fmt.Printf(`
Workspace: %s
Rack ID:   %s
Name: %s
Role: %s
Datacenter: %s
`,
			WorkspaceUUID.String(),
			RackUUID.String(),
			rack.Name,
			rack.Role,
			rack.Datacenter,
		)

		if *slotDetail {
			fmt.Println()

			slotNums := make([]int, 0, len(rack.Slots))
			for k := range rack.Slots {
				slotNums = append(slotNums, k)
			}
			sort.Sort(sort.Reverse(sort.IntSlice(slotNums)))

			table := util.GetMarkdownTable()
			table.SetHeader([]string{
				"RU",
				"Occupied",
				"Name",
				"Alias",
				"Vendor",
				"Occupied By",
				"Health",
			})

			for _, ru := range slotNums {
				slot := rack.Slots[ru]
				occupied := "X"

				occupantID := ""
				occupantHealth := ""

				if slot.Occupant.ID != "" {
					occupied = "+"
					occupantID = slot.Occupant.ID
					occupantHealth = slot.Occupant.Health
				}

				table.Append([]string{
					strconv.Itoa(ru),
					occupied,
					slot.Name,
					slot.Alias,
					slot.Vendor,
					occupantID,
					occupantHealth,
				})

			}
			table.Render()
		}
	}
}

func getRelays(app *cli.Cmd) {
	var (
		activeOnly = app.BoolOpt("active-only", false, "Only retrieve active relays")
		fullOutput = app.BoolOpt("full", false, "When global --json is used, provide full data about the devices rather than normal truncated data")
	)

	app.Action = func() {
		relays, err := util.API.GetWorkspaceRelays(WorkspaceUUID, *activeOnly)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON && *fullOutput {
			util.JSONOut(relays)
			return
		}

		type resultRow struct {
			ID         string        `json:"id"`
			Alias      string        `json:"asset_tag"`
			Created    pgtime.PgTime `json:"created, int"`
			IPAddr     string        `json:"ipaddr"`
			SSHPort    int           `json:"ssh_port"`
			Updated    pgtime.PgTime `json:"updated"`
			Version    string        `json:"version"`
			NumDevices int           `json:"num_devices"`
		}

		results := make([]resultRow, 0)

		for _, r := range relays {
			numDevices := len(r.Devices)
			results = append(results, resultRow{
				r.ID,
				r.Alias,
				r.Created,
				r.IPAddr,
				r.SSHPort,
				r.Updated,
				r.Version,
				numDevices,
			})
		}

		if util.JSON {
			util.JSONOut(results)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{
			"ID",
			"Alias",
			"Created",
			"IP Addr",
			"SSH Port",
			"Updated",
			"Version",
			"Number of Devices",
		})

		for _, r := range results {
			updated := ""
			if !r.Updated.IsZero() {
				updated = util.TimeStr(r.Updated)
			}

			table.Append([]string{
				r.ID,
				r.Alias,
				util.TimeStr(r.Created),
				r.IPAddr,
				strconv.Itoa(r.SSHPort),
				updated,
				r.Version,
				strconv.Itoa(r.NumDevices),
			})
		}

		table.Render()
	}
}

func getRooms(app *cli.Cmd) {
	app.Action = func() {
		rooms, err := util.API.GetWorkspaceRooms(WorkspaceUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(rooms)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{"ID", "AZ", "Alias", "Vendor Name"})

		for _, r := range rooms {
			table.Append([]string{r.ID, r.AZ, r.Alias, r.VendorName})
		}

		table.Render()
	}

}

func getSubs(app *cli.Cmd) {
	app.Action = func() {
		workspaces, err := util.API.GetSubWorkspaces(WorkspaceUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(workspaces)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{"Role", "Id", "Name", "Description"})

		for _, w := range workspaces {
			table.Append([]string{w.Role, w.ID.String(), w.Name, w.Description})
		}
		table.Render()
	}
}

func getRelayDevices(app *cli.Cmd) {
	var (
		fullOutput = app.BoolOpt("full", false, "When global --json is used, provide full data about the devices rather than normal truncated data")
	)

	app.Action = func() {
		relays, err := util.API.GetWorkspaceRelays(WorkspaceUUID, false)
		if err != nil {
			util.Bail(err)
		}
		var relay conch.Relay
		foundRelay := false
		for _, r := range relays {
			if r.ID == RelayID {
				relay = r
				foundRelay = true
			}
		}
		if !foundRelay {
			util.Bail(conch.ErrDataNotFound)
		}

		_ = util.DisplayDevices(relay.Devices, *fullOutput)

	}
}

func addRack(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.AddRackToWorkspace(WorkspaceUUID, RackUUID); err != nil {
			util.Bail(err)
		}
	}
}

func deleteRack(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.DeleteRackFromWorkspace(WorkspaceUUID, RackUUID); err != nil {
			util.Bail(err)
		}
	}
}
