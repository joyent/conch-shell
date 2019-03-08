// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tester

import (
	"fmt"
	"strconv"
	"time"

	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	uuid "gopkg.in/satori/go.uuid.v1"
)

const (
	ServerPlanName = "Conch v1 Legacy Plan: Server"
	SwitchPlanName = "Conch v1 Legacy Plan: Switch"
)

var (
	UserAgent    string
	API          *conch.Conch
	ServerPlanID uuid.UUID
	SwitchPlanID uuid.UUID
	Validations  map[uuid.UUID]conch.Validation
)

var (
	rootCmd = &cobra.Command{
		Use:     "tester",
		Version: util.Version,
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
		log.Fatal(err)
	}
}

func init() {
	initFlags()
	buildAPI()
	prepEnv()

	UserAgent = fmt.Sprintf("conch %s-%s / API Tester", util.Version, util.GitRev)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Run: func(cmd *cobra.Command, args []string) {
			buildTime := util.BuildTime
			t, err := strconv.ParseInt(util.BuildTime, 10, 64)
			if err == nil {
				buildTime = util.TimeStr(time.Unix(t, 0))
			}

			fmt.Printf(
				"Conch %s - API Tester\n"+
					"  Git Revision: %s\n"+
					"  Build Time: %s\n"+
					"  Build Host: %s\n",
				util.Version,
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

	flag.Bool(
		"verbose",
		false,
		"Verbose logging. Less chatty than debug or trace.",
	)

	flag.String(
		"interval",
		"1 hour",
		"Interval for the database query. Resolves to \"now() - interval '1 hour'\"",
	)

	flag.Bool(
		"json",
		false,
		"Log in json format",
	)

	flag.Int(
		"limit",
		20,
		"Submit a maximum of this many reports",
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

	if viper.GetBool("trace") {
		log.SetLevel(log.TraceLevel)
	} else if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	} else if viper.GetBool("verbose") {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if viper.GetBool("json") {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func prepEnv() {
	// Find the IDs for the One True Plans
	plans, err := API.GetValidationPlans()
	if err != nil {
		log.Fatalf("error getting validation plans: %s", err)
	}
	for _, plan := range plans {
		if plan.Name == ServerPlanName {
			ServerPlanID = plan.ID
		} else if plan.Name == SwitchPlanName {
			SwitchPlanID = plan.ID
		}
	}
	if uuid.Equal(SwitchPlanID, uuid.UUID{}) {
		log.Fatalf("failed to find validation plan '%s'", SwitchPlanName)
	}

	if uuid.Equal(ServerPlanID, uuid.UUID{}) {
		log.Fatalf("failed to find validation plan '%s'", ServerPlanName)
	}

	// Build a cache of Validation names and details
	Validations = make(map[uuid.UUID]conch.Validation)
	v, err := API.GetValidations()
	if err != nil {
		log.Fatalf("error getting list of validations: '%s'", err)
	}
	for _, validation := range v {
		Validations[validation.ID] = validation
	}

}
