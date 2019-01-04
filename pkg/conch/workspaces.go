// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// Workspace represents a Conch data partition which allows users to create
// custom lists of hardware
type Workspace struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Role        string    `json:"role"`
	ParentID    uuid.UUID `json:"parent_id,omitempty"`
}

// Room represents a physical area in a datacenter/AZ
type Room struct {
	ID         string `json:"id"`
	AZ         string `json:"az"`
	Alias      string `json:"alias"`
	VendorName string `json:"vendor_name"`
}

// WorkspaceAndRole ...
type WorkspaceAndRole struct {
	Workspace
	RoleVia uuid.UUID `json:"role_via"`
}

// WorkspaceUser ...
type WorkspaceUser struct {
	User
	RoleVia uuid.UUID `json:"role_via,omitempty"`
}

// GetWorkspaces returns the contents of /workspace, getting the list of all
// workspaces that the user has access to
func (c *Conch) GetWorkspaces() ([]Workspace, error) {
	workspaces := make([]Workspace, 0)
	return workspaces, c.get("/workspace", &workspaces)
}

// GetWorkspace returns the contents of /workspace/:uuid, getting information
// about a single workspace
// BUG(sungo): why is this returning a pointer
func (c *Conch) GetWorkspace(workspaceUUID fmt.Stringer) (*Workspace, error) {
	var workspace Workspace
	return &workspace, c.get("/workspace/"+workspaceUUID.String(), &workspace)
}

// GetSubWorkspaces returns the contents of /workspace/:uuid/child, getting
// a list of subworkspaces for the given workspace id
func (c *Conch) GetSubWorkspaces(workspaceUUID fmt.Stringer) ([]Workspace, error) {
	workspaces := make([]Workspace, 0)
	return workspaces, c.get(
		"/workspace/"+workspaceUUID.String()+"/child",
		&workspaces,
	)
}

// GetWorkspaceUsers returns the contents of /workspace/:uuid/users, getting
// a list of users for the given workspace id
func (c *Conch) GetWorkspaceUsers(workspaceUUID fmt.Stringer) ([]WorkspaceUser, error) {
	users := make([]WorkspaceUser, 0)
	return users, c.get(
		"/workspace/"+workspaceUUID.String()+"/user",
		&users,
	)
}

// GetWorkspaceRooms returns the contents of /workspace/:uuid/room, getting
// a list of rooms for the given workspace id
func (c *Conch) GetWorkspaceRooms(workspaceUUID fmt.Stringer) ([]Room, error) {
	rooms := make([]Room, 0)
	return rooms, c.get(
		"/workspace/"+workspaceUUID.String()+"/room",
		&rooms,
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
		"/workspace/"+parent.ID.String()+"/child",
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

	return c.post("/workspace/"+workspaceUUID.String()+"/rack", j, nil)
}

// DeleteRackFromWorkspace removes an existing rack from an existing workplace,
// via /workspace/:uuid/rack/:uuid
func (c *Conch) DeleteRackFromWorkspace(workspaceUUID fmt.Stringer, rackUUID fmt.Stringer) error {
	return c.httpDelete(
		"/workspace/" + workspaceUUID.String() + "/rack/" + rackUUID.String(),
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

	return c.post("/workspace/"+workspaceUUID.String()+"/user", body, nil)
}

// RemoveUserFromWorkspace ...
func (c *Conch) RemoveUserFromWorkspace(workspaceUUID fmt.Stringer, email string) error {
	return c.httpDelete("/workspace/" + workspaceUUID.String() + "/user/email=" + email)
}
