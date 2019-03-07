// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tester

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/joyent/conch-shell/pkg/util"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	uuid "gopkg.in/satori/go.uuid.v1"
)

var FailedCount = 0

type Report struct {
	Raw                string
	DeviceSerial       string
	ValidationPlanID   uuid.UUID
	ValidationPlanName string
	Parsed             map[string]interface{}
	Passed             bool
	Reasons            []string
	Exists             bool
}

type Reports []Report

type queryRow struct {
	deviceID string
	reportID uuid.UUID
	report   string
	created  time.Time
}
type queryRows []queryRow

type resultRow struct {
	deviceID string
	pass     bool
	reason   string
}

/************************/

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "destructive",
		Short: "Run the tester, in a potentially-destructive mode",
		Long:  "Submits the reports to /device/:id as would a real client. This updates many database tables and could possibly destroy or break live data. Devices do not have to exist in the database before these tests can run.",
		Run:   destructiveTest,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:     "test",
		Aliases: []string{"run"},
		Short:   "Run the tester, with no side effects",
		Long:    "Submits the reports to the validation endpoints, running the validations in a stateless mode. No data will be written to the database and all the database munging code in the device report processing code will NOT be exercised. This also requires that the device exist already in the database.",
		Run:     nonbindingTest,
	})

}

/************************/

func failMe(r Report) {
	FailedCount++

	logger := log.WithFields(log.Fields{
		"device":               r.DeviceSerial,
		"validation_plan_name": r.ValidationPlanName,
		"server":               viper.GetString("conch_api"),
	})

	if !r.Exists {
		logger.Error("device does not exist")
		return
	}

	logger.Error(strings.Join(r.Reasons, " || "))
}

func nonbindingTest(cmd *cobra.Command, args []string) {
	version, err := API.GetVersion()
	if err != nil {
		log.Fatalf("error retrieving API's version: %s", err)
	}
	log.Info(fmt.Sprintf(
		"Testing %s, API %s",
		viper.GetString("conch_api"),
		version,
	))

	reports := parseReports(extractReportsFromDB())

	for serial, report := range reports {

		_, err := API.GetDevice(serial)
		if err != nil {
			failMe(report)
			continue
		}
		report.Exists = true

		results, err := API.RunDeviceValidationPlan(
			report.DeviceSerial,
			report.ValidationPlanID,
			report.Raw,
		)

		if err != nil {
			report.Reasons = append(report.Reasons, fmt.Sprintf("%s", err))
			failMe(report)
			continue
		}

		validationPassed := true
		for _, result := range results {
			validationName := "[unknown]"
			if val, ok := Validations[result.ValidationID]; ok {
				validationName = val.Name
			}

			if result.Status != "pass" {
				validationPassed = false
				report.Passed = false
				report.Reasons = append(
					report.Reasons,
					fmt.Sprintf("%s : %s", validationName, result.Message),
				)
			}
		}
		if !validationPassed {
			failMe(report)
		}
	}

	log.Infof("Of %d results, %d failed", len(reports), FailedCount)
}

/************************/

func destructiveTest(cmd *cobra.Command, args []string) {
	version, err := API.GetVersion()
	if err != nil {
		log.Fatalf("error retrieving API's version: %s", err)
	}
	log.Info(fmt.Sprintf(
		"Testing %s, API %s",
		viper.GetString("conch_api"),
		version,
	))

	results := extractReportsFromDB()

	/**
	*** Submit reports to the API
	**/

	log.Info("Submitting reports")
	submitted := make([]*resultRow, 0)

	for i, result := range results {
		log.Info(fmt.Sprintf("Processing entry %d of %d", i, len(results)))

		status := &resultRow{result.deviceID, false, ""}
		submitted = append(submitted, status)

		state, err := API.SubmitDeviceReport(result.deviceID, result.report)

		if err != nil {
			status.pass = false
			status.reason = err.Error()

			log.WithFields(log.Fields{
				"error":  status.reason,
				"device": result.deviceID,
			}).Warn("error in device report submission")

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

/**
*** Grab reports from the database
*** Eventually this should be an API endpoint
**/
func extractReportsFromDB() queryRows {
	log.Debug("Attempting database connection")
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

	log.Debug("Database connection was successful")

	sql := fmt.Sprintf(
		"select distinct on (device_id) device_id, created, device_report_id from validation_state where created > now() - interval '%s' and status = 'pass'",
		viper.GetString("interval"),
	)
	log.Trace(sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	results := make(queryRows, 0)

	for rows.Next() {
		var row queryRow

		if err := rows.Scan(&row.deviceID, &row.created, &row.reportID); err != nil {
			log.Fatal(err)
		}

		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	rows.Close()

	resultsWithReports := make(queryRows, 0)
	for _, row := range results {
		sql = fmt.Sprintf(
			"select report from device_report where id = '%s'",
			row.reportID,
		)
		log.Trace(sql)

		rows, err = db.Query(sql)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {

			if err := rows.Scan(&row.report); err != nil {
				log.Fatal(err)
			}

			resultsWithReports = append(resultsWithReports, row)
		}
	}

	log.Info(fmt.Sprintf("Found %d device reports to submit", len(resultsWithReports)))

	log.Debug("Closing database connection")
	db.Close()

	return resultsWithReports
}

func parseReports(q queryRows) map[string]Report {
	reports := make(map[string]Report)

	for _, row := range q {
		log.Debug("Parsing " + row.deviceID)

		report := Report{
			DeviceSerial:       row.deviceID,
			ValidationPlanID:   ServerPlanID,
			ValidationPlanName: ServerPlanName,
			Raw:                row.report,
		}

		if err := json.Unmarshal([]byte(row.report), &report.Parsed); err != nil {
			log.Printf(
				"Report for device '%s' failed to parse: %s",
				row.deviceID,
				err.Error(),
			)
			continue
		}

		if val, ok := report.Parsed["device_type"]; ok {
			deviceType := val.(string)
			if deviceType == "switch" {
				report.ValidationPlanID = SwitchPlanID
				report.ValidationPlanName = SwitchPlanName
			}
		}

		reports[report.DeviceSerial] = report
	}

	return reports
}
