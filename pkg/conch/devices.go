// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"

	"github.com/joyent/conch-shell/pkg/conch/uuid"
)

// GetDevice returns a Device given a specific serial/id
func (c *Conch) GetDevice(serial string) (d Device, err error) {
	d.ID = serial

	return c.FillInDevice(d)
}

func (c *Conch) GetExtendedDevice(serial string) (ed ExtendedDevice, err error) {

	d, err := c.GetDevice(serial)
	if err != nil {
		return ExtendedDevice{}, err
	}

	enclosures := make(map[string]map[int]Disk)
	for _, disk := range d.Disks {
		enclosure, ok := enclosures[disk.Enclosure]
		if !ok {
			enclosure = make(map[int]Disk)
		}

		if _, ok := enclosure[disk.Slot]; !ok {
			enclosure[disk.Slot] = disk
		}

		enclosures[disk.Enclosure] = enclosure
	}

	/***********/

	var role RackRole
	if !d.Location.Rack.RoleID.IsZero() {
		role, err = c.GetRackRole(d.Location.Rack.RoleID)
		if err != nil {
			return ed, err
		}
	}

	ed = ExtendedDevice{
		Device:        d,
		IPMI:          "",
		HardwareName:  "",
		SKU:           "",
		RackRole:      role,
		Enclosures:    enclosures,
		IsGraduated:   !d.Graduated.IsZero(),
		IsTritonSetup: !d.TritonSetup.IsZero(),
		IsValidated:   !d.Validated.IsZero(),
		Validations:   make([]ValidationPlanExecution, 0),
	}

	allValidations := make(map[uuid.UUID]Validation)
	serverValidations, err := c.GetValidations()
	if err != nil {
		return ed, err
	}
	for _, v := range serverValidations {
		allValidations[v.ID] = v
	}

	if validationStates, err := c.DeviceValidationStates(d.ID); err == nil {

		plans := make([]ValidationPlanExecution, 0)

		for _, state := range validationStates {
			name := "[unknown]"
			validationPlan, err := c.GetValidationPlan(state.ValidationPlanID)
			if err == nil {
				name = validationPlan.Name
			}

			byValidationID := make(map[uuid.UUID][]ValidationResult)

			for _, result := range state.Results {
				results, ok := byValidationID[result.ValidationID]
				if !ok {
					results = make([]ValidationResult, 0)
				}

				results = append(results, result)
				byValidationID[result.ValidationID] = results
			}

			runs := make(ValidationRuns, 0)
			for id, results := range byValidationID {
				var run ValidationRun
				run.ID = id
				run.Name = "[unknown]"
				if v, ok := allValidations[id]; ok {
					run.Name = v.Name
				}
				passed := true
				for _, result := range results {
					if result.Status != "pass" {
						passed = false
					}
				}

				run.Passed = passed
				run.Results = results
				runs = append(runs, run)
			}

			sort.Sort(runs)
			plans = append(plans, ValidationPlanExecution{
				ID:          state.ValidationPlanID,
				Name:        name,
				Validations: runs,
			})

		}

		ed.Validations = plans
	}

	ipmi, err := c.GetDeviceIPMI(d.ID)
	if err == nil {
		ed.IPMI = ipmi
	}

	if !uuid.Equal(d.HardwareProduct, uuid.UUID{}) {
		hp, err := c.GetHardwareProduct(d.HardwareProduct)
		if err == nil {
			ed.HardwareName = hp.Name
			ed.SKU = hp.SKU
		}
	}

	return ed, nil
}

// FillInDevice takes an existing device and fills in its data using "/device"
//
// This exists because the API hands back partial devices in most cases. It's
// likely, though, that any client utility will eventually want all the data
// about a device and not just bits
func (c *Conch) FillInDevice(d Device) (Device, error) {
	escaped := url.PathEscape(d.ID)
	return d, c.get("/device/"+escaped, &d)
}

// GetDeviceLocation fetches the location for a device, via
// /device/:serial/location
func (c *Conch) GetDeviceLocation(serial string) (loc DeviceLocation, err error) {
	escaped := url.PathEscape(serial)
	return loc, c.get("/device/"+escaped+"/location", &loc)
}

// GraduateDevice sets the 'graduated' field for the given device, via
// /device/:serial/graduate
// WARNING: This is a one way operation and cannot currently be undone via the
// API
func (c *Conch) GraduateDevice(serial string) error {
	escaped := url.PathEscape(serial)
	return c.post("/device/"+escaped+"/graduate", nil, nil)
}

// DeviceTritonReboot sets the 'triton_reboot' field for the given device, via
// /device/:serial/triton_reboot
// WARNING: This is a one way operation and cannot currently be undone via the
// API
func (c *Conch) DeviceTritonReboot(serial string) error {
	escaped := url.PathEscape(serial)
	return c.post("/device/"+escaped+"/triton_reboot", nil, nil)
}

// SetDeviceTritonUUID sets the triton UUID via /device/:serial/triton_uuid
func (c *Conch) SetDeviceTritonUUID(serial string, id uuid.UUID) error {
	j := struct {
		TritonUUID string `json:"triton_uuid"`
	}{
		id.String(),
	}

	escaped := url.PathEscape(serial)
	return c.post("/device/"+escaped+"/triton_uuid", j, nil)
}

// MarkDeviceTritonSetup marks the device as setup for Triton
// For this action to succeed, the device must have its Triton UUID set and
// marked as rebooted into Triton. If these conditions are not met, this
// function will return ErrBadInput
func (c *Conch) MarkDeviceTritonSetup(serial string) error {
	escaped := url.PathEscape(serial)
	return c.post("/device/"+escaped+"/triton_setup", nil, nil)
}

// SetDeviceAssetTag sets the asset tag for the provided serial
func (c *Conch) SetDeviceAssetTag(serial string, tag string) error {
	j := struct {
		AssetTag string `json:"asset_tag"`
	}{
		tag,
	}

	escaped := url.PathEscape(serial)
	return c.post("/device/"+escaped+"/asset_tag", j, nil)
}

// GetDeviceIPMI retrieves "/device/:serial/interface/impi1/ipaddr"
func (c *Conch) GetDeviceIPMI(serial string) (string, error) {
	j := make(map[string]string)

	escaped := url.PathEscape(serial)
	if err := c.get("/device/"+escaped+"/interface/ipmi1/ipaddr", &j); err != nil {
		return "", err
	}

	return j["ipaddr"], nil
}

func (c *Conch) GetDevicesByField(key string, value string) (d Devices, err error) {
	escapedKey := url.PathEscape(key)
	escapedValue := url.PathEscape(value)
	url := fmt.Sprintf("/device?%s=%s", escapedKey, escapedValue)
	return d, c.get(url, &d)
}

func (c *Conch) SubmitDeviceReport(serial string, report string) (state ValidationState, err error) {
	reportReader := bytes.NewReader([]byte(report))

	escaped := url.PathEscape(serial)
	req, err := c.sling().New().Post("/device/"+escaped).
		Set("Content-Type", "application/json").Body(reportReader).Request()

	if err != nil {
		return state, err
	}

	_, err = c.httpDo(req, &state)
	return state, err
}

func (c *Conch) GetDevicePhase(serial string) (string, error) {
	ret := struct {
		DeviceID string `json:"id"`
		Phase    string `json:"phase"`
	}{}

	return ret.Phase, c.get("/device/"+url.PathEscape(serial)+"/phase", &ret)
}

func (c *Conch) SetDevicePhase(serial string, phase string) error {
	data := struct {
		Phase string `json:"phase"`
	}{phase}

	return c.post("/device/"+url.PathEscape(serial)+"/phase", data, nil)
}
