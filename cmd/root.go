package cmd

import (
	"github.com/mkideal/cli"
)

type GlobalArgs struct {
	Verbose    cli.Counter `cli:"v,verbose" usage:"Verbose mode (Multiple -v options increase the verbosity.)"`
	ConfigPath string      `cli:"c,config" usage:"Config file location" dft:"$HOME/.conch.json"`
	JSON       bool        `cli:"json" usage:"Results of the request operation are output in JSON. Overrides 'verbose'" dft:"false"`
}

var RootCmd = &cli.Command{
	Name:   "conch",
	Desc:   "Conch Shell - CLI for https://github.com/joyent/conch",
	Argv:   func() interface{} { return new(GlobalArgs) },
	Global: true,
	Fn: func(ctx *cli.Context) error {
		ctx.WriteUsage()
		return nil
	},
}
