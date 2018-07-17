// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hardware

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/util"
	uuid "gopkg.in/satori/go.uuid.v1"
)

const singleDBHWPTemplate = `ID: {{ .ID }}
    Name: {{ .Name }}
    Alias: {{ .Alias }}
    Prefix: {{ .Prefix }}
    Vendor: {{ .Vendor }}
    SKU: {{ .SKU }}
    Generation Name: {{ .GenerationName }}
    Legacy Product Name: {{ .LegacyProductName }}

    Specification:
    {{ marshal .Specification }}
`

func displayDBHardwareProduct(h conch.DBHardwareProduct) {
	funcMap := template.FuncMap{
		"marshal": func(v interface{}) string {
			if len(h.Specification) <= 1 {
				return ""
			}
			var j interface{}

			if err := json.Unmarshal([]byte(h.Specification), &j); err != nil {
				return err.Error()
			}

			s, err := json.MarshalIndent(j, "    ", "  ")
			if err != nil {
				return err.Error()
			}
			return string(s)
		},
	}

	t := template.Must(template.New("hw").Funcs(funcMap).Parse(singleDBHWPTemplate))
	if err := t.Execute(os.Stdout, h); err != nil {
		util.Bail(err)
	}
}

func getOneDB(app *cli.Cmd) {
	app.Action = func() {
		ret, err := util.API.GetDBHardwareProduct(ProductUUID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(ret)
			return
		}

		if err != nil {
			util.Bail(err)
		}
		displayDBHardwareProduct(ret)
	}
}

func getAllDB(app *cli.Cmd) {
	app.Action = func() {
		ret, err := util.API.GetDBHardwareProducts()
		if err != nil {
			util.Bail(err)
		}
		if util.JSON {
			util.JSONOut(ret)
			return
		}

		table := util.GetMarkdownTable()
		table.SetHeader([]string{
			"ID",
			"Name",
			"Alias",
			"SKU",
			"Prefix",
			"Generation",
			"Legacy Name",
		})
		for _, h := range ret {
			table.Append([]string{
				h.ID.String(),
				h.Name,
				h.Alias,
				h.SKU,
				h.Prefix,
				h.GenerationName,
				h.LegacyProductName,
			})
		}
		table.Render()
	}
}

func createOneDB(app *cli.Cmd) {
	var (
		nameOpt     = app.StringOpt("name", "", "Joyent's Name")
		aliasOpt    = app.StringOpt("alias", "", "Joyent's Name")
		vendorOpt   = app.StringOpt("vendor", "", "Vendor UUID")
		prefixOpt   = app.StringOpt("prefix", "", "Prefix")
		skuOpt      = app.StringOpt("sku", "", "SKU")
		genOpt      = app.StringOpt("generation-name generation gen", "", "Generation Name")
		legacyOpt   = app.StringOpt("legacy-name legacy", "", "Legacy Product Name")
		specOpt     = app.BoolOpt("specification spec", false, "Will provide specification as last arg")
		filePathArg = app.StringArg("FILE", "-", "Path to a JSON file to use as the specification. '-' indicates STDIN")
	)
	app.Spec = "--name --alias --vendor [OPTIONS] [FILE]"

	app.Action = func() {
		vendor, err := uuid.FromString(*vendorOpt)
		if err != nil {
			util.Bail(err)
		}

		h := conch.DBHardwareProduct{
			Name:              *nameOpt,
			Alias:             *aliasOpt,
			Vendor:            vendor,
			SKU:               *skuOpt,
			GenerationName:    *genOpt,
			LegacyProductName: *legacyOpt,
			Prefix:            *prefixOpt,
		}

		if *specOpt {
			var b []byte
			var err error
			if *filePathArg == "-" {
				b, err = ioutil.ReadAll(os.Stdin)
			} else {
				b, err = ioutil.ReadFile(*filePathArg)
			}
			if err != nil {
				util.Bail(err)
			}
			if len(string(b)) <= 1 {
				util.Bail(errors.New("No specification provided"))
			}
			h.Specification = string(b)
		}

		if err := util.API.SaveDBHardwareProduct(&h); err != nil {
			util.Bail(err)
		}

		ret, err := util.API.GetDBHardwareProduct(&h.ID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(ret)
			return
		}
		displayDBHardwareProduct(ret)
	}

}

func removeOneDB(app *cli.Cmd) {
	app.Action = func() {
		if err := util.API.DeleteDBHardwareProduct(ProductUUID); err != nil {
			util.Bail(err)
		}
	}

}

func updateOneDB(app *cli.Cmd) {
	var (
		nameOpt     = app.StringOpt("name", "", "Joyent's Name")
		aliasOpt    = app.StringOpt("alias", "", "Joyent's Name")
		vendorOpt   = app.StringOpt("vendor", "", "Vendor UUID")
		prefixOpt   = app.StringOpt("prefix", "", "Prefix")
		skuOpt      = app.StringOpt("sku", "", "SKU")
		genOpt      = app.StringOpt("generation-name generation gen", "", "Generation Name")
		legacyOpt   = app.StringOpt("legacy-name legacy", "", "Legacy Product Name")
		specOpt     = app.BoolOpt("specification spec", false, "Will provide specification as last arg")
		filePathArg = app.StringArg("FILE", "-", "Path to a JSON file to use as the specification. '-' indicates STDIN")
	)
	app.Spec = "[OPTIONS] [ --specification ] [FILE]"

	app.Action = func() {
		h, err := util.API.GetDBHardwareProduct(ProductUUID)
		if err != nil {
			util.Bail(err)
		}

		if *nameOpt != "" {
			h.Name = *nameOpt
		}

		if *aliasOpt != "" {
			h.Alias = *aliasOpt
		}

		if *vendorOpt != "" {
			vendor, err := uuid.FromString(*vendorOpt)
			if err != nil {
				util.Bail(err)
			}

			h.Vendor = vendor
		}

		if *prefixOpt != "" {
			h.Prefix = *prefixOpt
		}

		if *skuOpt != "" {
			h.SKU = *skuOpt
		}

		if *genOpt != "" {
			h.GenerationName = *genOpt
		}

		if *legacyOpt != "" {
			h.LegacyProductName = *legacyOpt
		}

		if *specOpt {
			var b []byte
			var err error
			if *filePathArg == "-" {
				b, err = ioutil.ReadAll(os.Stdin)
			} else {
				b, err = ioutil.ReadFile(*filePathArg)
			}
			if err != nil {
				util.Bail(err)
			}
			if len(string(b)) <= 1 {
				util.Bail(errors.New("No specification provided"))
			}
			h.Specification = string(b)
		}

		if err := util.API.SaveDBHardwareProduct(&h); err != nil {
			util.Bail(err)
		}

		ret, err := util.API.GetDBHardwareProduct(h.ID)
		if err != nil {
			util.Bail(err)
		}

		if util.JSON {
			util.JSONOut(ret)
			return
		}
		displayDBHardwareProduct(ret)
	}
}
