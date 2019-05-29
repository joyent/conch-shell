// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"net/url"

	"github.com/joyent/conch-shell/pkg/conch/uuid"
)

func (c *Conch) GetRooms() ([]Room, error) {
	r := make([]Room, 0)
	return r, c.get("/room", &r)
}

func (c *Conch) GetRoom(id uuid.UUID) (r Room, err error) {
	return r, c.get("/room/"+url.PathEscape(id.String()), &r)
}

func (c *Conch) SaveRoom(r *Room) error {
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

func (c *Conch) DeleteRoom(id uuid.UUID) error {
	return c.httpDelete("/room/" + url.PathEscape(id.String()))
}

func (c *Conch) GetRoomRacks(r Room) ([]Rack, error) {
	rs := make([]Rack, 0)
	return rs, c.get("/room/"+url.PathEscape(r.ID.String())+"/racks", &rs)
}
