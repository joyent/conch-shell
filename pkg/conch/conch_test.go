package conch_test

import (
	"fmt"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/joyent/conch-shell/pkg/conch"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"net/http/cookiejar"
	"os"
)

// Recorded builds a Conch client that uses go-vcr and passes the client to
// function 'runner' to perform requests . Any requests performed will be
// recorded to a yaml file in the fixtures/ directory, named by the string
// parameter 'name'. If the environment variable RE_RECORD is non-empty, all
// existing files are overwritten with new requests. The environment variables
// CONCH_USER and CONCH_PASS must be specified with valid Conch credentials
// when recording requests which require authorization
func Recorded(name string, runner func(conch.Conch)) {
	var mode recorder.Mode
	if os.Getenv("RE_RECORD") != "" {
		mode = recorder.ModeRecording
	} else {
		mode = recorder.ModeReplaying
	}

	r, err := recorder.NewAsMode("fixtures/"+name, mode, nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	defer r.Stop()

	cj, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cj,
	}

	c := conch.Conch{
		HTTPClient: client,
		CookieJar:  cj,
	}

	if os.Getenv("CONCH_URL") != "" {
		c.BaseURL = os.Getenv("CONCH_URL")
	}

	if os.Getenv("CONCH_USER") != "" && os.Getenv("CONCH_PASS") != "" {
		c.Login(os.Getenv("CONCH_USER"), os.Getenv("CONCH_PASS"))
	}

	// Login is conditionally executed before setting the transport to the
	// recorder, to avoid capturing login information
	c.HTTPClient.Transport = r

	runner(c)

}

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
