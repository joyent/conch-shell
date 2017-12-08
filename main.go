// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

// BUG(sungo): Errors do not understand the --json flag

import (
	"fmt"
	"github.com/joyent/conch-shell/cmd"
	"github.com/mkideal/cli"
	"os"
)

var (
	Version   string
	BuildTime string
	GitRev    string
)

var versionCmd = &cli.Command{
	Name: "version",
	Desc: "Display version information",
	Fn: func(ctx *cli.Context) error {
		fmt.Printf(
			"Conch Shell v%s\n"+
				"  Git Revision: %s\n"+
				"  Build Time: %s\n",
			Version,
			BuildTime,
			GitRev,
		)
		return nil
	},
}

func main() {
	/*
		nesting commands looks like:
		cli.Root(
			rootCmd,
			cli.Tree(
				levelOneCmd,
				cli.Tree(
					levelTwoCmd,
					cli.Tree(
						levelThreeCmd,
					),
				),
			),
			cli.Tree(anotherLevelOneCmd),
		).Run
	*/

	if err := cli.Root(
		cmd.RootCmd,
		cli.Tree(versionCmd),
		cli.Tree(cmd.LoginCmd),
		cli.Tree(cmd.GetWorkspacesCmd),
		cli.Tree(cmd.GetWorkspaceCmd),
		cli.Tree(cmd.GetSubWorkspacesCmd),
		cli.Tree(cmd.GetWorkspaceUsersCmd),
		cli.Tree(cmd.GetWorkspaceRoomsCmd),
		cli.Tree(cmd.GetSettingsCmd),
		cli.Tree(cmd.GetSettingCmd),
		cli.Tree(cmd.GetWorkspaceDevicesCmd),
		cli.Tree(cmd.GetWorkspaceRelaysCmd),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
