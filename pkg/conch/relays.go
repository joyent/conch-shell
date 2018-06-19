// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	"time"
)

// Relay represents a Conch Relay unit, a physical piece of hardware that
// mediates Livesys interactions in the field
type Relay struct {
	Alias   string    `json:"alias"`
	Created time.Time `json:"created"`
	Devices []Device  `json:"devices"`
	ID      string    `json:"id"`
	IPAddr  string    `json:"ipaddr"`
	SSHPort int       `json:"ssh_port"`
	Updated time.Time `json:"updated"`
	Version string    `json:"version"`
}

// GetWorkspaceRelays returns all Relays associated with the given workspace
func (c *Conch) GetWorkspaceRelays(workspaceUUID fmt.Stringer, activeOnly bool) ([]Relay, error) {
	var err error

	relays := make([]Relay, 0)

	var url string
	if activeOnly {
		url = "/workspace/" + workspaceUUID.String() + "/relay?active_only=1"
	} else {
		url = "/workspace/" + workspaceUUID.String() + "/relay"
	}

	aerr := &APIError{}
	res, err := c.sling().New().Get(url).Receive(&relays, aerr)

	return relays, c.isHTTPResOk(res, err, aerr)

}

// RegisterRelay registers/updates a Relay via /relay/:serial/register
// If the provided relay does not have an IP, SSHPort, and Version, ErrBadInput
// will be returned
func (c *Conch) RegisterRelay(r Relay) error {
	if (r.ID == "") || (r.SSHPort == 0) || (r.Version == "") {
		return ErrBadInput
	}

	d := struct {
		Alias   string `json:"alias,omitempty"`
		IPAddr  string `json:"ipaddr,omitempty"`
		SSHPort int    `json:"ssh_port"`
		Version string `json:"version"`
	}{
		r.Alias,
		r.IPAddr,
		r.SSHPort,
		r.Version,
	}

	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/relay/"+r.ID+"/register").
		BodyJSON(d).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// GetAllRelays uses the /relay endpoint to get a list of all relays devices
// known in the system. This endpoint is only supported in API versions >=2.1.0
// and will return ErrNotSupported if the API version is not compliant. This
// endpoint is also limited to Administrators on the GLOBAL workspace. If
// permissions are not met, ErrNotAuthorized will be returned
func (c *Conch) GetAllRelays() ([]Relay, error) {
	relays := make([]Relay, 0)

	aerr := &APIError{}
	res, err := c.sling().New().Get("/relay").Receive(&relays, aerr)
	return relays, c.isHTTPResOk(res, err, aerr)
}

// GetAllRelaysWithoutDevices uses the /relay endpoint to get a list of all
// relays, but without their assigned devices.
func (c *Conch) GetAllRelaysWithoutDevices() ([]Relay, error) {
	relays := make([]Relay, 0)

	aerr := &APIError{}
	res, err := c.sling().New().Get("/relay?no_devices=1").Receive(&relays, aerr)
	return relays, c.isHTTPResOk(res, err, aerr)
}
