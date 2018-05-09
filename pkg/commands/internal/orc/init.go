// +build ignore

// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package orc contains commands pertaining to orchestration
package orc

import (
	"github.com/joyent/conch-shell/pkg/util"
	"github.com/pkg/errors"
	"gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// ActiveUUID represents the detected UUID of whatever single entity is being
// processed currently
var ActiveUUID uuid.UUID

// Init loads up all the device related commands
func Init(app *cli.Cli) {
	app.Command(
		"orc",
		"Commands for dealing with orchestration",
		func(cmd *cli.Cmd) {
			cmd.Before = util.BuildAPIAndVerifyLogin

			cmd.Command(
				"lifecycles ls",
				"Commands for dealing with lifecycles",
				func(cmd *cli.Cmd) {
					cmd.Command(
						"get",
						"Get a list of all lifecycles",
						getAllOrcLifecycles,
					)

					cmd.Command(
						"create",
						"Create a new lifecycle",
						createOrcLifecycle,
					)
				},
			)

			cmd.Command(
				"lifecycle l",
				"Commands for dealing with a single lifecycle",
				func(cmd *cli.Cmd) {
					var lifecycleIDStr = cmd.StringArg(
						"ID",
						"",
						"The UUID of the lifecycle",
					)

					cmd.Spec = "ID"
					cmd.Before = func() {
						var err error
						ActiveUUID, err = util.MagicOrcLifecycleID(*lifecycleIDStr)
						if err != nil {

							util.Bail(errors.Wrap(err, "Lifecycle is not valid"))
						}
					}

					cmd.Command(
						"get",
						"Get a single lifecycle",
						getOrcLifecycle,
					)

					cmd.Command(
						"update",
						"Update a lifecycle",
						updateOrcLifecycle,
					)

					cmd.Command(
						"add-workflow aw",
						"Add a workflow to a lifecycle",
						addWorkflowToLifecycle,
					)

					cmd.Command(
						"remove-workflow rw",
						"Remove a workflow from a lifecycle",
						removeWorkflowFromLifecycle,
					)
				},
			)

			cmd.Command(
				"workflows ws",
				"Commands for dealing with workflows",
				func(cmd *cli.Cmd) {
					cmd.Command(
						"get",
						"Get a list of all workflows",
						getAllOrcWorkflows,
					)

					cmd.Command(
						"create",
						"Create a new workflow",
						createOrcWorkflow,
					)
				},
			)

			cmd.Command(
				"workflow w",
				"Commands for dealing with a single lifecycle",
				func(cmd *cli.Cmd) {
					var idStr = cmd.StringArg(
						"ID",
						"",
						"The UUID of the workflow",
					)

					cmd.Spec = "ID"
					cmd.Before = func() {
						var err error
						ActiveUUID, err = util.MagicOrcWorkflowID(*idStr)
						if err != nil {
							util.Bail(err)
						}
					}

					cmd.Command(
						"get",
						"Get a single workflow",
						getOrcWorkflow,
					)

					cmd.Command(
						"update",
						"Update a single workflow",
						updateOrcWorkflow,
					)

					cmd.Command(
						"steps",
						"Commands for dealing with workflow steps",
						func(cmd *cli.Cmd) {
							cmd.Command(
								"append",
								"Append a new step to this workflow",
								createOrcWorkflowStep,
							)
						},
					)
				},
			)

			cmd.Command(
				"step s",
				"Commands for dealing with an individual step (See 'workflow' for create",
				func(cmd *cli.Cmd) {
					var idStr = cmd.StringArg(
						"ID",
						"",
						"The UUID of the step",
					)

					cmd.Spec = "ID"
					cmd.Before = func() {
						var err error
						ActiveUUID, err = uuid.FromString(*idStr)
						if err != nil {
							util.Bail(errors.New("ID is not a valid UUID: " + err.Error()))
						}
					}

					cmd.Command(
						"get",
						"Get a single workflow step",
						getOrcWorkflowStep,
					)

					cmd.Command(
						"update",
						"Update a single workflow step",
						updateOrcWorkflowStep,
					)
				},
			)

		},
	)
}
