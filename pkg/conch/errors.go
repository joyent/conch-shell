// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"errors"
	"net/http"
)

//APIError matches the structure of request errors reported by the Conch API
type APIError struct {
	ErrorMsg string `json:"error"`
}

// Return a new error with Conch Error message as the text
func (ai *APIError) Error() error {
	return errors.New(ai.ErrorMsg)
}

var (
	// ErrLoginFailed indicates that the login process failed for unspecified
	// reasons
	ErrLoginFailed = errors.New("login failed")

	// ErrNoSessionData indicates that an auth related error occurred where
	// either the user did not provide session data or no data was returned
	// from the API
	ErrNoSessionData = errors.New("no session data provided")

	// ErrHTTPNotOk indicates that the API returned a non-200 status code that
	// we don't know how to handle
	ErrHTTPNotOk = errors.New("non-200 HTTP status code returned")

	// ErrDataNotFound inidicates that the API returned a status code
	// inidicating that the requested data does not exist or is not available.
	// NOTE: The API will also return this error if the user is not allowed to
	// access the data in question.
	ErrDataNotFound = errors.New("API could not find the data requested")

	// ErrBadInput indicates that the user passed incomplete or bad data to a
	// routine. This typicallly only occurs when a struct parameter isn't
	// filled out with enough data.
	ErrBadInput = errors.New("incomplete data passed to the routine")

	// ErrSemVerParse indicates that a semantic version string could not be
	// parsed
	ErrSemVerParse = errors.New("could not parse semantic version string")

	// ErrNotSupported indicates that the API server does not support this
	// command. This is typically determined via checks on conch.apiVersion
	ErrNotSupported = errors.New("this function is not supported")

	// ErrNotAuthorized indicates that the API server returned a 401
	ErrNotAuthorized = errors.New("not authorized for this endpoint")

	// ErrForbidden indicates that the API server returned a 403
	ErrForbidden = errors.New("access to this endpoint is forbidden")

	// ErrMustChangePassword is used to signal that the user must change their
	// password before proceeding. Typically, the existing auth credentials
	// will continue to work for a few minutes.
	ErrMustChangePassword = errors.New("password must be changed")
)

// isHTTPResOk is a convenience function to abstract out all the code to see
// if an API call came back ok.
func (c *Conch) isHTTPResOk(res *http.Response, err error, aerr *APIError) error {
	if res == nil {
		return err
	}

	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return ErrNotAuthorized
	}

	if res.StatusCode == http.StatusForbidden {
		return ErrForbidden
	}

	if res.StatusCode == http.StatusNotFound {
		return ErrDataNotFound
	}

	if aerr != nil {
		if aerr.ErrorMsg != "" {
			return errors.New(aerr.ErrorMsg)
		}
	}

	if res.StatusCode >= 400 {
		return ErrHTTPNotOk
	}

	return nil
}
