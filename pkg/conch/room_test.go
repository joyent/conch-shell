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

func TestRoomErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetRooms", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/room").Reply(400).JSON(ErrApi)

		ret, err := API.GetRooms()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.Room{})
	})

	t.Run("GetRoom", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/room/" + id.String()).Reply(400).JSON(ErrApi)

		ret, err := API.GetRoom(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Room{})
	})

	t.Run("CreateRoom", func(t *testing.T) {
		r := conch.Room{
			DatacenterID: uuid.NewV4(),
			AZ:           "a",
			Alias:        "l",
			VendorName:   "v",
		}

		gock.New(API.BaseURL).Post("/room").Reply(400).JSON(ErrApi)

		err := API.SaveRoom(&r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("UpdateRoom", func(t *testing.T) {
		id := uuid.NewV4()

		r := conch.Room{
			ID:           id,
			DatacenterID: uuid.NewV4(),
			AZ:           "a",
			Alias:        "l",
			VendorName:   "v",
		}

		gock.New(API.BaseURL).Post("/room/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.SaveRoom(&r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteRoom", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/room/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.DeleteRoom(id)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("GetRoomRacks", func(t *testing.T) {
		id := uuid.NewV4()

		r := conch.Room{
			ID:           id,
			DatacenterID: uuid.NewV4(),
			AZ:           "a",
			Alias:        "l",
			VendorName:   "v",
		}

		gock.New(API.BaseURL).Get("/room/" + id.String() + "/racks").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetRoomRacks(r)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.Rack{})
	})

}
