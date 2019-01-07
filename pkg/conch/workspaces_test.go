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

func TestWorkspaceErrors(t *testing.T) {
	BuildAPI()
	gock.Flush()

	aerr := struct {
		ErrorMsg string `json:"error"`
	}{"totally broken"}
	aerrUnpacked := errors.New(aerr.ErrorMsg)

	t.Run("GetWorkspaces", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/workspace").Reply(400).JSON(aerr)
		ret, err := API.GetWorkspaces()
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.Workspace{})
	})

	t.Run("GetWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Get("/workspace/" + id.String()).
			Reply(400).JSON(aerr)

		ret, err := API.GetWorkspace(id)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, conch.Workspace{})
	})

	t.Run("GetSubWorkspaces", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Get("/workspace/" + id.String() + "/child").
			Reply(400).JSON(aerr)

		ret, err := API.GetSubWorkspaces(id)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.Workspace{})
	})

	t.Run("GetWorkspaceUsers", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Get("/workspace/" + id.String() + "/user").
			Reply(400).JSON(aerr)

		ret, err := API.GetWorkspaceUsers(id)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.WorkspaceUser{})
	})

	t.Run("GetWorkspaceRooms", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Get("/workspace/" + id.String() + "/room").
			Reply(400).JSON(aerr)

		ret, err := API.GetWorkspaceRooms(id)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.Room{})
	})

	t.Run("CreateSubWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		id2 := uuid.NewV4()
		w := conch.Workspace{ID: id}
		s := conch.Workspace{ID: id2, Name: "test", Description: "test"}

		gock.New(API.BaseURL).Post("/workspace/" + id.String() + "/child").
			Reply(400).JSON(aerr)

		ret, err := API.CreateSubWorkspace(w, s)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, s)
	})

	t.Run("AddRackToWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		id2 := uuid.NewV4()

		gock.New(API.BaseURL).Post("/workspace/" + id.String() + "/rack").
			Reply(400).JSON(aerr)

		err := API.AddRackToWorkspace(id, id2)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("DeleteRackFromWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		id2 := uuid.NewV4()

		gock.New(API.BaseURL).Delete("/workspace/" + id.String() + "/rack").
			Reply(400).JSON(aerr)

		err := API.DeleteRackFromWorkspace(id, id2)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("AddUserToWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Post("/workspace/" + id.String() + "/user").
			Reply(400).JSON(aerr)
		err := API.AddUserToWorkspace(id, "user", "role")
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("RemoveUserFromWorkspace", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Delete("/workspace/" + id.String() + "/user").
			Reply(400).JSON(aerr)
		err := API.RemoveUserFromWorkspace(id, "user")
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("GetWorkspaceDevices", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/workspace/" + id.String() + "/device").
			Persist().Reply(400).JSON(aerr)
		defer gock.Flush()

		ret, err := API.GetWorkspaceDevices(id, false, "g", "h")
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.Device{})

		ret, err = API.GetWorkspaceDevices(id, true, "g", "h")
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.Device{})
	})

	t.Run("GetWorkspaceRacks", func(t *testing.T) {
		id := uuid.NewV4()

		gock.New(API.BaseURL).Get("/workspace/" + id.String() + "/rack").
			Reply(400).JSON(aerr)

		ret, err := API.GetWorkspaceRacks(id)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.Rack{})
	})

	t.Run("GetWorkspaceRack", func(t *testing.T) {
		id := uuid.NewV4()
		rID := uuid.NewV4()

		gock.New(API.BaseURL).
			Get("/workspace/" + id.String() + "/rack/" + rID.String()).
			Reply(400).JSON(aerr)

		ret, err := API.GetWorkspaceRack(id, rID)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, conch.Rack{})
	})

}
