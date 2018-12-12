// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package workspaces

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	gotree "github.com/DiSiqueira/GoTree"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
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

func buildWSTree(parents map[string][]conch.Workspace, parent uuid.UUID, tree *gotree.GTStructure) {

	for _, ws := range parents[parent.String()] {
		sub := gotree.GTStructure{}
		sub.Name = fmt.Sprintf("%s / %s (%s)", ws.Name, ws.Role, ws.ID.String())

		buildWSTree(parents, ws.ID, &sub)
		tree.Items = append(tree.Items, sub)
	}
}

func getOne(app *cli.Cmd) {

	var (
		treeOutput = app.BoolOpt("tree", false, "Show workspace membership as a tree, based on subworkspace relationships. Specifying a workspace changes the root. Has no affect on --json")
	)

	app.Action = func() {
		if *treeOutput {
			wss, err := util.API.GetWorkspaces()
			if err != nil {
				util.Bail(err)
			}

			workspaces := make(map[string]conch.Workspace)
			for _, ws := range wss {
				workspaces[ws.ID.String()] = ws
			}

			parents := make(map[string][]conch.Workspace)

			for _, ws := range workspaces {
				if !uuid.Equal(ws.ParentID, uuid.UUID{}) {
					if _, ok := parents[ws.ParentID.String()]; !ok {
						parents[ws.ParentID.String()] = make([]conch.Workspace, 0)
					}
					parents[ws.ParentID.String()] = append(
						parents[ws.ParentID.String()],
						ws,
					)
				}
			}

			tree := gotree.GTStructure{}
			root := workspaces[WorkspaceUUID.String()]
			tree.Name = fmt.Sprintf("%s / %s (%s)", root.Name, root.Role, root.ID.String())

			buildWSTree(parents, WorkspaceUUID, &tree)
			gotree.PrintTree(tree)
			return
		}

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
		table.SetHeader([]string{"Name", "Email", "Role", "Role Via"})

		for _, u := range users {
			roleVia := ""
			if !uuid.Equal(u.RoleVia, WorkspaceUUID) && !uuid.Equal(u.RoleVia, uuid.UUID{}) {
				ws, err := util.API.GetWorkspace(u.RoleVia)
				if err != nil {
					util.Bail(err)
				}
				roleVia = ws.Name
			}
			table.Append([]string{u.Name, u.Email, u.Role, roleVia})
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
			"Datacenter",
			"Name",
			"Role",
			"Size",
		})

		for _, r := range racks {
			table.Append([]string{
				r.ID.String(),
				r.Datacenter,
				r.Name,
				r.Role,
				strconv.Itoa(r.Size),
			})
		}

		table.Render()
	}
}

type slotByRackUnitStart []conch.RackSlot

func (b slotByRackUnitStart) Len() int {
	return len(b)
}
func (b slotByRackUnitStart) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b slotByRackUnitStart) Less(i, j int) bool {
	return b[i].RackUnitStart > b[j].RackUnitStart
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

			sort.Sort(slotByRackUnitStart(rack.Slots))

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

			for _, slot := range rack.Slots {
				occupied := "X"

				occupantID := ""
				occupantHealth := ""

				if slot.Occupant.ID != "" {
					occupied = "+"
					occupantID = slot.Occupant.ID
					occupantHealth = slot.Occupant.Health
				}

				table.Append([]string{
					strconv.Itoa(slot.RackUnitStart),
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
		activeOnly   = app.BoolOpt("active-only", false, "Only retrieve active relays")
		activeWithin = app.IntOpt("active-within", 5, "If active-only is used, this specifies the number of minutes in which a relay must have reported to be considered active")
		fullOutput   = app.BoolOpt("full", false, "When global --json is used, provide full data about the devices rather than normal truncated data")
	)

	app.Action = func() {
		var relays []conch.WorkspaceRelay
		var err error

		if *activeOnly {
			relays, err = util.API.GetActiveWorkspaceRelays(
				WorkspaceUUID,
				*activeWithin,
			)
		} else {
			relays, err = util.API.GetWorkspaceRelays(WorkspaceUUID)
		}

		if err != nil {
			util.Bail(err)
		}

		if util.JSON && *fullOutput {
			util.JSONOut(relays)
			return
		}

		type resultRow struct {
			ID         string    `json:"id"`
			Alias      string    `json:"asset_tag"`
			Created    time.Time `json:"created"`
			IPAddr     string    `json:"ipaddr"`
			SSHPort    int       `json:"ssh_port"`
			Updated    time.Time `json:"updated"`
			Version    string    `json:"version"`
			NumDevices int       `json:"num_devices"`
		}

		results := make([]resultRow, 0)

		for _, r := range relays {
			results = append(results, resultRow{
				r.ID,
				r.Alias,
				r.Created.Time,
				r.IPAddr,
				r.SSHPort,
				r.Updated.Time,
				r.Version,
				r.NumDevices,
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
		devices, err := util.API.GetWorkspaceRelayDevices(
			WorkspaceUUID,
			RelayID,
		)
		if err != nil {
			util.Bail(err)
		}
		_ = util.DisplayDevices(devices, *fullOutput)
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
