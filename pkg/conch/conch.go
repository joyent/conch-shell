// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package conch provides access to the Conch API
package conch

const (
	// MinimumAPIVersion sets the earliest API version that we support.
	MinimumAPIVersion = "2.24.0"
)

// GetVersion returns the API's version string, via /version
func (c *Conch) GetVersion() (string, error) {
	v := struct {
		Version string `json:"version"`
	}{}

	err := c.get("/version", &v)
	if err != nil {
		return "", err
	}

	return v.Version, nil
}
