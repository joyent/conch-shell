// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch_test

import (
	"testing"

	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func TestRackErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("SetRackPhase", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Post("/rack/" + id.String() + "/phase").
			Reply(400).JSON(ErrApi)

		err := API.SetRackPhase(id, "wat", true)
		st.Expect(t, err, ErrApiUnpacked)

		gock.New(API.BaseURL).Post("/rack/"+id.String()+"/phase").
			MatchParam("rack_only", "1").Reply(400).JSON(ErrApi)

		err = API.SetRackPhase(id, "wat", false)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("GetRackPhase", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/rack/" + id.String()).Reply(400).JSON(ErrApi)

		ret, err := API.GetRackPhase(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, "")
	})
}
