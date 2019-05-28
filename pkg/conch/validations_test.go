// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch_test

import (
	"testing"

	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/joyent/conch-shell/pkg/conch/uuid"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func TestValidationErrors(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetValidations", func(t *testing.T) {
		url := "/validation"
		gock.New(API.BaseURL).Get(url).Reply(400).JSON(ErrApi)
		ret, err := API.GetValidations()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.Validations{})
	})

	t.Run("GetValidationPlans", func(t *testing.T) {
		url := "/validation_plan"

		gock.New(API.BaseURL).Get(url).Reply(400).JSON(ErrApi)
		ret, err := API.GetValidationPlans()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.ValidationPlan{})
	})

	t.Run("GetValidationPlan", func(t *testing.T) {
		id := uuid.NewV4()
		url := "/validation_plan/" + id.String()

		gock.New(API.BaseURL).Get(url).Reply(400).JSON(ErrApi)
		ret, err := API.GetValidationPlan(id)
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, conch.ValidationPlan{})
	})

	t.Run("RunDeviceValidationPlan", func(t *testing.T) {
		dID := "test"
		vpID := uuid.NewV4()
		url := "/device/" + dID + "/validation_plan/" + vpID.String()

		gock.New(API.BaseURL).Post(url).Reply(400).JSON(ErrApi)
		ret, err := API.RunDeviceValidationPlan(dID, vpID, "{}")
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.ValidationResult{})
	})

	t.Run("RunDeviceValidation", func(t *testing.T) {
		dID := "test"
		vpID := uuid.NewV4()
		url := "/device/" + dID + "/validation/" + vpID.String()

		gock.New(API.BaseURL).Post(url).Reply(400).JSON(ErrApi)
		ret, err := API.RunDeviceValidation(dID, vpID, "{}")
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, []conch.ValidationResult{})
	})

}
