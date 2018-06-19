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

func TestGlobalRoomErrors(t *testing.T) {
	BuildAPI()
	gock.Flush()

	aerr := conch.APIError{ErrorMsg: "totally broken"}
	aerrUnpacked := errors.New(aerr.ErrorMsg)

	t.Run("GetGlobalRooms", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/room").Persist().Reply(400).JSON(aerr)

		defer gock.Flush()

		ret, err := API.GetGlobalRooms()
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.GlobalRoom{})
	})

	t.Run("GetGlobalRoom", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/room/" + id.String()).Reply(400).JSON(aerr)

		ret, err := API.GetGlobalRoom(id)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, conch.GlobalRoom{})
	})

	t.Run("CreateGlobalRoom", func(t *testing.T) {
		r := conch.GlobalRoom{
			DatacenterID: uuid.NewV4(),
			AZ:           "a",
			Alias:        "l",
			VendorName:   "v",
		}

		gock.New(API.BaseURL).Post("/room").Reply(400).JSON(aerr)

		err := API.SaveGlobalRoom(&r)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("UpdateGlobalRoom", func(t *testing.T) {
		id := uuid.NewV4()

		r := conch.GlobalRoom{
			ID:           id,
			DatacenterID: uuid.NewV4(),
			AZ:           "a",
			Alias:        "l",
			VendorName:   "v",
		}

		gock.New(API.BaseURL).Post("/room/" + id.String()).Reply(400).JSON(aerr)

		err := API.SaveGlobalRoom(&r)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("DeleteGlobalRoom", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/room/" + id.String()).Reply(400).JSON(aerr)

		err := API.DeleteGlobalRoom(id)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("GetGlobalRoomRacks", func(t *testing.T) {
		id := uuid.NewV4()

		r := conch.GlobalRoom{
			ID:           id,
			DatacenterID: uuid.NewV4(),
			AZ:           "a",
			Alias:        "l",
			VendorName:   "v",
		}

		gock.New(API.BaseURL).Get("/room/" + id.String() + "/racks").
			Persist().Reply(400).JSON(aerr)

		defer gock.Flush()

		ret, err := API.GetGlobalRoomRacks(r)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.GlobalRack{})
	})

}
