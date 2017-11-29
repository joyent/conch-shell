// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"fmt"
	config "github.com/joyent/conch-shell/config"
	conch "github.com/joyent/go-conch"
	"github.com/mkideal/cli"
	"os"
)

type loginArgs struct {
	cli.Helper
	ApiUrl   string `cli:"url,api,a" usage:"Conch API url" prompt:"Enter the Conch API URL" dft:"http://localhost:5001"`
	User     string `cli:"u,user,username" usage:"Conch user name" prompt:"Enter Conch user name" dft:"conch"`
	Password string `pw:"p,password" usage:"Conch password" prompt:"Enter password" dft:"conch"`
}

var LoginCmd = &cli.Command{
	Name: "login",
	Desc: "Get login credentials via the API. Best use is 'conch login 2> ~/.conch.json'",
	Argv: func() interface{} { return new(loginArgs) },
	Fn: func(ctx *cli.Context) error {
		argv := &loginArgs{}
		if err := ctx.GetArgvList(argv); err != nil {
			return err
		}

		api := &conch.Conch{
			BaseUrl: argv.ApiUrl,
			User:    argv.User,
		}

		if err := api.Login(argv.Password); err != nil {
			fmt.Printf("An error occurred during login: %s\n", err)
			os.Exit(1)
		}

		if api.Session == "" {
			fmt.Println("API login did not result in session data. Not sure why.")
			os.Exit(1)
		}

		cfg := &config.ConchConfig{
			Api:     api.BaseUrl,
			User:    api.User,
			Session: api.Session,
		}

		j, err := cfg.Serialize()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stderr, j)

		return nil
	},
}
