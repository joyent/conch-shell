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

func TestWorkspaceErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetWorkspaces", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/workspace").Reply(400).JSON(ErrApi)
		ret, err := API.GetWorkspaces()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Workspaces{})
	})

	t.Run("GetWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Get("/workspace/" + id.String()).
			Reply(400).JSON(ErrApi)

		ret, err := API.GetWorkspace(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Workspace{})
	})

	t.Run("GetWorkspaceByName", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/workspace/GLOBAL").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetWorkspaceByName("GLOBAL")
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Workspace{})
	})

	t.Run("GetSubWorkspaces", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Get("/workspace/" + id.String() + "/child").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetSubWorkspaces(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Workspaces{})
	})

	t.Run("GetWorkspaceUsers", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Get("/workspace/" + id.String() + "/user").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetWorkspaceUsers(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.WorkspaceUser{})
	})

	t.Run("CreateSubWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		id2 := uuid.NewV4()
		w := conch.Workspace{ID: id}
		s := conch.Workspace{ID: id2, Name: "test", Description: "test"}

		gock.New(API.BaseURL).Post("/workspace/" + id.String() + "/child").
			Reply(400).JSON(ErrApi)

		ret, err := API.CreateSubWorkspace(w, s)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, s)
	})

	t.Run("AddRackToWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		id2 := uuid.NewV4()

		gock.New(API.BaseURL).Post("/workspace/" + id.String() + "/rack").
			Reply(400).JSON(ErrApi)

		err := API.AddRackToWorkspace(id, id2)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteRackFromWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		id2 := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/workspace/" + id.String() + "/rack").
			Reply(400).JSON(ErrApi)

		err := API.DeleteRackFromWorkspace(id, id2)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("AddUserToWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Post("/workspace/" + id.String() + "/user").
			Reply(400).JSON(ErrApi)
		err := API.AddUserToWorkspace(id, "user", "role")
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("RemoveUserFromWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Delete("/workspace/" + id.String() + "/user").
			Reply(400).JSON(ErrApi)
		err := API.RemoveUserFromWorkspace(id, "user")
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("GetWorkspaceDevices", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/workspace/" + id.String() + "/device").
			Persist().Reply(400).JSON(ErrApi)

		ret, err := API.GetWorkspaceDevices(id, false, "g", "h", "T")
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Devices{})

		ret, err = API.GetWorkspaceDevices(id, true, "g", "h", "T")
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Devices{})

		ret, err = API.GetWorkspaceDevices(id, true, "g", "h", "T")
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Devices{})

		ret, err = API.GetWorkspaceDevices(id, true, "g", "h", "F")
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Devices{})

		gock.Flush()

	})

	t.Run("GetWorkspaceRacks", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/workspace/" + id.String() + "/rack").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetWorkspaceRacks(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.WorkspaceRack{})
	})

	t.Run("GetWorkspaceRack", func(t *testing.T) {
		id := uuid.NewV4()
		rID := uuid.NewV4()

		gock.New(API.BaseURL).
			Get("/workspace/" + id.String() + "/rack/" + rID.String()).
			Reply(400).JSON(ErrApi)

		ret, err := API.GetWorkspaceRack(id, rID)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.WorkspaceRack{})
	})

}
