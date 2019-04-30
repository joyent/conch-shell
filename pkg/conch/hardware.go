// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch

import (
	"encoding/json"
	"fmt"
	"net/url"

	uuid "gopkg.in/satori/go.uuid.v1"
)

// GetHardwareProduct fetches a single hardware product via
// /hardware/product/:uuid
func (c *Conch) GetHardwareProduct(
	hardwareProductUUID fmt.Stringer,
) (hp HardwareProduct, err error) {
	return hp, c.get(
		"/hardware_product/"+url.PathEscape(hardwareProductUUID.String()),
		&hp,
	)
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

	var specification string
	if h.Specification == nil {
		specification = ""
	} else {
		j, err := json.Marshal(h.Specification)

		if err != nil {
			return err
		}

		specification = string(j)
	}

	profile := struct {
		*HardwareProfile
		ID          omit `json:"id,omitempty"`
		Created     omit `json:"created,omitempty"`
		Updated     omit `json:"updated,omitempty"`
		Deactivated omit `json:"deactivated,omitempty"`
	}{HardwareProfile: &h.Profile}

	out := struct {
		*HardwareProduct
		ID            omit        `json:"id,omitempty"`
		Created       omit        `json:"created,omitempty"`
		Updated       omit        `json:"updated,omitempty"`
		Deactivated   omit        `json:"deactivated,omitempty"`
		Specification string      `json:"specification,omitempty"`
		Profile       interface{} `json:"hardware_product_profile,omitempty"`
	}{
		HardwareProduct: h,
		Specification:   specification,
		Profile:         profile,
	}

	if uuid.Equal(h.ID, uuid.UUID{}) {
		return c.post("/hardware_product", out, &h)
	} else {
		return c.post(
			"/hardware_product/"+url.PathEscape(h.ID.String()),
			out,
			&h,
		)
	}
}

// DeleteHardwareProduct deletes a hardware product by marking it as
// deactivated
func (c *Conch) DeleteHardwareProduct(hwUUID fmt.Stringer) error {
	return c.httpDelete("/hardware_product/" + url.PathEscape(hwUUID.String()))
}

// GetHardwareVendor ...
func (c *Conch) GetHardwareVendor(name string) (v HardwareVendor, err error) {
	return v, c.get("/hardware_vendor/"+url.PathEscape(name), &v)
}

func (c *Conch) GetHardwareVendorByID(id fmt.Stringer) (v HardwareVendor, err error) {
	return v, c.get("/hardware_vendor/"+url.PathEscape(id.String()), &v)
}

// GetHardwareVendors ...
func (c *Conch) GetHardwareVendors() ([]HardwareVendor, error) {
	vendors := make([]HardwareVendor, 0)
	return vendors, c.get("/hardware_vendor", &vendors)
}

// DeleteHardwareVendor ...
func (c *Conch) DeleteHardwareVendor(name string) error {
	return c.httpDelete("/hardware_vendor/" + url.PathEscape(name))
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
		"/hardware_vendor/"+url.PathEscape(v.Name),
		out,
		&v,
	)
}
