// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package util contains common routines used throughout the command base
package util

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/joyent/conch-shell/pkg/conch/uuid"
)

// MagicWorkspaceID takes a string and tries to find a valid UUID. If the
// string is a UUID, it doesn't get checked further. If not, we dig through
// GetWorkspaces() looking for UUIDs that match up to the first hyphen or where
// the workspace name matches the string
func MagicWorkspaceID(wat string) (id uuid.UUID, err error) {
	id, err = uuid.FromString(wat)
	if err == nil {
		return id, err
	}
	// So, it's not a UUID. Let's try for a string name or partial UUID
	workspaces, err := API.GetWorkspaces()
	if err != nil {
		return id, err
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s-", wat))
	for _, w := range workspaces {
		if (w.Name == wat) || re.MatchString(w.ID.String()) {
			return w.ID, nil
		}
	}

	return id, errors.New("Could not find workspace " + wat)
}

// MagicWorkspaceRackID takes a workspace UUID and a string and tries to find a
// valid rack UUID. If the string is a UUID, it doesn't get checked further. If
// it's not a UUID, we dig through GetWorkspaceRacks() looking for UUIDs that
// match up to the first hyphen or where the name matches the string.
func MagicWorkspaceRackID(workspace fmt.Stringer, wat string) (uuid.UUID, error) {
	id, err := uuid.FromString(wat)
	if err == nil {
		return id, err
	}

	// So, it's not a UUID. Let's try for a string name or partial UUID
	racks, err := API.GetWorkspaceRacks(workspace)
	if err != nil {
		return id, err
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s-", wat))
	for _, r := range racks {
		if (r.Name == wat) || re.MatchString(r.ID.String()) {
			return r.ID, nil
		}
	}

	return id, errors.New("Could not find rack " + wat)
}

// MagicRackID takes a string and tries to find a valid global rack UUID.
// If the string is a UUID, it doesn't get checked further. If it's not a UUID,
// we dig through GetRacks() looking for UUIDs that match up to the first
// hyphen.
func MagicRackID(wat string) (uuid.UUID, error) {
	id, err := uuid.FromString(wat)
	if err == nil {
		return id, err
	}

	// So, it's not a UUID. Let's try for a string name or partial UUID
	racks, err := API.GetRacks()
	if err != nil {
		return id, err
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s-", wat))
	for _, r := range racks {
		if re.MatchString(r.ID.String()) {
			return r.ID, nil
		}
	}

	return id, errors.New("Could not find rack " + wat)
}

// MagicProductID takes a string and tries to find a valid UUID. If the
// string is a UUID, it doesn't get checked further. If not, we dig through
// GetHardwareProducts() looking for UUIDs that match up to the first hyphen or
// where the product name or SKU matches the string
func MagicProductID(wat string) (uuid.UUID, error) {
	id, err := uuid.FromString(wat)
	if err == nil {
		return id, err
	}

	// So, it's not a UUID. Let's try for a string name or partial UUID
	d, err := API.GetHardwareProducts()
	if err != nil {
		return id, err
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s-", wat))
	for _, r := range d {
		if (r.Name == wat) || (r.SKU == wat) || re.MatchString(r.ID.String()) {
			return r.ID, nil
		}
	}

	return id, errors.New("Could not find product " + wat)
}

// MagicValidationID takes a string and tries to find a valid UUID. If the
// string is a UUID, it doesn't get checked further. Otherwise, we use
// FindShortUUID to see if the string matches an existing Validation ID
func MagicValidationID(s string) (uuid.UUID, error) {
	id, err := uuid.FromString(s)
	if err == nil {
		return id, err
	}

	vs, err := API.GetValidations()
	if err != nil {
		return id, err
	}
	ids := make([]uuid.UUID, len(vs))
	for i, v := range vs {
		ids[i] = v.ID
	}
	id, err = FindShortUUID(s, ids)

	return id, err

}

// MagicValidationPlanID takes a string and tries to find a valid UUID. If the
// string is a UUID, it doesn't get checked further. Otherwise, we use
// FindShortUUID to see if the string matches an existing Validation Plan ID
func MagicValidationPlanID(s string) (uuid.UUID, error) {
	id, err := uuid.FromString(s)
	if err == nil {
		return id, err
	}

	vs, err := API.GetValidationPlans()
	if err != nil {
		return id, err
	}
	ids := make([]uuid.UUID, len(vs))
	for i, v := range vs {
		ids[i] = v.ID
	}
	id, err = FindShortUUID(s, ids)

	return id, err

}

// FindShortUUID takes a string and tries to find a UUID in a list of UUIDs
// that match by prefix (first 4 bytes)
func FindShortUUID(s string, uuids []uuid.UUID) (uuid.UUID, error) {
	re := regexp.MustCompile(fmt.Sprintf("^%s-", s))
	for _, uuid := range uuids {
		if re.MatchString(uuid.String()) {
			return uuid, nil
		}
	}
	var id uuid.UUID
	return id, errors.New("Could not find short UUID " + s)

}

// MagicDatacenterID takes a string and tries to find a valid global
// datacenter UUID.  If the string is a UUID, it doesn't get checked further.
// If it's not a UUID, we dig through GetDatacenters() looking for UUIDs
// that match up to the first hyphen.
// *NOTE*: This will fail if the user is not a global admin
func MagicDatacenterID(wat string) (uuid.UUID, error) {
	id, err := uuid.FromString(wat)
	if err == nil {
		return id, err
	}

	// So, it's not a UUID. Let's try for a partial UUID
	ds, err := API.GetDatacenters()
	if err != nil {
		return id, err
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s-", wat))
	for _, d := range ds {
		if re.MatchString(d.ID.String()) {
			return d.ID, nil
		}
	}

	return id, errors.New("Could not find datacenter " + wat)
}

// MagicRoomID takes a string and tries to find a valid global UUID.  If
// the string is a UUID, it doesn't get checked further.  If it's not a UUID,
// we dig through GetRooms() looking for UUIDs that match up to the first
// hyphen.
func MagicRoomID(wat string) (uuid.UUID, error) {
	id, err := uuid.FromString(wat)
	if err == nil {
		return id, err
	}

	// So, it's not a UUID. Let's try for a partial UUID
	ds, err := API.GetRooms()
	if err != nil {
		return id, err
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s-", wat))
	for _, d := range ds {
		if re.MatchString(d.ID.String()) {
			return d.ID, nil
		}
	}

	return id, errors.New("Could not find room " + wat)
}

// MagicRackRoleID takes a string and tries to find a valid UUID. If the
// string is a UUID, it doesn't get checked further. If not, we dig through
// GetRackRoles() looking for UUIDs that match up to the first hyphen or
// where the role name matches the string
func MagicRackRoleID(wat string) (id uuid.UUID, err error) {
	id, err = uuid.FromString(wat)
	if err == nil {
		return id, err
	}
	// So, it's not a UUID. Let's try for a string name or partial UUID
	ret, err := API.GetRackRoles()
	if err != nil {
		return id, err
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s-", wat))
	for _, r := range ret {
		if (r.Name == wat) || re.MatchString(r.ID.String()) {
			return r.ID, nil
		}
	}

	return id, errors.New("Could not find rack role " + wat)
}

// MagicRackLayoutSlotID takes a string and tries to find a valid UUID.
// If the string is a UUID, it doesn't get checked further.  If it's not a
// UUID, we dig through GetRackLayoutSlots() looking for UUIDs that
// match up to the first hyphen.
func MagicRackLayoutSlotID(wat string) (uuid.UUID, error) {
	id, err := uuid.FromString(wat)
	if err == nil {
		return id, err
	}

	// So, it's not a UUID. Let's try for a partial UUID
	ds, err := API.GetRackLayoutSlots()
	if err != nil {
		return id, err
	}

	re := regexp.MustCompile(fmt.Sprintf("^%s-", wat))
	for _, d := range ds {
		if re.MatchString(d.ID.String()) {
			return d.ID, nil
		}
	}

	return id, errors.New("Could not find rack layout " + wat)
}
