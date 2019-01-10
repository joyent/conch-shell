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
	uuid "gopkg.in/satori/go.uuid.v1"
)

func TestHardwareVendorErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()
	name := "hardware vendor"

	t.Run("GetHardwareVendor", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/hardware_vendor/" + name).
			Reply(400).JSON(ErrApi)

		ret, err := API.GetHardwareVendor(name)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.HardwareVendor{})
	})

	t.Run("GetHardwareVendors", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/hardware_vendor").Reply(400).JSON(ErrApi)

		ret, err := API.GetHardwareVendors()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.HardwareVendor{})
	})

	t.Run("DeleteHardwareVendor", func(t *testing.T) {
		gock.New(API.BaseURL).Delete("/hardware_vendor/" + name).
			Reply(400).JSON(ErrApi)

		err := API.DeleteHardwareVendor(name)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("SaveHardwareVendor", func(t *testing.T) {
		v := conch.HardwareVendor{
			Name: name,
		}

		gock.New(API.BaseURL).Post("/hardware_vendor/" + name).
			Reply(400).JSON(ErrApi)

		err := API.SaveHardwareVendor(&v)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("GetHardwareProducts", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/hardware_product").Reply(400).JSON(ErrApi)

		ret, err := API.GetHardwareProducts()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.HardwareProduct{})
	})

	t.Run("GetHardwareProduct", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).
			Get("/hardware_product/" + id.String()).
			Reply(400).JSON(ErrApi)

		ret, err := API.GetHardwareProduct(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.HardwareProduct{})
	})

	// BUG(sungo): a lot of hardware product stuff is totally untested
}
