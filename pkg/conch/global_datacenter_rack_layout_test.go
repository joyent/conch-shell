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

func TestGlobalRackLayoutSlotErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetGlobalRackLayoutSlots", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/layout").Persist().Reply(400).JSON(ErrApi)

		ret, err := API.GetGlobalRackLayoutSlots()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.GlobalRackLayoutSlots{})
	})

	t.Run("GetGlobalRackLayoutSlot", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/layout/" + id.String()).Reply(400).JSON(ErrApi)

		ret, err := API.GetGlobalRackLayoutSlot(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, &conch.GlobalRackLayoutSlot{})
	})

	t.Run("CreateGlobalRackLayoutSlot", func(t *testing.T) {
		r := conch.GlobalRackLayoutSlot{
			RackID:    uuid.NewV4(),
			ProductID: uuid.NewV4(),
			RUStart:   2,
		}

		gock.New(API.BaseURL).Post("/layout").Reply(400).JSON(ErrApi)

		err := API.SaveGlobalRackLayoutSlot(&r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("UpdateGlobalRackLayoutSlot", func(t *testing.T) {
		id := uuid.NewV4()
		r := conch.GlobalRackLayoutSlot{
			ID:        id,
			RackID:    uuid.NewV4(),
			ProductID: uuid.NewV4(),
			RUStart:   3,
		}

		gock.New(API.BaseURL).Post("/layout/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.SaveGlobalRackLayoutSlot(&r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteGlobalRackLayoutSlot", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/layout/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.DeleteGlobalRackLayoutSlot(id)
		st.Expect(t, err, ErrApiUnpacked)
	})

}
