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

func (c *Conch) GetRackRoles() ([]RackRole, error) {
	r := make([]RackRole, 0)
	return r, c.get("/rack_role", &r)
}

func (c *Conch) GetRackRole(id uuid.UUID) (r RackRole, err error) {
	return r, c.get("/rack_role/"+url.PathEscape(id.String()), &r)
}

func (c *Conch) SaveRackRole(r *RackRole) error {
	if r.Name == "" {
		return ErrBadInput
	}

	if r.RackSize == 0 {
		return ErrBadInput
	}

	j := struct {
		Name     string `json:"name"`
		RackSize int    `json:"rack_size"`
	}{
		r.Name,
		r.RackSize,
	}

	if uuid.Equal(r.ID, uuid.UUID{}) {
		return c.post("/rack_role", j, &r)
	} else {
		return c.post("/rack_role/"+url.PathEscape(r.ID.String()), j, &r)
	}

}

func (c *Conch) DeleteRackRole(id uuid.UUID) error {
	return c.httpDelete("/rack_role/" + url.PathEscape(id.String()))
}
