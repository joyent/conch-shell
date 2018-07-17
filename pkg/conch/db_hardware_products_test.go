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

func TestDBHardwareProductsErrors(t *testing.T) {
	BuildAPI()
	gock.Flush()

	aerr := conch.APIError{ErrorMsg: "totally broken"}
	aerrUnpacked := errors.New(aerr.ErrorMsg)

	t.Run("GetDBHardwareProduct", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).
			Get("/db/hardware_product/" + id.String()).
			Reply(400).JSON(aerr)

		ret, err := API.GetDBHardwareProduct(id)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, conch.DBHardwareProduct{})
	})

	t.Run("GetDBHardwareProducts", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/db/hardware_product").Reply(400).JSON(aerr)

		ret, err := API.GetDBHardwareProducts()
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.DBHardwareProduct{})
	})

	t.Run("DeleteDBHardwareProduct", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).
			Delete("/db/hardware_product/" + id.String()).
			Reply(400).JSON(aerr)

		err := API.DeleteDBHardwareProduct(id)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("SaveDBHardwareProduct Create", func(t *testing.T) {
		h := conch.DBHardwareProduct{
			Name:   "tester",
			Alias:  "tester",
			Prefix: "tester",
			Vendor: uuid.NewV4(),
			SKU:    "tester",
		}

		gock.New(API.BaseURL).
			Post("/db/hardware_product").
			Reply(400).JSON(aerr)

		err := API.SaveDBHardwareProduct(&h)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("SaveDBHardwareProduct Update", func(t *testing.T) {
		id := uuid.NewV4()
		h := conch.DBHardwareProduct{
			ID:     id,
			Name:   "tester",
			Alias:  "tester",
			Prefix: "tester",
			Vendor: uuid.NewV4(),
			SKU:    "tester",
		}

		gock.New(API.BaseURL).
			Post("/db/hardware_product/" + id.String()).
			Reply(400).JSON(aerr)

		err := API.SaveDBHardwareProduct(&h)
		st.Expect(t, err, aerrUnpacked)
	})

}
