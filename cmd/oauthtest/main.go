package main

import (
	"net/http"

	"github.com/we7/go-mediagraft/pkg/mediagraft/oauth"
)

var testdomain = "api.stagingf.we7.com"

func main() {
	c := oauth.DefaultClient

	creds := oauth.DefaultCredentials()
	c.AddDomain(testdomain, creds)

	r, _ := http.NewRequest("GET", "http://slashdot.org", nil)
	c.Do(r)

	r, _ = http.NewRequest("GET", testdomain, nil)
	c.Do(r)

	return
}
