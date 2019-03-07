// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

// ddp pretty prints a structure to stderr. "Deep Data Printer"
func (c *Conch) ddp(v interface{}) {
	spew.Fdump(
		os.Stderr,
		v,
	)
}

// debugLog prints a string to stderr *if* the Debug flag is set
func (c *Conch) debugLog(out string) {
	if c.Debug {
		fmt.Fprintln(os.Stderr, out)
	}
}

// traceLog prints a string to stderr *if* the Trace flag is set
func (c *Conch) traceLog(out string) {
	if c.Trace {
		fmt.Fprintln(os.Stderr, out)
	}
}

// traceLogDDP is a convenience function for printing a data structure with a
// message, *if* the Trace flag is set
func (c *Conch) traceLogDDP(out string, v interface{}) {
	if c.Trace {
		c.traceLog(out)
		c.ddp(v)
	}
}

// debugLogDDP is a convenience function for printing a data structure with a
// message, *if* the Debug flag is set
//lint:ignore U1000 sure this will get used later
func (c *Conch) debugLogDDP(out string, v interface{}) {
	if c.Debug {
		c.debugLog(out)
		c.ddp(v)
	}
}

func init() {
	spew.Config = spew.ConfigState{
		Indent:                  "    ",
		SortKeys:                true,
		DisablePointerAddresses: true,
	}
}
