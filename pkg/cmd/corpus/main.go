// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package corpus

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	_ "github.com/lib/pq"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	uuid "gopkg.in/satori/go.uuid.v1"
)

var FailedCount = 0

type Report struct {
	ID           uuid.UUID
	SKU          string
	Completed    time.Time
	Raw          string
	DeviceSerial string
}

type Reports []Report

/************************/

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Dump reports",
		Run:   run,
	})

}

/************************/

func run(cmd *cobra.Command, args []string) {
	reports := extractReports()

	expandedPath, err := homedir.Expand(viper.GetString("data_directory"))
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Writing reports into '" + expandedPath + "'")

	for i, report := range reports {
		log.Info(fmt.Sprintf("Processing entry %d of %d", i, len(reports)))
		fileName := fmt.Sprintf(expandedPath + "/" + report.SKU + ".json")

		err := ioutil.WriteFile(fileName, []byte(report.Raw), 0644)
		if err != nil {
			log.Warn(err)
		} else {
			log.Info("Wrote " + fileName)
		}
	}
}

/************************/

func extractReports() Reports {
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
	log.Info("Querying database")

	sql := fmt.Sprintf(`with reports(device_id, device_report_id, completed, result_num) as (
	select
		device_id,
		device_report_id,
		completed,
		row_number() over (
			partition by device_id order by completed desc
		) as result_num
	from validation_state
	where status = 'pass'
) select
    distinct on (hp.name) hp.name as "hardware_product_name",
	d.id as "device_id",
	dr.id as "device_report_id",
	r.completed as "validation_completed",
	dr.report as "device_report"
	from hardware_product hp
		join device d on d.hardware_product_id = hp.id
		join reports r on r.device_id = d.id
		join device_report dr on dr.id = r.device_report_id
	where r.result_num = 1
;`)
	log.Trace(sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	reports := make(Reports, 0)

	for rows.Next() {
		report := Report{}

		if err := rows.Scan(&report.SKU, &report.DeviceSerial, &report.ID, &report.Completed, &report.Raw); err != nil {
			log.Fatal(err)
		}

		reports = append(reports, report)
	}
	rows.Close()

	log.Info(fmt.Sprintf("Found %d device reports", len(reports)))

	log.Debug("Closing database connection")
	db.Close()

	return reports
}
