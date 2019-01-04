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

// GetGlobalDatacenters fetches a list of all datacenters in the global domain
func (c *Conch) GetGlobalDatacenters() ([]GlobalDatacenter, error) {
	d := make([]GlobalDatacenter, 0)
	return d, c.get("/dc", &d)
}

// GetGlobalDatacenter fetches a single datacenter in the global domain, by its
// UUID
func (c *Conch) GetGlobalDatacenter(id fmt.Stringer) (d GlobalDatacenter, err error) {
	return d, c.get("/dc/"+id.String(), &d)
}

// SaveGlobalDatacenter creates or updates a datacenter in the global domain,
// based on the presence of an ID
func (c *Conch) SaveGlobalDatacenter(d *GlobalDatacenter) error {
	if d.Vendor == "" {
		return ErrBadInput
	}
	if d.Region == "" {
		return ErrBadInput
	}
	if d.Location == "" {
		return ErrBadInput
	}

	if uuid.Equal(d.ID, uuid.UUID{}) {
		j := struct {
			Vendor     string `json:"vendor"`
			Region     string `json:"region"`
			Location   string `json:"location"`
			VendorName string `json:"vendor_name,omitempty"`
		}{d.Vendor, d.Region, d.Location, d.VendorName}

		return c.post("/dc", j, &d)
	} else {
		j := struct {
			ID         string `json:"id"` // BUG(sungo): this is probably wrong
			Vendor     string `json:"vendor,omitempty"`
			Region     string `json:"region,omitempty"`
			Location   string `json:"location,omitempty"`
			VendorName string `json:"vendor_name,omitempty"`
		}{d.ID.String(), d.Vendor, d.Region, d.Location, d.VendorName}

		return c.post("/dc/"+d.ID.String(), j, &d)
	}
}

// DeleteGlobalDatacenter deletes a datacenter
func (c *Conch) DeleteGlobalDatacenter(id fmt.Stringer) error {
	return c.httpDelete("/dc/" + id.String())
}

// GetGlobalDatacenterRooms gets the global rooms assigned to a global datacenter
func (c *Conch) GetGlobalDatacenterRooms(d GlobalDatacenter) ([]GlobalRoom, error) {
	r := make([]GlobalRoom, 0)
	return r, c.get("/dc/"+d.ID.String()+"/rooms", &r)
}
