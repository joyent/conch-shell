// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package workspaces

import (
	"fmt"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
)

func createSubWorkspace(app *cli.Cmd) {
	var (
		nameArg        = app.StringArg("NAME", "", "The name for the new workspace")
		descriptionOpt = app.StringOpt("description desc", "", "The description of the new workspace")
	)

	app.Spec = "NAME [OPTIONS]"

	app.Action = func() {
		sub := conch.Workspace{
			Name:        *nameArg,
			Description: *descriptionOpt,
		}

		parent := conch.Workspace{
			ID: WorkspaceUUID,
		}

		ws, err := util.API.CreateSubWorkspace(parent, sub)

		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(ws)
		} else {
			fmt.Printf("Workspace '%s' created with ID '%s'\n", ws.Name, ws.ID)
		}

	}
}
