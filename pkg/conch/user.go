// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"github.com/joyent/conch-shell/pkg/pgtime"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// User represents a person able to access the Conch API or UI
type User struct {
	ID      string    `json:"id,omitempty"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Role    string    `json:"role"`
	RoleVia uuid.UUID `json:"role_via,omitempty"`
}

// UserDetailed ...
type UserDetailed struct {
	ID                  uuid.UUID          `json:"id"`
	Name                string             `json:"name"`
	Email               string             `json:"email"`
	Created             pgtime.PgTime      `json:"created"`
	LastLogin           pgtime.PgTime      `json:"last_login"`
	RefuseSessionAuth   bool               `json:"refuse_session_auth"`
	ForcePasswordChange bool               `json:"force_password_change"`
	Workspaces          []WorkspaceAndRole `json:"workspaces,omitempty"`
	IsAdmin             bool               `json:"is_admin"`
}

// GetUserSettings returns the results of /user/me/settings
// The return is a map[string]interface{} because the database structure is a
// string name and a jsonb data field.  There is no way for this library to
// know in advanace what's in that data so here there be dragons.
func (c *Conch) GetUserSettings() (map[string]interface{}, error) {
	settings := make(map[string]interface{})

	aerr := &APIError{}
	res, err := c.sling().New().Get("/user/me/settings").Receive(&settings, aerr)
	return settings, c.isHTTPResOk(res, err, aerr)
}

// GetUserSetting returns the results of /user/me/settings/:key
// The return is an interface{} because the database structure is a string name
// and a jsonb data field.  There is no way for this library to know in
// advanace what's in that data so here there be dragons.
func (c *Conch) GetUserSetting(key string) (interface{}, error) {
	var setting interface{}

	aerr := &APIError{}
	res, err := c.sling().New().Get("/user/me/settings/"+key).
		Receive(&setting, aerr)

	return setting, c.isHTTPResOk(res, err, aerr)
}

// SetUserSettings sets the value of *all* user settings via /user/me/settings
func (c *Conch) SetUserSettings(settings map[string]interface{}) error {
	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/user/me/settings").
		BodyJSON(settings).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// SetUserSetting sets the value of a user setting via /user/me/settings/:name
func (c *Conch) SetUserSetting(name string, value interface{}) error {
	aerr := &APIError{}
	res, err := c.sling().New().
		Post("/user/me/settings/"+name).
		BodyJSON(value).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// DeleteUserSetting deletes a user setting via /user/me/settings/:name
func (c *Conch) DeleteUserSetting(name string) error {
	aerr := &APIError{}
	res, err := c.sling().New().
		Delete("/user/me/settings/"+name).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// DeleteUser deletes a user and, optionally, clears their JWT credentials
func (c *Conch) DeleteUser(emailAddress string, clearTokens bool) error {
	url := "/user/email=" + emailAddress

	if clearTokens {
		url = url + "?clear_tokens=1"
	}

	aerr := &APIError{}
	res, err := c.sling().New().Delete(url).Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// CreateUser creates a new user. They are *not* added to a workspace.
// The 'name' argument is optional and will be omitted if set to ""
// The 'password' argument is optional and will be omitted if set to ""
func (c *Conch) CreateUser(email string, password string, name string) error {
	u := struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		Name     string `json:"name,omitempty"`
	}{email, password, name}

	aerr := &APIError{}
	res, err := c.sling().New().Post("/user").BodyJSON(u).Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// ResetUserPassword resets the password for the provided user, causing an
// email to be sent
func (c *Conch) ResetUserPassword(email string) error {
	aerr := &APIError{}
	res, err := c.sling().New().Delete("/user/email="+email+"/password").Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// GetAllUsers retrieves a list of all users, if the user has the right
// permissions, in the system. Returns UserDetailed structs
func (c *Conch) GetAllUsers() ([]UserDetailed, error) {
	u := make([]UserDetailed, 0)
	aerr := &APIError{}

	res, err := c.sling().New().Get("/user").Receive(&u, aerr)
	return u, c.isHTTPResOk(res, err, aerr)
}
