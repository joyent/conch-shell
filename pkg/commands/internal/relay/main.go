//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package relay

import (
	"github.com/joyent/conch-shell/pkg/util"
	conch "github.com/joyent/go-conch"
	"gopkg.in/jawher/mow.cli.v1"
)

func register(app *cli.Cmd) {
	var (
		versionOpt = app.StringOpt("version", "", "The version of the relay")
		sshPortOpt = app.IntOpt("ssh_port port", 22, "The SSH port for the relay")
		ipAddrOpt  = app.StringOpt("ipaddr ip", "", "The IP address for the relay")
		aliasOpt   = app.StringOpt("alias name", "", "The alias for the relay")
	)

	app.Spec = "--version [OPTIONS]"

	app.Action = func() {
		r := conch.Relay{
			ID:      RelaySerial,
			Version: *versionOpt,
			IPAddr:  *ipAddrOpt,
			Alias:   *aliasOpt,
			SSHPort: *sshPortOpt,
		}

		if err := util.API.RegisterRelay(r); err != nil {
			util.Bail(err)
		}
	}
}
