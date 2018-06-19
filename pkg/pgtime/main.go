// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package pgtime provides a wrapper around the raw Postgres 'timestamp with
// time zone' that comes back from the API, providing JSON marshalling and
// unmarshalling
package pgtime

import (
	"fmt"
	"time"
)

// PgTime can be used just like time.Time. When used in JSON
// unmarshalling, it parses the field into a standard time.Time structure. If
// the value in the JSON is null, this is set to the zero value of time.Time{}.
// When marshalled, the JSON value will be an epoch integer timestamp, set in
// the UTC timezone
type PgTime struct {
	time.Time
}

// UnmarshalJSON parses the incoming string value into something useful. In
// the case of null or error, the zero value of time.Time{} is used
func (p *PgTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	// Get rid of the quotes "" around the value.
	s = s[1 : len(s)-1]

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		// If this fails too, t will be the zero value and thus still valid
		t, _ = time.Parse("2006-01-02 15:04:05.999999999Z07", s)
	}
	p.Time = t
	return nil
}

// MarshalJSON translates the time struct into a string, specifically an
// epoch integer timestamp set in the UTC timezone. If the time is a zero
// value, the strign "null" is used instead
func (p *PgTime) MarshalJSON() (b []byte, err error) {
	if p.IsZero() {
		return []byte("null"), nil
	}
	ts := fmt.Sprintf("%d", p.UTC().Unix())
	return []byte(ts), nil
}

// AsUTC returns a time.Time in UTC timezone
func (p *PgTime) AsUTC() time.Time {
	return p.UTC()
}
