// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tester

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const Version = "v0.0.1"

var (
	API *conch.Conch
)

var (
	rootCmd = &cobra.Command{
		Use:     "tester",
		Version: Version,
		Short:   "tester is a tool to test the conch api using recent device reports, given a database connection",
	}
)

// Root returns the root command
func Root() *cobra.Command {
	return rootCmd
}

// Execute gets this party started
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		Bail(err)
	}
}

func init() {
	initFlags()
	buildAPI()

	UserAgent = fmt.Sprintf("conch tester %s-%s", Version, util.GitRev)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Run: func(cmd *cobra.Command, args []string) {
			buildTime := util.BuildTime
			t, err := strconv.ParseInt(util.BuildTime, 10, 64)
			if err == nil {
				buildTime = TimeStr(time.Unix(t, 0))
			}

			fmt.Printf(
				"Conch API Tester %s\n"+
					"  Git Revision: %s\n"+
					"  Build Time: %s\n"+
					"  Build Host: %s\n",
				rootCmd.Version,
				util.GitRev,
				buildTime,
				util.BuildHost,
			)
		},
	})
}

func buildAPI() {
	API = &conch.Conch{BaseURL: viper.GetString("conch_api")}
	err := API.Login(
		viper.GetString("conch_user"),
		viper.GetString("conch_password"),
	)

	if err != nil {
		log.Fatalf("error logging into %s : %s", viper.GetString("conch_api"), err)
	}
	API.Debug = viper.GetBool("debug")
	API.Trace = viper.GetBool("trace")

}

func initFlags() {
	flag.String(
		"conch_api",
		"https://staging.conch.joyent.us",
		"URL of the Conch API server to test",
	)

	flag.String(
		"conch_user",
		"",
		"Conch API user name",
	)

	flag.String(
		"conch_password",
		"",
		"Password for Conch API user",
	)

	flag.String(
		"db_host",
		"localhost",
		"Database Hostname",
	)

	flag.String(
		"db_name",
		"conch",
		"Database name",
	)

	flag.String(
		"db_user",
		"conch",
		"Database username",
	)

	flag.String(
		"db_password",
		"conch",
		"Database password",
	)

	flag.Bool(
		"debug",
		false,
		"Debug mode",
	)

	flag.Bool(
		"trace",
		false,
		"Trace mode. This is super loud",
	)

	flag.String(
		"interval",
		"1 hour",
		"Interval for the database query. Resolves to \"now() - interval '1 hour'\"",
	)

	viper.SetConfigName("conch_tester")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/usr/local/etc")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("conch_tester")
	viper.AutomaticEnv()

	viper.BindPFlags(flag.CommandLine)
	flag.Parse()

	viper.ReadInConfig()

	Debug = viper.GetBool("debug")
	Trace = viper.GetBool("trace")
}
