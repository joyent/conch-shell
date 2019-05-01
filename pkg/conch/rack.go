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
	r, err := c.GetGlobalRack(id)
	return r.Phase, err
}
