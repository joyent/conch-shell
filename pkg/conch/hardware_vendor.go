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

	aerr := &APIError{}
	res, err := c.sling().New().Get("/hardware_vendor/"+name).
		Receive(&vendor, aerr)

	return vendor, c.isHTTPResOk(res, err, aerr)
}

// GetHardwareVendors ...
func (c *Conch) GetHardwareVendors() ([]HardwareVendor, error) {
	vendors := make([]HardwareVendor, 0)

	aerr := &APIError{}
	res, err := c.sling().New().Get("/hardware_vendor").Receive(&vendors, aerr)

	return vendors, c.isHTTPResOk(res, err, aerr)
}

// DeleteHardwareVendor ...
func (c *Conch) DeleteHardwareVendor(name string) error {
	aerr := &APIError{}
	res, err := c.sling().New().Delete("/hardware_vendor/"+name).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// SaveHardwareVendor ...
func (c *Conch) SaveHardwareVendor(v *HardwareVendor) error {
	if v.Name == "" {
		return ErrBadInput
	}

	if !uuid.Equal(v.ID, uuid.UUID{}) {
		return ErrBadInput
	}

	aerr := &APIError{}

	out := struct {
		Name string `json:"name"`
	}{v.Name}

	res, err := c.sling().New().Post("/hardware_vendor/"+v.Name).BodyJSON(out).
		Receive(&v, aerr)

	return c.isHTTPResOk(res, err, aerr)
}
