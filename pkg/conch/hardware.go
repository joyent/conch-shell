// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"fmt"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// GetHardwareProduct fetches a single hardware product via
// /hardware/product/:uuid
func (c *Conch) GetHardwareProduct(
	hardwareProductUUID fmt.Stringer,
) (hp HardwareProduct, err error) {
	return hp, c.get("/hardware_product/"+hardwareProductUUID.String(), &hp)
}

// GetHardwareProducts fetches a single hardware product via
// /hardware_product
func (c *Conch) GetHardwareProducts() ([]HardwareProduct, error) {
	prods := make([]HardwareProduct, 0)
	return prods, c.get("/hardware_product", &prods)
}

// SaveHardwareProduct creates or saves s hardware product, based
// on the presence of an ID
func (c *Conch) SaveHardwareProduct(h *HardwareProduct) error {
	if h.Name == "" {
		return ErrBadInput
	}

	if h.Alias == "" {
		return ErrBadInput
	}

	if uuid.Equal(h.HardwareVendorID, uuid.UUID{}) {
		return ErrBadInput
	}

	out := struct {
		Name              string `json:"name"`
		Alias             string `json:"alias"`
		Prefix            string `json:"prefix"`
		HardwareVendorID  string `json:"hardware_vendor_id"`
		Specification     string `json:"specification,omitempty"`
		SKU               string `json:"sku"`
		GenerationName    string `json:"generation_name"`
		LegacyProductName string `json:"legacy_product_name"`
	}{
		h.Name,
		h.Alias,
		h.Prefix,
		h.HardwareVendorID.String(),
		h.Specification,
		h.SKU,
		h.GenerationName,
		h.LegacyProductName,
	}

	if uuid.Equal(h.ID, uuid.UUID{}) {
		return c.post("/hardware_product", out, &h)
	} else {
		return c.post(
			"/hardware_product/"+h.ID.String(), out, &h)
	}
}

// DeleteHardwareProduct deletes a hardware product by marking it as
// deactivated
func (c *Conch) DeleteHardwareProduct(hwUUID fmt.Stringer) error {
	return c.httpDelete("/hardware_product/" + hwUUID.String())
}

// GetHardwareVendor ...
func (c *Conch) GetHardwareVendor(name string) (v HardwareVendor, err error) {
	return v, c.get("/hardware_vendor/"+name, &v)
}

func (c *Conch) GetHardwareVendorByID(id fmt.Stringer) (v HardwareVendor, err error) {
	return v, c.get("/hardware_vendor/"+id.String(), &v)
}

// GetHardwareVendors ...
func (c *Conch) GetHardwareVendors() ([]HardwareVendor, error) {
	vendors := make([]HardwareVendor, 0)
	return vendors, c.get("/hardware_vendor", &vendors)
}

// DeleteHardwareVendor ...
func (c *Conch) DeleteHardwareVendor(name string) error {
	return c.httpDelete("/hardware_vendor/" + name)
}

// SaveHardwareVendor ...
func (c *Conch) SaveHardwareVendor(v *HardwareVendor) error {
	if v.Name == "" {
		return ErrBadInput
	}

	if !uuid.Equal(v.ID, uuid.UUID{}) {
		return ErrBadInput
	}

	out := struct {
		Name string `json:"name"`
	}{v.Name}

	return c.post(
		"/hardware_vendor/"+v.Name,
		out,
		&v,
	)
}
