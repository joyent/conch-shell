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
	"sort"
	"strings"

	"github.com/dghubble/sling"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	uuid "gopkg.in/satori/go.uuid.v1"
)

var FailedCount = 0

type Report struct {
	ID                 uuid.UUID
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

type mmField struct {
	Short bool   `json:"short,omitempty"`
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
}

type mmAttachment struct {
	Pretext  string    `json:"pretext,omitempty"`
	Text     string    `json:"text,omitempty"`
	Title    string    `json:"title,omitempty"`
	Color    string    `json:"color,omitempty"`
	Fallback string    `json:"fallback,omitempty"`
	Fields   []mmField `json:"fields,omitempty"`
}

type mmPayload struct {
	Text        string         `json:"text,omitempty"`
	Attachments []mmAttachment `json:"attachments,omitempty"`
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

func failMe(r Report, destructive bool) {
	FailedCount++

	sort.Strings(r.Reasons)

	submitType := "Validations Only"
	if destructive {
		submitType = "Full POST"
	}

	logger := log.WithFields(log.Fields{
		"device":               r.DeviceSerial,
		"report_id":            r.ID,
		"validation_plan_name": r.ValidationPlanName,
		"server":               viper.GetString("conch_api"),
		"submission_type":      submitType,
	})

	if !r.Exists {
		logger.Error("device does not exist")
		return
	}

	logger.Error(strings.Join(r.Reasons, " || "))

	/***/

	if !viper.GetBool("mattermost") {
		return
	}

	var msg string

	if !r.Exists {
		msg = fmt.Sprintf("Device %s does not exist in target API", r.DeviceSerial)
		r.Reasons = []string{msg}
	} else {
		msg = fmt.Sprintf(
			"%s failed (%s): %s",
			r.DeviceSerial,
			submitType,
			strings.Join(r.Reasons, " || "),
		)
	}
	payload := mmPayload{
		Attachments: []mmAttachment{
			{
				Color:    "#FF0000",
				Fallback: msg,
				Fields: []mmField{
					{
						Title: "API Server",
						Value: viper.GetString("conch_api"),
					},
					{
						Title: "Device ID",
						Value: r.DeviceSerial,
					},
					{
						Title: "Report ID",
						Value: r.ID.String(),
					},
					{
						Title: "Validation Plan",
						Value: r.ValidationPlanName,
					},
					{
						Title: "Submission Type",
						Value: submitType,
					},
					{
						Title: "Failure Reasons",
						Value: "* " + strings.Join(r.Reasons, "\n* "),
						Short: false,
					},
				},
			},
		},
	}

	sendToMM(payload)
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
			failMe(report, false)
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
			failMe(report, false)
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
					fmt.Sprintf(
						"%s : %s : %s -> %s",
						validationName,
						result.Category,
						result.Status,
						result.Message,
					),
				)
			}
		}
		if !validationPassed {
			failMe(report, false)
		}
	}

	msg := fmt.Sprintf(
		"Submitted %d reports to %s (validations only). %d failed",
		len(reports),
		viper.GetString("conch_api"),
		FailedCount,
	)

	log.Info(msg)
	sendToMM(mmPayload{
		Text: msg,
	})
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
			failMe(report, true)

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
							"%s : %s : %s -> %s",
							r.Validation.Name,
							r.Result.Category,
							r.Result.Status,
							r.Result.Message,
						),
					)
				}
			}

			failMe(report, true)
		}
	}

	msg := fmt.Sprintf(
		"Submitted %d reports to %s (full report process). %d failed",
		len(reports),
		viper.GetString("conch_api"),
		FailedCount,
	)

	log.Info(msg)
	sendToMM(mmPayload{
		Text: msg,
	})
}

/**
*** Grab reports from the database
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
	foo.device_id,
	foo.device_report_id,
	dr.report
	from (
		select
			device_id,
			device_report_id,
			row_number() over (
				partition by device_id order by created desc
			) as result_num
			from validation_state 
			where
				created > now() - interval '%s'
				and status = 'pass'
	) foo
	join device_report dr on dr.id = foo.device_report_id
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

	reports := make(Reports, 0)

	for rows.Next() {
		report := Report{
			ValidationPlanID:   ServerPlanID,
			ValidationPlanName: ServerPlanName,
		}

		if err := rows.Scan(&report.DeviceSerial, &report.ID, &report.Raw); err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal([]byte(report.Raw), &report.Parsed); err != nil {
			log.Errorf(
				"Report for device '%s' failed to parse: %s",
				report.DeviceSerial,
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
	rows.Close()

	log.Info(fmt.Sprintf("Found %d device reports to submit", len(reports)))

	log.Debug("Closing database connection")
	db.Close()

	return reports
}

func sendToMM(payload mmPayload) {
	_, err := sling.New().Set("User-Agent", UserAgent).
		Post(viper.GetString("mattermost_webhook")).
		BodyJSON(payload).
		ReceiveSuccess(nil)

	if err != nil {
		log.Warn(err)
		return
	}
}
