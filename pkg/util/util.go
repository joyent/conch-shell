// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package util contains common routines used throughout the command base
package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Bowery/prompt"
	"github.com/blang/semver"
	"github.com/dghubble/sling"
	cli "github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/config"
	"github.com/joyent/conch-shell/pkg/pgtime"
	"github.com/olekukonko/tablewriter"
	"unicode/utf8"
)

var (
	// UserAgent will be used as the http user agent when making API calls
	UserAgent string

	// JSON tells us if we should output JSON
	JSON bool

	// Config is a global Config object
	Config *config.ConchConfig

	// ActiveProfile represents, well, the active profile
	ActiveProfile *config.ConchProfile

	// API is a global Conch API object
	API *conch.Conch

	// Debug decides if we should put the API in debug mode
	// Yes, this is a bit of a kludge
	Debug bool

	// Trace decides if we should trace the HTTP transactions
	// Yes, this is a bit of a kludge
	Trace bool
)

// These variables are provided by the build environment
var (
	Version   string
	BuildTime string
	GitRev    string
	BuildHost string
)

// DateFormat should be used in date formatting calls to ensure uniformity of
// output
const DateFormat = "2006-01-02 15:04:05 -0700 MST"

// RefreshTokenTime represent when a JWT token will be refreshed, based on this
// many seconds left on the expiry time
const RefreshTokenTime = 86400

// TimeStr ensures that all Times are formatted using .Local() and DateFormat
func TimeStr(t time.Time) string {
	return t.Local().Format(DateFormat)
}

// MinimalDevice represents a limited subset of Device data, that which we are
// going to present to the user
type MinimalDevice struct {
	ID       string        `json:"id"`
	AssetTag string        `json:"asset_tag"`
	Created  pgtime.PgTime `json:"created"`
	LastSeen pgtime.PgTime `json:"last_seen"`
	Health   string        `json:"health"`
	Flags    string        `json:"flags"`
	AZ       string        `json:"az"`
	Rack     string        `json:"rack"`
}

// BuildAPIAndVerifyLogin builds a Conch object using the Config data and calls
// VerifyLogin
func BuildAPIAndVerifyLogin() {
	BuildAPI()
	if err := API.VerifyLogin(RefreshTokenTime, false); err != nil {
		Bail(err)
	}
	ActiveProfile.Session = API.Session
	ActiveProfile.JWToken = API.JWToken
	ActiveProfile.Expires = API.Expires
	WriteConfig()
}

// WriteConfig serializes the Config struct to disk
func WriteConfig() {
	if err := Config.SerializeToFile(Config.Path); err != nil {
		Bail(err)
	}
}

// BuildAPI builds a Conch object
func BuildAPI() {
	if ActiveProfile == nil {
		Bail(errors.New("no active profile. Please use 'conch profile' to create or set an active profile"))
	}

	API = &conch.Conch{
		BaseURL: ActiveProfile.BaseURL,
		Session: ActiveProfile.Session,
		JWToken: ActiveProfile.JWToken,
		Debug:   Debug,
		Trace:   Trace,
	}
	if UserAgent != "" {
		API.UA = UserAgent
	}
}

// GetMarkdownTable returns a tablewriter configured to output markdown
// compatible text
func GetMarkdownTable() (table *tablewriter.Table) {
	table = tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	return table
}

// Bail is a --json aware way of dying
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

// DisplayDevices is an abstraction to make sure that the output of
// Devices is uniform, be it tables, json, or full json
func DisplayDevices(devices []conch.Device, fullOutput bool) (err error) {
	minimals := make([]MinimalDevice, 0)
	for _, d := range devices {
		minimals = append(minimals, MinimalDevice{
			d.ID,
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
		if fullOutput {
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

	TableizeMinimalDevices(minimals, fullOutput, GetMarkdownTable()).Render()

	return nil
}

// TableizeMinimalDevices is an abstraction to make sure that tables of
// Devices-turned-MinimalDevices are uniform
func TableizeMinimalDevices(devices []MinimalDevice, fullOutput bool, table *tablewriter.Table) *tablewriter.Table {
	if fullOutput {
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
		lastSeen := ""
		if !d.LastSeen.IsZero() {
			lastSeen = TimeStr(d.LastSeen.AsUTC())
		}

		if fullOutput {
			table.Append([]string{
				d.AZ,
				d.Rack,
				d.ID,
				d.AssetTag,
				TimeStr(d.Created.AsUTC()),
				lastSeen,
				d.Health,
				d.Flags,
			})
		} else {
			table.Append([]string{
				d.ID,
				d.AssetTag,
				TimeStr(d.Created.AsUTC()),
				lastSeen,
				d.Health,
				d.Flags,
			})
		}
	}

	return table
}

// GenerateDeviceFlags is an abstraction to make sure that the 'flags' field
// for Devices remains uniform
func GenerateDeviceFlags(d conch.Device) (flags string) {
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

// JSONOut marshals an interface to JSON
func JSONOut(thingy interface{}) {
	j, err := json.Marshal(thingy)

	if err != nil {
		Bail(err)
	}

	fmt.Println(string(j))
}

// JSONOutIndent marshals an interface to indented JSON
func JSONOutIndent(thingy interface{}) {
	j, err := json.MarshalIndent(thingy, "", "     ")

	if err != nil {
		Bail(err)
	}

	fmt.Println(string(j))
}

// GithubRelease represents a 'release' for a Github project
type GithubRelease struct {
	URL     string         `json:"html_url"`
	TagName string         `json:"tag_name"`
	SemVer  semver.Version `json:"-"` // Will be set to 0.0.0 if no releases are found
	Body    string         `json:"body"`
	Name    string         `json:"name"`
	Assets  []GithubAsset  `json:"assets"`
}

// GithubAsset represents a file inside of a github release
type GithubAsset struct {
	URL                string `json:"url"`
	Name               string `json:"name"`
	State              string `json:"state"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// LatestGithubRelease returns some fields from the latest Github Release
// object for the given owner and repo via
// "https://api.github.com/repos/:owner/:repo/releases/latest"
func LatestGithubRelease(owner string, repo string) (*GithubRelease, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/releases/latest",
		owner,
		repo,
	)

	gh := &GithubRelease{}

	_, err := sling.New().
		Set("User-Agent", UserAgent).
		Get(url).Receive(&gh, nil)

	if err != nil {
		return gh, err
	}

	if gh.TagName == "" {
		gh.SemVer = semver.MustParse("0.0.0")
	} else {
		sem, err := semver.Make(
			strings.TrimLeft(gh.TagName, "v"),
		)
		if err != nil {
			return gh, err
		}
		gh.SemVer = sem
	}

	return gh, err
}

// IsPasswordSane verifies that the given password follows the current rules
// and restrictions
func IsPasswordSane(password string, profile *config.ConchProfile) error {
	if utf8.RuneCountInString(password) < 12 {
		return errors.New("length must be >= 12")
	}
	if profile != nil {
		if strings.EqualFold(password, profile.User) {
			return errors.New("password cannot match user name")
		}
	}
	return nil
}

// InteractiveForcePasswordChange is an abstraction for the process of
// prompting a user for a new password, validating it, and issuing the API
// calls to execute the change
func InteractiveForcePasswordChange() {
	fmt.Println("You must change your password to continue.")
	fmt.Println("The new password must:")
	fmt.Println("  * Be at least 12 characters long")
	fmt.Println("  * Cannot match the user name or email address on the account")
	fmt.Println()

	var password string
	for {
		s, err := prompt.Password("New Password:")
		if err != nil {
			Bail(err)
		}
		err = IsPasswordSane(s, ActiveProfile)
		if err != nil {
			fmt.Println(err)
		} else {
			password = s
			break
		}

	}
	if err := API.ChangePassword(password); err != nil {
		Bail(err)
	}

	if err := API.Login(ActiveProfile.User, password); err != nil {
		Bail(err)
	}

	ActiveProfile.Session = API.Session
	ActiveProfile.JWToken = API.JWToken

	WriteConfig()
}
