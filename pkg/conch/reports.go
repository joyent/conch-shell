// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	uuid "gopkg.in/satori/go.uuid.v1"
)

// ServerReport represents data obtained in the field about the current state
// of a Device
type ServerReport struct {
	BiosVersion  string                            `json:"bios_version"`
	Disks        map[string]map[string]interface{} `json:"disks"`
	Interfaces   map[string]map[string]interface{} `json:"interfaces"`
	Memory       map[string]interface{}            `json:"memory"`
	Processor    map[string]interface{}            `json:"processor"`
	ProductName  string                            `json:"product_name"`
	Relay        map[string]interface{}            `json:"relay"` // "Only key in hash is currently 'serial"
	SerialNumber string                            `json:"serial_number"`
	State        string                            `json:"state"`
	SystemUUID   uuid.UUID                         `json:"system_uuid"`
	Temp         map[string]interface{}            `json:"temp"`
	UptimeSince  string                            `json:"uptime_since"`
	Aux          interface{}                       `json:"aux"`
}
