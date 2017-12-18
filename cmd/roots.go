// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"github.com/mkideal/cli"
)

type GlobalArgs struct {
	Verbose    cli.Counter `cli:"v,verbose" usage:"Verbose mode (Multiple -v options increase the verbosity.)"`
	ConfigPath string      `cli:"c,config" usage:"Config file location" dft:"$HOME/.conch.json"`
	JSON       bool        `cli:"json" usage:"Results of the request operation are output in JSON. Overrides 'verbose'" dft:"false"`
}
type EmptyArgs struct{}

func writeUsage(ctx *cli.Context) error {
	ctx.WriteUsage()
	return nil
}

var RootCmd = &cli.Command{
	Name:   "conch",
	Desc:   "Conch Shell - CLI for https://github.com/joyent/conch",
	Argv:   func() interface{} { return new(GlobalArgs) },
	Global: true,
	Fn:     writeUsage,
}

