// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"github.com/joyent/conch-shell/pkg/pgtime"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// HardwareVendor ...
type HardwareVendor struct {
	ID      uuid.UUID     `json:"id"`
	Name    string        `json:"name"`
	Created pgtime.PgTime `json:"created"`
	Updated pgtime.PgTime `json:"updated"`
}

// GetHardwareVendor ...
func (c *Conch) GetHardwareVendor(name string) (HardwareVendor, error) {
	var vendor HardwareVendor
	return vendor, c.get("/hardware_vendor/"+name, &vendor)
}

// GetHardwareVendors ...
func (c *Conch) GetHardwareVendors() ([]HardwareVendor, error) {
	vendors := make([]HardwareVendor, 0)
	return vendors, c.get("/hardware_vendor", &vendors)
}

// DeleteHardwareVendor ...
func (c *Conch) DeleteHardwareVendor(name string) error {
	return c.httpDelete("/hardware_vendor/" + name)
}

// SaveHardwareVendor ...
func (c *Conch) SaveHardwareVendor(v *HardwareVendor) error {
	if v.Name == "" {
		return ErrBadInput
	}

	if !uuid.Equal(v.ID, uuid.UUID{}) {
		return ErrBadInput
	}

	out := struct {
		Name string `json:"name"`
	}{v.Name}

	return c.post(
		"/hardware_vendor/"+v.Name,
		out,
		&v,
	)
}
