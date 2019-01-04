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
	uuid "gopkg.in/satori/go.uuid.v1"
	"testing"
)

func TestDevicesErrors(t *testing.T) {
	BuildAPI()
	gock.Flush()

	aerr := struct {
		ErrorMsg string `json:"error"`
	}{"totally broken"}
	aerrUnpacked := errors.New(aerr.ErrorMsg)

	t.Run("GetDevice", func(t *testing.T) {
		serial := "test"

		gock.New(API.BaseURL).Get("/device/" + serial).
			Reply(400).JSON(aerr)

		ret, err := API.GetDevice(serial)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, conch.Device{ID: serial})
	})

	t.Run("FillInDevice", func(t *testing.T) {
		serial := "test"
		d := conch.Device{ID: serial}

		gock.New(API.BaseURL).Get("/device/" + serial).
			Reply(400).JSON(aerr)

		ret, err := API.FillInDevice(d)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, d)
	})

	t.Run("GetDeviceLocation", func(t *testing.T) {
		serial := "test"

		gock.New(API.BaseURL).Get("/device/" + serial + "/location").
			Reply(400).JSON(aerr)

		ret, err := API.GetDeviceLocation(serial)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, conch.DeviceLocation{})
	})

	t.Run("GraduateDevice", func(t *testing.T) {
		serial := "test"
		gock.New(API.BaseURL).Post("/device/" + serial + "/graduate").
			Reply(400).JSON(aerr)

		err := API.GraduateDevice(serial)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("DeviceTritonReboot", func(t *testing.T) {
		serial := "test"
		gock.New(API.BaseURL).Post("/device/" + serial + "/triton_reboot").
			Reply(400).JSON(aerr)

		err := API.DeviceTritonReboot(serial)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("SetDeviceTritonUUID", func(t *testing.T) {
		serial := "test"
		id := uuid.NewV4()
		gock.New(API.BaseURL).Post("/device/" + serial + "/triton_uuid").
			Reply(400).JSON(aerr)

		err := API.SetDeviceTritonUUID(serial, id)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("MarkDeviceTritonSetup", func(t *testing.T) {
		serial := "test"
		gock.New(API.BaseURL).Post("/device/" + serial + "/triton_setup").
			Reply(400).JSON(aerr)

		err := API.MarkDeviceTritonSetup(serial)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("SetDeviceAssetTag", func(t *testing.T) {
		serial := "test"
		tag := "tag"
		gock.New(API.BaseURL).Post("/device/" + serial + "/asset_tag").
			Reply(400).JSON(aerr)

		err := API.SetDeviceAssetTag(serial, tag)
		st.Expect(t, err, aerrUnpacked)
	})
}
