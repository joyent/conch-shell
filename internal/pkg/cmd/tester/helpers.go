// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tester

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

var (
	// UserAgent will be used as the http user agent when making API calls
	UserAgent string

	// JSON tells us if we should output JSON
	JSON bool
)

// DateFormat should be used in date formatting calls to ensure uniformity of
// output
const DateFormat = "2006-01-02 15:04:05 -0700 MST"

// Bail is a --json aware way of dying
func Bail(err error) {
	TraceLogDDP("Bailing out", err)

	msg := err.Error()

	if JSON {
		j, _ := json.Marshal(struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}{
			true,
			msg,
		})

		fmt.Println(string(j))
	} else {
		fmt.Fprintln(os.Stderr, msg)
	}
	os.Exit(1)
}

// GetMarkdownTable returns a tablewriter configured to output markdown
// compatible text
func GetMarkdownTable() (table *tablewriter.Table) {
	table = tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	return table
}

// JSONOut marshals an interface to indented JSON
func JSONOut(thingy interface{}) {
	TraceLogDDP("Outputting the JSON version of: ", thingy)

	j, err := json.MarshalIndent(thingy, "", "     ")

	if err != nil {
		Bail(err)
	}

	fmt.Println(string(j))
}

// TimeStr ensures that all Times are formatted using .Local() and DateFormat
func TimeStr(t time.Time) string {
	return t.Local().Format(DateFormat)
}
