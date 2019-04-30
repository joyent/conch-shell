// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	"net/url"

	uuid "gopkg.in/satori/go.uuid.v1"
)

// GetWorkspaceRacks fetchest the list of racks for a workspace, via
// /workspace/:uuid/rack
// NOTE: The API currently returns a hash of arrays where the key is the
// datacenter/az. This routine copies that key into the Datacenter field in the
// Rack struct.
func (c *Conch) GetWorkspaceRacks(workspaceUUID fmt.Stringer) ([]Rack, error) {
	racks := make([]Rack, 0)
	j := make(map[string][]Rack)

	if err := c.get("/workspace/"+url.PathEscape(workspaceUUID.String())+"/rack", &j); err != nil {
		return racks, err
	}

	for az, loc := range j {
		for _, rack := range loc {
			rack.Datacenter = az
			racks = append(racks, rack)
		}
	}

	return racks, nil
}

// GetWorkspaceRack fetches a single rack for a workspace, via
// /workspace/:uuid/rack/:id
func (c *Conch) GetWorkspaceRack(
	workspaceUUID fmt.Stringer,
	rackUUID fmt.Stringer,
) (rack Rack, err error) {
	return rack, c.get(
		"/workspace/"+
			url.PathEscape(workspaceUUID.String())+
			"/rack/"+
			url.PathEscape(rackUUID.String()),
		&rack,
	)
}

// GetWorkspaceDevices retrieves a list of Devices for the given
// workspace.
// Pass true for 'IDsOnly' to get Devices with only the ID field populated
// Pass a string for 'graduated' to filter devices by graduated value, as per https://conch.joyent.us/doc#getdevices
// Pass a string for 'health' to filter devices by health value, as per https://conch.joyent.us/doc#getdevices
func (c *Conch) GetWorkspaceDevices(
	workspaceUUID fmt.Stringer,
	idsOnly bool,
	graduated string,
	health string,
	validated string,
) (Devices, error) {

	devices := make([]Device, 0)

	opts := struct {
		IDsOnly   bool   `url:"ids_only,omitempty"`
		Graduated string `url:"graduated,omitempty"`
		Health    string `url:"health,omitempty"`
		Validated string `url:"validated,omitempty"`
	}{
		idsOnly,
		graduated,
		health,
		validated,
	}

	url := "/workspace/" + url.PathEscape(workspaceUUID.String()) + "/device"
	if idsOnly {
		ids := make([]string, 0)

		if err := c.getWithQuery(url, opts, &ids); err != nil {
			return devices, err
		}

		for _, v := range ids {
			device := Device{ID: v}
			devices = append(devices, device)
		}
		return devices, nil
	}
	return devices, c.getWithQuery(url, opts, &devices)
}

// GetWorkspaces returns the contents of /workspace, getting the list of all
// workspaces that the user has access to
func (c *Conch) GetWorkspaces() (Workspaces, error) {
	workspaces := make([]Workspace, 0)
	return workspaces, c.get("/workspace", &workspaces)
}

// GetWorkspace returns the contents of /workspace/:uuid, getting information
// about a single workspace
func (c *Conch) GetWorkspace(workspaceUUID fmt.Stringer) (w Workspace, e error) {
	return w, c.get("/workspace/"+url.PathEscape(workspaceUUID.String()), &w)
}

// GetWorkspaceByName returns the contents of /workspace/:name, getting
// information about a single workspace
func (c *Conch) GetWorkspaceByName(name string) (w Workspace, e error) {

	return w, c.get("/workspace/"+url.PathEscape(name), &w)
}

// GetSubWorkspaces returns the contents of /workspace/:uuid/child, getting
// a list of subworkspaces for the given workspace id
func (c *Conch) GetSubWorkspaces(workspaceUUID fmt.Stringer) (Workspaces, error) {
	workspaces := make(Workspaces, 0)
	return workspaces, c.get(
		"/workspace/"+url.PathEscape(workspaceUUID.String())+"/child",
		&workspaces,
	)
}

// GetWorkspaceUsers returns the contents of /workspace/:uuid/users, getting
// a list of users for the given workspace id
func (c *Conch) GetWorkspaceUsers(workspaceUUID fmt.Stringer) ([]WorkspaceUser, error) {
	users := make([]WorkspaceUser, 0)
	return users, c.get(
		"/workspace/"+url.PathEscape(workspaceUUID.String())+"/user",
		&users,
	)
}

// CreateSubWorkspace creates a sub workspace under the parent, via
// /workspace/:uuid/child
// If the provided parent lacks an ID, ErrBadInput is returned
// Currently, if an attempt to create a workspace with a conflicting name
// happens, the API returns a 500 rather than something useful. The routine
// will return ErrHTTPNotOk in that case.
func (c *Conch) CreateSubWorkspace(parent Workspace, sub Workspace) (Workspace, error) {
	if uuid.Equal(parent.ID, uuid.UUID{}) {
		return sub, ErrBadInput
	}
	j := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{
		sub.Name,
		sub.Description,
	}

	return sub, c.post(
		"/workspace/"+url.PathEscape(parent.ID.String())+"/child",
		j,
		&sub,
	)
}

// AddRackToWorkspace adds an existing rack to an existing workspace, via
// /workspace/:uuid/rack
func (c *Conch) AddRackToWorkspace(workspaceUUID fmt.Stringer, rackUUID fmt.Stringer) error {
	j := struct {
		ID string `json:"id"`
	}{
		rackUUID.String(),
	}

	return c.post("/workspace/"+url.PathEscape(workspaceUUID.String())+"/rack", j, nil)
}

// DeleteRackFromWorkspace removes an existing rack from an existing workplace,
// via /workspace/:uuid/rack/:uuid
func (c *Conch) DeleteRackFromWorkspace(workspaceUUID fmt.Stringer, rackUUID fmt.Stringer) error {
	return c.httpDelete(
		"/workspace/" +
			url.PathEscape(workspaceUUID.String()) +
			"/rack/" +
			url.PathEscape(rackUUID.String()),
	)
}

// AddUserToWorkspace adds a user to a workspace via /workspace/:uuid/user
func (c *Conch) AddUserToWorkspace(workspaceUUID fmt.Stringer, user string, role string) error {
	body := struct {
		User string `json:"user"`
		Role string `json:"role"`
	}{
		user,
		role,
	}

	return c.post("/workspace/"+url.PathEscape(workspaceUUID.String())+"/user", body, nil)
}

// RemoveUserFromWorkspace ...
func (c *Conch) RemoveUserFromWorkspace(workspaceUUID fmt.Stringer, email string) error {
	return c.httpDelete("/workspace/" +
		url.PathEscape(workspaceUUID.String()) +
		"/user/email=" +
		url.PathEscape(email),
	)
}

func (c *Conch) AssignDevicesToRackSlots(
	workspaceID fmt.Stringer,
	rackID fmt.Stringer,
	assignments WorkspaceRackLayoutAssignments,
) error {
	return c.post(
		"/workspace/"+
			url.PathEscape(workspaceID.String())+
			"/rack/"+
			url.PathEscape(rackID.String())+
			"/layout",
		assignments,
		nil,
	)
}
