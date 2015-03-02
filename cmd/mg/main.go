package main

import (
	"fmt"
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

	r, _ := c.SimpleSearch("jimi hendrix purple haze", []string{"artists"})

	a := r.Artists[0]
	log.Println(a.Id)
	log.Println(a.Images)

	sid := mg.StationIdent(fmt.Sprintf("a%v", a.Id))
	s, err := c.GetStation(sid)
	log.Println(err)
	log.Println(s)

	/*
		s, err := c.StreamInfo(t.Id, "RADIO", 0, []string{"MP3"})
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(s)
		log.Println(s.Unique)
		log.Println(s.Location)

		d := (time.Second * 15)
		time.Sleep(d)

		err = c.StreamEnd(s.Unique, d, 0)
		log.Println(err)
	*/

	return
}
