// Copyright 2018 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package validation contains commands for validation related commands
package validation

import (
	"github.com/joyent/conch-shell/pkg/util"
	"gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
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

			var validationPlanID = cmd.StringArg("ID", "", "The UUID of the validation plan")

			cmd.Spec = "ID"

			cmd.Before = func() {
				util.BuildAPIAndVerifyLogin()
				var err error
				validationUUID, err = uuid.FromString(*validationPlanID)
				if err != nil {
					util.Bail(err)
				}
				return
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
			cmd.Command(
				"create",
				"Create a new validation plan",
				createValidationPlan,
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
				validationPlanUUID, err = uuid.FromString(*validationPlanID)
				if err != nil {

					util.Bail(err)
				}
				return
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
				"add-validation",
				"Associate a validation with a validation plan",
				addValidationToPlan,
			)
			cmd.Command(
				"remove-validation",
				"Remove an associated validation from a validation plan",
				removeValidationFromPlan,
			)

			cmd.Command(
				"test",
				"Test a validation plan against a given device with input data from STDIN",
				testValidationPlan,
			)
		},
	)
}
