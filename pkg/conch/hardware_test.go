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

func TestHardwareErrors(t *testing.T) {
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

	t.Run("SaveHardwareProduct", func(t *testing.T) {
		gock.New(API.BaseURL).Persist().Post("/hardware_product").
			MatchType("json").Reply(400).JSON(ErrApi)

		hp := conch.HardwareProduct{}
		err := API.SaveHardwareProduct(&hp)
		st.Expect(t, err, conch.ErrBadInput)

		hp.Name = "test"
		err = API.SaveHardwareProduct(&hp)
		st.Expect(t, err, conch.ErrBadInput)

		hp.Alias = "test"
		err = API.SaveHardwareProduct(&hp)
		st.Expect(t, err, conch.ErrBadInput)

		hp.HardwareVendorID = uuid.NewV4()

		err = API.SaveHardwareProduct(&hp)
		st.Expect(t, err, ErrApiUnpacked)

		hp.ID = uuid.NewV4()

		err = API.SaveHardwareProduct(&hp)
		st.Expect(t, err, ErrApiUnpacked)

		gock.Flush()
	})

	t.Run("SaveHardwareVendor", func(t *testing.T) {
		gock.New(API.BaseURL).Persist().Post("/hardware_vendor").
			MatchType("json").Reply(400).JSON(ErrApi)

		hv := conch.HardwareVendor{}
		err := API.SaveHardwareVendor(&hv)
		st.Expect(t, err, conch.ErrBadInput)

		hv.Name = "test"
		err = API.SaveHardwareVendor(&hv)
		st.Expect(t, err, ErrApiUnpacked)

		hv2 := conch.HardwareVendor{ID: uuid.NewV4()}
		err = API.SaveHardwareVendor(&hv2)
		st.Expect(t, err, conch.ErrBadInput)

		gock.Flush()
	})

	t.Run("DeleteHardwareProduct", func(t *testing.T) {
		gock.New(API.BaseURL).Persist().Delete("/hardware_product").
			Reply(400).JSON(ErrApi)

		err := API.DeleteHardwareProduct(uuid.NewV4())
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteHardwareVendor", func(t *testing.T) {
		gock.New(API.BaseURL).Persist().Delete("/hardware_vendor").
			Reply(400).JSON(ErrApi)

		err := API.DeleteHardwareVendor("vendor")
		st.Expect(t, err, ErrApiUnpacked)
	})

}
