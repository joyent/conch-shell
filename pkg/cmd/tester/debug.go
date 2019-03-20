// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tester

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

// DDP pretty prints a structure to stderr
// See Data::Printer in perl for the origin of the name.
func DDP(v interface{}) {
	if viper.GetBool("trace") {
		spew.Fdump(
			os.Stderr,
			v,
		)
	}
}

func init() {
	spew.Config = spew.ConfigState{
		Indent:                  "    ",
		SortKeys:                true,
		DisablePointerAddresses: true,
	}

	log.SetOutput(os.Stderr)
}
