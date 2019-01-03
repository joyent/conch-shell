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

// GetGlobalRacks fetches a list of all racks in the global domain
func (c *Conch) GetGlobalRacks() ([]GlobalRack, error) {
	r := make([]GlobalRack, 0)
	return r, c.get("/rack", &r)
}

// GetGlobalRack fetches a single rack in the global domain, by its
// UUID
func (c *Conch) GetGlobalRack(id fmt.Stringer) (GlobalRack, error) {
	var r GlobalRack
	return r, c.get("/rack/"+id.String(), &r)
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

		return c.post("/rack", j, &r)
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
		return c.post("/rack/"+r.ID.String(), j, &r)
	}

}

// DeleteGlobalRack deletes a rack
func (c *Conch) DeleteGlobalRack(id fmt.Stringer) error {
	return c.httpDelete("/rack/" + id.String())
}

// GetGlobalRackLayout fetches the layout entries for a rack in the global domain
func (c *Conch) GetGlobalRackLayout(r GlobalRack) ([]GlobalRackLayoutSlot, error) {
	rs := make([]GlobalRackLayoutSlot, 0)
	return rs, c.get("/rack/"+r.ID.String()+"/layouts", &rs)
}
