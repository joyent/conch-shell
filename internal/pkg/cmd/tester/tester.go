// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tester

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/joyent/conch-shell/pkg/util"
	_ "github.com/lib/pq"
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
	version, err := API.GetVersion()
	if err != nil {
		log.Fatalf("error retrieving API's version: %s", err)
	}
	DebugLog(fmt.Sprintf(
		"Testing %s, API %s\n\n",
		viper.GetString("conch_api"),
		version,
	))

	/**
	*** Grab reports from the database
	*** Eventually this should be an API endpoint
	**/

	DebugLog("Attempting database connection")
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s sslmode=disable",
		viper.GetString("db_user"),
		viper.GetString("db_password"),
		viper.GetString("db_hostname"),
		viper.GetString("db_name"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	DebugLog("Database connection was successful")

	sql := fmt.Sprintf(
		"select distinct on (device_id) device_id, created, report from device_report where created > now() - interval '%s' and invalid_report is null  order by device_id, created desc",
		viper.GetString("interval"),
	)
	TraceLog(sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	results := make(queryRows, 0)

	for rows.Next() {
		var row queryRow

		if err := rows.Scan(&row.deviceID, &row.created, &row.report); err != nil {
			log.Fatal(err)
		}

		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	rows.Close()

	DebugLog(fmt.Sprintf("Found %d device reports to submit", len(results)))

	DebugLog("Closing database connection")
	db.Close()

	/**
	*** Submit reports to the API
	**/

	DebugLog("Submitting reports")
	submitted := make([]*resultRow, 0)

	for i, result := range results {
		DebugLog(fmt.Sprintf("Processing entry %d of %d", i, len(results)))

		status := &resultRow{result.deviceID, false, ""}
		submitted = append(submitted, status)

		state, err := API.SubmitDeviceReport(result.deviceID, result.report)

		if err != nil {
			status.pass = false
			status.reason = err.Error()

			DebugLog(fmt.Sprintf("Error: %s", err))
			continue
		}
		if state.State.Status == "pass" {
			status.pass = true

		} else {
			status.pass = false

			msg := fmt.Sprintf("Validation plan '%s' failed:\n", state.Plan.Name)

			for _, r := range state.Results {
				if r.Result.Status != "pass" {
					submsg := fmt.Sprintf(
						"- %s : %s : %s\n   %s\n",
						r.Validation.Name,
						r.Result.Category,
						r.Result.Status,
						r.Result.Message,
					)

					msg = msg + submsg
				}
			}

			status.reason = msg
		}
	}
	DDP(submitted)

	table := util.GetMarkdownTable()
	table.SetHeader([]string{"Device", "Status", "Reason"})

	for _, s := range submitted {
		status := "FAIL"
		if s.pass {
			status = "pass"
		}
		table.Append([]string{s.deviceID, status, s.reason})
	}
	table.Render()

}

type queryRow struct {
	deviceID string
	report   string
	created  time.Time
}
type queryRows []queryRow

type resultRow struct {
	deviceID string
	pass     bool
	reason   string
}
