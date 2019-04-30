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

// GetGlobalRooms fetches a list of all rooms in the global domain
func (c *Conch) GetGlobalRooms() ([]GlobalRoom, error) {
	r := make([]GlobalRoom, 0)
	return r, c.get("/room", &r)
}

// GetGlobalRoom fetches a single room in the global domain, by its
// UUID
func (c *Conch) GetGlobalRoom(id fmt.Stringer) (r GlobalRoom, err error) {
	return r, c.get("/room/"+url.PathEscape(id.String()), &r)
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

	if r.Alias == "" {
		return ErrBadInput
	}

	j := struct {
		Datacenter string `json:"datacenter"`
		AZ         string `json:"az"`
		Alias      string `json:"alias"`
		VendorName string `json:"vendor_name,omitempty"`
	}{r.DatacenterID.String(), r.AZ, r.Alias, r.VendorName}

	if uuid.Equal(r.ID, uuid.UUID{}) {
		return c.post("/room", j, &r)
	} else {
		return c.post("/room/"+url.PathEscape(r.ID.String()), j, &r)
	}
}

// DeleteGlobalRoom deletes a room
func (c *Conch) DeleteGlobalRoom(id fmt.Stringer) error {
	return c.httpDelete("/room/" + url.PathEscape(id.String()))
}

// GetGlobalRoomRacks retrieves the racks assigned to a room in the global
// domain
func (c *Conch) GetGlobalRoomRacks(r GlobalRoom) ([]GlobalRack, error) {
	rs := make([]GlobalRack, 0)
	return rs, c.get("/room/"+url.PathEscape(r.ID.String())+"/racks", &rs)
}
