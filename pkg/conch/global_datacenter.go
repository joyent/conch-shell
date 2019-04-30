// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"net/url"

	uuid "gopkg.in/satori/go.uuid.v1"
)

func (c *Conch) GetDatacenters() ([]Datacenter, error) {
	d := make([]Datacenter, 0)
	return d, c.get("/dc", &d)
}

func (c *Conch) GetDatacenter(id uuid.UUID) (d Datacenter, err error) {

	return d, c.get("/dc/"+id.String(), &d)
}

// SaveDatacenter creates or updates a datacenter in the global domain,
// based on the presence of an ID
func (c *Conch) SaveDatacenter(d *Datacenter) error {
	if d.Vendor == "" {
		return ErrBadInput
	}
	if d.Region == "" {
		return ErrBadInput
	}
	if d.Location == "" {
		return ErrBadInput
	}
	j := struct {
		Vendor     string `json:"vendor"`
		Region     string `json:"region"`
		Location   string `json:"location"`
		VendorName string `json:"vendor_name,omitempty"`
	}{d.Vendor, d.Region, d.Location, d.VendorName}

	if uuid.Equal(d.ID, uuid.UUID{}) {
		return c.post("/dc", j, &d)
	} else {
		escaped := url.PathEscape(d.ID.String())
		return c.post("/dc/"+escaped, j, &d)
	}
}

// DeleteDatacenter deletes a datacenter
func (c *Conch) DeleteDatacenter(id uuid.UUID) error {
	escaped := url.PathEscape(id.String())
	return c.httpDelete("/dc/" + escaped)
}

// GetDatacenterRooms gets the global rooms assigned to a global datacenter
func (c *Conch) GetDatacenterRooms(d Datacenter) ([]GlobalRoom, error) {
	r := make([]GlobalRoom, 0)
	escaped := url.PathEscape(d.ID.String())
	return r, c.get("/dc/"+escaped+"/rooms", &r)
}
