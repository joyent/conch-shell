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
)

func TestUserErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetUserSettings", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/user/me/settings").Reply(400).JSON(ErrApi)
		ret, err := API.GetUserSettings()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, make(map[string]interface{}))
	})

	t.Run("GetUserSetting", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/user/me/settings/test").
			Reply(400).JSON(ErrApi)

		ret, err := API.GetUserSetting("test")
		st.Expect(t, err, ErrApiUnpacked)
		var f interface{}
		st.Expect(t, ret, f)
	})

	t.Run("SetUserSettings", func(t *testing.T) {
		s := make(map[string]interface{})

		gock.New(API.BaseURL).Post("/user/me/settings").
			Reply(400).JSON(ErrApi)

		err := API.SetUserSettings(s)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("SetUserSetting", func(t *testing.T) {
		gock.New(API.BaseURL).Post("/user/me/settings/test").
			Reply(400).JSON(ErrApi)

		err := API.SetUserSetting("test", "wat")
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteUserSetting", func(t *testing.T) {
		gock.New(API.BaseURL).Delete("/user/me/settings/test").
			Reply(400).JSON(ErrApi)
		err := API.DeleteUserSetting("test")
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		gock.New(API.BaseURL).Delete("/user/email=foo@bar.bat").
			Reply(400).JSON(ErrApi)
		err := API.DeleteUser("foo@bar.bat", false)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("CreateUser", func(t *testing.T) {
		gock.New(API.BaseURL).Post("/user").Reply(400).JSON(ErrApi)
		err := API.CreateUser("foo@bar.bat", "", "", false)
		st.Expect(t, err, ErrApiUnpacked)

		gock.New(API.BaseURL).Post("/user").Reply(400).JSON(ErrApi)
		err = API.CreateUser("foo@bar.bat", "", "", true)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("ResetUserPassword", func(t *testing.T) {
		gock.New(API.BaseURL).Delete("/user/email=foo@bar.bat").
			MatchParam("clear_tokens", "login_only").Reply(400).JSON(ErrApi)

		err := API.ResetUserPassword("foo@bar.bat", false)
		st.Expect(t, err, ErrApiUnpacked)

		gock.New(API.BaseURL).Delete("/user/email=foo@bar.bat").
			MatchParam("clear_tokens", "all").Reply(400).JSON(ErrApi)

		err = API.ResetUserPassword("foo@bar.bat", true)
		st.Expect(t, err, ErrApiUnpacked)

	})

	t.Run("GetAllUsers", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/user").Reply(400).JSON(ErrApi)
		users, err := API.GetAllUsers()
		st.Expect(t, users, make(conch.UsersDetailed, 0))
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("GetUserProfile", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/user/me").Reply(400).JSON(ErrApi)
		profile, err := API.GetUserProfile()

		st.Expect(t, profile, conch.UserProfile{})
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("GetMyTokens", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/user/me/token").Reply(400).JSON(ErrApi)
		tokens, err := API.GetMyTokens()
		st.Expect(t, tokens, make(conch.UserTokens, 0))
		st.Expect(t, err, ErrApiUnpacked)

	})

	t.Run("GetMyToken", func(t *testing.T) {
		tokenName := "token_test"
		gock.New(API.BaseURL).Get("/user/me/token/" + tokenName).Reply(400).JSON(ErrApi)

		token, err := API.GetMyToken(tokenName)
		st.Expect(t, token, conch.UserToken{})
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("CreateMyToken", func(t *testing.T) {
		tokenName := "token_test"
		gock.New(API.BaseURL).Post("/user/me/token").Reply(400).JSON(ErrApi)

		token, err := API.CreateMyToken(tokenName)
		st.Expect(t, token, conch.NewUserToken{})
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("DeleteMyToken", func(t *testing.T) {
		tokenName := "token_test"
		gock.New(API.BaseURL).Delete("/user/me/token/" + tokenName).Reply(400).JSON(ErrApi)

		err := API.DeleteMyToken(tokenName)
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("RevokeMyLogins", func(t *testing.T) {
		gock.New(API.BaseURL).Post("/user/me/revoke").
			MatchParam("auth_only", "1").Reply(400).JSON(ErrApi)

		err := API.RevokeMyLogins()
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("RevokeMyTokens", func(t *testing.T) {
		gock.New(API.BaseURL).Post("/user/me/revoke").
			MatchParam("api_only", "1").Reply(400).JSON(ErrApi)

		err := API.RevokeMyTokens()
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("RevokeMyTokensAndLogins", func(t *testing.T) {
		gock.New(API.BaseURL).Post("/user/me/revoke").Reply(400).JSON(ErrApi)
		err := API.RevokeMyTokensAndLogins()
		st.Expect(t, err, ErrApiUnpacked)
	})

	t.Run("ChangeMyPassword", func(t *testing.T) {
		gock.New(API.BaseURL).Post("/user/me/password").
			MatchParam("clear_tokens", "login_only").Reply(400).JSON(ErrApi)
		err := API.ChangeMyPassword("pants", false)
		st.Expect(t, err, ErrApiUnpacked)

		gock.New(API.BaseURL).Post("/user/me/password").
			MatchParam("clear_tokens", "all").Reply(400).JSON(ErrApi)
		err = API.ChangeMyPassword("pants", true)
		st.Expect(t, err, ErrApiUnpacked)

	})

}
