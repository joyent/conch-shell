// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package uuid

import (
	"encoding/json"

	gofrs "github.com/gofrs/uuid"
)

type UUID struct {
	uuid gofrs.UUID
}

func (u UUID) String() string {
	return u.uuid.String()
}

func (u UUID) Equal(u2 UUID) bool {
	return u.String() == u2.String()
}

func (u UUID) IsZero() bool {
	return u.uuid == gofrs.UUID{}
}

func NewV4() UUID {
	u, _ := gofrs.NewV4()
	return UUID{
		uuid: u,
	}
}

func New() UUID {
	return UUID{}
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + u.String() + "\""), nil
}

func (u *UUID) UnmarshalJSON(b []byte) error {
	var frs gofrs.UUID

	if err := json.Unmarshal(b, &frs); err != nil {
		return err
	}

	u.uuid = frs

	return nil
}
