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

func TestRelayErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetWorkspaceRelays", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/workspace/" + id.String() + "/relay").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetWorkspaceRelays(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.WorkspaceRelay{})
	})

	t.Run("RegisterRelay", func(t *testing.T) {
		id := uuid.NewV4()
		r := conch.WorkspaceRelay{ID: id.String(), SSHPort: 22, Version: "wat"}

		gock.New(API.BaseURL).Post("/relay/" + id.String() + "/register").
			Reply(400).JSON(ErrApi)

		err := API.RegisterRelay(r)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("GetAllRelays", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/relay").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetAllRelays()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.WorkspaceRelay{})
	})

}
