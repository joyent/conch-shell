// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package orc contains commands pertaining to orchestration
package orc

import (
	"fmt"
	"github.com/joyent/conch-shell/pkg/util"
	conch "github.com/joyent/go-conch"
	"github.com/pkg/errors"
	"gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
)

func createOrcWorkflowStep(app *cli.Cmd) {
	var (
		nameOpt = app.StringOpt(
			"name n",
			"",
			"The string name for the new lifecycle",
		)

		retryOpt = app.BoolOpt(
			"allow-retry ar",
			false,
			"Should retries be allowed?",
		)

		maxRetriesOpt = app.IntOpt(
			"max-retries mr",
			0,
			"The maximum amount of retries that are allowed",
		)

		validationOpt = app.StringOpt(
			"validation-plan vp",
			"",
			"The UUID of the validation plan to be run when a device completes this step",
		)
	)

	app.Spec = "--name [OPTIONS]"
	app.Action = func() {
		validationID, err := uuid.FromString(*validationOpt)
		if err != nil {
			util.Bail(errors.New("'validation-plan' is not a valid UUID: " + err.Error()))
		}

		_, err = util.API.GetValidationPlan(validationID)
		if err != nil {
			util.Bail(errors.Wrap(err, "Error retrieving validation plan"))
		}

		s := &conch.OrcWorkflowStep{
			Name:             *nameOpt,
			ValidationPlanID: validationID,
			WorkflowID:       ActiveUUID,
		}

		if *retryOpt {
			s.Retry = 1
		}

		if *maxRetriesOpt != 0 {
			s.MaxRetries = *maxRetriesOpt
		}

		err = util.API.CreateOrcWorkflowStep(ActiveUUID, s)
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(s)
			return
		}

		fmt.Printf("New Workflow Step ID: %s", s.ID.String())
	}
}

func printOrcWorkflowStep(s conch.OrcWorkflowStep) {
	w, err := util.API.GetOrcWorkflow(s.WorkflowID)
	if err != nil {
		util.Bail(err)
	}

	v, err := util.API.GetValidationPlan(s.ValidationPlanID)
	if err != nil {
		util.Bail(err)
	}

	canRetry := "No"
	if s.Retry == 1 {
		canRetry = "Yes"
	}

	fmt.Printf(`Name: %s
ID: %s

Workflow: %s (%s)
Order Of Execution: %d

Validation Plan: %s (%s)

Can Be Retried: %s
Maximum Number Of Retries: %d
`,
		s.Name,
		s.ID.String(),
		w.Name,
		w.ID.String(),
		s.Order,
		v.Name,
		s.ValidationPlanID,
		canRetry,
		s.MaxRetries,
	)
}

func getOrcWorkflowStep(app *cli.Cmd) {
	app.Action = func() {
		s, err := util.API.GetOrcWorkflowStep(ActiveUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(s)
			return
		}
		printOrcWorkflowStep(s)
	}
}

func updateOrcWorkflowStep(app *cli.Cmd) {
	var (
		nameOpt = app.StringOpt(
			"name n",
			"",
			"The string name for the new lifecycle",
		)

		retryOptByUser bool
		retryOpt       = app.Bool(cli.BoolOpt{
			Name:      "allow-retry ar",
			Value:     false,
			Desc:      "Should retries be allowed?",
			SetByUser: &retryOptByUser,
		})

		maxRetriesByUser bool
		maxRetriesOpt    = app.Int(cli.IntOpt{
			Name:      "max-retries mr",
			Value:     0,
			Desc:      "The maximum amount of retries that are allowed",
			SetByUser: &maxRetriesByUser,
		})

		validationOpt = app.StringOpt(
			"validation-plan vp",
			"",
			"The UUID of the validation plan to be run when a device completes this step",
		)
	)

	app.Action = func() {
		s, err := util.API.GetOrcWorkflowStep(ActiveUUID)
		if err != nil {
			util.Bail(err)
		}

		if *nameOpt != "" {
			s.Name = *nameOpt
		}

		if retryOptByUser {
			if *retryOpt {
				s.Retry = 1
			} else {
				s.Retry = 0
			}
		}

		if maxRetriesByUser {
			s.MaxRetries = *maxRetriesOpt
		}

		if *validationOpt != "" {
			validationID, err := uuid.FromString(*validationOpt)
			if err != nil {
				util.Bail(errors.Wrap(err, "Validation ID is not a valid UUID"))
			}

			_, err = util.API.GetValidationPlan(validationID)
			if err != nil {
				util.Bail(errors.Wrap(err, "Error retrieving validation plan"))
			}

			s.ValidationPlanID = validationID
		}

		err = util.API.UpdateOrcWorkflowStep(&s)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(s)
			return
		}
		printOrcWorkflowStep(s)
	}
}
