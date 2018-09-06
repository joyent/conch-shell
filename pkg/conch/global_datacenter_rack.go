// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	uuid "gopkg.in/satori/go.uuid.v1"
	"net/http"
)

// GetGlobalRacks fetches a list of all racks in the global domain
func (c *Conch) GetGlobalRacks() ([]GlobalRack, error) {
	r := make([]GlobalRack, 0)

	aerr := &APIError{}
	res, err := c.sling().New().Get("/rack").Receive(&r, aerr)

	return r, c.isHTTPResOk(res, err, aerr)
}

// GetGlobalRack fetches a single rack in the global domain, by its
// UUID
func (c *Conch) GetGlobalRack(id fmt.Stringer) (GlobalRack, error) {
	r := GlobalRack{}

	aerr := &APIError{}
	res, err := c.sling().New().Get("/rack/"+id.String()).Receive(&r, aerr)

	return r, c.isHTTPResOk(res, err, aerr)
}

// SaveGlobalRack creates or updates a rack in the global domain,
// based on the presence of an ID
func (c *Conch) SaveGlobalRack(r *GlobalRack) error {
	if uuid.Equal(r.DatacenterRoomID, uuid.UUID{}) {
		return ErrBadInput
	}
	if uuid.Equal(r.RoleID, uuid.UUID{}) {
		return ErrBadInput
	}
	if r.Name == "" {
		return ErrBadInput
	}

	var err error
	var res *http.Response
	aerr := &APIError{}

	if uuid.Equal(r.ID, uuid.UUID{}) {

		j := struct {
			DatacenterRoomID string `json:"datacenter_room_id"`
			Name             string `json:"name"`
			RoleID           string `json:"role"`
			SerialNumber     string `json:"serial_number,omitempty"`
			AssetTag         string `json:"asset_tag,omitempty"`
		}{
			r.DatacenterRoomID.String(),
			r.Name,
			r.RoleID.String(),
			r.SerialNumber,
			r.AssetTag,
		}

		res, err = c.sling().New().Post("/rack").BodyJSON(j).Receive(&r, aerr)
	} else {
		j := struct {
			DatacenterRoomID string `json:"datacenter_room_id"`
			Name             string `json:"name"`
			RoleID           string `json:"role"`
			SerialNumber     string `json:"serial_number,omitempty"`
			AssetTag         string `json:"asset_tag,omitempty"`
		}{
			r.DatacenterRoomID.String(),
			r.Name,
			r.RoleID.String(),
			r.SerialNumber,
			r.AssetTag,
		}
		res, err = c.sling().New().Post("/rack/"+r.ID.String()).
			BodyJSON(j).Receive(&r, aerr)
	}

	return c.isHTTPResOk(res, err, aerr)
}

// DeleteGlobalRack deletes a rack
func (c *Conch) DeleteGlobalRack(id fmt.Stringer) error {
	aerr := &APIError{}
	res, err := c.sling().New().Delete("/rack/"+id.String()).Receive(nil, aerr)
	return c.isHTTPResOk(res, err, aerr)
}

// GetGlobalRackLayout fetches the layout entries for a rack in the global domain
func (c *Conch) GetGlobalRackLayout(r GlobalRack) ([]GlobalRackLayoutSlot, error) {
	rs := make([]GlobalRackLayoutSlot, 0)

	aerr := &APIError{}
	res, err := c.sling().New().Get("/rack/"+r.ID.String()+"/layouts").
		Receive(&rs, aerr)

	return rs, c.isHTTPResOk(res, err, aerr)
}
