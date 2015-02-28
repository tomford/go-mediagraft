package main

import (
	"fmt"
	"io"
	"log"
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

	rurl := fmt.Sprintf("http://%s/api/0.1/simpleSearch?apiKey=%s&appVersion=%s&format=json&type=genres&query=blues",
		testdomain,
		creds.ApiKey,
		"1",
	)

	log.Println("QUERY 1...")
	r, _ := http.NewRequest("GET", rurl, nil)
	resp, err := c.Do(r)
	if err == nil {
		io.Copy(os.Stdout, resp.Body)
	}

	log.Println(err)
	log.Println(resp)

	log.Println("QUERY 2...")
	r, _ = http.NewRequest("GET", rurl, nil)
	resp, err = c.Do(r)
	if err == nil {
		io.Copy(os.Stdout, resp.Body)
	}

	log.Println(err)
	log.Println(resp)

	return
}
