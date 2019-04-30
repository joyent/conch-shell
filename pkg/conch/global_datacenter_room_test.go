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

func TestGlobalRoomErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetGlobalRooms", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/room").Reply(400).JSON(ErrApi)

		ret, err := API.GetGlobalRooms()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.GlobalRoom{})
	})

	t.Run("GetGlobalRoom", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/room/" + id.String()).Reply(400).JSON(ErrApi)

		ret, err := API.GetGlobalRoom(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.GlobalRoom{})
	})

	t.Run("CreateGlobalRoom", func(t *testing.T) {
		r := conch.GlobalRoom{
			DatacenterID: uuid.NewV4(),
			AZ:           "a",
			Alias:        "l",
			VendorName:   "v",
		}

		gock.New(API.BaseURL).Post("/room").Reply(400).JSON(ErrApi)

		err := API.SaveGlobalRoom(&r)
		st.Expect(t, err, ErrApiUnpacked)
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

		gock.New(API.BaseURL).Post("/room/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.SaveGlobalRoom(&r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteGlobalRoom", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/room/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.DeleteGlobalRoom(id)
		st.Expect(t, err, ErrApiUnpacked)
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
			Reply(400).JSON(ErrApi)

		ret, err := API.GetGlobalRoomRacks(r)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.GlobalRack{})
	})

}
