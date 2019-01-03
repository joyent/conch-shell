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

// GetGlobalRackRoles fetches a list of all rack roles in the global domain
func (c *Conch) GetGlobalRackRoles() ([]GlobalRackRole, error) {
	r := make([]GlobalRackRole, 0)
	return r, c.get("/rack_role", &r)
}

// GetGlobalRackRole fetches a single rack role in the global domain, by its
// UUID
func (c *Conch) GetGlobalRackRole(id fmt.Stringer) (GlobalRackRole, error) {
	var r GlobalRackRole
	return r, c.get("/rack_role/"+id.String(), &r)
}

// SaveGlobalRackRole creates or updates a rack role in the global domain,
// based on the presence of an ID
func (c *Conch) SaveGlobalRackRole(r *GlobalRackRole) error {
	if r.Name == "" {
		return ErrBadInput
	}

	if r.RackSize == 0 {
		return ErrBadInput
	}

	if uuid.Equal(r.ID, uuid.UUID{}) {
		j := struct {
			Name     string `json:"name"`
			RackSize int    `json:"rack_size"`
		}{
			r.Name,
			r.RackSize,
		}

		return c.post("/rack_role", j, &r)
	} else {
		j := struct {
			ID       string `json:"id"` // BUG(sungo): this is probably wrong
			Name     string `json:"name"`
			RackSize int    `json:"rack_size"`
		}{
			r.ID.String(),
			r.Name,
			r.RackSize,
		}

		return c.post(
			"/rack_role/"+r.ID.String(),
			j,
			&r,
		)
	}

}

// DeleteGlobalRackRole deletes a rack role
func (c *Conch) DeleteGlobalRackRole(id fmt.Stringer) error {
	return c.httpDelete("/rack_role/" + id.String())
}
