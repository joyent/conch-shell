// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joyent/conch-shell/config"
	conch "github.com/joyent/go-conch"
	pgtime "github.com/joyent/go-conch/pg_time"
	"github.com/olekukonko/tablewriter"
	cli "gopkg.in/jawher/mow.cli.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
	"os"
	"regexp"
	"time"
)

var (
	JSON   bool
	Config *config.ConchConfig
	API    *conch.Conch
)

type MinimalDevice struct {
	Id       string             `json:"id"`
	AssetTag string             `json:"asset_tag"`
	Created  pgtime.ConchPgTime `json:"created,int"`
	LastSeen pgtime.ConchPgTime `json:"last_seen,int"`
	Health   string             `json:"health"`
	Flags    string             `json:"flags"`
	AZ       string             `json:"az"`
	Rack     string             `json:"rack"`
}

func BuildApiAndVerifyLogin() {
	API = &conch.Conch{
		BaseUrl: Config.Api,
		User:    Config.Api,
		Session: Config.Session,
	}

	if err := API.VerifyLogin(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func BuildApi() {
	API = &conch.Conch{
		BaseUrl: Config.Api,
		User:    Config.Api,
		Session: Config.Session,
	}
}

func GetMarkdownTable() (table *tablewriter.Table) {
	table = tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	return table
}

func Bail(err error) {
	if JSON {
		j, _ := json.Marshal(struct {
			Error   bool   `json:"error"`
			Message string `json:"message"`
		}{
			true,
			fmt.Sprintf("%v", err),
		})

		fmt.Println(string(j))
	} else {
		fmt.Println(err)
	}
	cli.Exit(1)
}

// DisplayDevices() is an abstraction to make sure that the output of
// ConchDevices is uniform, be it tables, json, or full json
func DisplayDevices(devices []conch.ConchDevice, full_output bool) (err error) {
	minimals := make([]MinimalDevice, 0)
	for _, d := range devices {
		minimals = append(minimals, MinimalDevice{
			d.Id,
			d.AssetTag,
			d.Created,
			d.LastSeen,
			d.Health,
			GenerateDeviceFlags(d),
			d.Location.Datacenter.Name,
			d.Location.Rack.Name,
		})
	}

	if JSON {
		var j []byte
		if full_output {
			j, err = json.Marshal(devices)
		} else {
			j, err = json.Marshal(minimals)
		}
		if err != nil {
			return err
		}
		fmt.Println(string(j))
		return nil
	}

	TableizeMinimalDevices(minimals, full_output, GetMarkdownTable()).Render()

	return nil
}

// TableizeMinimalDevices() is an abstraction to make sure that tables of
// ConchDevices-turned-MinimalDevices are uniform
func TableizeMinimalDevices(devices []MinimalDevice, full_output bool, table *tablewriter.Table) *tablewriter.Table {
	if full_output {
		table.SetHeader([]string{
			"AZ",
			"Rack",
			"ID",
			"Asset Tag",
			"Created",
			"Last Seen",
			"Health",
			"Flags",
		})
	} else {
		table.SetHeader([]string{
			"ID",
			"Asset Tag",
			"Created",
			"Last Seen",
			"Health",
			"Flags",
		})
	}

	for _, d := range devices {
		last_seen := ""
		if !d.LastSeen.IsZero() {
			last_seen = d.LastSeen.Format(time.UnixDate)
		}

		if full_output {
			table.Append([]string{
				d.AZ,
				d.Rack,
				d.Id,
				d.AssetTag,
				d.Created.Format(time.UnixDate),
				last_seen,
				d.Health,
				d.Flags,
			})
		} else {
			table.Append([]string{
				d.Id,
				d.AssetTag,
				d.Created.Format(time.UnixDate),
				last_seen,
				d.Health,
				d.Flags,
			})
		}
	}

	return table
}

// GenerateDeviceFlags() is an abstraction to make sure that the 'flags' field
// for ConchDevices remains uniform
func GenerateDeviceFlags(d conch.ConchDevice) (flags string) {
	flags = ""

	if !d.Deactivated.IsZero() {
		flags += "X"
	}

	if !d.Validated.IsZero() {
		flags += "v"
	}

	if !d.Graduated.IsZero() {
		flags += "g"
	}
	return flags
}

func JsonOut(thingy interface{}) {
	j, err := json.Marshal(thingy)

	if err != nil {
		Bail(err)
	}

	fmt.Println(string(j))
}

func MagicWorkspaceId(wat string) (id uuid.UUID, err error) {
	id, err = uuid.FromString(wat)
	if err == nil {
		return id, err
	}
	// So, it's not a UUID. Let's try for a string name or partial UUID
	workspaces, err := API.GetWorkspaces()
	if err != nil {
		return id, err
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s-", wat))
	for _, w := range workspaces {
		if (w.Name == wat) || re.MatchString(w.Id.String()) {
			return w.Id, nil
		}
	}

	return id, errors.New("Could not find workspace " + wat)
}
