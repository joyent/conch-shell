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

	reports := extractReportsFromDB()

	for i, report := range reports {
		log.Info(fmt.Sprintf("Processing entry %d of %d", i, len(reports)))

		_, err := API.GetDevice(report.DeviceSerial)
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

	reports := extractReportsFromDB()

	/**
	*** Submit reports to the API
	**/

	log.Info("Submitting reports")

	for i, report := range reports {
		log.Info(fmt.Sprintf("Processing entry %d of %d", i, len(reports)))
		report.Exists = true

		state, err := API.SubmitDeviceReport(report.DeviceSerial, report.Raw)

		if err != nil {
			report.Reasons = append(report.Reasons, err.Error())
			failMe(report)

			continue
		}
		report.ValidationPlanID = state.Plan.ID
		report.ValidationPlanName = state.Plan.Name

		if state.State.Status != "pass" {
			for _, r := range state.Results {
				if r.Result.Status != "pass" {
					report.Reasons = append(
						report.Reasons,
						fmt.Sprintf(
							"- %s : %s : %s\n   %s\n",
							r.Validation.Name,
							r.Result.Category,
							r.Result.Status,
							r.Result.Message,
						),
					)
				}
			}

			failMe(report)
		}
	}
}

/**
*** Grab reports from the database
*** Eventually this should be an API endpoint
**/
func extractReportsFromDB() Reports {
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

	sql := fmt.Sprintf(`select
	device_id,
	created,
	device_report_id
	from (
		select
			device_id,
			device_report_id,
			created,
			row_number() over (
				partition by device_id order by created desc
			) as result_num
			from validation_state 
			where
				created > now() - interval '%s'
				and status = 'pass'
	) foo
	where result_num = 1
	order by random()
	limit %d;`,
		viper.GetString("interval"),
		viper.GetInt("limit"),
	)
	log.Trace(sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type queryRow struct {
		deviceID string
		reportID uuid.UUID
		report   string
		created  time.Time
	}

	results := make([]queryRow, 0)

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

	reports := make(Reports, 0)

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

			reports = append(reports, report)
		}
	}

	log.Info(fmt.Sprintf("Found %d device reports to submit", len(reports)))

	log.Debug("Closing database connection")
	db.Close()

	return reports
}
