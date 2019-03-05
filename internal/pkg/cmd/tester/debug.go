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
)

// Debug triggers lots and lots of output to stderr for use in debugging
var Debug bool

// Trace causes the API struct to activate http tracing
var Trace bool

// DDP pretty prints a structure to stderr
// See Data::Printer in perl for the origin of the name.
func DDP(v interface{}) {
	if Trace {
		spew.Fdump(
			os.Stderr,
			v,
		)
	}
}

// TraceLog prints a string to stderr *if* the Trace flag is set
func TraceLog(out string) {
	if Trace {
		log.Println(out)
	}
}

// DebugLog prints a string to stderr *if* the Debug flag is set
func DebugLog(out string) {
	if Debug {
		log.Println(out)
	}
}

// TraceLogDDP is a convenience function for printing a data structure with a
// message, *if* the Trace flag is set
func TraceLogDDP(out string, v interface{}) {
	if Trace {
		TraceLog(out)
		DDP(v)
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
