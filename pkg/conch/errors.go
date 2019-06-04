// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"errors"
)

var (

	// ErrHTTPNotOk indicates that the API returned a non-200 status code that
	// we don't know how to handle
	ErrHTTPNotOk = errors.New("non-200 HTTP status code returned")

	// ErrDataNotFound inidicates that the API returned a status code
	// inidicating that the requested data does not exist or is not available.
	// NOTE: The API will also return this error if the user is not allowed to
	// access the data in question.
	ErrDataNotFound = errors.New("server could not find the data requested")

	// ErrBadInput indicates that the user passed incomplete or bad data to a
	// routine. This typicallly only occurs when a struct parameter isn't
	// filled out with enough data.
	ErrBadInput = errors.New("incomplete data passed to the routine")

	// ErrNotSupported indicates that the API server does not support this
	// command. This is typically determined via checks on conch.apiVersion
	ErrNotSupported = errors.New("this function is not supported")

	// ErrNotAuthorized indicates that the API server returned a 401
	ErrNotAuthorized = errors.New("invalid or expired auth credentials")

	// ErrForbidden indicates that the API server returned a 403
	ErrForbidden = errors.New("access to this data is forbidden")
)
