// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	"regexp"
	"strings"
)

func isTag(str string) bool {
	// Settings that start with 'tag.' are special cased and only available
	// in the device tag commands
	return regexp.MustCompile(`^tag\.`).MatchString(str)
}

// GetDeviceSettings fetches settings for a device, via
// /device/:serial/settings
// Device settings that begin with 'tag.' are filtered out.
func (c *Conch) GetDeviceSettings(serial string) (map[string]string, error) {
	settings := make(map[string]string)
	filtered := make(map[string]string)

	if err := c.get("/device/"+serial+"/settings", &settings); err != nil {
		return filtered, err
	}

	for k, v := range settings {
		if !isTag(k) {
			filtered[k] = v
		}
	}

	return filtered, nil
}

// GetDeviceSetting fetches a single setting for a device, via
// /device/:serial/settings/:key
// Device settings that begin with 'tag.' are filtered out.
func (c *Conch) GetDeviceSetting(serial string, key string) (string, error) {

	if isTag(key) {
		return "", ErrDataNotFound
	}

	var setting string
	j := make(map[string]string)

	if err := c.get("/device/"+serial+"/settings/"+key, &j); err != nil {
		return setting, err
	}

	if _, ok := j[key]; ok {
		setting = j[key]
	}

	return setting, nil
}

// SetDeviceSetting sets a single setting for a device via /device/:deviceID/settings/:key
// Settings that begin with "tag." cannot be processed by this routine and will
// always return ErrDataNotFound
func (c *Conch) SetDeviceSetting(deviceID string, key string, value string) error {
	if isTag(key) {
		return ErrDataNotFound
	}

	j := make(map[string]string)
	j[key] = value

	return c.post(
		"/device/"+deviceID+"/settings/"+key,
		j,
		nil,
	)
}

// DeleteDeviceSetting deletes a single setting for a device via
// /device/:deviceID/settings/:key
// Settings that begin with "tag." cannot be processed by this routine and will
// always return ErrDataNotFound
func (c *Conch) DeleteDeviceSetting(deviceID string, key string) error {
	if isTag(key) {
		return ErrDataNotFound
	}
	return c.httpDelete("/device/" + deviceID + "/settings/" + key)
}

// GetDeviceTags fetches tags for a device, via /device/:serial/settings
// Device settings that do NOT begin with 'tag.' are filtered out.
func (c *Conch) GetDeviceTags(serial string) (map[string]string, error) {
	settings := make(map[string]string)
	filtered := make(map[string]string)

	if err := c.get("/device/"+serial+"/settings", &settings); err != nil {
		return filtered, err
	}

	for k, v := range settings {
		if isTag(k) {
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

	if !isTag(key) {
		key = "tag." + key
	}

	var setting string
	j := make(map[string]string)

	if err := c.get("/device/"+serial+"/settings/"+key, &j); err != nil {
		return setting, err
	}

	if _, ok := j[key]; ok {
		setting = j[key]
	}

	return setting, nil
}

// SetDeviceTag sets a single tag for a device via /device/:deviceID/settings/:key
// The key must either begin with 'tag.' or it will be prepended
func (c *Conch) SetDeviceTag(deviceID string, key string, value string) error {
	if !isTag(key) {
		key = "tag." + key
	}

	j := make(map[string]string)
	j[key] = value

	return c.post("/device/"+deviceID+"/settings/"+key, j, nil)
}

// DeleteDeviceTag deletes a single tag for a device via
// /device/:deviceID/settings/:key
// Settings that do NOT begin with "tag." cannot be processed by this routine
// and will always return ErrDataNotFound
func (c *Conch) DeleteDeviceTag(deviceID string, key string) error {
	if !isTag(key) {
		key = "tag." + key
	}

	return c.httpDelete("/device/" + deviceID + "/settings/" + key)
}

func (c *Conch) GetDevicesBySetting(key string, value string) (d Devices, err error) {
	url := fmt.Sprintf("/device?%s=%s", key, value)
	return d, c.get(url, &d)
}

func (c *Conch) GetDevicesByTag(key string, value string) (Devices, error) {
	return c.GetDevicesBySetting("tag."+key, value)

}
