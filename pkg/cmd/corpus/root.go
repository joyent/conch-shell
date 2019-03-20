// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package corpus

import (
	"fmt"

	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	uuid "gopkg.in/satori/go.uuid.v1"
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
		Use:     "report-corpus",
		Version: util.Version,
		Short:   "report-corpus dumps out device reports relevant to every known sku",
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

	UserAgent = fmt.Sprintf("conch %s-%s / Report Corpus", util.Version, util.GitRev)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(
				"Conch %s - API Tester\n"+
					"  Git Revision: %s\n",
				util.Version,
				util.GitRev,
			)
		},
	})
}

func initFlags() {
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
		"data_directory",
		"",
		"A directory full of device reports",
	)

	viper.SetConfigName("conch_corpus")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("/usr/local/etc")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("conch_corpus")
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

	if viper.GetString("data_directory") == "" {
		log.Fatal("Please provide the data_directory parameter")
	}
}
