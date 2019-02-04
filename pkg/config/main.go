// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package config wraps a conch shell config. Typically this is either coming
// from and/or becoming JSON on disk.
package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/joyent/conch-shell/pkg/conch"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// ErrConfigNoPath is issued when a file operation is attempted on a
// ConchConfig that lacks a path
var ErrConfigNoPath = errors.New("no path found in config data")

// ConchConfig represents the configuration information for the shell, mostly
// just a profile list
type ConchConfig struct {
	Path     string                   `json:"path"`
	Profiles map[string]*ConchProfile `json:"profiles"`
}

// ConchProfile is an individual environment, consisting of login data, API
// settings, and an optional default workspace
type ConchProfile struct {
	Name             string         `json:"name"`
	User             string         `json:"user"`
	Session          string         `json:"session,omitempty"`
	WorkspaceUUID    uuid.UUID      `json:"workspace_id"`
	WorkspaceName    string         `json:"workspace_name"`
	BaseURL          string         `json:"api_url"`
	Active           bool           `json:"active"`
	JWT              conch.ConchJWT `json:"jwt"`
	Expires          time.Time      `json:"expires,omitempty"`
	SkipVersionCheck bool           `json:"skip_version_check"`
}

// New provides an initialized struct with default values geared towards a
// dev environment. For instance, the default Api value is
// "http://localhost:5001".
func New() (c *ConchConfig) {
	c = &ConchConfig{
		Path:     "~/.conch.json",
		Profiles: make(map[string]*ConchProfile),
	}

	return c
}

// NewFromJSON unmarshals a JSON blob into a ConchConfig struct
func NewFromJSON(j string) (c *ConchConfig, err error) {

	// BUG(sungo): This transition code is a mess but necessary for
	// compatbility. Need to give it a release or two in production before
	// removing this grossness.
	type conchProfileTransition struct {
		Name             string    `json:"name"`
		User             string    `json:"user"`
		Session          string    `json:"session,omitempty"`
		WorkspaceUUID    uuid.UUID `json:"workspace_id"`
		WorkspaceName    string    `json:"workspace_name"`
		BaseURL          string    `json:"api_url"`
		Active           bool      `json:"active"`
		JWT              string    `json:"jwt"`
		Expires          int64     `json:"expires,omitempty"`
		SkipVersionCheck bool      `json:"skip_version_check"`
	}

	type conchConfigTransition struct {
		Path     string                             `json:"path"`
		Profiles map[string]*conchProfileTransition `json:"profiles"`
	}

	ct := &conchConfigTransition{
		Path:     "~/.conch.json",
		Profiles: make(map[string]*conchProfileTransition),
	}

	c = New()
	err = json.Unmarshal([]byte(j), ct)

	if err != nil {
		err = json.Unmarshal([]byte(j), c)
		if err != nil {
			return c, err
		}

		if c.Profiles == nil {
			c.Profiles = make(map[string]*ConchProfile)
		}

	} else {
		for _, p := range ct.Profiles {

			jwt := conch.ConchJWT{}

			bits := strings.Split(p.JWT, ".")
			fmt.Println(bits)
			if len(bits) == 3 {
				token := bits[0] + "." + bits[1]
				sig := bits[2]
				jwt, _ = parseJWT(token, sig)
			}

			pNew := &ConchProfile{
				Name:             p.Name,
				User:             p.User,
				Session:          p.Session,
				WorkspaceUUID:    p.WorkspaceUUID,
				WorkspaceName:    p.WorkspaceName,
				BaseURL:          p.BaseURL,
				Active:           p.Active,
				JWT:              jwt,
				Expires:          time.Unix(p.Expires, 0),
				SkipVersionCheck: p.SkipVersionCheck,
			}
			c.Profiles[pNew.Name] = pNew
		}
	}

	return c, nil
}

// NewFromJSONFile reads a file off disk and unmarshals it into ConchConfig
// struct.
func NewFromJSONFile(path string) (c *ConchConfig, err error) {

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return New(), err
	}

	return NewFromJSON(string(raw))
}

// Serialize marshals a ConchConfig struct into a JSON string
func (c *ConchConfig) Serialize() (s string, err error) {

	j, err := json.MarshalIndent(c, "", "	")

	if err != nil {
		return "", err
	}

	return string(j), nil
}

// SerializeToFile marshals a ConchConfig struct into a JSON string and
// writes it out to the provided path
func (c *ConchConfig) SerializeToFile(path string) (err error) {
	if c.Path == "" {
		return ErrConfigNoPath
	}

	j, err := json.MarshalIndent(c, "", "	")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, j, 0644)
	return err
}

// BUG(sungo): entirely for config backcompat
func decodeJWTsegment(seg string) (map[string]interface{}, error) {
	var payload map[string]interface{}

	b, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		return payload, err
	}

	err = json.Unmarshal(b, &payload)

	return payload, err
}

// BUG(sungo): entirely for config backcompat
func parseJWT(token string, signature string) (conch.ConchJWT, error) {
	var jwt conch.ConchJWT
	var err error

	jwt.Token = token
	jwt.Signature = signature

	bits := strings.Split(token, ".")
	if len(bits) != 2 {
		return jwt, conch.ErrMalformedJWT
	}

	jwt.Header, err = decodeJWTsegment(bits[0])
	if err != nil {
		return jwt, err
	}

	jwt.Claims, err = decodeJWTsegment(bits[1])
	if err != nil {
		return jwt, err
	}

	if val, ok := jwt.Claims["exp"]; ok {
		jwt.Expires = time.Unix(int64(val.(float64)), 0)
	}

	return jwt, nil
}
