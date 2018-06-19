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

// GetGlobalRooms fetches a list of all rooms in the global domain
func (c *Conch) GetGlobalRooms() ([]GlobalRoom, error) {
	r := make([]GlobalRoom, 0)

	aerr := &APIError{}
	res, err := c.sling().New().Get("/room").Receive(&r, aerr)

	return r, c.isHTTPResOk(res, err, aerr)
}

// GetGlobalRoom fetches a single room in the global domain, by its
// UUID
func (c *Conch) GetGlobalRoom(id fmt.Stringer) (GlobalRoom, error) {
	r := GlobalRoom{}

	aerr := &APIError{}
	res, err := c.sling().New().Get("/room/"+id.String()).Receive(&r, aerr)

	return r, c.isHTTPResOk(res, err, aerr)
}

// SaveGlobalRoom creates or updates a room in the global domain,
// based on the presence of an ID
func (c *Conch) SaveGlobalRoom(r *GlobalRoom) error {
	if uuid.Equal(r.DatacenterID, uuid.UUID{}) {
		return ErrBadInput
	}
	if r.AZ == "" {
		return ErrBadInput
	}

	var err error
	var res *http.Response
	aerr := &APIError{}

	if uuid.Equal(r.ID, uuid.UUID{}) {
		j := struct {
			Datacenter string `json:"datacenter"`
			AZ         string `json:"az"`
			Alias      string `json:"alias,omitempty"`
			VendorName string `json:"vendor_name,omitempty"`
		}{r.DatacenterID.String(), r.AZ, r.Alias, r.VendorName}

		res, err = c.sling().New().Post("/room").BodyJSON(j).Receive(&r, aerr)
	} else {
		j := struct {
			ID         string `json:"id"`
			Datacenter string `json:"datacenter"`
			AZ         string `json:"az"`
			Alias      string `json:"alias,omitempty"`
			VendorName string `json:"vendor_name,omitempty"`
		}{r.ID.String(), r.DatacenterID.String(), r.AZ, r.Alias, r.VendorName}

		res, err = c.sling().New().Post("/room/"+r.ID.String()).
			BodyJSON(j).Receive(&r, aerr)
	}

	return c.isHTTPResOk(res, err, aerr)
}

// DeleteGlobalRoom deletes a room
func (c *Conch) DeleteGlobalRoom(id fmt.Stringer) error {
	aerr := &APIError{}
	res, err := c.sling().New().Delete("/room/"+id.String()).Receive(nil, aerr)
	return c.isHTTPResOk(res, err, aerr)
}

// GetGlobalRoomRacks retrieves the racks assigned to a room in the global
// domain
func (c *Conch) GetGlobalRoomRacks(r GlobalRoom) ([]GlobalRack, error) {
	rs := make([]GlobalRack, 0)

	aerr := &APIError{}
	res, err := c.sling().New().Get("/room/"+r.ID.String()+"/racks").
		Receive(&rs, aerr)

	return rs, c.isHTTPResOk(res, err, aerr)
}
