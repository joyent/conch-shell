// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch_test

import (
	"errors"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func TestErrors(t *testing.T) {
	BuildAPI()
	gock.Flush()

	type APIError struct {
		ErrorMsg string `json:"error"`
	}

	aerr := struct {
		ErrorMsg string `json:"error"`
	}{"totally broken"}
	aerrUnpacked := errors.New(aerr.ErrorMsg)

	fourohfour := APIError{ErrorMsg: "Not found"}
	fourohthree := APIError{ErrorMsg: "Forbidden"}
	fourohone := APIError{ErrorMsg: "Not Authorized"}

	url := "/user/me/settings"
	gock.New(API.BaseURL).Get(url).Reply(404).JSON(fourohfour)
	_, err := API.GetUserSettings()
	st.Expect(t, err, conch.ErrDataNotFound)

	gock.New(API.BaseURL).Get(url).Reply(403).JSON(fourohthree)
	_, err = API.GetUserSettings()
	st.Expect(t, err, conch.ErrForbidden)

	gock.New(API.BaseURL).Get(url).Reply(401).JSON(fourohone)
	_, err = API.GetUserSettings()
	st.Expect(t, err, conch.ErrNotAuthorized)

	gock.New(API.BaseURL).Get(url).Reply(400).JSON(aerr)
	_, err = API.GetUserSettings()
	st.Expect(t, err, aerrUnpacked)

	gock.New(API.BaseURL).Get(url).Reply(500).JSON(aerr)
	_, err = API.GetUserSettings()
	st.Expect(t, err, aerrUnpacked)
}
