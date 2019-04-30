// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	"net/url"

	"github.com/joyent/conch-shell/pkg/conch/uuid"
)

// GetGlobalRackLayoutSlots fetches a list of all rack layouts in the global domain
func (c *Conch) GetGlobalRackLayoutSlots() (GlobalRackLayoutSlots, error) {
	r := make([]GlobalRackLayoutSlot, 0)
	return r, c.get("/layout", &r)
}

// GetGlobalRackLayoutSlot fetches a single rack layout in the global domain, by its
// UUID
func (c *Conch) GetGlobalRackLayoutSlot(id fmt.Stringer) (*GlobalRackLayoutSlot, error) {
	r := &GlobalRackLayoutSlot{}
	escaped := url.PathEscape(id.String())
	return r, c.get("/layout/"+escaped, &r)
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

// DeleteGlobalRackLayoutSlot deletes a rack layout
func (c *Conch) DeleteGlobalRackLayoutSlot(id fmt.Stringer) error {
	escaped := url.PathEscape(id.String())
	return c.httpDelete("/layout/" + escaped)
}
