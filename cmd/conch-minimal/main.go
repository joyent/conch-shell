// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"os"

	"github.com/joyent/conch-shell/pkg/cmd/conch1"
	"github.com/joyent/conch-shell/pkg/commands/api"
	"github.com/joyent/conch-shell/pkg/commands/profile"
	"github.com/joyent/conch-shell/pkg/commands/update"
)

func main() {
	app := conch1.Init()

	api.Init(app)
	profile.Init(app)
	update.Init(app)

	_ = app.Run(os.Args)
}
