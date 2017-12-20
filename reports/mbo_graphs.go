// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package reports

import (
	"fmt"
	"github.com/joyent/conch-shell/util"
	chart "github.com/wcharczuk/go-chart"
	"gopkg.in/jawher/mow.cli.v1"
	"net/http"
	"sort"
)

func mboHardwareFailureGraphListener(app *cli.Cmd) {
	var (
		manta_report_path = app.StringOpt("manta-report path", "", "Path to Manta job output file")
		manta_report_url  = app.StringOpt("manta-report-url url", "", "The url for manta report output")
		datacenter_choice = app.StringOpt("datacenter az", "", "Limit the output to a particular datacenter by UUID, partial UUID, or string name")
		remediation_min   = app.IntOpt("remediation-minimum", 90, "For a failure to be considered, its remediation time must be greater than or equal to this number")

		port = app.IntOpt("port", 1337, "Port to listen on")
	)

	app.Action = func() {
		manta_report := &mboMantaReport{}

		if *manta_report_path != "" {
			fmt.Println("Opening file " + *manta_report_path)
			if err := manta_report.NewFromFile(*manta_report_path); err != nil {
				util.Bail(err)
			}
		} else {
			fmt.Println("Downloading URL " + *manta_report_url)
			if err := manta_report.NewFromUrl(*manta_report_url); err != nil {
				util.Bail(err)
			}
		}

		fmt.Println("Parsing complete. Processing...")
		fmt.Println()
		manta_report.Process(*datacenter_choice, *remediation_min)
		report := manta_report.Processed

		az_names := make([]string, 0)
		for name := range report {
			az_names = append(az_names, name)
		}
		sort.Strings(az_names)

		fmt.Printf("Opening listener on port %d\n", *port)

		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("content-type", "text/html")
			fmt.Fprintf(w, "<html><body><h1>Conch : MBO Hardware Failures</h1>")
			fmt.Fprintf(w, "<h2>By Type</h2>")
			fmt.Fprintf(w, "<ul>")
			for _, name := range az_names {
				fmt.Fprintf(w, "<li><a href='/by_type.png?az=%s'>%s</a></li>", name, name)
			}
			fmt.Fprintf(w, "</ul>")
			fmt.Fprintf(w, "<h2>By Vendor</h2>")
			fmt.Fprintf(w, "<ul>")
			for _, name := range az_names {
				fmt.Fprintf(w, "<li><a href='/by_vendor.png?az=%s'>%s</a></li>", name, name)
			}
			fmt.Fprintf(w, "</ul>")

			fmt.Fprintf(w, "</body></html>")
		})

		http.HandleFunc("/by_type.png", func(w http.ResponseWriter, req *http.Request) {
			az_params, ok := req.URL.Query()["az"]
			if !ok || (len(az_params) == 0) {
				http.Error(w, "", 404)
				return
			}
			az_param := az_params[0]

			az, ok := report[az_param]
			if !ok {
				http.Error(w, "", 404)
				return
			}
			values := make([]chart.Value, 0)
			if len(az.TimesByType) == 0 {
				http.Error(w, "No data available", 500)
				return
			}

			for data_type, data := range az.TimesByType {
				values = append(values, chart.Value{
					Value: float64(data.Count),
					Label: fmt.Sprintf("%s : %d", data_type, data.Count),
				})
			}

			pie := chart.PieChart{
				Width:  512,
				Height: 512,
				Values: values,
			}

			w.Header().Set("Content-Type", "image/png")
			if err := pie.Render(chart.PNG, w); err != nil {
				fmt.Printf("Error rendering pie chart: %v\n", err)
				http.Error(w, err.Error(), 500)
				return
			}
		})

		http.HandleFunc("/by_vendor.png", func(w http.ResponseWriter, req *http.Request) {
			az_params, ok := req.URL.Query()["az"]
			if !ok || (len(az_params) == 0) {
				http.Error(w, "", 404)
				return
			}
			az_param := az_params[0]

			az, ok := report[az_param]
			if !ok {
				http.Error(w, "", 404)
				return
			}
			values := make([]chart.Value, 0)
			if len(az.TimesByVendorAndType) == 0 {
				http.Error(w, "No data available", 500)
				return
			}

			for vendor_name, vendor_data := range az.TimesByVendorAndType {
				var count int64

				for _, type_data := range vendor_data {
					count = count + type_data.Count
				}

				values = append(values, chart.Value{
					Value: float64(count),
					Label: fmt.Sprintf("%s : %d", vendor_name, count),
				})
			}

			pie := chart.PieChart{
				Width:  512,
				Height: 512,
				Values: values,
			}

			w.Header().Set("Content-Type", "image/png")
			if err := pie.Render(chart.PNG, w); err != nil {
				fmt.Printf("Error rendering pie chart: %v\n", err)
				http.Error(w, err.Error(), 500)
				return
			}
		})
		util.Bail(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))

	}
}
