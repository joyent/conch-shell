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

func TestGlobalRackErrors(t *testing.T) {
	BuildAPI()
	gock.Flush()

	aerr := struct {
		ErrorMsg string `json:"error"`
	}{"totally broken"}
	aerrUnpacked := errors.New(aerr.ErrorMsg)

	t.Run("GetGlobalRacks", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/rack").Persist().Reply(400).JSON(aerr)

		defer gock.Flush()

		ret, err := API.GetGlobalRacks()
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.GlobalRack{})
	})

	t.Run("GetGlobalRack", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/rack/" + id.String()).Reply(400).JSON(aerr)

		ret, err := API.GetGlobalRack(id)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, conch.GlobalRack{})
	})

	t.Run("CreateGlobalRack", func(t *testing.T) {
		r := conch.GlobalRack{
			DatacenterRoomID: uuid.NewV4(),
			RoleID:           uuid.NewV4(),
			Name:             "n",
		}

		gock.New(API.BaseURL).Post("/rack").Reply(400).JSON(aerr)

		err := API.SaveGlobalRack(&r)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("UpdateGlobalRack", func(t *testing.T) {
		id := uuid.NewV4()
		r := conch.GlobalRack{
			ID:               id,
			DatacenterRoomID: uuid.NewV4(),
			RoleID:           uuid.NewV4(),
			Name:             "n",
		}

		gock.New(API.BaseURL).Post("/rack/" + id.String()).Reply(400).JSON(aerr)

		err := API.SaveGlobalRack(&r)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("DeleteGlobalRack", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/rack/" + id.String()).Reply(400).JSON(aerr)

		err := API.DeleteGlobalRack(id)
		st.Expect(t, err, aerrUnpacked)
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
			Reply(400).JSON(aerr)

		ret, err := API.GetGlobalRackLayout(r)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.GlobalRackLayoutSlot{})
	})

}
