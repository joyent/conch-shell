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

func (c *Conch) GetRacks() ([]Rack, error) {
	r := make([]Rack, 0)
	return r, c.get("/rack", &r)
}

func (c *Conch) GetRack(id uuid.UUID) (r Rack, err error) {
	escaped := url.PathEscape(id.String())
	return r, c.get("/rack/"+escaped, &r)
}

func (c *Conch) SaveRack(r *Rack) error {
	if uuid.Equal(r.DatacenterRoomID, uuid.UUID{}) {
		return ErrBadInput
	}
	if uuid.Equal(r.RoleID, uuid.UUID{}) {
		return ErrBadInput
	}
	if r.Name == "" {
		return ErrBadInput
	}

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

	if uuid.Equal(r.ID, uuid.UUID{}) {
		return c.post("/rack", j, &r)
	} else {
		escaped := url.PathEscape(r.ID.String())
		return c.post("/rack/"+escaped, j, &r)
	}

}

func (c *Conch) DeleteRack(id uuid.UUID) error {
	escaped := url.PathEscape(id.String())
	return c.httpDelete("/rack/" + escaped)
}

// GetRackLayout fetches the layout entries for a rack in the global domain
func (c *Conch) GetRackLayout(r Rack) (RackLayoutSlots, error) {
	rs := make([]RackLayoutSlot, 0)
	escaped := url.PathEscape(r.ID.String())
	return rs, c.get("/rack/"+escaped+"/layouts", &rs)
}
func (c *Conch) SetRackPhase(id uuid.UUID, phase string, withDevices bool) error {
	data := struct {
		Phase string `json:"phase"`
	}{phase}
	escaped := url.PathEscape(id.String())

	url := "/rack/" + escaped + "/phase"
	if !withDevices {
		url = url + "?rack_only=1"
	}

	return c.post(url, data, nil)
}

func (c *Conch) GetRackPhase(id uuid.UUID) (string, error) {
	r, err := c.GetRack(id)
	return r.Phase, err
}
