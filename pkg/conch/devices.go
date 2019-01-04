// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	uuid "gopkg.in/satori/go.uuid.v1"
)

// GetDevice returns a Device given a specific serial/id
func (c *Conch) GetDevice(serial string) (d Device, err error) {
	d.ID = serial

	return c.FillInDevice(d)
}

// FillInDevice takes an existing device and fills in its data using "/device"
//
// This exists because the API hands back partial devices in most cases. It's
// likely, though, that any client utility will eventually want all the data
// about a device and not just bits
func (c *Conch) FillInDevice(d Device) (Device, error) {
	return d, c.get("/device/"+d.ID, &d)
}

// GetDeviceLocation fetches the location for a device, via
// /device/:serial/location
func (c *Conch) GetDeviceLocation(serial string) (loc DeviceLocation, err error) {
	return loc, c.get("/device/"+serial+"/location", &loc)
}

// GraduateDevice sets the 'graduated' field for the given device, via
// /device/:serial/graduate
// WARNING: This is a one way operation and cannot currently be undone via the
// API
func (c *Conch) GraduateDevice(serial string) error {
	return c.post("/device/"+serial+"/graduate", nil, nil)
}

// DeviceTritonReboot sets the 'triton_reboot' field for the given device, via
// /device/:serial/triton_reboot
// WARNING: This is a one way operation and cannot currently be undone via the
// API
func (c *Conch) DeviceTritonReboot(serial string) error {
	return c.post("/device/"+serial+"/triton_reboot", nil, nil)
}

// SetDeviceTritonUUID sets the triton UUID via /device/:serial/triton_uuid
func (c *Conch) SetDeviceTritonUUID(serial string, id uuid.UUID) error {
	j := struct {
		TritonUUID string `json:"triton_uuid"`
	}{
		id.String(),
	}

	return c.post("/device/"+serial+"/triton_uuid", j, nil)
}

// MarkDeviceTritonSetup marks the device as setup for Triton
// For this action to succeed, the device must have its Triton UUID set and
// marked as rebooted into Triton. If these conditions are not met, this
// function will return ErrBadInput
func (c *Conch) MarkDeviceTritonSetup(serial string) error {
	return c.post("/device/"+serial+"/triton_setup", nil, nil)
}

// SetDeviceAssetTag sets the asset tag for the provided serial
func (c *Conch) SetDeviceAssetTag(serial string, tag string) error {
	j := struct {
		AssetTag string `json:"asset_tag"`
	}{
		tag,
	}

	return c.post("/device/"+serial+"/asset_tag", j, nil)
}

// GetDeviceIPMI retrieves "/device/:serial/interface/impi1/ipaddr"
func (c *Conch) GetDeviceIPMI(serial string) (string, error) {
	j := make(map[string]string)

	if err := c.get("/device/"+serial+"/interface/ipmi1/ipaddr", &j); err != nil {
		return "", err
	}

	return j["ipaddr"], nil
}
