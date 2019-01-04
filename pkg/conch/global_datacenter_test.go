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

func TestGlobalDatacenterErrors(t *testing.T) {
	BuildAPI()
	gock.Flush()

	aerr := struct {
		ErrorMsg string `json:"error"`
	}{"totally broken"}
	aerrUnpacked := errors.New(aerr.ErrorMsg)

	t.Run("GetGlobalDatacenters", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/dc").Persist().Reply(400).JSON(aerr)

		defer gock.Flush()

		ret, err := API.GetGlobalDatacenters()
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.GlobalDatacenter{})
	})

	t.Run("GetGlobalDatacenter", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/dc/" + id.String()).Reply(400).JSON(aerr)

		ret, err := API.GetGlobalDatacenter(id)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, conch.GlobalDatacenter{})
	})

	t.Run("CreateGlobalDatacenter", func(t *testing.T) {
		d := conch.GlobalDatacenter{
			Region:   "r",
			Vendor:   "v",
			Location: "l",
		}

		gock.New(API.BaseURL).Post("/dc").Reply(400).JSON(aerr)

		err := API.SaveGlobalDatacenter(&d)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("UpdateGlobalDatacenter", func(t *testing.T) {
		id := uuid.NewV4()
		d := conch.GlobalDatacenter{
			ID:       id,
			Region:   "r",
			Vendor:   "v",
			Location: "l",
		}

		gock.New(API.BaseURL).Post("/dc/" + id.String()).Reply(400).JSON(aerr)

		err := API.SaveGlobalDatacenter(&d)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("DeleteGlobalDatacenter", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/dc/" + id.String()).Reply(400).JSON(aerr)

		err := API.DeleteGlobalDatacenter(id)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("GetGlobalDatacenterRooms", func(t *testing.T) {
		id := uuid.NewV4()
		d := conch.GlobalDatacenter{
			ID:       id,
			Region:   "r",
			Vendor:   "v",
			Location: "l",
		}

		gock.New(API.BaseURL).Get("/dc/" + id.String() + "/rooms").
			Reply(400).JSON(aerr)

		ret, err := API.GetGlobalDatacenterRooms(d)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.GlobalRoom{})
	})

}
