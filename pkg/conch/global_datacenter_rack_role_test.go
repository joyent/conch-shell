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

func TestGlobalRackRoleErrors(t *testing.T) {
	BuildAPI()
	gock.Flush()

	aerr := struct {
		ErrorMsg string `json:"error"`
	}{"totally broken"}
	aerrUnpacked := errors.New(aerr.ErrorMsg)

	t.Run("GetGlobalRackRoles", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/rack_role").Persist().Reply(400).JSON(aerr)

		defer gock.Flush()

		ret, err := API.GetGlobalRackRoles()
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.GlobalRackRole{})
	})

	t.Run("GetGlobalRackRole", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/rack_role/" + id.String()).Reply(400).JSON(aerr)

		ret, err := API.GetGlobalRackRole(id)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, conch.GlobalRackRole{})
	})

	t.Run("CreateGlobalRackRole", func(t *testing.T) {
		r := conch.GlobalRackRole{
			Name:     "n",
			RackSize: 2,
		}

		gock.New(API.BaseURL).Post("/rack_role").Reply(400).JSON(aerr)

		err := API.SaveGlobalRackRole(&r)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("UpdateGlobalRackRole", func(t *testing.T) {
		id := uuid.NewV4()
		r := conch.GlobalRackRole{
			ID:       id,
			Name:     "n",
			RackSize: 2,
		}

		gock.New(API.BaseURL).Post("/rack_role/" + id.String()).Reply(400).JSON(aerr)

		err := API.SaveGlobalRackRole(&r)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("DeleteGlobalRackRole", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/rack_role/" + id.String()).Reply(400).JSON(aerr)

		err := API.DeleteGlobalRackRole(id)
		st.Expect(t, err, aerrUnpacked)
	})

}
