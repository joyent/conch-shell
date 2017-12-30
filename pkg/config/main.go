// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package config wraps a conch shell config. Typically this is either coming
// from and/or becoming JSON on disk.
package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// ErrConfigNoPath is issued when a file operation is attempted on a
// ConchConfig that lacks a path
var ErrConfigNoPath = errors.New("No path found in config data")

// ConchConfig represents the configuration information for the shell,
// including the api server, user auth data, and a key/value store for
// future-proofing
type ConchConfig struct {
	Path    string                 `json:"path"`
	API     string                 `json:"api"`
	User    string                 `json:"user"`
	Session string                 `json:"session"`
	KV      map[string]interface{} `json:"kv"`
}

// New provides an initialized struct with default values geared towards a
// dev environment. For instance, the default Api value is
// "http://localhost:5001".
func New() (c *ConchConfig) {
	c = &ConchConfig{
		API: "http://localhost:5001",
	}
	c.KV = make(map[string]interface{})
	return c
}

// NewFromJSON unmarshals a JSON blob into a ConchConfig struct. It does
// *not* fill in any default values.
func NewFromJSON(j string) (c *ConchConfig, err error) {
	c = &ConchConfig{}

	err = json.Unmarshal([]byte(j), c)
	return c, err
}

// NewFromJSONFile reads a file off disk and unmarshals it into ConchConfig
// struct. It does *not* fill in any default values
func NewFromJSONFile(path string) (c *ConchConfig, err error) {
	c = &ConchConfig{}

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(raw, c)
	return c, err
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
