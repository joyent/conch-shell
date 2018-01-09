// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joyent/conch-shell/pkg/config"
	"github.com/joyent/conch-shell/pkg/reports/mbo"
	c_templates "github.com/joyent/conch-shell/pkg/templates"
	"github.com/joyent/conch-shell/pkg/util"
	homedir "github.com/mitchellh/go-homedir"
	chart "github.com/wcharczuk/go-chart"
	"gopkg.in/jawher/mow.cli.v1"
	"html/template"
	"net/http"
	"os"
	"sort"
	"time"
)

// VERSION is the application version
const VERSION = "0.0.0"

func main() {
	app := cli.App("conch-mbo", "HTTP interface for MBO hardware failure reports")
	app.Version("version", VERSION)

	var (
		configFile       = app.StringOpt("config c", "~/.conch.json", "Path to config file")
		mantaReportPath  = app.StringOpt("manta-report path", "", "Path to Manta job output file")
		mantaReportURL   = app.StringOpt("manta-report-url url", "", "The url for manta report output")
		datacenterChoice = app.StringOpt("datacenter az", "", "Limit the output to a particular datacenter by UUID, partial UUID, or string name")
		remediationMin   = app.IntOpt("remediation-minimum", 90, "For a failure to be considered, its remediation time must be greater than or equal to this number")

		port = app.IntOpt("port", 1337, "Port to listen on")
	)

	app.Before = func() {

		util.Pretty = true
		util.Spin = spinner.New(spinner.CharSets[10], 100*time.Millisecond)
		util.Spin.FinalMSG = "Complete.\n"

		configFilePath, err := homedir.Expand(*configFile)
		if err != nil {
			util.Bail(err)
		}

		cfg, err := config.NewFromJSONFile(configFilePath)
		if err != nil {
			fmt.Println("A login error occurred. Please use 'conch' to login...")
			util.Bail(err)
		}
		cfg.Path = configFilePath
		util.Config = cfg

		for _, prof := range cfg.Profiles {
			if prof.Active {
				util.ActiveProfile = prof
				break
			}
		}
		util.BuildAPI()

		if err := util.API.VerifyLogin(); err != nil {
			fmt.Println("A login error occurred. Please use 'conch' to login...")
			util.Bail(err)
		}

	}

	app.Action = func() {
		mantaReport := &mbo.MantaReport{}

		if *mantaReportPath != "" {
			fmt.Println("Opening file " + *mantaReportPath)
			if err := mantaReport.NewFromFile(*mantaReportPath); err != nil {
				util.Bail(err)
			}
		} else {
			fmt.Println("Downloading URL " + *mantaReportURL)
			if err := mantaReport.NewFromURL(*mantaReportURL); err != nil {
				util.Bail(err)
			}
		}

		fmt.Println("Parsing complete. Processing...")
		fmt.Println()

		mantaReport.Process(*datacenterChoice, *remediationMin)
		report := mantaReport.Processed

		azNames := make([]string, 0)
		for name := range report {
			azNames = append(azNames, name)
		}
		sort.Strings(azNames)

		fmt.Printf("Opening listener on port %d\n", *port)

		gorilla := mux.NewRouter()

		gorilla.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			tmpl, err := template.New("index").Parse(c_templates.MboGraphsIndex)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("content-type", "text/html")
			_ = tmpl.Execute(w,
				struct {
					AzNames []string
				}{
					azNames,
				},
			)
		})
		gorilla.HandleFunc("/full", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "text/plain")
			fmt.Fprintf(w, mantaReport.AsText(true, true, true))
		})

		gorilla.HandleFunc("/full.csv", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "text/csv")
			fmt.Fprintf(w, mantaReport.AsCsv())
		})

		gorilla.HandleFunc("/style.css", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "text/css")
			fmt.Fprintf(w, c_templates.MboGraphsReportsCSS)
		})

		gorilla.HandleFunc("/reports/times/{az}", func(w http.ResponseWriter, req *http.Request) {
			params := mux.Vars(req)
			azParam := params["az"]
			if len(azParam) == 0 {
				http.Error(w, "", 404)
				return
			}

			if _, ok := mantaReport.Processed[azParam]; !ok {
				http.Error(w, "No data found for "+azParam, 404)
				return
			}
			azData := mantaReport.Processed[azParam]

			tmpl, err := template.New("index").Parse(c_templates.MboGraphsReportsIndex)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("content-type", "text/html")
			_ = tmpl.Execute(w,
				struct {
					Name string
					Data mbo.DatacenterReport
				}{
					azParam,
					azData,
				},
			)
		})

		gorilla.HandleFunc("/reports/times/{az}/{component}", func(w http.ResponseWriter, req *http.Request) {
			params := mux.Vars(req)
			azParam := params["az"]
			if len(azParam) == 0 {
				http.Error(w, "", 404)
				return
			}
			componentParam := params["component"]
			if len(componentParam) == 0 {
				http.Error(w, "", 404)
				return
			}

			if _, ok := mantaReport.Processed[azParam]; !ok {
				http.Error(w, "No data found for "+azParam, 404)
				return
			}
			azData := mantaReport.Processed[azParam]

			if _, ok := azData.TimesBySubType[componentParam]; !ok {
				http.Error(w, fmt.Sprintf("No data found for AZ %s, type %s", azParam, componentParam), 404)
				return
			}
			subtypeData := azData.TimesBySubType[componentParam]

			tmpl, err := template.New("index").Parse(c_templates.MboGraphsReportsBySubtype)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("content-type", "text/html")
			_ = tmpl.Execute(w,
				struct {
					Az   string
					Name string
					Data map[string]*mbo.TypeReport
				}{
					azParam,
					componentParam,
					subtypeData,
				},
			)
		})

		gorilla.HandleFunc("/reports/times/{az}/{component}/{subtype}", func(w http.ResponseWriter, req *http.Request) {
			params := mux.Vars(req)
			azParam := params["az"]
			if len(azParam) == 0 {
				http.Error(w, "", 404)
				return
			}

			if _, ok := mantaReport.Processed[azParam]; !ok {
				http.Error(w, "No data found for "+azParam, 404)
				return
			}
			azData := mantaReport.Processed[azParam]

			componentParam := params["component"]
			if len(componentParam) == 0 {
				http.Error(w, "", 404)
				return
			}
			if _, ok := azData.TimesBySubType[componentParam]; !ok {
				http.Error(w, fmt.Sprintf("No data found for AZ %s, type %s", azParam, componentParam), 404)
				return
			}
			componentData := azData.TimesBySubType[componentParam]

			/**/

			subtypeParam := params["subtype"]
			if len(subtypeParam) == 0 {
				http.Error(w, "", 404)
				return
			}

			if _, ok := componentData[subtypeParam]; !ok {
				http.Error(w, fmt.Sprintf(
					"No data found for AZ %s, type %s, subtype %s",
					azParam,
					componentParam,
					subtypeParam,
				), 404)
				return
			}
			subtypeData := componentData[subtypeParam]

			tmpl, err := template.New("index").Parse(c_templates.MboGraphsReportsByComponentAndSubtype)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("content-type", "text/html")
			_ = tmpl.Execute(w,
				struct {
					Az        string
					Component string
					Subtype   string
					Data      *mbo.TypeReport
				}{
					azParam,
					componentParam,
					subtypeParam,
					subtypeData,
				},
			)
		})

		gorilla.HandleFunc("/graphics/{az}/by_type.png", func(w http.ResponseWriter, req *http.Request) {
			params := mux.Vars(req)
			azParam := params["az"]
			if len(azParam) == 0 {
				http.Error(w, "", 404)
				return
			}

			az, ok := report[azParam]
			if !ok {
				http.Error(w, "", 404)
				return
			}
			values := make([]chart.Value, 0)
			if len(az.TimesByType) == 0 {
				http.Error(w, "No data available", 500)
				return
			}

			for dataType, data := range az.TimesByType {
				values = append(values, chart.Value{
					Value: float64(data.Count),
					Label: fmt.Sprintf("%s : %d", dataType, data.Count),
				})
			}
			barChart := chart.BarChart{
				Height:   512,
				BarWidth: 60,
				XAxis: chart.Style{
					Show: true,
				},
				YAxis: chart.YAxis{
					Style: chart.Style{
						Show: true,
					},
				},
				Bars: values,
			}

			w.Header().Set("Content-Type", "image/png")
			if err := barChart.Render(chart.PNG, w); err != nil {
				fmt.Printf("Error rendering pie chart: %v\n", err)
				http.Error(w, err.Error(), 500)
				return
			}
		})

		gorilla.HandleFunc("/graphics/{az}/by_vendor.png", func(w http.ResponseWriter, req *http.Request) {
			params := mux.Vars(req)
			azParam := params["az"]
			if len(azParam) == 0 {
				http.Error(w, "", 404)
				return
			}

			az, ok := report[azParam]
			if !ok {
				http.Error(w, "", 404)
				return
			}
			values := make([]chart.Value, 0)
			if len(az.TimesByVendorAndType) == 0 {
				http.Error(w, "No data available", 500)
				return
			}

			for name, data := range az.TimesByVendorAndType {
				var count int64

				for _, typeData := range data {
					count = count + typeData.Count
				}

				values = append(values, chart.Value{
					Value: float64(count),
					Label: fmt.Sprintf("%s : %d", name, count),
				})
			}

			barChart := chart.BarChart{
				Height: 512,
				XAxis: chart.Style{
					Show: true,
				},
				YAxis: chart.YAxis{
					Style: chart.Style{
						Show: true,
					},
				},
				Bars: values,
			}

			w.Header().Set("Content-Type", "image/png")
			if err := barChart.Render(chart.PNG, w); err != nil {
				fmt.Printf("Error rendering pie chart: %v\n", err)
				http.Error(w, err.Error(), 500)
				return
			}
		})

		logger := handlers.CombinedLoggingHandler(os.Stdout, gorilla)
		http.Handle("/", gorilla)
		util.Bail(http.ListenAndServe(
			fmt.Sprintf(":%d", *port),
			logger,
		))
	}

	_ = app.Run(os.Args)
}
