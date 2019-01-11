// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	uuid "gopkg.in/satori/go.uuid.v1"
)

// RevokeUserTokens revokes all auth tokens for a the given user. This action
// is typically limited server-side to admins.
func (c *Conch) RevokeUserTokens(user string) error {
	var uPart string
	_, err := uuid.FromString(user)
	if err == nil {
		uPart = user
	} else {
		uPart = "email=" + user
	}

	return c.post("/user/"+uPart+"/revoke", nil, nil)
}

// RevokeOwnTokens revokes all auth tokens for the current user. Login() is
// required after to generate new tokens. Clears the Session, JWToken, and
// Expires attributes
func (c *Conch) RevokeOwnTokens() error {
	if err := c.post("/user/me/revoke", nil, nil); err != nil {
		return err
	}

	c.Session = ""
	c.JWToken = ""
	c.Expires = 0

	return nil
}

// VerifyLogin determines if the user's session data is still valid. If
// available, it uses the refresh API, falling back to plain cookie auth.
//
// One can pass in an integer value, representing when to force a token
// refresh, based on the number of seconds left until expiry. Pass in 0 to
// prevent refreshing
//
// If the second parameter is true, a JWT refresh is forced, regardless of any
// other parameters.
//
// NOTE: If the Conch struct contains cookie session data, it will be
// automatically upgraded to JWT and the Session data will no longer function
func (c *Conch) VerifyLogin(refreshTime int, forceJWT bool) error {
	u, _ := url.Parse(c.BaseURL)
	if c.JWToken != "" {
		if err := c.recordJWTExpiry(); err != nil {
			return ErrLoginFailed
		}
	}

	if !forceJWT {
		if refreshTime > 0 {
			now := int(time.Now().Unix())
			if c.Expires-now > refreshTime {
				return nil
			}
		}
	}

	jwtAuth := struct {
		Token string `json:"jwt_token,omitempty"`
	}{}

	if err := c.post("/refresh_token", nil, &jwtAuth); err != nil {
		return err
	}

	if jwtAuth.Token == "" {
		return ErrLoginFailed
	}

	signature := ""
	for _, cookie := range c.HTTPClient.Jar.Cookies(u) {
		if cookie.Name == "jwt_sig" {
			signature = cookie.Value
		}
	}
	if signature == "" {
		return ErrLoginFailed
	}

	c.JWToken = jwtAuth.Token + "." + signature
	if err := c.recordJWTExpiry(); err != nil {
		return ErrLoginFailed
	}

	c.Session = ""
	return nil
}

// Login uses the User, as listed in the Conch struct, and the provided
// password to log into the Conch API and populate the Session entry in the
// Conch struct
func (c *Conch) Login(user string, password string) error {

	payload := struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}{
		user,
		password,
	}

	jwtAuth := struct {
		Token string `json:"jwt_token,omitempty"`
	}{}

	res, err := c.postNeedsResponse("/login", payload, &jwtAuth)
	if err != nil {
		return err
	}

	u, _ := url.Parse(c.BaseURL)

	if jwtAuth.Token != "" {
		signature := ""
		for _, cookie := range c.HTTPClient.Jar.Cookies(u) {
			if cookie.Name == "jwt_sig" {
				signature = cookie.Value
			}
		}
		if signature == "" {
			return ErrLoginFailed
		}

		c.JWToken = jwtAuth.Token + "." + signature

		if err := c.recordJWTExpiry(); err != nil {
			return ErrLoginFailed
		}

	} else {
		for _, cookie := range c.HTTPClient.Jar.Cookies(u) {
			if cookie.Name == "conch" {
				c.Session = cookie.Value
			}
		}

		if c.Session == "" {
			return ErrLoginFailed
		}
	}

	location, err := res.Location()

	if err != nil {
		if err != http.ErrNoLocation {
			return err
		}
	}

	if location != nil {
		return ErrMustChangePassword
	}

	return nil
}

// ChangePassword changes the password for the currently active profile
func (c *Conch) ChangePassword(password string) error {
	b := struct {
		Password string `json:"password"`
	}{password}

	return c.post("/user/me/password", b, nil)

}

func (c *Conch) recordJWTExpiry() error {
	bits := strings.Split(c.JWToken, ".")
	if len(bits) != 3 {
		return ErrLoginFailed
	}

	b, err := base64.RawURLEncoding.DecodeString(bits[1])
	if err != nil {
		return err
	}

	jp := struct {
		Exp int `json:"exp"`
	}{}

	err = json.Unmarshal(b, &jp)
	if err != nil {
		return err
	}
	c.Expires = jp.Exp

	return nil
}