// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package util contains common routines used throughout the command base
package util

import (
	"os/user"
	"strconv"
)

func UserIsRoot() bool {
	current, err := user.Current()
	if err != nil {
		Bail(err)
	}

	uid, err := strconv.Atoi(current.Uid)
	if err != nil {
		Bail(err)
	}

	return (uid == 0)
}
