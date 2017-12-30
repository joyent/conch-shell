// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/joyent/conch-shell/pkg/commands"
	"github.com/joyent/conch-shell/pkg/config"
	"github.com/joyent/conch-shell/pkg/util"
	homedir "github.com/mitchellh/go-homedir"
	"gopkg.in/jawher/mow.cli.v1"
	"os"
	"time"
)

var (
	Version   string
	BuildTime string
	GitRev    string
)

func main() {
	app := cli.App("conch", "Command line interface for Conch")
	app.Version("version", Version)

	app.Command(
		"version",
		"Get more detailed version info than --version",
		func(cmd *cli.Cmd) {
			cmd.Action = func() {
				fmt.Printf(
					"Conch Shell v%s\n"+
						"  Git Revision: %s\n"+
						"  Build Time: %s\n",
					Version,
					GitRev,
					BuildTime,
				)
			}
		},
	)

	var (
		use_json    = app.BoolOpt("json", false, "Output JSON")
		config_file = app.StringOpt("config c", "~/.conch.json", "Path to config file")
		pretty      = app.BoolOpt("pretty", false, "Pretty CLI output, including spinners")
	)

	app.Before = func() {
		if *use_json {
			util.JSON = true
		} else {
			util.JSON = false
		}

		util.Pretty = *pretty
		if *pretty {
			util.Spin = spinner.New(spinner.CharSets[10], 100*time.Millisecond)
			util.Spin.FinalMSG = "Complete.\n"
		}

		config_file_path, err := homedir.Expand(*config_file)
		if err != nil {
			util.Bail(err)
		}

		cfg, err := config.NewFromJsonFile(config_file_path)
		if err != nil {
			cfg.Path = config_file_path
		}
		util.Config = cfg
	}

	commands.Init(app)

	app.Run(os.Args)
}