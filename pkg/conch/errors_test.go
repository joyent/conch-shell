// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch_test

import (
	"testing"

	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func TestErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	url := "/user/me/settings"
	gock.New(API.BaseURL).Get(url).Reply(404).JSON(ErrApi)
	_, err := API.GetUserSettings()
	st.Expect(t, err, conch.ErrDataNotFound)

	gock.New(API.BaseURL).Get(url).Reply(403).JSON(ErrApi)
	_, err = API.GetUserSettings()
	st.Expect(t, err, conch.ErrForbidden)

	gock.New(API.BaseURL).Get(url).Reply(401).JSON(ErrApi)
	_, err = API.GetUserSettings()
	st.Expect(t, err, conch.ErrNotAuthorized)

	gock.New(API.BaseURL).Get(url).Reply(400).JSON(ErrApi)
	_, err = API.GetUserSettings()
	st.Expect(t, err, ErrApiUnpacked)

	gock.New(API.BaseURL).Get(url).Reply(500).JSON(ErrApi)
	_, err = API.GetUserSettings()
	st.Expect(t, err, ErrApiUnpacked)
}
