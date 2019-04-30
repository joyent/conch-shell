// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	"net/url"

	uuid "gopkg.in/satori/go.uuid.v1"
)

func (c *Conch) GetMyTokens() (UserTokens, error) {
	u := make(UserTokens, 0)
	return u, c.get("/user/me/token", &u)
}

func (c *Conch) GetMyToken(name string) (u UserToken, err error) {
	escapedName := url.PathEscape(name)
	return u, c.get("/user/me/token/"+escapedName, &u)
}

func (c *Conch) CreateMyToken(name string) (u NewUserToken, err error) {
	return u, c.post(
		"/user/me/token",
		CreateNewUserToken{Name: name},
		&u,
	)
}

func (c *Conch) DeleteMyToken(name string) error {
	escapedName := url.PathEscape(name)
	return c.httpDelete("/user/me/token/" + escapedName)
}

func (c *Conch) RevokeMyLogins() error {
	return c.post("/user/me/revoke?auth_only=1", nil, nil)
}

func (c *Conch) RevokeMyTokens() error {
	return c.post("/user/me/revoke?api_only=1", nil, nil)
}

func (c *Conch) RevokeMyTokensAndLogins() error {
	if err := c.post("/user/me/revoke", nil, nil); err != nil {
		return err
	}

	c.JWT = ConchJWT{}
	return nil
}

func (c *Conch) ChangeMyPassword(password string, revokeTokens bool) error {
	b := struct {
		Password string `json:"password"`
	}{password}

	url := "/user/me/password?"
	if revokeTokens {
		url = url + "clear_tokens=all"
	} else {
		// This is a bit opinionated of me. Changing your password will always
		// clear your login tokens, but never your API tokens unless you ask.
		//
		// While the API allows one to not clear any tokens, I can't figure a
		// use case where you'd want to change your password but let existing
		// sessions to just keep working.
		url = url + "clear_tokens=login_only"
	}

	return c.post(url, b, nil)

}

// GetUserSettings returns the results of /user/me/settings
// The return is a map[string]interface{} because the database structure is a
// string name and a jsonb data field.  There is no way for this library to
// know in advanace what's in that data so here there be dragons.
func (c *Conch) GetUserSettings() (map[string]interface{}, error) {
	settings := make(map[string]interface{})
	return settings, c.get("/user/me/settings", &settings)
}

// GetUserSetting returns the results of /user/me/settings/:key
// The return is an interface{} because the database structure is a string name
// and a jsonb data field.  There is no way for this library to know in
// advanace what's in that data so here there be dragons.
func (c *Conch) GetUserSetting(key string) (setting interface{}, err error) {
	return setting, c.get("/user/me/settings/"+url.PathEscape(key), &setting)
}

// SetUserSettings sets the value of *all* user settings via /user/me/settings
func (c *Conch) SetUserSettings(settings map[string]interface{}) error {
	return c.post("/user/me/settings", settings, nil)
}

// SetUserSetting sets the value of a user setting via /user/me/settings/:name
func (c *Conch) SetUserSetting(name string, value interface{}) error {
	return c.post("/user/me/settings/"+url.PathEscape(name), value, nil)
}

// DeleteUserSetting deletes a user setting via /user/me/settings/:name
func (c *Conch) DeleteUserSetting(name string) error {
	return c.httpDelete("/user/me/settings/" + url.PathEscape(name))
}

// DeleteUser deletes a user and, optionally, clears their JWT credentials
func (c *Conch) DeleteUser(emailAddress string, clearTokens bool) error {
	url := "/user/email=" + url.PathEscape(emailAddress)

	if clearTokens {
		url = url + "?clear_tokens=1"
	}

	return c.httpDelete(url)
}

// CreateUser creates a new user. They are *not* added to a workspace.
// The 'email' argument is required.
// The 'name' argument is optional
// The 'password' argument is optional
// The 'isAdmin' argument sets the user to be an admin. Defaults to false.
func (c *Conch) CreateUser(email string, password string, name string, isAdmin bool) error {
	if email == "" {
		return ErrBadInput
	}

	u := struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		Name     string `json:"name,omitempty"`
		IsAdmin  bool   `json:"is_admin"`
	}{email, password, name, isAdmin}

	return c.post("/user", u, nil)
}

// ResetUserPassword resets the password for the provided user, causing an
// email to be sent
func (c *Conch) ResetUserPassword(email string, revokeTokens bool) error {
	url := "/user/email=" + url.PathEscape(email) + "/password?"
	if revokeTokens {
		url = url + "clear_tokens=all"
	} else {
		// This is a bit opinionated of me. Changing someone's password will
		// always clear their login tokens, but never their API tokens unless you
		// ask.
		//
		// While the API allows one to not clear any tokens, I can't figure a
		// use case where you'd want to change someone's password but let
		// existing sessions to just keep working.
		url = url + "clear_tokens=login_only"
	}

	return c.httpDelete(url)
}

// GetAllUsers retrieves a list of all users, if the user has the right
// permissions, in the system. Returns UserDetailed structs
func (c *Conch) GetAllUsers() (UsersDetailed, error) {
	u := make(UsersDetailed, 0)
	return u, c.get("/user", &u)
}

func (c *Conch) GetUserProfile() (profile UserProfile, err error) {
	return profile, c.get("/user/me", &profile)
}

func (c *Conch) GetUser(id uuid.UUID) (user UserDetailed, err error) {
	return user, c.get("/user/"+url.PathEscape(id.String()), &user)
}

func (c *Conch) GetUserByEmail(email string) (user UserDetailed, err error) {
	return user, c.get("/user/email="+url.PathEscape(email), &user)
}

// UpdateUser updates properties of a user. No workspace permissions are
// changed.
// The 'userID' argument is required
// The 'email' argument is optional
// The 'name' argument is optional
// The 'isAdmin' argument sets the user to be an admin. Defaults to false.
func (c *Conch) UpdateUser(
	userID uuid.UUID,
	email string,
	name string,
	isAdmin bool,
) error {
	if uuid.Equal(userID, uuid.UUID{}) {
		return ErrBadInput
	}

	u := struct {
		Email   string `json:"email,omitempty"`
		Name    string `json:"name,omitempty"`
		IsAdmin bool   `json:"is_admin"`
	}{email, name, isAdmin}

	return c.post(
		fmt.Sprintf("/user/%s", url.PathEscape(userID.String())),
		u,
		nil,
	)
}
