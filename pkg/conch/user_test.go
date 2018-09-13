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
	// "time"
)

func TestUserErrors(t *testing.T) {
	BuildAPI()
	gock.Flush()

	aerr := conch.APIError{ErrorMsg: "totally broken"}
	aerrUnpacked := errors.New(aerr.ErrorMsg)

	t.Run("GetUserSettings", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/user/me/settings").Reply(400).JSON(aerr)
		ret, err := API.GetUserSettings()
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, make(map[string]interface{}))
	})

	t.Run("GetUserSetting", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/user/me/settings/test").
			Reply(400).JSON(aerr)
		ret, err := API.GetUserSetting("test")
		st.Expect(t, err, aerrUnpacked)
		var f interface{}
		st.Expect(t, ret, f)
	})

	t.Run("SetUserSettings", func(t *testing.T) {
		s := make(map[string]interface{})

		gock.New(API.BaseURL).Post("/user/me/settings").
			Reply(400).JSON(aerr)
		err := API.SetUserSettings(s)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("SetUserSetting", func(t *testing.T) {
		gock.New(API.BaseURL).Post("/user/me/settings/test").
			Reply(400).JSON(aerr)
		err := API.SetUserSetting("test", "wat")
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("DeleteUserSetting", func(t *testing.T) {
		gock.New(API.BaseURL).Delete("/user/me/settings/test").
			Reply(400).JSON(aerr)
		err := API.DeleteUserSetting("test")
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("InviteUser", func(t *testing.T) {
		id := uuid.NewV4()
		gock.New(API.BaseURL).Post("/workspace/" + id.String() + "/user").
			Reply(400).JSON(aerr)
		err := API.InviteUser(id, "user", "role")
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		gock.New(API.BaseURL).Delete("/user/email=foo@bar.bat").
			Reply(400).JSON(aerr)
		err := API.DeleteUser("foo@bar.bat", false)
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("CreateUser", func(t *testing.T) {
		gock.New(API.BaseURL).Post("/user").Reply(400).JSON(aerr)
		err := API.CreateUser("foo@bar.bat", "", "")
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("ResetUserPassword", func(t *testing.T) {
		gock.New(API.BaseURL).Delete("/user/email=foo@bar.bat").Reply(400).JSON(aerr)
		err := API.ResetUserPassword("foo@bar.bat")
		st.Expect(t, err, aerrUnpacked)
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/user").Reply(400).JSON(aerr)
		users, err := API.GetAllUsers()
		st.Expect(t, users, make([]conch.UserDetailed, 0))
		st.Expect(t, err, aerrUnpacked)
	})

}
