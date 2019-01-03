// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	uuid "gopkg.in/satori/go.uuid.v1"
	"io"
	"time"
)

// Validation represents device validations loaded into Conch
type Validation struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Version     int       `json:"version"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// ValidationPlan represents an organized association of Validations
type ValidationPlan struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}

// ValidationResult is a result of running a validation on a device
type ValidationResult struct {
	ID              uuid.UUID `json:"id"`
	Category        string    `json:"category"`
	ComponentID     string    `json:"component_id"`
	DeviceID        string    `json:"device_id"`
	HardwareProduct uuid.UUID `json:"hardware_product_id"`
	Hint            string    `json:"hint"`
	Message         string    `json:"message"`
	Status          string    `json:"status"`
	ValidationID    uuid.UUID `json:"validation_id"`
}

// ValidationState is the result of running a validation plan on a device
type ValidationState struct {
	ID               uuid.UUID          `json:"id"`
	Created          time.Time          `json:"created"`
	Completed        time.Time          `json:"completed"`
	DeviceID         string             `json:"device_id"`
	Results          []ValidationResult `json:"results"`
	Status           string             `json:"status"`
	ValidationPlanID uuid.UUID          `json:"validation_plan_id"`
}

// GetValidations returns the contents of /validation, getting the list of all
// validations loaded in the system
func (c *Conch) GetValidations() ([]Validation, error) {
	validations := make([]Validation, 0)
	return validations, c.get("/validation", &validations)
}

// GetValidationPlans returns the contents of /validation_plan, getting the
// list of all validations plans loaded in the system
func (c *Conch) GetValidationPlans() ([]ValidationPlan, error) {
	validationPlans := make([]ValidationPlan, 0)
	return validationPlans, c.get("/validation_plan", &validationPlans)
}

// GetValidationPlan returns the contents of /validation_plan/:uuid, getting information
// about a single validation plan
// BUG(sungo): why is this returning a pointer?
func (c *Conch) GetValidationPlan(validationPlanUUID fmt.Stringer) (*ValidationPlan, error) {
	var validationPlan ValidationPlan
	return &validationPlan, c.get(
		"/validation_plan/"+validationPlanUUID.String(),
		&validationPlan,
	)
}

// CreateValidationPlan creates a new validation plan in Conch
func (c *Conch) CreateValidationPlan(newValidationPlan ValidationPlan) (ValidationPlan, error) {

	j := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{
		newValidationPlan.Name,
		newValidationPlan.Description,
	}

	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/validation_plan").
		BodyJSON(j).
		Receive(&newValidationPlan, aerr)

	return newValidationPlan, c.isHTTPResOk(res, err, aerr)
}

// AddValidationToPlan associates a validation with a validation plan
func (c *Conch) AddValidationToPlan(validationPlanUUID fmt.Stringer, validationUUID fmt.Stringer) error {
	j := struct {
		ID string `json:"id"`
	}{
		validationUUID.String(),
	}

	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/validation_plan/"+validationPlanUUID.String()+"/validation").
		BodyJSON(j).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// RemoveValidationFromPlan removes a validation from a validation plan
func (c *Conch) RemoveValidationFromPlan(validationPlanUUID fmt.Stringer, validationUUID fmt.Stringer) error {

	return c.httpDelete(
		"/validation_plan/" +
			validationPlanUUID.String() +
			"/validation/" +
			validationUUID.String(),
	)
}

// GetValidationPlanValidations gets the list of validations associated with a validation plan
func (c *Conch) GetValidationPlanValidations(validationPlanUUID fmt.Stringer) ([]Validation, error) {
	validations := make([]Validation, 0)
	return validations, c.get(
		"/validation_plan/"+validationPlanUUID.String()+"/validation",
		&validations,
	)
}

// RunDeviceValidation runs a validation against given a device and returns the results
func (c *Conch) RunDeviceValidation(deviceSerial string, validationUUID fmt.Stringer, body io.Reader) ([]ValidationResult, error) {
	results := make([]ValidationResult, 0)

	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/device/"+deviceSerial+"/validation/"+validationUUID.String()).
		Body(body).
		Receive(&results, aerr)

	return results, c.isHTTPResOk(res, err, aerr)
}

// RunDeviceValidationPlan runs a validation plan against a given device and returns the results
func (c *Conch) RunDeviceValidationPlan(deviceSerial string, validationPlanUUID fmt.Stringer, body io.Reader) ([]ValidationResult, error) {
	results := make([]ValidationResult, 0)

	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/device/"+deviceSerial+"/validation_plan/"+validationPlanUUID.String()).
		Body(body).
		Receive(&results, aerr)

	return results, c.isHTTPResOk(res, err, aerr)
}

// DeviceValidationStates returns the stored validation states for a device
func (c *Conch) DeviceValidationStates(deviceSerial string) ([]ValidationState, error) {
	states := make([]ValidationState, 0)
	return states, c.get("/device/"+deviceSerial+"/validation_state", &states)
}

// WorkspaceValidationStates returns the stored validation states for all devices in a workspace
func (c *Conch) WorkspaceValidationStates(workspaceUUID fmt.Stringer) ([]ValidationState, error) {
	states := make([]ValidationState, 0)
	return states, c.get(
		"/workspace/"+workspaceUUID.String()+"/validation_state",
		&states,
	)
}
