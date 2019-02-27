// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
)

// GetActiveWorkspaceRelays ...
func (c *Conch) GetActiveWorkspaceRelays(
	workspaceUUID fmt.Stringer,
	minutes int,
) ([]WorkspaceRelay, error) {
	if minutes == 0 {
		minutes = 5
	}

	relays := make([]WorkspaceRelay, 0)

	url := fmt.Sprintf("/workspace/%s/relay?active_within=%d",
		workspaceUUID.String(),
		minutes,
	)

	return relays, c.get(url, &relays)
}

// GetWorkspaceRelays returns all Relays associated with the given workspace
func (c *Conch) GetWorkspaceRelays(workspaceUUID fmt.Stringer) (WorkspaceRelays, error) {
	relays := make([]WorkspaceRelay, 0)

	url := "/workspace/" + workspaceUUID.String() + "/relay"
	return relays, c.get(url, &relays)
}

// GetWorkspaceRelayDevices ...
func (c *Conch) GetWorkspaceRelayDevices(
	workspaceUUID fmt.Stringer,
	relayName string,
) ([]Device, error) {

	devices := make([]Device, 0)
	url := "/workspace/" + workspaceUUID.String() + "/relay/" + relayName + "/device"
	return devices, c.get(url, &devices)
}

// RegisterRelay registers/updates a Relay via /relay/:serial/register
// If the provided relay does not have an IP, SSHPort, and Version, ErrBadInput
// will be returned
func (c *Conch) RegisterRelay(r WorkspaceRelay) error {
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

	return c.post(
		"/relay/"+r.ID+"/register",
		d,
		nil,
	)
}

// GetAllRelays uses the /relay endpoint to get a list of all
// relays, but without their assigned devices.
func (c *Conch) GetAllRelays() (WorkspaceRelays, error) {
	relays := make([]WorkspaceRelay, 0)
	return relays, c.get("/relay?no_devices=1", &relays)
}
