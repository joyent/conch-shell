// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package corpus

import (
	"log"
	"os"

	"github.com/joyent/conch-shell/pkg/util"
	"github.com/spf13/viper"
)

// DDP pretty prints a structure to stderr
// See Data::Printer in perl for the origin of the name.
func DDP(v interface{}) {
	if viper.GetBool("trace") {
		util.DDP(v)
	}
}

func init() {
	log.SetOutput(os.Stderr)
}
