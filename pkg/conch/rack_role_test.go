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

func TestRackRoleErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetRackRoles", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/rack_role").Persist().Reply(400).JSON(ErrApi)

		defer gock.Flush()

		ret, err := API.GetRackRoles()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.RackRole{})
	})

	t.Run("GetRackRole", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/rack_role/" + id.String()).Reply(400).JSON(ErrApi)

		ret, err := API.GetRackRole(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.RackRole{})
	})

	t.Run("CreateRackRole", func(t *testing.T) {
		r := conch.RackRole{
			Name:     "n",
			RackSize: 2,
		}

		gock.New(API.BaseURL).Post("/rack_role").Reply(400).JSON(ErrApi)

		err := API.SaveRackRole(&r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("UpdateRackRole", func(t *testing.T) {
		id := uuid.NewV4()
		r := conch.RackRole{
			ID:       id,
			Name:     "n",
			RackSize: 2,
		}

		gock.New(API.BaseURL).Post("/rack_role/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.SaveRackRole(&r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteRackRole", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/rack_role/" + id.String()).Reply(400).JSON(ErrApi)

		err := API.DeleteRackRole(id)
		st.Expect(t, err, ErrApiUnpacked)
	})

}
