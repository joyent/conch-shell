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
	"github.com/joyent/go-conch"
	"gopkg.in/jawher/mow.cli.v1"
	"strconv"
)

func getAllOrcWorkflows(app *cli.Cmd) {
	app.Action = func() {
		ws, err := util.API.GetOrcWorkflows()
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(ws)
			return
		}
		table := util.GetMarkdownTable()
		table.SetHeader([]string{
			"ID",
			"Name",
			"Locked",
			"Preflight",
			"Step Count",
			"Created",
			"Updated",
		})
		for _, w := range ws {
			locked := ""
			if w.Locked == 1 {
				locked = "X"
			}

			preflight := ""
			if w.Preflight == 1 {
				preflight = "X"
			}

			table.Append([]string{
				w.ID.String(),
				w.Name,
				locked,
				preflight,
				strconv.Itoa(len(w.Steps)),
				w.Created.String(),
				w.Updated.String(),
			})

		}
		table.Render()
	}
}

func createOrcWorkflow(app *cli.Cmd) {
	var (
		nameOpt = app.StringOpt(
			"name n",
			"",
			"The string name for the new workflow",
		)

		lockedOpt = app.BoolOpt(
			"locked",
			false,
			"The locked status for the new workflow",
		)

		preflightOpt = app.BoolOpt(
			"preflight",
			false,
			"The preflight status for the new workflow",
		)
	)

	app.Spec = "--name [OPTIONS]"
	app.Action = func() {

		w := &conch.OrcWorkflow{Name: *nameOpt}
		if *lockedOpt {
			w.Locked = 1
		}

		if *preflightOpt {
			w.Preflight = 1
		}

		err := util.API.CreateOrcWorkflow(w)
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(w)
			return
		}

		fmt.Printf("New Workflow ID: %s", w.ID)
	}
}

func printOrcWorkflow(w conch.OrcWorkflow) {
	locked := "No"
	if w.Locked == 1 {
		locked = "Yes"
	}
	preflight := "No"
	if w.Preflight == 1 {
		preflight = "Yes"
	}

	fmt.Printf(`Name: %s
ID: %s
Locked: %s
Preflight: %s

Created: %s
Updated: %s

Steps:
`,
		w.Name,
		w.ID.String(),
		locked,
		preflight,
		w.Created,
		w.Updated,
	)

	for _, id := range w.Steps {
		s, err := util.API.GetOrcWorkflowStep(id)
		if err != nil {
			util.Bail(err)
		}
		fmt.Printf(" * [ %d ] - %s (%s)\n", s.Order, s.Name, s.ID.String())
	}
	fmt.Println()
}

func getOrcWorkflow(app *cli.Cmd) {
	app.Action = func() {
		w, err := util.API.GetOrcWorkflow(ActiveUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(w)
			return
		}
		printOrcWorkflow(w)
	}
}

func updateOrcWorkflow(app *cli.Cmd) {
	var (
		lockedOptByUser    bool
		preflightOptByUser bool
		nameOpt            = app.StringOpt(
			"name n",
			"",
			"The string name for the new workflow",
		)

		lockedOpt = app.Bool(cli.BoolOpt{
			Name:      "locked",
			Value:     false,
			Desc:      "The locked status for the new workflow. (Use locked=false to unset)",
			SetByUser: &lockedOptByUser,
		})

		preflightOpt = app.Bool(cli.BoolOpt{
			Name:      "preflight",
			Value:     false,
			Desc:      "The preflight status for the new workflow. (Use preflight=false to unset)",
			SetByUser: &preflightOptByUser,
		})
	)

	app.Action = func() {
		w, err := util.API.GetOrcWorkflow(ActiveUUID)
		if err != nil {
			util.Bail(err)
		}
		if *nameOpt != "" {
			w.Name = *nameOpt
		}

		if lockedOptByUser {
			if *lockedOpt {
				w.Locked = 1
			} else {
				w.Locked = 0
			}
		}

		if preflightOptByUser {
			if *preflightOpt {
				w.Preflight = 1
			} else {
				w.Preflight = 0
			}
		}

		err = util.API.UpdateOrcWorkflow(&w)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(w)
			return
		}
		printOrcWorkflow(w)
	}
}
