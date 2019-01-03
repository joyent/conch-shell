// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/dghubble/sling"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	defaultUA      = "go-conch"
	defaultBaseURL = "https://conch.joyent.us"
)

var defaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
		DualStack: true,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

func (c *Conch) sling() *sling.Sling {
	if c.UA == "" {
		c.UA = defaultUA
	}

	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}

	if c.CookieJar == nil {
		c.CookieJar, _ = cookiejar.New(nil)
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Transport: defaultTransport,
			Jar:       c.CookieJar,

			// Preserve auth header on redirect
			// Inspired by: https://github.com/michiwend/gomusicbrainz/pull/4/files?utf8=%E2%9C%93&diff=unified
			// Under MIT License
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) > 30 {
					return fmt.Errorf("%d > 30 consecutive requests(redirects)", len(via))
				}
				if len(via) == 0 {
					// No redirects
					return nil
				}

				// This is a massive hack. In theory, go should already see
				// that these have the same host and copy the Authorization
				// header over on its own. Until I can track down why that's not
				// happening, this will get us back in business.
				// sungo [ 2018-06-21 ]
				if req.URL.Host == via[0].URL.Host {
					h, ok := via[0].Header["Authorization"]
					if ok {
						req.Header["Authorization"] = h
					}
				}
				return nil
			},
		}
	}

	s := sling.New().
		Client(c.HTTPClient).
		Base(c.BaseURL).
		Set("User-Agent", c.UA)

	if c.apiVersion == nil {
		sem, _ := semver.New(MinimumAPIVersion)
		c.apiVersion = sem

		body := struct {
			Version string `json:"version"`
		}{}

		_, err := s.Get("/version").Receive(&body, nil)
		if err != nil {
			return s
		}
		apiVer, err := semver.New(strings.TrimLeft(body.Version, "v"))
		if err == nil {
			c.apiVersion = apiVer
		}
	}

	u, _ := url.Parse(c.BaseURL)
	if c.JWToken != "" {
		if c.Expires == 0 {
			_ = c.recordJWTExpiry
		}

		s = s.Set("Authorization", "Bearer "+c.JWToken)

	} else if c.Session != "" {

		cookie := &http.Cookie{
			Name:  "conch",
			Value: c.Session,
		}
		c.CookieJar.SetCookies(
			u,
			[]*http.Cookie{cookie},
		)
	}

	return s
}

func (c *Conch) get(url string, data interface{}) error {
	aerr := &APIError{}
	res, err := c.sling().New().Get(url).Receive(data, aerr)
	return c.isHTTPResOk(res, err, aerr)
}

func (c *Conch) getWithQuery(url string, query interface{}, data interface{}) error {
	aerr := &APIError{}
	res, err := c.sling().New().Get(url).QueryStruct(query).Receive(data, aerr)
	return c.isHTTPResOk(res, err, aerr)
}

func (c *Conch) httpDelete(url string) error {
	aerr := &APIError{}
	res, err := c.sling().New().Delete(url).Receive(nil, aerr)
	return c.isHTTPResOk(res, err, aerr)
}

// RawGet allows the user to perform an HTTP GET against the API, with the
// library handling all auth but *not* processing the response.
func (c *Conch) RawGet(url string) (*http.Response, error) {
	req, err := c.sling().New().Get(url).Request()
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}

// RawDelete allows the user to perform an HTTP DELETE against the API, with the
// library handling all auth but *not* processing the response.
func (c *Conch) RawDelete(url string) (*http.Response, error) {
	req, err := c.sling().New().Delete(url).Request()
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}

// RawPost allows the user to perform an HTTP POST against the API, with the
// library handling all auth but *not* processing the response.
// The provided body *must* be JSON for the server to accept it.
func (c *Conch) RawPost(url string, body io.Reader) (*http.Response, error) {
	req, err := c.sling().New().Post(url).
		Set("Content-Type", "application/json").Body(body).Request()
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}
