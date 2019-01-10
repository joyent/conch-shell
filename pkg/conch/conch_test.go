package conch_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/joyent/conch-shell/pkg/conch"
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

var ErrApi = struct {
	ErrorMsg string `json:"error"`
}{"totally broken"}

var ErrApiUnpacked = errors.New(ErrApi.ErrorMsg)

var API = &conch.Conch{
	BaseURL:    "http://localhost",
	HTTPClient: http.DefaultClient,
}

func BuildAPI() {
	// BUG(sungo): noop until the whole test suite migrates
}

func TestConch(t *testing.T) {
	gock.Flush()
	defer gock.Flush()

	t.Run("GetVersion", func(t *testing.T) {

		good := struct {
			Version string `json:"version"`
		}{"99.99.99"}

		gock.New(API.BaseURL).Get("/version").Reply(200).JSON(good)

		ret, err := API.GetVersion()
		st.Expect(t, err, nil)
		st.Expect(t, ret, "99.99.99")
	})

	t.Run("GetVersionErrors", func(t *testing.T) {
		gock.New(API.BaseURL).Get("/version").Reply(400).JSON(ErrApi)

		ret, err := API.GetVersion()
		st.Expect(t, err, ErrApiUnpacked)
		st.Expect(t, ret, "")
	})
}
