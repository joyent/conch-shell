// Copyright 2018 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package validation contains commands for validation related commands
package validation

import (
	"github.com/joyent/conch-shell/pkg/util"
	conch "github.com/joyent/go-conch"
	"gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
	"strconv"
)

type validationStates []conch.ValidationState

func (vs validationStates) renderTable(validationPlans []conch.ValidationPlan) {
	table := util.GetMarkdownTable()

	planNameMap := make(map[uuid.UUID]string)
	for _, vp := range validationPlans {
		planNameMap[vp.ID] = vp.Name
	}

	table.SetHeader([]string{
		"Device ID",
		"Status",
		"Completed",
		"Validation Plan",
		"Pass / Fail / Error",
	})

	for _, v := range vs {
		fails := 0
		passes := 0
		errors := 0
		for _, r := range v.Results {
			if r.Status == "pass" {
				passes++
			}
			if r.Status == "fail" {
				fails++
			}
			if r.Status == "error" {
				errors++
			}
		}
		table.Append([]string{
			v.DeviceID,
			v.Status,
			v.Completed.String(),
			planNameMap[v.ValidationPlanID],
			strconv.Itoa(passes) + " / " + strconv.Itoa(fails) + " / " + strconv.Itoa(errors),
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
		validationStates.renderTable(validationPlans)
	}
}

func getWorkspaceValidationStates(app *cli.Cmd) {
	var workspaceStrID = app.StringArg("WORKSPACE_ID", "", "The Workspace UUID")

	app.Spec = "WORKSPACE_ID"

	app.Action = func() {
		workspaceUUID, err := util.MagicWorkspaceID(*workspaceStrID)
		if err != nil {
			util.Bail(err)
		}

		var validationStates validationStates
		validationStates, err = util.API.WorkspaceValidationStates(workspaceUUID)
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
		validationStates.renderTable(validationPlans)
	}
}
