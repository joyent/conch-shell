// Copyright 2018 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package validation contains commands for validation related commands
package validation

import (
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/joyent/conch-shell/pkg/util"
)

var validationUUID uuid.UUID
var validationPlanUUID uuid.UUID

// Init loads up the commands dealing with validations and validation plans
func Init(app *cli.Cli) {
	app.Command(
		"validations vs",
		"List available validations",
		getValidations,
	)
	app.Command(
		"validation v",
		"Commands for operating on a validation",
		func(cmd *cli.Cmd) {

			var validationID = cmd.StringArg("ID", "", "The UUID of the validation")

			cmd.Spec = "ID"

			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
				var err error
				validationUUID, err = util.MagicValidationID(*validationID)
				if err != nil {
					util.Bail(err)
				}
			}

			cmd.Command(
				"test",
				"Test a validation against a given device with input data from STDIN",
				testValidation,
			)
		},
	)
	app.Command(
		"validation-plans vps",
		"Manage validation plans",
		func(cmd *cli.Cmd) {
			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
			}
			cmd.Command(
				"get",
				"List all active validation plans",
				getValidationPlans,
			)
		},
	)
	app.Command(
		"validation-plan vp",
		"Commands for operating on a validation plan",
		func(cmd *cli.Cmd) {

			var validationPlanID = cmd.StringArg("ID", "", "The UUID of the validation plan")

			cmd.Spec = "ID"

			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
				var err error
				validationPlanUUID, err = util.MagicValidationPlanID(*validationPlanID)
				if err != nil {

					util.Bail(err)
				}
			}

			cmd.Command(
				"get",
				"Get details of a validation plan",
				getValidationPlan,
			)

			cmd.Command(
				"validations",
				"Show a validation plan's associated validations",
				showValidationPlanValidations,
			)

			cmd.Command(
				"test",
				"Test a validation plan against a given device with input data from STDIN",
				testValidationPlan,
			)
		},
	)
	app.Command(
		"validation-states vss",
		"Commands for validation states",
		func(cmd *cli.Cmd) {

			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
			}

			cmd.Command(
				"device",
				"Get validation states for a device",
				getDeviceValidationStates,
			)
		},
	)
}
