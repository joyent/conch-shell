// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package workspaces

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"

	gotree "github.com/DiSiqueira/GoTree"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/joyent/conch-shell/pkg/util"
)

type rackAssignedSlot struct {
	RackUnitStart       int    `json:"ru_start"`
	HardwareProductName string `json:"hardware_product_name"`
	DeviceID            string `json:"device_id"`
}

type rackAssignments []rackAssignedSlot

func (r rackAssignments) Len() int {
	return len(r)
}
func (r rackAssignments) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r rackAssignments) Less(i, j int) bool {
	return r[i].RackUnitStart > r[j].RackUnitStart
}

/******/

func getAll(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin

	app.Action = func() {
		workspaces, err := util.API.GetWorkspaces()
		if err != nil {
			util.Bail(err)
		}
		sort.Sort(workspaces)
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
		validated  = app.StringOpt("validated", "", "Filter by the 'validated' field")
	)

	app.Action = func() {
		devices, err := util.API.GetWorkspaceDevices(
			WorkspaceUUID,
			*idsOnly,
			*graduated,
			*health,
			*validated,
		)
		if err != nil {
			util.Bail(err)
		}

		sort.Sort(devices)

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
			locs := make(map[uuid.UUID]conch.DeviceLocation)

			dLocs := make([]conch.Device, 0)

			for _, d := range devices {
				if uuid.Equal(d.RackID, uuid.UUID{}) {
					continue
				}
				if loc, ok := locs[d.RackID]; ok {
					d.Location = loc
				} else {
					if loc, err := util.API.GetDeviceLocation(d.ID); err == nil {
						loc.TargetHardwareProduct = conch.HardwareProductTarget{}
						locs[loc.Rack.ID] = loc
						d.Location = loc
					}
				}

				dLocs = append(dLocs, d)
			}
			devices = dLocs
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

func getRack(app *cli.Cmd) {
	app.LongDesc = "The validation status in this command does *not* correspond to the 'validated' properly of a device. Rather, the app retrieves the real validation status."

	app.Action = func() {
		rack, err := util.API.GetWorkspaceRack(WorkspaceUUID, RackUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(rack)
			return
		}

		workspace, err := util.API.GetWorkspace(WorkspaceUUID)
		if err != nil {
			util.Bail(err)
		}

		fmt.Printf(`
Workspace:  %s
Datacenter: %s

Name: %s
Role: %s
Rack ID: %s
`,
			workspace.Name,
			rack.Datacenter,
			rack.Name,
			rack.Role,
			rack.ID,
		)

		fmt.Println()

		sort.Sort(rack.Slots)

		table := util.GetMarkdownTable()
		table.SetHeader([]string{
			"RU",
			"Occupied",
			"Validated",
			"Name",
			"Alias",
			"Vendor",
			"Occupied By",
			"Health",
		})

		for _, slot := range rack.Slots {
			occupied := "X"
			validated := "?"

			occupantID := ""
			occupantHealth := ""

			if slot.Occupant.ID != "" {
				occupied = "+"
				occupantID = slot.Occupant.ID
				occupantHealth = slot.Occupant.Health

				vstates, err := util.API.DeviceValidationStates(slot.Occupant.ID)
				if err != nil {
					util.Bail(err)
				}

				if len(vstates) > 0 {
					validated = "+"
					for _, vstate := range vstates {
						if vstate.Status != "pass" {
							validated = "X"
						}
					}
				}
			}

			table.Append([]string{
				strconv.Itoa(slot.RackUnitStart),
				occupied,
				validated,
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

func getSubs(app *cli.Cmd) {
	app.Action = func() {
		workspaces, err := util.API.GetSubWorkspaces(WorkspaceUUID)
		if err != nil {
			util.Bail(err)
		}
		sort.Sort(workspaces)

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

		if err := util.DisplayDevices(devices, *fullOutput); err != nil {
			util.Bail(err)
		}
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

func assignRack(app *cli.Cmd) {
	var (
		filePathArg = app.StringArg("FILE", "-", "Path to a JSON file to use as the data source. '-' indicates STDIN")
	)
	app.Spec = "FILE"
	app.Action = func() {
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
		if len(string(b)) <= 1 {
			util.Bail(errors.New("no data provided"))
		}

		a := make(rackAssignments, 0)
		if err := json.Unmarshal(b, &a); err != nil {
			util.Bail(err)
		}
		fmt.Println(a)

		keepers := make(rackAssignments, 0)
		for _, slot := range a {
			if slot.DeviceID != "" {
				keepers = append(keepers, slot)
			}
		}

		if len(keepers) == 0 {
			util.Bail(errors.New("no devices found. no changes to make"))
		}

		assignments := make(conch.WorkspaceRackLayoutAssignments)

		for _, keeper := range keepers {
			assignments[keeper.DeviceID] = keeper.RackUnitStart
		}

		if err := util.API.AssignDevicesToRackSlots(
			WorkspaceUUID,
			RackUUID,
			assignments,
		); err != nil {
			util.Bail(err)
		}
	}
}

func assignmentsRack(app *cli.Cmd) {
	app.Action = func() {
		rack, err := util.API.GetWorkspaceRack(WorkspaceUUID, RackUUID)
		if err != nil {
			util.Bail(err)
		}
		if (len(rack.Slots) == 1) && (rack.Slots[0].RackUnitStart == 0) {
			util.Bail(errors.New("rack has no layout"))
		}

		a := make(rackAssignments, 0)
		for _, slot := range rack.Slots {
			if slot.Occupant.ID == "" {
				a = append(a, rackAssignedSlot{
					RackUnitStart:       slot.RackUnitStart,
					HardwareProductName: slot.Name,
				})
				continue
			}

			a = append(a, rackAssignedSlot{
				RackUnitStart:       slot.RackUnitStart,
				HardwareProductName: slot.Name,
				DeviceID:            slot.Occupant.ID,
			})
		}

		sort.Sort(a)

		util.JSONOutIndent(a)
	}
}
