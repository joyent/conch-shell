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

// GetGlobalRackLayoutSlots fetches a list of all rack layouts in the global domain
func (c *Conch) GetGlobalRackLayoutSlots() ([]GlobalRackLayoutSlot, error) {
	r := make([]GlobalRackLayoutSlot, 0)

	aerr := &APIError{}
	res, err := c.sling().New().Get("/layout").Receive(&r, aerr)

	return r, c.isHTTPResOk(res, err, aerr)
}

// GetGlobalRackLayoutSlot fetches a single rack layout in the global domain, by its
// UUID
func (c *Conch) GetGlobalRackLayoutSlot(id fmt.Stringer) (*GlobalRackLayoutSlot, error) {
	r := &GlobalRackLayoutSlot{}

	aerr := &APIError{}
	res, err := c.sling().New().Get("/layout/"+id.String()).Receive(&r, aerr)

	return r, c.isHTTPResOk(res, err, aerr)
}

// SaveGlobalRackLayoutSlot creates or updates a rack layout in the global domain,
// based on the presence of an ID
func (c *Conch) SaveGlobalRackLayoutSlot(r *GlobalRackLayoutSlot) error {
	if uuid.Equal(r.RackID, uuid.UUID{}) {
		return ErrBadInput
	}
	if uuid.Equal(r.ProductID, uuid.UUID{}) {
		return ErrBadInput
	}
	if r.RUStart == 0 {
		return ErrBadInput
	}

	var err error
	var res *http.Response
	aerr := &APIError{}

	j := struct {
		RackID    string `json:"rack_id"`
		ProductID string `json:"product_id"`
		RUStart   int    `json:"ru_start"`
	}{
		r.RackID.String(),
		r.ProductID.String(),
		r.RUStart,
	}

	if uuid.Equal(r.ID, uuid.UUID{}) {
		res, err = c.sling().New().Post("/layout").BodyJSON(j).Receive(&r, aerr)
	} else {
		res, err = c.sling().New().Post("/layout/"+r.ID.String()).
			BodyJSON(j).Receive(&r, aerr)
	}

	return c.isHTTPResOk(res, err, aerr)
}

// DeleteGlobalRackLayoutSlot deletes a rack layout
func (c *Conch) DeleteGlobalRackLayoutSlot(id fmt.Stringer) error {
	aerr := &APIError{}
	res, err := c.sling().New().Delete("/layout/"+id.String()).Receive(nil, aerr)
	return c.isHTTPResOk(res, err, aerr)
}
