//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package relay

import (
	"strconv"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
	conch "github.com/joyent/go-conch"
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

func getAllRelays(app *cli.Cmd) {
	app.Before = util.BuildAPIAndVerifyLogin
	app.Action = func() {
		relays, err := util.API.GetAllRelaysWithoutDevices()
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(relays)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{
			"ID",
			"Alias",
			"IP Addr",
			"SSH Port",
			"Version",
			"Created",
			"Updated",
		})

		for _, r := range relays {
			table.Append([]string{
				r.ID,
				r.Alias,
				r.IPAddr,
				strconv.Itoa(r.SSHPort),
				r.Version,
				util.TimeStr(r.Created),
				util.TimeStr(r.Updated),
			})
		}

		table.Render()

	}
}
