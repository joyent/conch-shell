// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"encoding/json"
	"fmt"
)

// GetValidations returns the contents of /validation, getting the list of all
// validations loaded in the system
func (c *Conch) GetValidations() ([]Validation, error) {
	validations := make([]Validation, 0)
	return validations, c.get("/validation", &validations)
}
func (c *Conch) GetValidation(id fmt.Stringer) (v Validation, err error) {
	return v, c.get("/validation/"+id.String(), &v)
}

// GetValidationPlans returns the contents of /validation_plan, getting the
// list of all validations plans loaded in the system
func (c *Conch) GetValidationPlans() ([]ValidationPlan, error) {
	validationPlans := make([]ValidationPlan, 0)
	return validationPlans, c.get("/validation_plan", &validationPlans)
}

// GetValidationPlan returns the contents of /validation_plan/:uuid, getting information
// about a single validation plan
func (c *Conch) GetValidationPlan(
	validationPlanUUID fmt.Stringer,
) (vp ValidationPlan, err error) {

	return vp, c.get(
		"/validation_plan/"+validationPlanUUID.String(),
		&vp,
	)
}

// GetValidationPlanValidations gets the list of validations associated with a validation plan
func (c *Conch) GetValidationPlanValidations(
	validationPlanUUID fmt.Stringer,
) ([]Validation, error) {

	validations := make([]Validation, 0)
	return validations, c.get(
		"/validation_plan/"+validationPlanUUID.String()+"/validation",
		&validations,
	)
}

// RunDeviceValidation runs a validation against given a device and returns the results
func (c *Conch) RunDeviceValidation(
	deviceSerial string,
	validationUUID fmt.Stringer,
	body string,
) ([]ValidationResult, error) {

	results := make([]ValidationResult, 0)

	return results, c.post(
		"/device/"+deviceSerial+"/validation/"+validationUUID.String(),
		body,
		&results,
	)
}

// RunDeviceValidationPlan runs a validation plan against a given device and returns the results
func (c *Conch) RunDeviceValidationPlan(
	deviceSerial string,
	validationPlanUUID fmt.Stringer,
	body string,
) ([]ValidationResult, error) {

	results := make([]ValidationResult, 0)

	var j interface{}
	err := json.Unmarshal([]byte(body), &j)
	if err != nil {
		return results, err
	}

	return results, c.post(
		"/device/"+deviceSerial+"/validation_plan/"+validationPlanUUID.String(),
		j,
		&results,
	)
}

// DeviceValidationStates returns the stored validation states for a device
func (c *Conch) DeviceValidationStates(
	deviceSerial string,
) ([]ValidationState, error) {

	states := make([]ValidationState, 0)
	return states, c.get("/device/"+deviceSerial+"/validation_state", &states)
}

// WorkspaceValidationStates returns the stored validation states for all devices in a workspace
func (c *Conch) WorkspaceValidationStates(
	workspaceUUID fmt.Stringer,
) ([]ValidationState, error) {

	states := make([]ValidationState, 0)
	return states, c.get(
		"/workspace/"+workspaceUUID.String()+"/validation_state",
		&states,
	)
}
