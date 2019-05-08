// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch_test

import (
	"testing"

	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func TestDevices(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	serial := "test"
	d := conch.Device{ID: serial}

	t.Run("GetDevice", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/device/" + serial).Reply(200).JSON(d)

		ret, err := API.GetDevice(serial)
		st.Expect(t, err, nil)
		st.Expect(t, ret, conch.Device{ID: serial})
	})

	t.Run("GetDeviceErrors", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/device/" + serial).Reply(400).JSON(ErrApi)

		ret, err := API.GetDevice(serial)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Device{ID: serial})
	})

	t.Run("FillInDeviceErrors", func(t *testing.T) {
		serial := "test"
		d := conch.Device{ID: serial}

		gock.New(API.BaseURL).Get("/device/" + serial).Reply(400).JSON(ErrApi)

		ret, err := API.FillInDevice(d)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, d)
	})

	t.Run("GetDeviceLocationErrors", func(t *testing.T) {
		serial := "test"

		gock.New(API.BaseURL).Get("/device/" + serial + "/location").Reply(400).JSON(ErrApi)

		ret, err := API.GetDeviceLocation(serial)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.DeviceLocation{})
	})

	t.Run("GraduateDeviceErrors", func(t *testing.T) {
		serial := "test"
		gock.New(API.BaseURL).Post("/device/" + serial + "/graduate").Reply(400).JSON(ErrApi)

		err := API.GraduateDevice(serial)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeviceTritonRebootErrors", func(t *testing.T) {
		serial := "test"
		gock.New(API.BaseURL).Post("/device/" + serial + "/triton_reboot").Reply(400).JSON(ErrApi)

		err := API.DeviceTritonReboot(serial)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("SetDeviceTritonUUIDErrors", func(t *testing.T) {
		serial := "test"
		id := uuid.NewV4()
		gock.New(API.BaseURL).Post("/device/" + serial + "/triton_uuid").Reply(400).JSON(ErrApi)

		err := API.SetDeviceTritonUUID(serial, id)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("MarkDeviceTritonSetupErrors", func(t *testing.T) {
		serial := "test"
		gock.New(API.BaseURL).Post("/device/" + serial + "/triton_setup").Reply(400).JSON(ErrApi)

		err := API.MarkDeviceTritonSetup(serial)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("SetDeviceAssetTagErrors", func(t *testing.T) {
		serial := "test"
		tag := "tag"
		gock.New(API.BaseURL).Post("/device/" + serial + "/asset_tag").Reply(400).JSON(ErrApi)

		err := API.SetDeviceAssetTag(serial, tag)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("GetDevicesByField", func(t *testing.T) {
		var d conch.Devices

		gock.New(API.BaseURL).Get("/device").MatchParam("hostname", "bar").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetDevicesByField("hostname", "bar")
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, d)
	})

	t.Run("SubmitDeviceReport", func(t *testing.T) {
		serial := "test"
		gock.New(API.BaseURL).Post("/device/" + serial).Reply(400).JSON(ErrApi)

		ret, err := API.SubmitDeviceReport(serial, "")
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.ValidationState{})
	})

	t.Run("GetDevicePhase", func(t *testing.T) {
		serial := "test"
		gock.New(API.BaseURL).Get("/device/" + serial + "/phase").Reply(400).JSON(ErrApi)

		ret, err := API.GetDevicePhase(serial)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, "")
	})

	t.Run("SetDevicePhase", func(t *testing.T) {
		serial := "test"
		gock.New(API.BaseURL).Post("/device/" + serial + "/phase").Reply(400).JSON(ErrApi)

		err := API.SetDevicePhase(serial, "production")
		st.Expect(t, err, ErrApiUnpacked)
	})

}
