// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	"github.com/joyent/conch-shell/pkg/pgtime"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// WorkspaceRelays ...
type WorkspaceRelays []WorkspaceRelay

// WorkspaceRelay represents a Conch Relay unit, a physical piece of hardware that
// mediates Livesys interactions in the field
type WorkspaceRelay struct {
	ID         string                 `json:"id"` // *not* a UUID
	Created    pgtime.PgTime          `json:"created"`
	Updated    pgtime.PgTime          `json:"updated"`
	Alias      string                 `json:"alias"`
	IPAddr     string                 `json:"ipaddr"`
	SSHPort    int                    `json:"ssh_port"`
	Version    string                 `json:"version"`
	LastSeen   pgtime.PgTime          `json:"last_seen"`
	NumDevices int                    `json:"num_devices"`
	Location   WorkspaceRelayLocation `json:"location"`
}

// WorkspaceRelayLocation ...
type WorkspaceRelayLocation struct {
	Az            string    `json:"az"`
	RackID        uuid.UUID `json:"rack_id"`
	RackName      string    `json:"rack_name"`
	RackUnitStart int       `json:"rack_unit_start"`
	RoleName      string    `json:"role_name"`
}

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
func (c *Conch) GetWorkspaceRelays(workspaceUUID fmt.Stringer) ([]WorkspaceRelay, error) {
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

	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/relay/"+r.ID+"/register").
		BodyJSON(d).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// GetAllRelays uses the /relay endpoint to get a list of all
// relays, but without their assigned devices.
func (c *Conch) GetAllRelays() ([]WorkspaceRelay, error) {
	relays := make([]WorkspaceRelay, 0)
	return relays, c.get("/relay?no_devices=1", &relays)
}
