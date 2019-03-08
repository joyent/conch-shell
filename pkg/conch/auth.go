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
// required after to generate new tokens. Clears the JWT, and
// Expires attributes
func (c *Conch) RevokeOwnTokens() error {
	if err := c.post("/user/me/revoke", nil, nil); err != nil {
		return err
	}

	c.JWT = ConchJWT{}
	return nil
}

// VerifyLogin determines if the user's session data is still valid.
//
// One can pass in an integer value, representing when to force a token
// refresh, based on the number of seconds left until expiry. Pass in 0 to
// prevent refreshing
//
// If the second parameter is true, a JWT refresh is forced, regardless of any
// other parameters.
//
func (c *Conch) VerifyLogin(refreshTime int, forceJWT bool) error {
	u, _ := url.Parse(c.BaseURL)

	if !forceJWT {
		if (refreshTime > 0) && !c.JWT.Expires.IsZero() {
			now := time.Now()
			if c.JWT.Expires.Sub(now).Seconds() > float64(refreshTime) {
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
		return ErrMalformedJWT
	}

	signature := ""
	for _, cookie := range c.HTTPClient.Jar.Cookies(u) {
		if cookie.Name == "jwt_sig" {
			signature = cookie.Value
		}
	}
	if signature == "" {
		return ErrMalformedJWT
	}

	jwt, err := c.ParseJWT(jwtAuth.Token, signature)
	if err != nil {
		return err
	}

	c.JWT = jwt

	return nil
}

// Login uses the User, as listed in the Conch struct, and the provided
// password to log into the Conch API and populate the JWT entry in the
// Conch struct
func (c *Conch) Login(user string, password string) error {
	u, _ := url.Parse(c.BaseURL)

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

	signature := ""
	for _, cookie := range c.HTTPClient.Jar.Cookies(u) {
		if cookie.Name == "jwt_sig" {
			signature = cookie.Value
		}
	}
	if signature == "" {
		return ErrMalformedJWT
	}

	jwt, err := c.ParseJWT(jwtAuth.Token, signature)
	if err != nil {
		return err
	}

	c.JWT = jwt

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

func decodeJWTsegment(seg string) (map[string]interface{}, error) {
	var payload map[string]interface{}

	b, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		return payload, err
	}

	err = json.Unmarshal(b, &payload)

	return payload, err
}

func (c *Conch) ParseJWT(token string, signature string) (ConchJWT, error) {
	var jwt ConchJWT
	var err error

	jwt.Token = token
	jwt.Signature = signature

	if c.Trace {
		c.ddp(jwt)
	}

	bits := strings.Split(token, ".")
	if len(bits) != 2 {
		return jwt, ErrMalformedJWT
	}

	jwt.Header, err = decodeJWTsegment(bits[0])
	if err != nil {
		return jwt, ErrMalformedJWT
	}

	jwt.Claims, err = decodeJWTsegment(bits[1])
	if err != nil {
		return jwt, err
	}

	if c.Trace {
		c.ddp(jwt)
	}

	if val, ok := jwt.Claims["exp"]; ok {
		jwt.Expires = time.Unix(int64(val.(float64)), 0)
	}

	return jwt, nil
}

// ChangePassword changes the password for the currently active profile
func (c *Conch) ChangePassword(password string) error {
	b := struct {
		Password string `json:"password"`
	}{password}

	return c.post("/user/me/password", b, nil)

}
