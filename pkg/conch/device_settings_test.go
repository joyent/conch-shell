// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch_test

import (
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func TestDeviceSettingsErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetDeviceSettings", func(t *testing.T) {
		serial := "test"

		gock.New(API.BaseURL).Get("/device/" + serial + "/settings").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetDeviceSettings(serial)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, make(map[string]string))
	})

	t.Run("GetDeviceSetting", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Get("/device/" + serial + "/settings/" + key).
			Reply(400).JSON(ErrApi)

		ret, err := API.GetDeviceSetting(serial, key)
		st.Expect(t, err, ErrApiUnpacked)
		var setting string
		st.Expect(t, ret, setting)
	})

	t.Run("SetDeviceSetting", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Post("/device/" + serial + "/settings/" + key).
			Reply(400).JSON(ErrApi)

		err := API.SetDeviceSetting(serial, key, "val")
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteDeviceSetting", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Delete("/device/" + serial + "/settings/" + key).
			Reply(400).JSON(ErrApi)

		err := API.DeleteDeviceSetting(serial, key)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("GetDeviceTags", func(t *testing.T) {
		serial := "test"

		gock.New(API.BaseURL).Get("/device/" + serial + "/settings").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetDeviceTags(serial)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, make(map[string]string))
	})

	t.Run("GetDeviceTag", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Get("/device/" + serial + "/settings/tag." + key).
			Reply(400).JSON(ErrApi)

		ret, err := API.GetDeviceTag(serial, key)
		st.Expect(t, err, ErrApiUnpacked)
		var setting string
		st.Expect(t, ret, setting)
	})

	t.Run("SetDeviceTag", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Post("/device/" + serial + "/settings/tag." + key).
			Reply(400).JSON(ErrApi)

		err := API.SetDeviceTag(serial, key, "val")
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteDeviceTag", func(t *testing.T) {
		serial := "test"
		key := "key"

		gock.New(API.BaseURL).Delete("/device/" + serial + "/settings/tag." + key).
			Reply(400).JSON(ErrApi)

		err := API.DeleteDeviceTag(serial, key)
		st.Expect(t, err, ErrApiUnpacked)
	})

}
