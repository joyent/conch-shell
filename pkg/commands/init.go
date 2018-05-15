// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package commands is the parent that loads up the full command set
package commands

import (
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/commands/internal/admin"
	"github.com/joyent/conch-shell/pkg/commands/internal/devices"
	"github.com/joyent/conch-shell/pkg/commands/internal/hardware"
	"github.com/joyent/conch-shell/pkg/commands/internal/profile"
	"github.com/joyent/conch-shell/pkg/commands/internal/relay"
	"github.com/joyent/conch-shell/pkg/commands/internal/update"
	"github.com/joyent/conch-shell/pkg/commands/internal/user"
	"github.com/joyent/conch-shell/pkg/commands/internal/validation"
	"github.com/joyent/conch-shell/pkg/commands/internal/workspaces"
)

// Init loads up all the commands
func Init(app *cli.Cli) {
	admin.Init(app)
	devices.Init(app)
	hardware.Init(app)
	profile.Init(app)
	relay.Init(app)
	user.Init(app)
	workspaces.Init(app)
	validation.Init(app)
	update.Init(app)
}
