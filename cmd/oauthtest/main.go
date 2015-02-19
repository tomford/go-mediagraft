package main

import (
	"net/http"
	"os"

	"github.com/we7/go-mediagraft/pkg/mediagraft/oauth"
)

var testdomain = "api.stagingf.we7.com"

func main() {

	c := oauth.DefaultClient

	creds := oauth.DefaultCredentials()
	creds.ClientID = os.Getenv("OAUTH_CLIENT_ID")
	creds.ClientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
	creds.Username = os.Getenv("OAUTH_USERNAME")
	creds.Password = os.Getenv("OAUTH_PASSWORD")

	c.AddDomain(testdomain, creds)

	r, _ := http.NewRequest("GET", "http://"+testdomain, nil)
	resp, err := c.Do(r)

	return
}
