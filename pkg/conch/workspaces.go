// Copyright 2017 Joyent, Inc.
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
	Description string    `json:"description"`
	Role        string    `json:"role"`
	ParentID    uuid.UUID `json:"parent_id"`
}

// Room represents a physical area in a datacenter/AZ
type Room struct {
	ID         string `json:"id"`
	AZ         string `json:"az"`
	Alias      string `json:"alias"`
	VendorName string `json:"vendor_name"`
}

// GetWorkspaces returns the contents of /workspace, getting the list of all
// workspaces that the user has access to
func (c *Conch) GetWorkspaces() ([]Workspace, error) {
	workspaces := make([]Workspace, 0)

	aerr := &APIError{}
	res, err := c.sling().New().Get("/workspace").Receive(&workspaces, aerr)

	return workspaces, c.isHTTPResOk(res, err, aerr)
}

// GetWorkspace returns the contents of /workspace/:uuid, getting information
// about a single workspace
func (c *Conch) GetWorkspace(workspaceUUID fmt.Stringer) (*Workspace, error) {
	workspace := &Workspace{}

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/workspace/"+workspaceUUID.String()).
		Receive(&workspace, aerr)

	return workspace, c.isHTTPResOk(res, err, aerr)
}

// GetSubWorkspaces returns the contents of /workspace/:uuid/child, getting
// a list of subworkspaces for the given workspace id
func (c *Conch) GetSubWorkspaces(workspaceUUID fmt.Stringer) ([]Workspace, error) {
	workspaces := make([]Workspace, 0)

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/workspace/"+workspaceUUID.String()+"/child").
		Receive(&workspaces, aerr)

	return workspaces, c.isHTTPResOk(res, err, aerr)
}

// GetWorkspaceUsers returns the contents of /workspace/:uuid/users, getting
// a list of users for the given workspace id
func (c *Conch) GetWorkspaceUsers(workspaceUUID fmt.Stringer) ([]User, error) {
	users := make([]User, 0)

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/workspace/"+workspaceUUID.String()+"/user").
		Receive(&users, aerr)

	return users, c.isHTTPResOk(res, err, aerr)
}

// GetWorkspaceRooms returns the contents of /workspace/:uuid/room, getting
// a list of rooms for the given workspace id
func (c *Conch) GetWorkspaceRooms(workspaceUUID fmt.Stringer) ([]Room, error) {
	rooms := make([]Room, 0)

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/workspace/"+workspaceUUID.String()+"/room").
		Receive(&rooms, aerr)

	return rooms, c.isHTTPResOk(res, err, aerr)
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

	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/workspace/"+parent.ID.String()+"/child").
		BodyJSON(j).
		Receive(&sub, aerr)

	return sub, c.isHTTPResOk(res, err, aerr)
}

// AddRackToWorkspace adds an existing rack to an existing workspace, via
// /workspace/:uuid/rack
func (c *Conch) AddRackToWorkspace(workspaceUUID fmt.Stringer, rackUUID fmt.Stringer) error {
	j := struct {
		ID string `json:"id"`
	}{
		rackUUID.String(),
	}

	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/workspace/"+workspaceUUID.String()+"/rack").
		BodyJSON(j).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// DeleteRackFromWorkspace removes an existing rack from an existing workplace,
// via /workspace/:uuid/rack/:uuid
func (c *Conch) DeleteRackFromWorkspace(workspaceUUID fmt.Stringer, rackUUID fmt.Stringer) error {
	aerr := &APIError{}
	res, err := c.sling().New().
		Delete("/workspace/"+workspaceUUID.String()+"/rack/"+rackUUID.String()).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}
