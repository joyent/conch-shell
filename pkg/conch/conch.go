// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package conch provides access to the Conch API
package conch

/*
This is a bit of trickery used elsewhere to help make it clear that we are
omitting fields from json output.

Specifically, you'll see something like

out := type writableFoo {
	*Foo
	ID omit `json:"id,omitempty"`
}{ myFoo }

This gives us a special version of Foo that leaves out the ID field unless we
explicitly set it in the new structure.

It's a bit ugly and hacky but it's the best solution I've come up with yet,
rather than repeating the full contents of the structure and just leaving those
fields out. This at least lets us grow the structure without remembering to
update it in multiple places and it lets us document explicitly via this 'omit'
data type that we're leaving stuff out.

*/
type omit bool

const (
	// MinimumAPIVersion sets the earliest API version that we support.
	MinimumAPIVersion  = "2.27.0"
	BreakingAPIVersion = "3.0.0"
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
