// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package api

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
	"io/ioutil"
	"os"
)

func get(cmd *cli.Cmd) {
	var cmdArg = cmd.StringArg("CMD", "", "The API path to GET. Must *not* include the hostname or port")
	cmd.Spec = "CMD"
	cmd.Action = func() {
		util.JSON = true
		if *cmdArg != "" {
			res, err := util.API.RawGet(*cmdArg)
			if err != nil {
				util.Bail(err)
			}
			if res == nil {
				util.Bail(errors.New("Empty response"))
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				util.Bail(err)
			}
			bodyStr := string(body)

			if res.StatusCode != 200 {
				errStr := fmt.Sprintf(
					"HTTP Error: Status: %s\nBody: %s\n",
					res.Status,
					bodyStr,
				)
				util.Bail(errors.New(errStr))
			}
			fmt.Println(bodyStr)
		}
	}
}

func deleteAPI(cmd *cli.Cmd) {
	var cmdArg = cmd.StringArg("CMD", "", "The API path to DELETE. Must *not* include the hostname or port")
	cmd.Spec = "CMD"
	cmd.Action = func() {
		util.JSON = true
		if *cmdArg != "" {
			res, err := util.API.RawDelete(*cmdArg)
			if err != nil {
				util.Bail(err)
			}
			if res == nil {
				util.Bail(errors.New("Empty response"))
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				util.Bail(err)
			}
			bodyStr := string(body)

			if res.StatusCode != 200 {
				if res.StatusCode != 204 {
					errStr := fmt.Sprintf(
						"HTTP Error: Status: %s\nBody: %s\n",
						res.Status,
						bodyStr,
					)
					util.Bail(errors.New(errStr))
				}
			}
			fmt.Println(bodyStr)

		}
	}
}

func postAPI(cmd *cli.Cmd) {
	var cmdArg = cmd.StringArg("API", "", "The API path to GET. Must *not* include the hostname or port")
	var filePathArg = cmd.StringArg("FILE", "-", "Path to a JSON file to use as the request body. '-' indicates STDIN")
	cmd.Spec = "API [FILE]"
	cmd.Action = func() {
		util.JSON = true
		var b []byte
		var err error
		if *filePathArg == "-" {
			b, err = ioutil.ReadAll(os.Stdin)
		} else {
			b, err = ioutil.ReadFile(*filePathArg)
		}
		if err != nil {
			util.Bail(err)
		}

		res, err := util.API.RawPost(*cmdArg, bytes.NewReader(b))
		if err != nil {
			util.Bail(err)
		}
		if res == nil {
			util.Bail(errors.New("Empty response"))
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			util.Bail(err)
		}
		bodyStr := string(body)

		if res.StatusCode != 200 {
			errStr := fmt.Sprintf(
				"HTTP Error: Status: %s\nBody: %s\n",
				res.Status,
				bodyStr,
			)
			util.Bail(errors.New(errStr))
		}
		fmt.Println(bodyStr)
	}
}
