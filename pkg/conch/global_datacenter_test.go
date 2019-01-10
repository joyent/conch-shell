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

func TestGlobalDatacenterErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetGlobalDatacenters", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/dc").Reply(400).JSON(ErrApi)

		ret, err := API.GetGlobalDatacenters()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.GlobalDatacenter{})
	})

	t.Run("GetGlobalDatacenter", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/dc/" + id.String()).Reply(400).JSON(ErrApi)

		ret, err := API.GetGlobalDatacenter(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.GlobalDatacenter{})
	})

	t.Run("CreateGlobalDatacenter", func(t *testing.T) {
		d := conch.GlobalDatacenter{
			Region:   "r",
			Vendor:   "v",
			Location: "l",
		}

		gock.New(API.BaseURL).Post("/dc").Reply(400).JSON(ErrApi)

		err := API.SaveGlobalDatacenter(&d)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("UpdateGlobalDatacenter", func(t *testing.T) {
		id := uuid.NewV4()
		d := conch.GlobalDatacenter{
			ID:       id,
			Region:   "r",
			Vendor:   "v",
			Location: "l",
		}

		gock.New(API.BaseURL).Post("/dc/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.SaveGlobalDatacenter(&d)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteGlobalDatacenter", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/dc/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.DeleteGlobalDatacenter(id)
		st.Expect(t, err, ErrApiUnpacked)
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
			Reply(400).JSON(ErrApi)

		ret, err := API.GetGlobalDatacenterRooms(d)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.GlobalRoom{})
	})

}
