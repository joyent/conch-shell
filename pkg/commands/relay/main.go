//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package relay

import (
	"errors"
	"regexp"
	"sort"
	"strconv"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
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
		r := conch.WorkspaceRelay{
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
		relays, err := util.API.GetAllRelays()
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

		sort.Sort(sortRelaysByUpdated(relays))
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

type sortRelaysByUpdated []conch.WorkspaceRelay

func (s sortRelaysByUpdated) Len() int {
	return len(s)
}

func (s sortRelaysByUpdated) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRelaysByUpdated) Less(i, j int) bool {
	return s[i].Updated.Before(s[j].Updated)
}

func findRelaysByName(app *cli.Cmd) {
	var (
		relays = app.StringsArg("RELAYS", nil, "List of regular expressions to match against relay IDs")
		andOpt = app.BoolOpt("and", false, "Match the list as a logical AND")
	)

	app.Spec = "[OPTIONS] RELAYS..."
	app.LongDesc = `
Takes a list of regular expressions and matches those against the IDs of all known relays.

The default behavior is to match as a logical OR but this behavior can be changed by providing the --and flag

For instance:

* "conch relay find drd" will find all relays with 'drd' in their ID. For perl folks, this is essentially 'm/drd/'
* "conch relay find '^ams-'" will find all relays with IDs that begin with 'ams-'
* "conch relay find drd '^ams-' will find all relays with IDs that contain 'drd' OR begin with 'ams-'
* "conch relay find --and drd '^ams-' will find all relays with IDs that contain 'drd' AND begin with '^ams-'`

	app.Action = func() {
		if *relays == nil {
			util.Bail(errors.New("please provide a list of regular expressions"))
		}

		// If a user for some strange reason gives us a relay name of "", the
		// cli lib will pass it on to us. That name is obviously useless so
		// let's filter it out.
		relayREs := make([]*regexp.Regexp, 0)
		for _, matcher := range *relays {
			if matcher == "" {
				continue
			}
			re, err := regexp.Compile(matcher)
			if err != nil {
				util.Bail(err)
			}
			relayREs = append(relayREs, re)
		}
		if len(relayREs) == 0 {
			util.Bail(errors.New("please provide a list of regular expressions"))
		}

		relays, err := util.API.GetAllRelays()
		if err != nil {
			util.Bail(err)
		}

		keepers := make(sortRelaysByUpdated, 0)
		for _, relay := range relays {
			matched := 0
			for _, re := range relayREs {
				if re.MatchString(relay.ID) {
					if *andOpt {
						matched++
					} else {
						keepers = append(keepers, relay)
						continue
					}
				}
			}
			if *andOpt {
				if matched == len(relayREs) {
					keepers = append(keepers, relay)
				}
			}
		}
		sort.Sort(keepers)

		if util.JSON {
			util.JSONOut(keepers)
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

		for _, r := range keepers {
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
