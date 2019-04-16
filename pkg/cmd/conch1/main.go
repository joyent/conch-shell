// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch1

import (
	"errors"
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
				if util.DisableApiVersionCheck() {
					fmt.Println("\n** API version checking is disabled. Functionality cannot be guaranteed **")
				}
			}
		},
	)

	var (
		tokenOpt = app.String(cli.StringOpt{
			Name:   "token",
			Value:  "",
			Desc:   "API token",
			EnvVar: "CONCH_TOKEN",
		})

		environmentOpt = app.String(cli.StringOpt{
			Name:   "environment env",
			Value:  "production",
			Desc:   "Specify the environment to be used: production, staging, development (provide URL in the --url parameter)",
			EnvVar: "CONCH_ENV",
		})

		urlOpt = app.String(cli.StringOpt{
			Name:   "url",
			Value:  "",
			Desc:   "If the environment is 'development', this specifies the API URL. Ignored if --environment is 'production' or 'staging'",
			EnvVar: "CONCH_URL",
		})

		useJSON         = app.BoolOpt("json j", false, "Output JSON")
		configFile      = app.StringOpt("config c", "~/.conch.json", "Path to config file")
		noVersion       = app.BoolOpt("no-version-check", false, "Does nothing. Included for backwards compatibility.") // TODO(sungo): remove back compat
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

		if *noVersion {
			fmt.Fprintf(os.Stderr, "--no-version-check is deprecated and no longer functional")
		}

		if (*profileOverride != "") && (len(*tokenOpt) > 0) {
			util.IgnoreConfig = true
			util.Token = *tokenOpt

			if len(*environmentOpt) > 0 {
				if (*environmentOpt == "development") && (len(*urlOpt) == 0) {
					util.Bail(errors.New("--url must be provied if --environment=development is set"))
				}
			}

			switch *environmentOpt {
			case "staging":
				util.BaseURL = config.StagingURL
			case "development":
				util.BaseURL = *urlOpt
			default:
				util.BaseURL = config.ProductionURL
			}

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
					if prof.Token != "" {
						util.Token = prof.Token
					}

					break
				}
			} else if prof.Active {
				util.ActiveProfile = prof
				if prof.Token != "" {
					util.Token = prof.Token
				}

				break
			}
		}

		if !util.IgnoreConfig {
			if (*profileOverride != "") && (util.ActiveProfile == nil) {
				util.Bail(fmt.Errorf("could not find a profile named '%s'", *profileOverride))
			}
		}

		// There is no way to avoid the version check, save piping stderr to
		// /dev/null.  The API is changing too much and introducing too much
		// breakage on the regular for users to stick using old versions.
		util.GithubReleaseCheck()
	}

	return app
}
