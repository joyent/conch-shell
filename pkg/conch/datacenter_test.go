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

func TestDatacenterErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetDatacenters", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/dc").Reply(400).JSON(ErrApi)

		ret, err := API.GetDatacenters()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.Datacenter{})
	})

	t.Run("GetDatacenter", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/dc/" + id.String()).Reply(400).JSON(ErrApi)

		ret, err := API.GetDatacenter(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Datacenter{})
	})

	t.Run("CreateDatacenter", func(t *testing.T) {
		d := conch.Datacenter{
			Region:   "r",
			Vendor:   "v",
			Location: "l",
		}

		gock.New(API.BaseURL).Post("/dc").Reply(400).JSON(ErrApi)

		err := API.SaveDatacenter(&d)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("UpdateDatacenter", func(t *testing.T) {
		id := uuid.NewV4()
		d := conch.Datacenter{
			ID:       id,
			Region:   "r",
			Vendor:   "v",
			Location: "l",
		}

		gock.New(API.BaseURL).Post("/dc/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.SaveDatacenter(&d)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteDatacenter", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/dc/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.DeleteDatacenter(id)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("GetDatacenterRooms", func(t *testing.T) {
		id := uuid.NewV4()
		d := conch.Datacenter{
			ID:       id,
			Region:   "r",
			Vendor:   "v",
			Location: "l",
		}

		gock.New(API.BaseURL).Get("/dc/" + id.String() + "/rooms").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetDatacenterRooms(d)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.Room{})
	})

}
