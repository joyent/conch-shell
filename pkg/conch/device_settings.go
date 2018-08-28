// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"regexp"
)

// GetDeviceSettings fetches settings for a device, via
// /device/:serial/settings
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
