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
	"sort"
	"strconv"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
)

func renderTableValidations(vs conch.Validations, showDeactivated bool) {
	sort.Sort(vs)
	table := util.GetMarkdownTable()

	if showDeactivated {
		table.SetHeader([]string{"Id", "Name", "Version", "Active", "Description"})
	} else {
		table.SetHeader([]string{"Id", "Name", "Version", "Description"})
	}

	for _, v := range vs {
		if showDeactivated {
			active := ""
			if v.Deactivated.IsZero() {
				active = "X"
			}
			table.Append([]string{
				v.ID.String(),
				v.Name,
				strconv.Itoa(v.Version),
				active,
				v.Description,
			})
		} else {
			table.Append([]string{
				v.ID.String(),
				v.Name,
				strconv.Itoa(v.Version),
				v.Description,
			})
		}
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
	var showDeactivated = app.BoolOpt("deactivated", false, "Show deactivated (old) versions of validations")

	app.Action = func() {
		validations, err := util.API.GetValidations()
		if err != nil {
			util.Bail(err)
		}

		if !*showDeactivated {
			v := make(conch.Validations, 0)
			for _, validation := range validations {
				if validation.Deactivated.IsZero() {
					v = append(v, validation)
				}
			}
			validations = v
		}

		if util.JSON {
			util.JSONOut(validations)
			return
		}

		renderTableValidations(validations, *showDeactivated)
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
