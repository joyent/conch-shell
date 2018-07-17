// Copyright 2017 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	uuid "gopkg.in/satori/go.uuid.v1"
	"net/http"
)

// DBHardwareProduct represents the specification for a specific piece of hardware
type DBHardwareProduct struct {
	ID                uuid.UUID `json:"id,omitempty"`
	Name              string    `json:"name"`
	Alias             string    `json:"alias"`
	Prefix            string    `json:"prefix"`
	Vendor            uuid.UUID `json:"vendor"`
	Specification     string    `json:"specification,omitempty"`
	SKU               string    `json:"sku,omitempty"`
	GenerationName    string    `json:"generation_name,omitempty"`
	LegacyProductName string    `json:"legacy_product_name,omitempty"`
}

// GetDBHardwareProduct fetches a single new-style hardware product
func (c *Conch) GetDBHardwareProduct(hardwareProductUUID fmt.Stringer) (DBHardwareProduct, error) {
	var prod DBHardwareProduct

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/db/hardware_product/"+hardwareProductUUID.String()).
		Receive(&prod, aerr)

	return prod, c.isHTTPResOk(res, err, aerr)
}

// GetDBHardwareProducts fetches all new-style hardware products
func (c *Conch) GetDBHardwareProducts() ([]DBHardwareProduct, error) {
	prods := make([]DBHardwareProduct, 0)

	aerr := &APIError{}
	res, err := c.sling().New().
		Get("/db/hardware_product").
		Receive(&prods, aerr)

	return prods, c.isHTTPResOk(res, err, aerr)
}

// DeleteDBHardwareProduct deletes a hardware product by marking it as
// deactivated
func (c *Conch) DeleteDBHardwareProduct(hwUUID fmt.Stringer) error {
	aerr := &APIError{}
	res, err := c.sling().New().
		Delete("/db/hardware_product/"+hwUUID.String()).
		Receive(nil, aerr)

	return c.isHTTPResOk(res, err, aerr)
}

// SaveDBHardwareProduct creates or saves a new-style hardware product, based
// on the presence of an ID
func (c *Conch) SaveDBHardwareProduct(h *DBHardwareProduct) error {
	if h.Name == "" {
		return ErrBadInput
	}

	if h.Alias == "" {
		return ErrBadInput
	}

	if uuid.Equal(h.Vendor, uuid.UUID{}) {
		return ErrBadInput
	}

	var err error
	var res *http.Response
	aerr := &APIError{}

	if uuid.Equal(h.ID, uuid.UUID{}) {
		out := struct {
			Name              string `json:"name"`
			Alias             string `json:"alias"`
			Prefix            string `json:"prefix"`
			Vendor            string `json:"vendor"`
			Specification     string `json:"specification,omitempty"`
			SKU               string `json:"sku"`
			GenerationName    string `json:"generation_name"`
			LegacyProductName string `json:"legacy_product_name"`
		}{
			h.Name,
			h.Alias,
			h.Prefix,
			h.Vendor.String(),
			h.Specification,
			h.SKU,
			h.GenerationName,
			h.LegacyProductName,
		}

		res, err = c.sling().New().Post("/db/hardware_product").
			BodyJSON(out).Receive(&h, aerr)
	} else {
		res, err = c.sling().New().Post("/db/hardware_product/"+h.ID.String()).
			BodyJSON(h).Receive(&h, aerr)
	}

	return c.isHTTPResOk(res, err, aerr)
}
