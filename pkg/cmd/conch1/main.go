// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch1

import (
	"fmt"
	"os"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/config"
	"github.com/joyent/conch-shell/pkg/util"

	homedir "github.com/mitchellh/go-homedir"
)

func Init() *cli.Cli {
	util.UserAgent = fmt.Sprintf("conch shell v%s-%s", util.Version, util.GitRev)

	app := cli.App("conch", "Command line interface for Conch")

	app.Version("version", util.Version)

	app.Command(
		"version",
		"Get more detailed version info than --version",
		func(cmd *cli.Cmd) {

			cmd.Action = func() {
				fmt.Printf(
					"Conch Shell v%s\n"+
						"  Git Revision: %s\n"+
						"  Requires API version: >= %s and < %s\n",
					util.Version,
					util.GitRev,
					conch.MinimumAPIVersion,
					conch.BreakingAPIVersion,
				)
				if util.NoApiVersionCheck {
					fmt.Println("\n** API version checking is disabled. Functionality cannot be guaranteed **")
				}
			}
		},
	)
	var (
		useJSON         = app.BoolOpt("json j", false, "Output JSON")
		configFile      = app.StringOpt("config c", "~/.conch.json", "Path to config file")
		noVersion       = app.BoolOpt("no-version-check", false, "Skip Github version check")
		profileOverride = app.StringOpt("profile p", "", "Override the active profile")
		debugMode       = app.BoolOpt("debug", false, "Debug mode")
		traceMode       = app.BoolOpt("trace", false, "Trace http requests. Warning: this is super loud")
	)

	app.Before = func() {
		util.Debug = *debugMode
		util.Trace = *traceMode

		if *useJSON {
			util.JSON = true
		} else {
			util.JSON = false
		}

		expandedPath, err := homedir.Expand(*configFile)
		if err != nil {
			util.Bail(err)
		}

		cfg, _ := config.NewFromJSONFile(expandedPath)
		cfg.Path = expandedPath
		util.Config = cfg

		for _, prof := range cfg.Profiles {
			if *profileOverride != "" {
				if prof.Name == *profileOverride {
					util.ActiveProfile = prof
					break
				}
			} else if prof.Active {
				util.ActiveProfile = prof
				break
			}
		}
		if (*profileOverride != "") && (util.ActiveProfile == nil) {
			util.Bail(fmt.Errorf("could not find a profile named '%s'", *profileOverride))
		}

		checkVersion := true
		if *noVersion || cfg.SkipVersionCheck {
			checkVersion = false
		}

		if checkVersion {
			gh, err := util.LatestGithubRelease()
			if (err != nil) && (err != util.ErrNoGithubRelease) {
				util.Bail(err)
			}
			if gh.Upgrade {
				os.Stderr.WriteString(fmt.Sprintf(`
A new release is available! You have v%s but %s is available.
The changelog can be viewed via 'conch update changelog'

You can obtain the new release by:
  * Running 'conch update self', which will attempt to overwrite the current application
  * Manually download the new release at %s

`,
					util.Version,
					gh.TagName,
					gh.URL,
				))
			}
		}

	}

	return app
}
