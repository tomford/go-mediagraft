package main

import (
	"log"
	"os"

	mg "github.com/we7/go-mediagraft/pkg/mediagraft"
	"github.com/we7/go-mediagraft/pkg/mediagraft/oauth"
)

func main() {
	testdomain := "api.stagingf.we7.com"
	creds := oauth.DefaultCredentials()
	creds.ClientID = os.Getenv("OAUTH_CLIENT_ID")
	creds.ClientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
	creds.Username = os.Getenv("OAUTH_USERNAME")
	creds.Password = os.Getenv("OAUTH_PASSWORD")

	c := mg.DefaultClient
	c.ApiKey = "sonos"
	c.Host = testdomain
	c.OAuthClient().AddDomain(testdomain, creds)

	r, err := c.SimpleSearch("jimi", []string{"artists"})
	log.Println(err)
	log.Println(r)

	return
}
