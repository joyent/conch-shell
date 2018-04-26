// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package orc contains commands pertaining to orchestration
package orc

import (
	"errors"
	"fmt"
	"github.com/joyent/conch-shell/pkg/util"
	conch "github.com/joyent/go-conch"
	"gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
	"strconv"
)

func printOrcLifecycle(l conch.OrcLifecycle) {
	locked := "No"
	if l.Locked == 1 {
		locked = "Yes"
	}
	fmt.Printf(`
ID: %s
Name: %s
Locked: %s
Role ID: %s
Created: %s
Updated: %s

Workflow Plan:
`,
		l.ID.String(),
		l.Name,
		locked,
		l.RoleID.String(),
		l.Created.String(),
		l.Updated.String(),
	)

	for idx, id := range l.Plan {
		w, err := util.API.GetOrcWorkflow(id)
		if err != nil {
			util.Bail(err)
		}
		fmt.Printf("    * [ %d ] %s (%s)\n", idx, w.Name, w.ID)
	}

	fmt.Println()
}

func getAllOrcLifecycles(app *cli.Cmd) {
	app.Action = func() {
		ls, err := util.API.GetOrcLifecycles()
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(ls)
			return
		}
		table := util.GetMarkdownTable()
		table.SetHeader([]string{
			"ID",
			"Name",
			"Locked",
			"Role ID",
			"Workflow Count",
			"Created",
			"Updated",
		})

		for _, l := range ls {
			locked := ""
			if l.Locked == 1 {
				locked = "X"
			}

			table.Append([]string{
				l.ID.String(),
				l.Name,
				locked,
				l.RoleID.String(),
				strconv.Itoa(len(l.Plan)),
				l.Created.String(),
				l.Updated.String(),
			})
		}
		table.Render()
	}
}

func createOrcLifecycle(app *cli.Cmd) {
	var (
		nameOpt = app.StringOpt(
			"name n",
			"",
			"The string name for the new lifecycle",
		)

		roleOpt = app.StringOpt(
			"role r",
			"",
			"The UUID of a device role",
		)

		lockedOpt = app.BoolOpt(
			"locked",
			false,
			"The locked status",
		)
	)

	app.Spec = "--name --role [OPTIONS]"
	app.Action = func() {
		roleID, err := uuid.FromString(*roleOpt)
		if err != nil {
			util.Bail(errors.New("'role' is not a valid UUID: " + err.Error()))
		}

		l := &conch.OrcLifecycle{Name: *nameOpt, RoleID: roleID}
		if *lockedOpt {
			l.Locked = 1
		}

		err = util.API.CreateOrcLifecycle(l)
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(l)
			return
		}

		printOrcLifecycle(*l)
	}
}
func getOrcLifecycle(app *cli.Cmd) {
	app.Action = func() {
		l, err := util.API.GetOrcLifecycle(ActiveUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(l)
			return
		}
		printOrcLifecycle(l)
	}
}

func addWorkflowToLifecycle(app *cli.Cmd) {
	var (
		wIDOpt = app.StringArg("WORKFLOW", "", "Name or UUID of the workflow to add")
	)

	app.Spec = "WORKFLOW"

	app.Action = func() {
		wID, err := util.MagicOrcWorkflowID(*wIDOpt)
		if err != nil {
			util.Bail(errors.New("'workflow' is not valid: " + err.Error()))
		}

		l, err := util.API.GetOrcLifecycle(ActiveUUID)
		if err != nil {
			util.Bail(err)
		}

		w, err := util.API.GetOrcWorkflow(wID)
		if err != nil {
			util.Bail(err)
		}

		err = util.API.AddWorkflowToOrcLifecycle(&l, w.ID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(l)
		}
		printOrcLifecycle(l)
	}

}

func removeWorkflowFromLifecycle(app *cli.Cmd) {
	var (
		wIDOpt = app.StringArg("WORKFLOW", "", "Name or UUID of the workflow to add")
	)

	app.Spec = "WORKFLOW"

	app.Action = func() {
		wID, err := util.MagicOrcWorkflowID(*wIDOpt)
		if err != nil {
			util.Bail(errors.New("'workflow' is not valid: " + err.Error()))
		}

		l, err := util.API.GetOrcLifecycle(ActiveUUID)
		if err != nil {
			util.Bail(err)
		}

		w, err := util.API.GetOrcWorkflow(wID)
		if err != nil {
			util.Bail(err)
		}

		err = util.API.RemoveWorkflowFromOrcLifecycle(&l, w.ID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(l)
		}
		printOrcLifecycle(l)
	}
}

func updateOrcLifecycle(app *cli.Cmd) {
	var (
		lockedOptByUser bool
		roleOptByUser   bool
		nameOpt         = app.StringOpt(
			"name n",
			"",
			"The string name for the lifecycle",
		)

		lockedOpt = app.Bool(cli.BoolOpt{
			Name:      "locked",
			Value:     false,
			Desc:      "The locked status for the lifecycle. (Use locked=false to unset)",
			SetByUser: &lockedOptByUser,
		})

		roleOpt = app.String(cli.StringOpt{
			Name:      "role-id r",
			Value:     "",
			Desc:      "The role ID assigned to the lifecycle",
			SetByUser: &roleOptByUser,
		})
	)

	app.Action = func() {
		l, err := util.API.GetOrcLifecycle(ActiveUUID)
		if err != nil {
			util.Bail(err)
		}
		if *nameOpt != "" {
			l.Name = *nameOpt
		}

		if lockedOptByUser {
			if *lockedOpt {
				l.Locked = 1
			} else {
				l.Locked = 0
			}
		}
		if roleOptByUser {
			roleID, err := uuid.FromString(*roleOpt)
			if err != nil {
				util.Bail(err)
			}
			l.RoleID = roleID
		}

		err = util.API.UpdateOrcLifecycle(&l)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(l)
			return
		}
		printOrcLifecycle(l)
	}
}
