// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package workspaces

import (
	"fmt"
	"github.com/joyent/conch-shell/util"
	"github.com/joyent/go-conch"
	pgtime "github.com/joyent/go-conch/pg_time"
	"gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
	"sort"
	"strconv"
	"time"
)

func getAll(app *cli.Cmd) {
	app.Before = util.BuildApiAndVerifyLogin

	app.Action = func() {
		workspaces, err := util.API.GetWorkspaces()
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JsonOut(workspaces)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{"Role", "Id", "Name", "Description"})

		for _, w := range workspaces {
			table.Append([]string{w.Role, w.Id.String(), w.Name, w.Description})
		}

		table.Render()
	}
}

func getOne(app *cli.Cmd) {
	app.Action = func() {
		workspace, err := util.API.GetWorkspace(WorkspaceUuid)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JsonOut(workspace)
			return
		}

		fmt.Printf(
			"Role: %s\nID: %s\nName: %s\nDescription: %s\n",
			workspace.Role,
			workspace.Id.String(),
			workspace.Name,
			workspace.Description,
		)
	}
}

func getUsers(app *cli.Cmd) {
	app.Action = func() {
		users, err := util.API.GetWorkspaceUsers(WorkspaceUuid)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JsonOut(users)
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
		full_output = app.BoolOpt("full", false, "When --ids-only is *not* used, provide additional data about the devices rather than normal truncated data. Note: this slows things down immensely")
		ids_only    = app.BoolOpt("ids-only", false, "Only retrieve device IDs")
		graduated   = app.StringOpt("graduated", "", "Filter by the 'graduated' field")
		health      = app.StringOpt("health", "", "Filter by the 'health' field")
	)

	app.Action = func() {
		devices, err := util.API.GetWorkspaceDevices(
			WorkspaceUuid,
			*ids_only,
			*graduated,
			*health,
		)
		if err != nil {
			util.Bail(err)
		}

		if *ids_only {
			ids := make([]string, 0)
			if util.JSON {
				for _, d := range devices {
					ids = append(ids, d.Id)
				}
				util.JsonOut(ids)
				return
			} else {
				for _, d := range devices {
					fmt.Println(d.Id)
				}
				return
			}
		}

		if *full_output {
			filled_in := make([]conch.ConchDevice, 0)
			for _, d := range devices {
				full_d, err := util.API.FillInDevice(d)
				if err != nil {
					util.Bail(err)
				}
				filled_in = append(filled_in, full_d)
			}
			devices = filled_in
		}

		if err := util.DisplayDevices(devices, *full_output); err != nil {
			util.Bail(err)
		}
	}
}

func getRacks(app *cli.Cmd) {
	app.Action = func() {
		racks, err := util.API.GetWorkspaceRacks(WorkspaceUuid)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JsonOut(racks)
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
				r.Id.String(),
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
		rack_id     = app.StringArg("RACK", "", "Rack UUID")
		slot_detail = app.BoolOpt("slots", false, "Show details about each rack slot")
	)

	app.Spec = "RACK [OPTIONS]"
	app.Action = func() {
		rack_uuid, err := uuid.FromString(*rack_id)
		if err != nil {
			util.Bail(err)
		}

		rack, err := util.API.GetWorkspaceRack(WorkspaceUuid, rack_uuid)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JsonOut(rack)
			return
		}

		fmt.Printf(`
Workspace: %s
Rack ID:   %s
Name: %s
Role: %s
Datacenter: %s
`,
			WorkspaceUuid.String(),
			rack_uuid.String(),
			rack.Name,
			rack.Role,
			rack.Datacenter,
		)

		if *slot_detail {
			fmt.Println()

			slot_nums := make([]int, 0, len(rack.Slots))
			for k := range rack.Slots {
				slot_nums = append(slot_nums, k)
			}
			sort.Sort(sort.Reverse(sort.IntSlice(slot_nums)))

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

			for _, ru := range slot_nums {
				slot := rack.Slots[ru]
				occupied := "X"

				occupant_id := ""
				occupant_health := ""

				if slot.Occupant.Id != "" {
					occupied = "+"
					occupant_id = slot.Occupant.Id
					occupant_health = slot.Occupant.Health
				}

				table.Append([]string{
					strconv.Itoa(ru),
					occupied,
					slot.Name,
					slot.Alias,
					slot.Vendor,
					occupant_id,
					occupant_health,
				})

			}
			table.Render()
		}
	}
}

func getRelays(app *cli.Cmd) {
	var (
		active_only = app.BoolOpt("active-only", false, "Only retrieve active relays")
		full_output = app.BoolOpt("full", false, "When global --json is used, provide full data about the devices rather than normal truncated data")
	)

	app.Action = func() {
		relays, err := util.API.GetWorkspaceRelays(WorkspaceUuid, *active_only)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON && *full_output {
			util.JsonOut(relays)
			return
		}

		type resultRow struct {
			Id         string             `json:"id"`
			Alias      string             `json:"asset_tag"`
			Created    pgtime.ConchPgTime `json:"created, int"`
			IpAddr     string             `json:"ipaddr"`
			SshPort    int                `json:"ssh_port"`
			Updated    pgtime.ConchPgTime `json:"updated"`
			Version    string             `json:"version"`
			NumDevices int                `json:"num_devices"`
		}

		results := make([]resultRow, 0)

		for _, r := range relays {
			num_devices := len(r.Devices)
			results = append(results, resultRow{
				r.Id,
				r.Alias,
				r.Created,
				r.IpAddr,
				r.SshPort,
				r.Updated,
				r.Version,
				num_devices,
			})
		}

		if util.JSON {
			util.JsonOut(results)
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
				updated = r.Updated.Format(time.UnixDate)
			}

			table.Append([]string{
				r.Id,
				r.Alias,
				r.Created.Format(time.UnixDate),
				r.IpAddr,
				strconv.Itoa(r.SshPort),
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
		rooms, err := util.API.GetWorkspaceRooms(WorkspaceUuid)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JsonOut(rooms)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{"ID", "AZ", "Alias", "Vendor Name"})

		for _, r := range rooms {
			table.Append([]string{r.Id, r.Az, r.Alias, r.VendorName})
		}

		table.Render()
	}

}

func getSubs(app *cli.Cmd) {
	app.Action = func() {
		workspaces, err := util.API.GetSubWorkspaces(WorkspaceUuid)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON == true {
			util.JsonOut(workspaces)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{"Role", "Id", "Name", "Description"})

		for _, w := range workspaces {
			table.Append([]string{w.Role, w.Id.String(), w.Name, w.Description})
		}
		table.Render()
	}
}

func getRelayDevices(app *cli.Cmd) {
	var (
		full_output = app.BoolOpt("full", false, "When global --json is used, provide full data about the devices rather than normal truncated data")
	)

	app.Action = func() {
		relays, err := util.API.GetWorkspaceRelays(WorkspaceUuid, false)
		if err != nil {
			util.Bail(err)
		}
		var relay conch.ConchRelay
		found_relay := false
		for _, r := range relays {
			if r.Id == RelayId {
				relay = r
				found_relay = true
			}
		}
		if found_relay == false {
			util.Bail(conch.ConchDataNotFound)
		}

		util.DisplayDevices(relay.Devices, *full_output)

	}
}
