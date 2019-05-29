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

func (c *Conch) GetRackLayoutSlots() (RackLayoutSlots, error) {
	r := make([]RackLayoutSlot, 0)
	return r, c.get("/layout", &r)
}

func (c *Conch) GetRackLayoutSlot(id uuid.UUID) (*RackLayoutSlot, error) {
	r := &RackLayoutSlot{}
	escaped := url.PathEscape(id.String())
	return r, c.get("/layout/"+escaped, &r)
}

func (c *Conch) SaveRackLayoutSlot(r *RackLayoutSlot) error {
	if uuid.Equal(r.RackID, uuid.UUID{}) {
		return ErrBadInput
	}
	if uuid.Equal(r.ProductID, uuid.UUID{}) {
		return ErrBadInput
	}
	if r.RUStart == 0 {
		return ErrBadInput
	}

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
		return c.post("/layout", j, &r)
	} else {
		escaped := url.PathEscape(r.ID.String())
		return c.post("/layout/"+escaped, j, &r)
	}
}

func (c *Conch) DeleteRackLayoutSlot(id uuid.UUID) error {
	escaped := url.PathEscape(id.String())
	return c.httpDelete("/layout/" + escaped)
}
