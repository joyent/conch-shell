package conch_test

import (
	"github.com/joyent/conch-shell/pkg/conch"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"net/http/cookiejar"
)

var API *conch.Conch

func BuildAPI() {
	cj, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar:       cj,
		Transport: &http.Transport{},
	}

	gock.InterceptClient(client)
	c := &conch.Conch{
		BaseURL:    "http://localhost:5001",
		HTTPClient: client,
	}
	API = c

	gock.New(c.BaseURL).Get("/version").
		Reply(200).BodyString("{ \"version\": \"99.99.99\" }")
}
