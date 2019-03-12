// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package update

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/jawher/mow.cli"
	"github.com/joyent/conch-shell/pkg/util"
)

func status(cmd *cli.Cmd) {
	cmd.Action = func() {
		gh, err := util.LatestGithubRelease()
		if err != nil {
			if err == util.ErrNoGithubRelease {
				fmt.Printf(
					"This is v%s. No upgrade is available.\n",
					util.Version,
				)
				return
			}

			util.Bail(err)
		}

		if gh.Upgrade {
			fmt.Printf(
				"This is v%s. An upgrade to %s is available\n",
				util.Version,
				gh.TagName,
			)
		} else {
			fmt.Printf(
				"This is v%s. No upgrade is available.\n",
				util.Version,
			)
		}
	}
}

func changelog(cmd *cli.Cmd) {
	cmd.Action = func() {
		gh, err := util.LatestGithubRelease("joyent", "conch-shell")
		if err != nil {
			if err == util.ErrNoGithubRelease {
				fmt.Println("No changelog found")
				return
			}
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
		gh, err := util.LatestGithubRelease()

		if err != nil {
			if err == util.ErrNoGithubRelease {
				fmt.Fprintln(os.Stderr, "no upgrade available")
				return
			}

			util.Bail(err)
		}

		if !*force {
			if !gh.Upgrade {
				util.Bail(errors.New("no upgrade required"))
			}
		}

		if !util.JSON {
			fmt.Fprintf(
				os.Stderr,
				"Attempting to upgrade from %s to %s...\n",
				util.SemVersion,
				gh.SemVer,
			)

			fmt.Fprintf(
				os.Stderr,
				"Detected OS to be '%s' and arch to be '%s'\n",
				runtime.GOOS,
				runtime.GOARCH,
			)
		}

		// What platform are we on?
		lookingFor := fmt.Sprintf("conch-%s-%s", runtime.GOOS, runtime.GOARCH)
		downloadURL := ""

		// Is this a supported platform
		for _, a := range gh.Assets {
			if a.Name == lookingFor {
				downloadURL = a.BrowserDownloadURL
			}
		}

		if downloadURL == "" {
			util.Bail(fmt.Errorf(
				"could not find an appropriate binary for %s-%s",
				runtime.GOOS,
				runtime.GOARCH,
			))
		}
		/// Download the binary
		conchBin, err := updaterDownloadFile(downloadURL)
		if err != nil {
			util.Bail(err)
		}

		/// Verify checksum

		// This assumes our build system is being sensible about file names.
		// At time of writing, it is.
		shaURL := downloadURL + ".sha256"
		shaBin, err := updaterDownloadFile(shaURL)
		if err != nil {
			util.Bail(err)
		}

		// The checksum file looks like "thisisahexstring ./conch-os-arch"
		bits := strings.Split(string(shaBin[:]), " ")
		remoteSum := bits[0]

		if !util.JSON {
			fmt.Fprintf(
				os.Stderr,
				"Server-side SHA256 sum: %s\n",
				remoteSum,
			)
		}

		h := sha256.New()
		h.Write(conchBin)
		sum := hex.EncodeToString(h.Sum(nil))

		if !util.JSON {
			fmt.Fprintf(
				os.Stderr,
				"SHA256 sum of downloaded binary: %s\n",
				sum,
			)
		}

		if sum == remoteSum {
			if !util.JSON {
				fmt.Fprintf(
					os.Stderr,
					"SHA256 checksums match\n",
				)
			}
		} else {
			util.Bail(fmt.Errorf(
				"!!! SHA of downloaded file does not match the provided SHA sum: '%s' != '%s'",
				sum,
				remoteSum,
			))
		}

		/// Write out the binary
		binPath, err := os.Executable()
		if err != nil {
			util.Bail(err)
		}

		fullPath, err := filepath.EvalSymlinks(binPath)
		if err != nil {
			util.Bail(err)
		}
		if !util.JSON {
			fmt.Fprintf(
				os.Stderr,
				"Detected local binary path: %s\n",
				fullPath,
			)
		}
		existingStat, err := os.Lstat(fullPath)
		if err != nil {
			util.Bail(err)
		}
		// On sensible operating systems, we can't open and write to our
		// own binary, because it's in use. We can, however, move a file
		// into that place.

		newPath := fmt.Sprintf("%s-%s", fullPath, gh.SemVer)
		if !util.JSON {
			fmt.Fprintf(
				os.Stderr,
				"Writing to temp file '%s'\n",
				newPath,
			)
		}
		if err := ioutil.WriteFile(newPath, conchBin, existingStat.Mode()); err != nil {
			util.Bail(err)
		}

		if !util.JSON {
			fmt.Fprintf(
				os.Stderr,
				"Renaming '%s' to '%s'\n",
				newPath,
				fullPath,
			)
		}

		if err := os.Rename(newPath, fullPath); err != nil {
			util.Bail(err)
		}

		if !util.JSON {
			fmt.Fprintf(
				os.Stderr,
				"Successfully upgraded from %s to %s\n",
				util.SemVersion,
				gh.SemVer,
			)
		}

	}
}

func updaterDownloadFile(downloadURL string) (data []byte, err error) {
	if !util.JSON {
		fmt.Fprintf(
			os.Stderr,
			"Downloading '%s'\n",
			downloadURL,
		)
	}

	resp, err := http.Get(downloadURL)
	if err != nil {
		return data, err
	}

	if resp.StatusCode != 200 {
		return data, fmt.Errorf(
			"could not download '%s' (status %d)",
			downloadURL,
			resp.StatusCode,
		)
	}

	data, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return data, err
}
