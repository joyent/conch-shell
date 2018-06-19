// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package conch_test

import (
	"errors"
	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
	uuid "gopkg.in/satori/go.uuid.v1"
	"strings"
	"testing"
)

func TestValidationErrors(t *testing.T) {
	BuildAPI()
	gock.Flush()

	aerr := conch.APIError{ErrorMsg: "totally broken"}
	aerrUnpacked := errors.New(aerr.ErrorMsg)

	t.Run("GetValidations", func(t *testing.T) {
		url := "/validation"
		gock.New(API.BaseURL).Get(url).Reply(400).JSON(aerr)
		ret, err := API.GetValidations()
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.Validation{})
	})

	t.Run("GetValidationPlans", func(t *testing.T) {
		url := "/validation_plan"

		gock.New(API.BaseURL).Get(url).Reply(400).JSON(aerr)
		ret, err := API.GetValidationPlans()
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.ValidationPlan{})
	})

	t.Run("GetValidationPlan", func(t *testing.T) {
		id := uuid.NewV4()
		url := "/validation_plan/" + id.String()

		gock.New(API.BaseURL).Get(url).Reply(400).JSON(aerr)
		ret, err := API.GetValidationPlan(id)
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, &conch.ValidationPlan{})
	})

	t.Run("RunDeviceValidationPlan", func(t *testing.T) {
		dID := "test"
		vpID := uuid.NewV4()
		url := "/device/" + dID + "/validation_plan/" + vpID.String()

		gock.New(API.BaseURL).Post(url).Reply(400).JSON(aerr)
		ret, err := API.RunDeviceValidationPlan(dID, vpID, strings.NewReader("{ }"))
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.ValidationResult{})
	})

	t.Run("RunDeviceValidation", func(t *testing.T) {
		dID := "test"
		vpID := uuid.NewV4()
		url := "/device/" + dID + "/validation/" + vpID.String()

		gock.New(API.BaseURL).Post(url).Reply(400).JSON(aerr)
		ret, err := API.RunDeviceValidation(dID, vpID, strings.NewReader("{ }"))
		st.Expect(t, err, aerrUnpacked)
		st.Expect(t, ret, []conch.ValidationResult{})
	})

}

func TestGetValidations(t *testing.T) {
	Recorded("get_validations",
		func(c conch.Conch) {
			vs, err := c.GetValidations()
			st.Expect(t, err, nil)
			st.Reject(t, len(vs), 0)
		},
	)
}

func TestGetValidationPlans(t *testing.T) {
	Recorded("get_validation_plans",
		func(c conch.Conch) {
			vps, err := c.GetValidationPlans()
			st.Expect(t, err, nil)
			st.Reject(t, len(vps), 0)
		},
	)
}

func TestGetValidationPlan(t *testing.T) {
	Recorded("get_validation_plan",
		func(c conch.Conch) {
			vps, err := c.GetValidationPlans()
			st.Assert(t, err, nil)
			st.Refute(t, len(vps), 0)

			vp0 := vps[0]

			vp1, err := c.GetValidationPlan(vp0.ID)
			st.Expect(t, err, nil)
			st.Expect(t, vp1.ID, vp0.ID)
		},
	)
}

func TestRunValidationPlan(t *testing.T) {
	Recorded("run_validation_plan",
		func(c conch.Conch) {
			vps, err := c.GetValidationPlans()
			st.Assert(t, err, nil)
			st.Refute(t, len(vps), 0)

			vp := vps[0]

			ws, err := c.GetWorkspaces()
			st.Assert(t, err, nil)
			st.Refute(t, len(ws), 0)

			w := ws[0]
			ds, err := c.GetWorkspaceDevices(w.ID, true, "", "")
			st.Assert(t, err, nil)
			st.Refute(t, len(ds), 0)

			d := ds[0]
			results, err := c.RunDeviceValidationPlan(d.ID, vp.ID, strings.NewReader("{ }"))
			st.Assert(t, err, nil)
			st.Refute(t, len(results), 0)
		},
	)
}

func TestRunValidation(t *testing.T) {
	Recorded("run_validation",
		func(c conch.Conch) {
			vs, err := c.GetValidations()
			st.Assert(t, err, nil)
			st.Refute(t, len(vs), 0)

			v := vs[0]

			ws, err := c.GetWorkspaces()
			st.Assert(t, err, nil)
			st.Refute(t, len(ws), 0)

			w := ws[0]

			ds, err := c.GetWorkspaceDevices(w.ID, true, "", "")
			st.Assert(t, err, nil)
			st.Refute(t, len(ds), 0)

			d := ds[0]

			results, err := c.RunDeviceValidation(d.ID, v.ID, strings.NewReader("{ }"))
			st.Assert(t, err, nil)
			st.Refute(t, len(results), 0)

		},
	)
}
