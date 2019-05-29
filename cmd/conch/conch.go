// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"os"

	"github.com/joyent/conch-shell/pkg/cmd/conch1"
	"github.com/joyent/conch-shell/pkg/commands/admin"
	"github.com/joyent/conch-shell/pkg/commands/api"
	"github.com/joyent/conch-shell/pkg/commands/datacenter"
	"github.com/joyent/conch-shell/pkg/commands/devices"
	"github.com/joyent/conch-shell/pkg/commands/global"
	"github.com/joyent/conch-shell/pkg/commands/hardware"
	"github.com/joyent/conch-shell/pkg/commands/profile"
	"github.com/joyent/conch-shell/pkg/commands/rack"
	"github.com/joyent/conch-shell/pkg/commands/relay"
	"github.com/joyent/conch-shell/pkg/commands/update"
	"github.com/joyent/conch-shell/pkg/commands/user"
	"github.com/joyent/conch-shell/pkg/commands/validation"
	"github.com/joyent/conch-shell/pkg/commands/workspaces"
)

func main() {
	app := conch1.Init()

	api.Init(app)
	admin.Init(app)
	datacenter.Init(app)
	devices.Init(app)
	global.Init(app)
	hardware.Init(app)
	profile.Init(app)
	rack.Init(app)
	relay.Init(app)
	user.Init(app)
	workspaces.Init(app)
	validation.Init(app)
	update.Init(app)

	_ = app.Run(os.Args)
}
