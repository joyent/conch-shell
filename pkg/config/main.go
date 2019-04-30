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
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/joyent/conch-shell/pkg/config/obfuscate"
)

const (
	ProductionURL = "https://conch.joyent.us"
	StagingURL    = "https://staging.conch.joyent.us"
)

var (
	// For the love of Eris, override this default via the Makefile
	ObfuscationKey = "shies9rohz1beigheyoish1viovohWachohw7aithee9apheez"
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

// We're going to obfuscate the token itself. I'm aware this is krypto and not
// even remotely secure. But it will prevent the tokens from being just c&p'd
// out of the configs on a remote box.
type Token string

func (t Token) String() string {
	return string(t)
}

func (t Token) MarshalJSON() ([]byte, error) {
	if len(string(t)) == 0 {
		return []byte("\"\""), nil
	}
	str, err := obfuscate.Obfuscate(string(t), ObfuscationKey)
	return []byte(fmt.Sprintf("\"%s\"", str)), err
}

func (t *Token) UnmarshalJSON(b []byte) error {
	if string(b) == "\"\"" {
		*t = Token("")
		return nil
	}

	str := strings.ReplaceAll(string(b), "\"", "")

	token, _ := obfuscate.Deobfuscate(str, ObfuscationKey)
	*t = Token(token)
	return nil
}

// ConchProfile is an individual environment, consisting of login data, API
// settings, and an optional default workspace
type ConchProfile struct {
	Name          string         `json:"name"`
	User          string         `json:"user"`
	WorkspaceUUID uuid.UUID      `json:"workspace_id"`
	WorkspaceName string         `json:"workspace_name"`
	BaseURL       string         `json:"api_url"`
	Active        bool           `json:"active"`
	JWT           conch.ConchJWT `json:"jwt"`               // TODO(sungo): DEPRECATED
	Expires       time.Time      `json:"expires,omitempty"` // TODO(sungo): DEPRECATED
	Token         Token          `json:"token"`
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

	// Keeping all these structures local because they're old, crufty, and
	// there's no need to pollute the global space with this noise
	type conchProfileTransition struct {
		Name          string    `json:"name"`
		User          string    `json:"user"`
		Session       string    `json:"session,omitempty"`
		WorkspaceUUID uuid.UUID `json:"workspace_id"`
		WorkspaceName string    `json:"workspace_name"`
		BaseURL       string    `json:"api_url"`
		Active        bool      `json:"active"`
		JWT           string    `json:"jwt"`
		Expires       int64     `json:"expires,omitempty"`
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

	// First we see if the JSON parses as the old profile structure
	err = json.Unmarshal([]byte(j), ct)

	if err != nil {
		// Well, that didn't work. Let's try again with the current structure
		err = json.Unmarshal([]byte(j), c)
		if err != nil {
			// That didn't work either. Not much we can do except bail
			return c, err
		}

		// Great. We have a current profile
		if c.Profiles == nil {
			// Except we don't have any profiles
			c.Profiles = make(map[string]*ConchProfile)
		}

		// If we have a token, zero out the old JWT structure because who cares
		// about that if we have a token
		for _, profile := range c.Profiles {
			if string(profile.Token) != "" {
				profile.JWT = conch.ConchJWT{}
			}
		}

		return c, nil
	}

	// Oh joy. We have old data.
	//
	// This mostly just pulls apart the old JWT string into a fancy structure.
	// This will also be totally replaced by the Token string which is... well,
	// it's a JWT string. Circle of life, I guess. Or something.
	for _, p := range ct.Profiles {

		jwt := conch.ConchJWT{}

		bits := strings.Split(p.JWT, ".")
		if len(bits) == 3 {
			token := bits[0] + "." + bits[1]
			sig := bits[2]
			jwt, _ = parseJWT(token, sig)
		}

		pNew := &ConchProfile{
			Name:          p.Name,
			User:          p.User,
			WorkspaceUUID: p.WorkspaceUUID,
			WorkspaceName: p.WorkspaceName,
			BaseURL:       p.BaseURL,
			Active:        p.Active,
			JWT:           jwt,
			Expires:       time.Unix(p.Expires, 0),
		}
		c.Profiles[pNew.Name] = pNew
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
