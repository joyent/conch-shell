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

func TestRackLayoutSlotErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetRackLayoutSlots", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/layout").Persist().Reply(400).JSON(ErrApi)

		ret, err := API.GetRackLayoutSlots()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.RackLayoutSlots{})
	})

	t.Run("GetRackLayoutSlot", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/layout/" + id.String()).Reply(400).JSON(ErrApi)

		ret, err := API.GetRackLayoutSlot(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, &conch.RackLayoutSlot{})
	})

	t.Run("CreateRackLayoutSlot", func(t *testing.T) {
		r := conch.RackLayoutSlot{
			RackID:    uuid.NewV4(),
			ProductID: uuid.NewV4(),
			RUStart:   2,
		}

		gock.New(API.BaseURL).Post("/layout").Reply(400).JSON(ErrApi)

		err := API.SaveRackLayoutSlot(&r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("UpdateRackLayoutSlot", func(t *testing.T) {
		id := uuid.NewV4()
		r := conch.RackLayoutSlot{
			ID:        id,
			RackID:    uuid.NewV4(),
			ProductID: uuid.NewV4(),
			RUStart:   3,
		}

		gock.New(API.BaseURL).Post("/layout/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.SaveRackLayoutSlot(&r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteRackLayoutSlot", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/layout/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.DeleteRackLayoutSlot(id)
		st.Expect(t, err, ErrApiUnpacked)
	})

}
