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

func TestDeviceSettingsErrors(t *testing.T) {
	BuildAPI()
	gock.Flush()

	aerr := conch.APIError{ErrorMsg: "totally broken"}
	aerrUnpacked := errors.New(aerr.ErrorMsg)

	t.Run("GetDeviceSettings", func(t *testing.T) {
		serial := "test"

		gock.New(API.BaseURL).Get("/device/" + serial + "/settings").
			Reply(400).JSON(aerr)

		ret, err := API.GetDeviceSettings(serial)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, make(map[string]string))
	})

	t.Run("GetDeviceSetting", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Get("/device/" + serial + "/settings/" + key).
			Reply(400).JSON(aerr)

		ret, err := API.GetDeviceSetting(serial, key)
		st.Expect(t, err, aerrUnpacked)
		var setting string
		st.Expect(t, ret, setting)
	})

	t.Run("SetDeviceSetting", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Post("/device/" + serial + "/settings/" + key).
			Reply(400).JSON(aerr)

		err := API.SetDeviceSetting(serial, key, "val")
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("DeleteDeviceSetting", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Delete("/device/" + serial + "/settings/" + key).
			Reply(400).JSON(aerr)

		err := API.DeleteDeviceSetting(serial, key)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("GetDeviceTags", func(t *testing.T) {
		serial := "test"

		gock.New(API.BaseURL).Get("/device/" + serial + "/settings").
			Reply(400).JSON(aerr)

		ret, err := API.GetDeviceTags(serial)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, make(map[string]string))
	})

	t.Run("GetDeviceTag", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Get("/device/" + serial + "/settings/tag." + key).
			Reply(400).JSON(aerr)

		ret, err := API.GetDeviceTag(serial, key)
		st.Expect(t, err, aerrUnpacked)
		var setting string
		st.Expect(t, ret, setting)
	})

	t.Run("SetDeviceTag", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Post("/device/" + serial + "/settings/tag." + key).
			Reply(400).JSON(aerr)

		err := API.SetDeviceTag(serial, key, "val")
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("DeleteDeviceTag", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Delete("/device/" + serial + "/settings/tag." + key).
			Reply(400).JSON(aerr)

		err := API.DeleteDeviceTag(serial, key)
		st.Expect(t, err, aerrUnpacked)
	})

}
