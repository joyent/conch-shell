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
	required := []string{
		"conch_user",
		"conch_password",
	}
	for _, r := range required {
		if viper.GetString(r) == "" {
			log.Fatalf("please provide --%s", r)
		}
	}

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

	DDP(db)

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

	type queryRow struct {
		deviceID string
		report   string
		created  time.Time
	}
	type queryRows []queryRow

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

	DDP(results)
	DebugLog(fmt.Sprintf("Found %d device reports to submit", len(results)))

	DebugLog("Closing database connection")
	db.Close()

}
