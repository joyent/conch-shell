// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// config wraps a conch shell config. Typically this is either coming from
// and/or becoming JSON on disk.
package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

var ConchConfigNoPath = errors.New("No path found in config data")

type ConchConfig struct {
	Path    string                 `json:"path"`
	Api     string                 `json:"api"`
	User    string                 `json:"user"`
	Session string                 `json:"session"`
	KV      map[string]interface{} `json:"kv"`
}

// New() provides an initialized struct with default values geared towards a
// dev environment. For instance, the default Api value is
// "http://localhost:5001".
func New() (c *ConchConfig) {
	c = &ConchConfig{
		Api: "http://localhost:5001",
	}
	c.KV = make(map[string]interface{})
	return c
}

// NewFromJson() unmarshals a JSON blob into a ConchConfig struct. It does
// *not* fill in any default values.
func NewFromJson(j string) (c *ConchConfig, err error) {
	c = &ConchConfig{}
	err = json.Unmarshal([]byte(j), c)

	if err != nil {
		return nil, err
	}
	return c, nil
}

// NewFromJsonFile reads a file off disk and unmarshals it into ConchConfig
// struct. It does *not* fill in any default values
func NewFromJsonFile(path string) (c *ConchConfig, err error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c = &ConchConfig{}
	err = json.Unmarshal(raw, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Serialize() marshals a ConchConfig struct into a JSON string
func (c *ConchConfig) Serialize() (s string, err error) {

	j, err := json.MarshalIndent(c, "", "	")

	if err != nil {
		return "", err
	}

	return string(j), nil
}

// SerializeToFile() marshals a ConchConfig struct into a JSON string and
// writes it out to the provided path
func (c *ConchConfig) SerializeToFile(path string) (err error) {
	if c.Path == "" {
		return ConchConfigNoPath
	}

	j, err := json.MarshalIndent(c, "", "	")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, j, 0644)
	if err != nil {
		return err
	}

	return nil
}
