// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// config wraps a conch shell config. Typically this is either coming from
// and/or becoming JSON on disk.
package config

import (
	"encoding/json"
)

type ConchConfig struct {
	Api     string
	User    string
	Session string
	KV      map[string]interface{}
}



// Serialize() marshals a ConchConfig struct into a JSON string
func (c *ConchConfig) Serialize() (s string, err error) {

	j, err := json.MarshalIndent(c, "", "	")

	if err != nil {
		return "", err
	}

	return string(j), nil
}
