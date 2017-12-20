// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"os"
	"github.com/joyent/conch-shell/config"
	"github.com/joyent/conch-shell/util"
	"github.com/joyent/conch-shell/workspaces"
	"github.com/joyent/conch-shell/user"
	"github.com/joyent/conch-shell/reports"
	"github.com/joyent/conch-shell/devices"
	homedir "github.com/mitchellh/go-homedir"
	"gopkg.in/jawher/mow.cli.v1"
)


var (
	Version   string
	BuildTime string
	GitRev    string
)

func main() {
	app := cli.App("conch", "Command line interface for Conch")
	app.Version("version", Version)

	var (
		use_json    = app.BoolOpt("json", false, "Output JSON")
		config_file = app.StringOpt("config c", "~/.conch.json", "Path to config file")
	)

	app.Before = func() {
		if *use_json {
			util.JSON = true
		} else {
			util.JSON = false
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


	workspaces.Init(app)
	devices.Init(app)
	user.Init(app)
	reports.Init(app)

	app.Command("login", "Log in", user.Login)

	app.Run(os.Args)
}


