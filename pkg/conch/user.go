// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"

	uuid "gopkg.in/satori/go.uuid.v1"
)

func (c *Conch) GetMyTokens() (UserTokens, error) {
	u := make(UserTokens, 0)
	return u, c.get("/user/me/token", &u)
}

func (c *Conch) GetMyToken(name string) (u UserToken, err error) {
	return u, c.get("/user/me/token/"+name, &u)
}

func (c *Conch) CreateMyToken(name string) (u NewUserToken, err error) {
	return u, c.post(
		"/user/me/token",
		CreateNewUserToken{Name: name},
		&u,
	)
}

func (c *Conch) DeleteMyToken(name string) error {
	return c.httpDelete("/user/me/token/" + name)
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

func (c *Conch) ChangeMyPassword(password string) error {
	return c.ChangePassword(password)
}

// ChangePassword changes the password for the currently active profile
func (c *Conch) ChangePassword(password string) error {
	b := struct {
		Password string `json:"password"`
	}{password}

	return c.post("/user/me/password", b, nil)

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
	return setting, c.get("/user/me/settings/"+key, &setting)
}

// SetUserSettings sets the value of *all* user settings via /user/me/settings
func (c *Conch) SetUserSettings(settings map[string]interface{}) error {
	return c.post("/user/me/settings", settings, nil)
}

// SetUserSetting sets the value of a user setting via /user/me/settings/:name
func (c *Conch) SetUserSetting(name string, value interface{}) error {
	return c.post("/user/me/settings/"+name, value, nil)
}

// DeleteUserSetting deletes a user setting via /user/me/settings/:name
func (c *Conch) DeleteUserSetting(name string) error {
	return c.httpDelete("/user/me/settings/" + name)
}

// DeleteUser deletes a user and, optionally, clears their JWT credentials
func (c *Conch) DeleteUser(emailAddress string, clearTokens bool) error {
	url := "/user/email=" + emailAddress

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
func (c *Conch) ResetUserPassword(email string) error {
	return c.httpDelete("/user/email=" + email + "/password")
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
	return user, c.get("/user/"+id.String(), &user)
}

func (c *Conch) GetUserByEmail(email string) (user UserDetailed, err error) {
	return user, c.get("/user/email="+email, &user)
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
		fmt.Sprintf("/user/%s", userID.String()),
		u,
		nil,
	)
}
