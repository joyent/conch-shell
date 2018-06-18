// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/blang/semver"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/commands"
	"github.com/joyent/conch-shell/pkg/config"
	"github.com/joyent/conch-shell/pkg/util"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	util.UserAgent = fmt.Sprintf("conch shell v%s-%s", util.Version, util.GitRev)
	app := cli.App("conch", "Command line interface for Conch")
	app.Version("version", util.Version)

	app.Command(
		"version",
		"Get more detailed version info than --version",
		func(cmd *cli.Cmd) {
			buildTime := util.BuildTime
			t, err := strconv.ParseInt(util.BuildTime, 10, 64)
			if err == nil {
				buildTime = util.TimeStr(time.Unix(t, 0))
			}

			cmd.Action = func() {
				fmt.Printf(
					"Conch Shell v%s\n"+
						"  Git Revision: %s\n"+
						"  Build Time: %s\n"+
						"  Build Host: %s\n",
					util.Version,
					util.GitRev,
					buildTime,
					util.BuildHost,
				)
			}
		},
	)
	var (
		useJSON    = app.BoolOpt("json j", false, "Output JSON")
		configFile = app.StringOpt("config c", "~/.conch.json", "Path to config file")
		noVersion  = app.BoolOpt("no-version-check", false, "Skip Github version check")
	)

	app.Before = func() {
		if *useJSON {
			util.JSON = true
		} else {
			util.JSON = false
		}

		if !*noVersion {
			gh, err := util.LatestGithubRelease("joyent", "conch-shell")
			if err != nil {
				util.Bail(err)
			}

			if gh.SemVer.GT(semver.MustParse(util.Version)) {
				os.Stderr.WriteString(fmt.Sprintf(
					"** A new release is available! You have v%s and %s is available.\n",
					util.Version,
					gh.TagName,
				))
				os.Stderr.WriteString(fmt.Sprintf("   The changelog can be viewed via 'conch update changelog'\n\n"))
				os.Stderr.WriteString(fmt.Sprintln("   You can obtain the new release by:"))
				os.Stderr.WriteString(fmt.Sprintln("     * Running 'conch update self', which will attempt to overwrite the current application"))
				os.Stderr.WriteString(fmt.Sprintf("     * Download the new release at %s and manually install it\n\n", gh.URL))
			}
		}

		expandedPath, err := homedir.Expand(*configFile)
		if err != nil {
			util.Bail(err)
		}

		cfg, err := config.NewFromJSONFile(expandedPath)
		if err != nil {
			cfg.Path = expandedPath
		}
		util.Config = cfg

		for _, prof := range cfg.Profiles {
			if prof.Active {
				util.ActiveProfile = prof
				break
			}
		}
	}

	commands.Init(app)

	_ = app.Run(os.Args)
}
