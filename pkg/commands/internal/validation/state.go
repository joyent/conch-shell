// Copyright 2018 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package validation contains commands for validation related commands
package validation

import (
	"strings"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
)

type validationStates []conch.ValidationState

func (vs validationStates) renderTable(validationPlans []conch.ValidationPlan, validations []conch.Validation) {
	table := util.GetMarkdownTable()

	planNameMap := make(map[uuid.UUID]string)
	for _, vp := range validationPlans {
		planNameMap[vp.ID] = vp.Name
	}

	validationNameMap := make(map[uuid.UUID]string)
	for _, v := range validations {
		validationNameMap[v.ID] = v.Name
	}

	table.SetHeader([]string{
		"Device ID",
		"Status",
		"Completed",
		"Validation Plan",
		"Results",
	})

	/* group the validations by name. For each name group, output "pass" if all
	* results pass or create a string with each failing validation results */
	for _, v := range vs {
		resultGroup := make(map[string]string)
		for _, r := range v.Results {
			vName := validationNameMap[r.ValidationID]
			if _, ok := resultGroup[vName]; !ok {
				resultGroup[vName] = "pass" // Default to pass
			}

			if r.Status != "pass" {
				message := r.Status + "\n  " + r.Message
				if r.Hint != "" {
					message = message + " (" + r.Hint + ")"
				}
				resultGroup[vName] = message
			}
		}
		results := make([]string, 0, len(resultGroup))
		for vName, vResult := range resultGroup {
			results = append(results, vName+": "+vResult)
		}
		table.Append([]string{
			v.DeviceID,
			v.Status,
			v.Completed.String(),
			planNameMap[v.ValidationPlanID],
			strings.Join(results, "\n"),
		})
	}

	table.Render()
}

func getDeviceValidationStates(app *cli.Cmd) {
	var deviceSerial = app.StringArg("DEVICE_ID", "", "The Device ID (serial number)")

	app.Spec = "DEVICE_ID"

	app.Action = func() {
		var validationStates validationStates
		validationStates, err := util.API.DeviceValidationStates(*deviceSerial)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(validationStates)
			return
		}
		validationPlans, err := util.API.GetValidationPlans()
		if err != nil {
			util.Bail(err)
		}
		validations, err := util.API.GetValidations()
		if err != nil {
			util.Bail(err)
		}
		validationStates.renderTable(validationPlans, validations)
	}
}
