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
	"time"
)

// GlobalDatacenter represents a datacenter in the global domain
type GlobalDatacenter struct {
	ID         uuid.UUID `json:"id"`
	Vendor     string    `json:"vendor"`
	VendorName string    `json:"vendor_name"`
	Region     string    `json:"region"`
	Location   string    `json:"location"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

// GlobalRoom represents a datacenter room in the global domain
type GlobalRoom struct {
	ID           uuid.UUID `json:"id"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	DatacenterID uuid.UUID `json:"datacenter"`
	AZ           string    `json:"az"`
	Alias        string    `json:"alias"`
	VendorName   string    `json:"vendor_name"`
}

// GlobalRack represents a datacenter rack in the global domain
type GlobalRack struct {
	ID               uuid.UUID `json:"id"`
	Created          time.Time `json:"created"`
	Updated          time.Time `json:"updated"`
	DatacenterRoomID uuid.UUID `json:"datacenter_room_id"`
	Name             string    `json:"name"`
	RoleID           uuid.UUID `json:"role"`
	SerialNumber     string    `json:"serial_number"`
	AssetTag         string    `json:"asset_tag"`
}

// GlobalRackRole represents a rack role in the global domain
type GlobalRackRole struct {
	ID       uuid.UUID `json:"id"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Name     string    `json:"name"`
	RackSize int       `json:"rack_size"`
}

// GlobalRackLayoutSlot represents an individual rack layout entry
type GlobalRackLayoutSlot struct {
	ID        uuid.UUID `json:"id"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	RackID    uuid.UUID `json:"rack_id"`
	ProductID uuid.UUID `json:"product_id"`
	RUStart   int       `json:"ru_start"`
}

// GetGlobalDatacenters fetches a list of all datacenters in the global domain
func (c *Conch) GetGlobalDatacenters() ([]GlobalDatacenter, error) {
	d := make([]GlobalDatacenter, 0)
	return d, c.get("/dc", &d)
}

// GetGlobalDatacenter fetches a single datacenter in the global domain, by its
// UUID
func (c *Conch) GetGlobalDatacenter(id fmt.Stringer) (GlobalDatacenter, error) {
	var d GlobalDatacenter
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

	var err error
	var res *http.Response
	aerr := &APIError{}

	if uuid.Equal(d.ID, uuid.UUID{}) {
		j := struct {
			Vendor     string `json:"vendor"`
			Region     string `json:"region"`
			Location   string `json:"location"`
			VendorName string `json:"vendor_name,omitempty"`
		}{d.Vendor, d.Region, d.Location, d.VendorName}

		res, err = c.sling().New().Post("/dc").BodyJSON(j).Receive(&d, aerr)
	} else {
		j := struct {
			ID         string `json:"id"`
			Vendor     string `json:"vendor,omitempty"`
			Region     string `json:"region,omitempty"`
			Location   string `json:"location,omitempty"`
			VendorName string `json:"vendor_name,omitempty"`
		}{d.ID.String(), d.Vendor, d.Region, d.Location, d.VendorName}

		res, err = c.sling().New().Post("/dc/"+d.ID.String()).
			BodyJSON(j).Receive(&d, aerr)
	}

	return c.isHTTPResOk(res, err, aerr)
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
