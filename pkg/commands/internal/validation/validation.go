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
	"strconv"
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
