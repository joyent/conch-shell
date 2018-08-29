// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"regexp"
	"strings"
)

// GetDeviceSettings fetches settings for a device, via
// /device/:serial/settings
// Device settings that begin with 'tag.' are filtered out.
func (c *Conch) GetDeviceSettings(serial string) (map[string]string, error) {
	settings := make(map[string]string)

	aerr := &APIError{}

	res, err := c.sling().New().
		Get("/device/"+serial+"/settings").
		Receive(&settings, aerr)

	filtered := make(map[string]string)

	if ret := c.isHTTPResOk(res, err, aerr); ret != nil {
		return filtered, ret
	}

	// Settings that start with 'tag.' are special cased and only availabe
	// in the device tag commands
	re := regexp.MustCompile("^tag\\.")
	for k, v := range settings {
		if !re.MatchString(k) {
			filtered[k] = v
		}
	}

	return filtered, nil
}

// GetDeviceSetting fetches a single setting for a device, via
// /device/:serial/settings/:key
// Device settings that begin with 'tag.' are filtered out.
func (c *Conch) GetDeviceSetting(serial string, key string) (string, error) {

	// Settings that start with 'tag.' are special cased and only available
	// in the device tag interface
	re := regexp.MustCompile("^tag\\.")
	if re.MatchString(key) {
		return "", ErrDataNotFound
	}

	var setting string
	j := make(map[string]string)

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/device/"+serial+"/settings/"+key).
		Receive(&j, aerr)

	if _, ok := j[key]; ok {
		setting = j[key]
	}

	return setting, c.isHTTPResOk(res, err, aerr)
}

// SetDeviceSetting sets a single setting for a device via /device/:deviceID/settings/:key
// Settings that begin with "tag." cannot be processed by this routine and will
// always return ErrDataNotFound
func (c *Conch) SetDeviceSetting(deviceID string, key string, value string) error {
	// Settings that start with 'tag.' are special cased and only available
	// in the device tag interface
	re := regexp.MustCompile("^tag\\.")
	if re.MatchString(key) {
		return ErrDataNotFound
	}

	j := make(map[string]string)
	j[key] = value

	aerr := &APIError{}
	res, err := c.sling().New().Post("/device/"+deviceID+"/settings/"+key).
		BodyJSON(j).Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// DeleteDeviceSetting deletes a single setting for a device via
// /device/:deviceID/settings/:key
// Settings that begin with "tag." cannot be processed by this routine and will
// always return ErrDataNotFound
func (c *Conch) DeleteDeviceSetting(deviceID string, key string) error {
	// Settings that start with 'tag.' are special cased and only available
	// in the device tag interface
	re := regexp.MustCompile("^tag\\.")
	if re.MatchString(key) {
		return ErrDataNotFound
	}

	aerr := &APIError{}
	res, err := c.sling().New().Delete("/device/"+deviceID+"/settings/"+key).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// GetDeviceTags fetches tags for a device, via /device/:serial/settings
// Device settings that do NOT begin with 'tag.' are filtered out.
func (c *Conch) GetDeviceTags(serial string) (map[string]string, error) {
	settings := make(map[string]string)

	aerr := &APIError{}

	res, err := c.sling().New().
		Get("/device/"+serial+"/settings").
		Receive(&settings, aerr)

	filtered := make(map[string]string)

	if ret := c.isHTTPResOk(res, err, aerr); ret != nil {
		return filtered, ret
	}

	// Settings that start with 'tag.' are special cased and only availabe
	// in the device tag commands
	re := regexp.MustCompile("^tag\\.")
	for k, v := range settings {
		if re.MatchString(k) {
			filtered[strings.TrimPrefix(k, "tag.")] = v
		}
	}

	return filtered, nil
}

//////////////////

// GetDeviceTag fetches a single tag for a device, via
// /device/:serial/settings/:key
// The key must either begin with 'tag.' or it will be prepended
func (c *Conch) GetDeviceTag(serial string, key string) (string, error) {

	re := regexp.MustCompile("^tag\\.")
	if !re.MatchString(key) {
		key = "tag." + key
	}

	var setting string
	j := make(map[string]string)

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/device/"+serial+"/settings/"+key).
		Receive(&j, aerr)

	if _, ok := j[key]; ok {
		setting = j[key]
	}

	return setting, c.isHTTPResOk(res, err, aerr)
}

// SetDeviceTag sets a single tag for a device via /device/:deviceID/settings/:key
// The key must either begin with 'tag.' or it will be prepended
func (c *Conch) SetDeviceTag(deviceID string, key string, value string) error {
	// Settings that start with 'tag.' are special cased and only available
	// in the device tag interface
	re := regexp.MustCompile("^tag\\.")
	if !re.MatchString(key) {
		key = "tag." + key
	}

	j := make(map[string]string)
	j[key] = value

	aerr := &APIError{}
	res, err := c.sling().New().Post("/device/"+deviceID+"/settings/"+key).
		BodyJSON(j).Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// DeleteDeviceTag deletes a single tag for a device via
// /device/:deviceID/settings/:key
// Settings that do NOT begin with "tag." cannot be processed by this routine
// and will always return ErrDataNotFound
func (c *Conch) DeleteDeviceTag(deviceID string, key string) error {
	// Settings that start with 'tag.' are special cased and only available
	// in the device tag interface
	re := regexp.MustCompile("^tag\\.")
	if !re.MatchString(key) {
		key = "tag." + key
	}

	aerr := &APIError{}
	res, err := c.sling().New().Delete("/device/"+deviceID+"/settings/"+key).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}
