// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"net/url"

	"github.com/joyent/conch-shell/pkg/conch/uuid"
)

func (c *Conch) RevokeUserTokensAndLogins(user string) error {
	var uPart string
	_, err := uuid.FromString(user)
	if err == nil {
		uPart = user
	} else {
		uPart = "email=" + url.PathEscape(user)
	}

	return c.post("/user/"+uPart+"/revoke", nil, nil)
}

func (c *Conch) RevokeUserLogins(user string) error {
	var uPart string
	_, err := uuid.FromString(user)
	if err == nil {
		uPart = user
	} else {
		uPart = "email=" + url.PathEscape(user)
	}

	return c.post("/user/"+uPart+"/revoke?auth_only=1", nil, nil)
}

func (c *Conch) RevokeUserTokens(user string) error {
	var uPart string
	_, err := uuid.FromString(user)
	if err == nil {
		uPart = user
	} else {
		uPart = "email=" + url.PathEscape(user)
	}

	return c.post("/user/"+uPart+"/revoke?api_only=1", nil, nil)
}

func (c *Conch) GetUserToken(user string, name string) (u UserToken, err error) {
	escapedName := url.PathEscape(name)
	return u, c.get("/user/email="+user+"/token/"+escapedName, &u)
}

func (c *Conch) GetUserTokens(user string) (UserTokens, error) {
	u := make(UserTokens, 0)
	escaped := url.PathEscape(user)
	return u, c.get("/user/email="+escaped+"/token", &u)
}

func (c *Conch) DeleteUserToken(user string, name string) error {
	escapedName := url.PathEscape(name)
	return c.httpDelete("/user/email=" + user + "/token/" + escapedName)
}

func (c *Conch) VerifyToken() (bool, error) {
	if c.Token == "" {
		return false, ErrBadInput
	}

	_, err := c.GetUserSettings()
	if err != nil {
		return false, err
	}

	return true, nil
}
