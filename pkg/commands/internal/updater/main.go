// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package updater

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/joyent/conch-shell/pkg/util"
	"gopkg.in/jawher/mow.cli.v1"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

func status(cmd *cli.Cmd) {
	cmd.Action = func() {
		gh, err := util.LatestGithubRelease("joyent", "conch-shell")
		if err != nil {
			util.Bail(err)
		}
		fmt.Printf("This is v%s. Current release is %s.\n",
			util.Version,
			gh.TagName,
		)
	}
}

func changelog(cmd *cli.Cmd) {
	cmd.Action = func() {
		gh, err := util.LatestGithubRelease("joyent", "conch-shell")
		if err != nil {
			util.Bail(err)
		}

		// I'm not going to try and fully sanitize the output
		// for a shell environment but removing the markdown
		// backticks seems like a no-brainer for safety.
		re := regexp.MustCompile("`")
		body := gh.Body
		re.ReplaceAllLiteralString(body, "'")
		fmt.Printf("Version %s Changelog:\n\n", gh.TagName)
		fmt.Println(body)
	}
}

func selfUpdate(cmd *cli.Cmd) {
	var force = cmd.BoolOpt(
		"force",
		false,
		"Update the binary even if it appears we are on the current release",
	)
	cmd.Action = func() {
		gh, err := util.LatestGithubRelease("joyent", "conch-shell")
		if err != nil {
			util.Bail(err)
		}
		sem := semver.MustParse(util.Version)
		if !*force {
			if gh.SemVer.LTE(sem) {
				fmt.Println("Already at the latest release")
				return
			}
		}
		fmt.Printf(
			"=> Attempting to upgrade from %s to %s...\n",
			sem,
			gh.SemVer,
		)

		fmt.Printf("===> Detected OS to be '%s' and arch to be '%s'\n", runtime.GOOS, runtime.GOARCH)
		lookingFor := fmt.Sprintf("conch-%s-%s", runtime.GOOS, runtime.GOARCH)
		downloadURL := ""

		for _, a := range gh.Assets {
			if a.Name == lookingFor {
				downloadURL = a.BrowserDownloadURL
			}
		}
		if downloadURL == "" {
			fmt.Println("XX Could not find an appropriate binary")
			return
		}
		fmt.Printf("===> Found new binary URL: %s\n", downloadURL)
		fmt.Println("=====> Downloading binary...")
		resp, err := http.Get(downloadURL)
		if err != nil {
			util.Bail(err)
		}
		if resp.StatusCode != 200 {
			fmt.Printf("XX Could not download binary (status %d)\n", resp.StatusCode)
			return
		}

		newBinary, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			util.Bail(err)
		}
		fmt.Println("=====> Done.")

		binPath, err := os.Executable()
		if err != nil {
			util.Bail(err)
		}

		fullPath, err := filepath.EvalSymlinks(binPath)
		if err != nil {
			util.Bail(err)
		}
		fmt.Printf("===> Detected local binary path: %s\n", fullPath)
		existingStat, err := os.Lstat(fullPath)
		if err != nil {
			util.Bail(err)
		}

		fmt.Println("===> Overwriting local binary...")
		err = ioutil.WriteFile(fullPath, newBinary, existingStat.Mode())
		if err != nil {
			util.Bail(err)
		}
		fmt.Println("===> Done.")

	}
}
