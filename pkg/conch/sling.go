// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

const (
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
	userAgent := fmt.Sprintf("Conch/%s", VersionStr)
	if len(c.UserAgent) > 0 {
		for k, v := range c.UserAgent {
			userAgent = fmt.Sprintf("%s %s/%s", userAgent, k, v)
		}
	}

	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Transport: defaultTransport,
		}
	}

	s := sling.New().
		Client(c.HTTPClient).
		Base(c.BaseURL).
		Set("User-Agent", userAgent)

	if c.Token != "" {
		s = s.Set("Authorization", "Bearer "+c.Token)
	}

	return s
}

func (c *Conch) get(url string, data interface{}) error {
	req, err := c.sling().New().Get(url).Request()
	if err != nil {
		return err
	}

	_, err = c.httpDo(req, data)
	return err
}

func (c *Conch) httpDo(req *http.Request, data interface{}) (*http.Response, error) {

	c.debugLog(fmt.Sprintf(
		"Request: %s %s",
		req.Method,
		req.URL,
	))

	if (req.Method == "POST") && (req.Body != nil) {
		if read, err := req.GetBody(); err == nil {
			if bodyBytes, err := ioutil.ReadAll(read); err == nil {
				c.traceLog(
					fmt.Sprintf(
						"  Request Body: %s",
						string(bodyBytes),
					),
				)
			}
		}
	}

	res, err := c.HTTPClient.Do(req)
	if (res == nil) || (err != nil) {
		return res, err
	}

	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}

	if c.Trace {
		c.traceLogDDP(
			fmt.Sprintf(
				"Response: HTTP %d",
				res.StatusCode,
			),
			string(bodyBytes),
		)
	} else {
		c.debugLog(
			fmt.Sprintf(
				"Response: HTTP %d",
				res.StatusCode,
			),
		)
	}

	if res.StatusCode == http.StatusUnauthorized {
		return res, ErrNotAuthorized
	}

	if res.StatusCode == http.StatusForbidden {
		return res, ErrForbidden
	}

	if res.StatusCode == http.StatusNotFound {
		return res, ErrDataNotFound
	}

	// BUG(sungo): an awfully simplistic view of the world
	if code := res.StatusCode; code >= 200 && code < 300 {
		if data != nil {
			// BUG(sungo): do we really want to throw away parse errors?
			json.Unmarshal(bodyBytes, data)

			if c.Trace {
				c.ddp(data)
			}
		}
		return res, nil
	}

	aerr := struct {
		Error string `json:"error"`
	}{""}
	if err := json.Unmarshal(bodyBytes, &aerr); err == nil {
		if c.Trace {
			c.ddp(aerr)
		}
		return res, errors.New(aerr.Error)
	}

	// In general, we should expect the API to give us error structures when
	// things go awry, but just in case not...
	return res, ErrHTTPNotOk
}

func (c *Conch) getWithQuery(url string, query interface{}, data interface{}) error {
	req, err := c.sling().New().Get(url).QueryStruct(query).Request()
	if err != nil {
		return err
	}
	_, err = c.httpDo(req, data)
	return err
}

func (c *Conch) httpDelete(url string) error {
	req, err := c.sling().New().Delete(url).Request()
	if err != nil {
		return err
	}
	_, err = c.httpDo(req, nil)
	return err
}

func (c *Conch) httpDeleteWithPayload(url string, payload interface{}) error {
	req, err := c.sling().New().Delete(url).BodyJSON(payload).Request()
	if err != nil {
		return err
	}
	_, err = c.httpDo(req, nil)
	return err
}

func (c *Conch) post(url string, payload interface{}, response interface{}) error {
	req, err := c.sling().New().
		Post(url).
		BodyJSON(payload).
		Request()

	if err != nil {
		return err
	}

	_, err = c.httpDo(req, response)
	return err
}

//lint:ignore U1000 keeping for later
func (c *Conch) postNeedsResponse(
	url string,
	payload interface{},
	response interface{},

) (*http.Response, error) {
	req, err := c.sling().New().
		Post(url).
		BodyJSON(payload).
		Request()

	if err != nil {
		return nil, err
	}
	res, err := c.httpDo(req, response)
	return res, err
}

//////

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
