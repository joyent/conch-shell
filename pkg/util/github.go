// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package util contains common routines used throughout the command base
package util

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/blang/semver"
	"github.com/dghubble/sling"
)

const GhOrg = "joyent"
const GhRepo = "conch-shell"

// GithubRelease represents a 'release' for a Github project
type GithubRelease struct {
	URL        string         `json:"html_url"`
	TagName    string         `json:"tag_name"`
	SemVer     semver.Version `json:"-"` // Will be set to 0.0.0 if no releases are found
	Body       string         `json:"body"`
	Name       string         `json:"name"`
	Assets     []GithubAsset  `json:"assets"`
	PreRelease bool           `json:"prerelease"`
	Upgrade    bool           `json:"-"`
}

type GithubReleases []GithubRelease

func (g GithubReleases) Len() int {
	return len(g)
}

func (g GithubReleases) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g GithubReleases) Less(i, j int) bool {
	var iSem, jSem semver.Version

	if g[i].TagName == "" {
		iSem = semver.MustParse("0.0.0")
	} else {
		iSem = semver.MustParse(
			strings.TrimLeft(g[i].TagName, "v"),
		)
	}

	if g[j].TagName == "" {
		jSem = semver.MustParse("0.0.0")
	} else {
		jSem = semver.MustParse(
			strings.TrimLeft(g[j].TagName, "v"),
		)
	}

	return iSem.GT(jSem) // reversing sort
}

// GithubAsset represents a file inside of a github release
type GithubAsset struct {
	URL                string `json:"url"`
	Name               string `json:"name"`
	State              string `json:"state"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var ErrNoGithubRelease = errors.New("no appropriate github release found")

// LatestGithubRelease returns some fields from the latest Github Release
// that matches our major version
func LatestGithubRelease() (gh GithubRelease, err error) {
	releases := make(GithubReleases, 0)

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/releases",
		GhOrg,
		GhRepo,
	)

	_, err = sling.New().
		Set("User-Agent", fmt.Sprintf("Conch/%s", Version)).
		Get(url).Receive(&releases, nil)

	if err != nil {
		return gh, err
	}

	sort.Sort(releases)

	for _, r := range releases {
		if r.PreRelease {
			continue
		}
		if r.TagName == "" {
			continue
		}
		r.SemVer = CleanVersion(
			strings.TrimLeft(r.TagName, "v"),
		)

		if r.SemVer.Major == SemVersion.Major {
			if r.SemVer.GT(SemVersion) {
				r.Upgrade = true
			}
			return r, nil
		}
	}

	return gh, ErrNoGithubRelease
}

func GithubReleasesSince(start semver.Version) GithubReleases {
	releases := make(GithubReleases, 0)

	diff := make(GithubReleases, 0)

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/releases",
		GhOrg,
		GhRepo,
	)

	_, err := sling.New().
		Set("User-Agent", fmt.Sprintf("Conch/%s", Version)).
		Get(url).Receive(&releases, nil)

	if err != nil {
		return diff
	}

	sort.Sort(releases)

	for _, r := range releases {
		if r.PreRelease {
			continue
		}
		if r.TagName == "" {
			continue
		}
		r.SemVer = CleanVersion(
			strings.TrimLeft(r.TagName, "v"),
		)

		if r.SemVer.Major == SemVersion.Major {
			if r.SemVer.GT(start) {
				diff = append(diff, r)
			}
		}
	}

	sort.Sort(diff)

	return diff
}

// CleanVersion removes a "v" prefix, and anything after a dash
// For example, pass in v2.99.10-abcde-dirty and get back a semver containing
// 2.29.10
// Why? Git and Semver differ in their notions of what those extra bits mean.
// In Git, they mean "v2.99.10, plus some other stuff that happend". In semver,
// they indicate that this is a prerelease of v2.99.10. Obviously this screws
// up comparisions. This function lets us clean that stuff out so we can get a
// clean comparison
func CleanVersion(version string) semver.Version {
	bits := strings.Split(strings.TrimLeft(version, "v"), "-")
	return semver.MustParse(bits[0])
}
