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

func TestGlobalRackErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetGlobalRacks", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/rack").Reply(400).JSON(ErrApi)

		ret, err := API.GetGlobalRacks()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.GlobalRack{})
	})

	t.Run("GetGlobalRack", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/rack/" + id.String()).Reply(400).JSON(ErrApi)

		ret, err := API.GetGlobalRack(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.GlobalRack{})
	})

	t.Run("CreateGlobalRack", func(t *testing.T) {
		r := conch.GlobalRack{
			DatacenterRoomID: uuid.NewV4(),
			RoleID:           uuid.NewV4(),
			Name:             "n",
		}

		gock.New(API.BaseURL).Post("/rack").Reply(400).JSON(ErrApi)

		err := API.SaveGlobalRack(&r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("UpdateGlobalRack", func(t *testing.T) {
		id := uuid.NewV4()
		r := conch.GlobalRack{
			ID:               id,
			DatacenterRoomID: uuid.NewV4(),
			RoleID:           uuid.NewV4(),
			Name:             "n",
		}

		gock.New(API.BaseURL).Post("/rack/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.SaveGlobalRack(&r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteGlobalRack", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/rack/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.DeleteGlobalRack(id)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("GetGlobalRackLayout", func(t *testing.T) {
		id := uuid.NewV4()
		r := conch.GlobalRack{
			ID:               id,
			DatacenterRoomID: uuid.NewV4(),
			RoleID:           uuid.NewV4(),
			Name:             "n",
		}

		gock.New(API.BaseURL).Get("/rack/" + id.String() + "/layouts").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetGlobalRackLayout(r)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.GlobalRackLayoutSlots{})
	})

}
