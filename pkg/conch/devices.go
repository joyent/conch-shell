// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// GetWorkspaceDevices retrieves a list of Devices for the given
// workspace.
// Pass true for 'IDsOnly' to get Devices with only the ID field populated
// Pass a string for 'graduated' to filter devices by graduated value, as per https://conch.joyent.us/doc#getdevices
// Pass a string for 'health' to filter devices by health value, as per https://conch.joyent.us/doc#getdevices
func (c *Conch) GetWorkspaceDevices(workspaceUUID fmt.Stringer, idsOnly bool, graduated string, health string) ([]Device, error) {

	devices := make([]Device, 0)

	opts := struct {
		IDsOnly   bool   `url:"ids_only,omitempty"`
		Graduated string `url:"graduated,omitempty"`
		Health    string `url:"health,omitempty"`
	}{
		idsOnly,
		graduated,
		health,
	}

	url := "/workspace/" + workspaceUUID.String() + "/device"
	if idsOnly {
		ids := make([]string, 0)

		if err := c.getWithQuery(url, opts, &ids); err != nil {
			return devices, err
		}

		for _, v := range ids {
			device := Device{ID: v}
			devices = append(devices, device)
		}
		return devices, nil
	}
	return devices, c.getWithQuery(url, opts, &devices)
}

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

// GetWorkspaceRacks fetchest the list of racks for a workspace, via
// /workspace/:uuid/rack
// NOTE: The API currently returns a hash of arrays where the key is the
// datacenter/az. This routine copies that key into the Datacenter field in the
// Rack struct.
func (c *Conch) GetWorkspaceRacks(workspaceUUID fmt.Stringer) ([]Rack, error) {
	racks := make([]Rack, 0)
	j := make(map[string][]Rack)

	if err := c.get("/workspace/"+workspaceUUID.String()+"/rack", &j); err != nil {
		return racks, err
	}

	for az, loc := range j {
		for _, rack := range loc {
			rack.Datacenter = az
			racks = append(racks, rack)
		}
	}

	return racks, nil
}

// GetWorkspaceRack fetches a single rack for a workspace, via
// /workspace/:uuid/rack/:id
func (c *Conch) GetWorkspaceRack(
	workspaceUUID fmt.Stringer,
	rackUUID fmt.Stringer,
) (rack Rack, err error) {
	return rack, c.get(
		"/workspace/"+
			workspaceUUID.String()+
			"/rack/"+
			rackUUID.String(),
		&rack,
	)
}

// GetHardwareProduct fetches a single hardware product via
// /hardware/product/:uuid
func (c *Conch) GetHardwareProduct(
	hardwareProductUUID fmt.Stringer,
) (hp HardwareProduct, err error) {
	return hp, c.get("/hardware_product/"+hardwareProductUUID.String(), &hp)
}

// GetHardwareProducts fetches a single hardware product via
// /hardware_product
func (c *Conch) GetHardwareProducts() ([]HardwareProduct, error) {
	prods := make([]HardwareProduct, 0)
	return prods, c.get("/hardware_product", &prods)
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

// SaveHardwareProduct creates or saves s hardware product, based
// on the presence of an ID
func (c *Conch) SaveHardwareProduct(h *HardwareProduct) error {
	if h.Name == "" {
		return ErrBadInput
	}

	if h.Alias == "" {
		return ErrBadInput
	}

	if uuid.Equal(h.HardwareVendorID, uuid.UUID{}) {
		return ErrBadInput
	}

	out := struct {
		Name              string `json:"name"`
		Alias             string `json:"alias"`
		Prefix            string `json:"prefix"`
		HardwareVendorID  string `json:"hardware_vendor_id"`
		Specification     string `json:"specification,omitempty"`
		SKU               string `json:"sku"`
		GenerationName    string `json:"generation_name"`
		LegacyProductName string `json:"legacy_product_name"`
	}{
		h.Name,
		h.Alias,
		h.Prefix,
		h.HardwareVendorID.String(),
		h.Specification,
		h.SKU,
		h.GenerationName,
		h.LegacyProductName,
	}

	if uuid.Equal(h.ID, uuid.UUID{}) {
		return c.post("/hardware_product", out, &h)
	} else {
		return c.post(
			"/hardware_product/"+h.ID.String(),
			out,
			&h,
		)
	}
}

// DeleteHardwareProduct deletes a hardware product by marking it as
// deactivated
func (c *Conch) DeleteHardwareProduct(hwUUID fmt.Stringer) error {
	return c.httpDelete("/hardware_product/" + hwUUID.String())
}
