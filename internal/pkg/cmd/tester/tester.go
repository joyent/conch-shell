// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tester

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/************************/

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:     "test",
		Aliases: []string{"run"},
		Short:   "Run the tester",
		Run:     testAPI,
	})
}

/************************/

func testAPI(cmd *cobra.Command, args []string) {
	required := []string{
		"conch_user",
		"conch_password",
	}
	for _, r := range required {
		if viper.GetString(r) == "" {
			log.Fatalf("please provide --%s", r)
		}
	}

}
