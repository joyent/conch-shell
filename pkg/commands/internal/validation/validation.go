// Copyright 2018 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package validation contains commands for validation related commands
package validation

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
)

type validations []conch.Validation

func (vs validations) renderTable() {
	table := util.GetMarkdownTable()

	table.SetHeader([]string{"Id", "Name", "Version", "Description"})

	for _, v := range vs {
		table.Append([]string{v.ID.String(), v.Name, strconv.Itoa(v.Version), v.Description})
	}

	table.Render()
}

type validationResults []conch.ValidationResult

func (rs validationResults) renderTable() {
	table := util.GetMarkdownTable()

	table.SetHeader([]string{"Status", "Category", "Message", "Hint", "Component ID"})

	for _, r := range rs {
		table.Append([]string{r.Status, r.Category, r.Message, r.Hint, r.ComponentID})
	}

	table.Render()
}

func getValidations(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin

	app.Action = func() {
		var validations validations
		validations, err := util.API.GetValidations()
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(validations)
			return
		}
		validations.renderTable()
	}
}

func testValidation(app *cli.Cmd) {
	var deviceSerial = app.StringArg("DEVICE_ID", "", "The Device ID (serial number) to test the validation against")

	app.Spec = "DEVICE_ID"

	app.Action = func() {
		bodyBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			util.Bail(err)
		}

		body := string(bodyBytes)
		if len(body) <= 1 {
			util.Bail(errors.New("no device report provided on stdin"))
		}

		var validationResults validationResults
		validationResults, err = util.API.RunDeviceValidation(
			*deviceSerial,
			validationUUID,
			body,
		)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(validationResults)
			return
		}
		validationResults.renderTable()
	}
}
