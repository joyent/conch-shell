// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package reports

import (
	"fmt"
	"github.com/joyent/conch-shell/pkg/reports/mbo"
	"github.com/joyent/conch-shell/pkg/util"
	"gopkg.in/jawher/mow.cli.v1"
)

func mboHardwareFailures(app *cli.Cmd) {
	var (
		mantaReportPath   = app.StringOpt("manta-report path", "", "Path to Manta job output file")
		mantaReportURL    = app.StringOpt("manta-report-url url", "", "The url for manta report output")
		fullOutput        = app.BoolOpt("full", false, "Instead of just presenting a datacenter summary, break results out by rack as well. Has no effect on --json")
		datacenterChoice  = app.StringOpt("datacenter az", "", "Limit the output to a particular datacenter by UUID, partial UUID, or string name")
		csvOutput         = app.BoolOpt("csv", false, "Output report as CSV. Assumes --full and overrides --json")
		includeVendors    = app.BoolOpt("include-vendors", false, "Include vendor data")
		includeComponents = app.BoolOpt("include-components", false, "Break out failures by components")
		remediationMin    = app.IntOpt("remediation-minimum", 90, "For a failure to be considered, its remediation time must be greater than or equal to this number")
	)

	app.Spec = "--manta-report=<manta-report.json> | --manta-report-url=<https://manta/report.json> [OPTIONS]"

	app.Action = func() {
		report := &mbo.MantaReport{}

		if *csvOutput {
			util.Pretty = false
		}

		if *mantaReportPath != "" {
			if !*csvOutput {
				fmt.Println("Opening file " + *mantaReportPath)
			}
			if err := report.NewFromFile(*mantaReportPath); err != nil {
				util.Bail(err)
			}
		} else {
			if !*csvOutput {
				fmt.Println("Downloading URL " + *mantaReportURL)
			}
			if err := report.NewFromURL(*mantaReportURL); err != nil {
				util.Bail(err)
			}
		}

		if !*csvOutput {
			fmt.Println("Parsing complete. Processing...")
			fmt.Println()
		}

		report.Process(*datacenterChoice, *remediationMin)
		if *csvOutput {
			fmt.Println(report.AsCsv())
		} else {
			fmt.Println(report.AsText(
				*fullOutput,
				*includeVendors,
				*includeComponents,
			))
		}
	}
}
