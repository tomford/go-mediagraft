package mediagraft

import (
	"encoding/json"
	"net/url"
)

type StationIdent string
type StationMood struct{}
type StationInfluence struct{}
type StationSeed struct {
	ID     int    `json:"id,string"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Images Images `json:"images"`
	Full   bool   `json:"full"`
}

type Station struct {
	ID         StationIdent `json:"id"`
	Cookie     string       `json:"cookie"`
	Moods      StationMood  `json:"moods"`
	LinkText   string       `json:"linkText"`
	Influences struct {
		Positive []StationInfluence `json:"positive"`
		Negative []StationInfluence `json:"negative"`
	} `json:"influences"`
	Tracks            []Track           `json:"tracks"`
	Artists           []string          `json:"artists"`
	Description       string            `json:"description"`
	Name              string            `json:"name"`
	Images            Images            `json:"images"`
	CategorizedImages map[string]Images `json:"categorizedImages"`
	URL               URL               `json:"url"`
	Searchable        bool              `json:"searchable,string"`
	Popularity        float64           `json:"popularity,string"`
	Subtitle          string            `json:"subtitle"`
	Badge             string            `json:"badge"`
	ExplicitCount     int               `json:"explicitCount"`
	TrackCount        int               `json:"trackCount"`
	Promoted          bool              `json:"promoted"`
	Tags              []string          `json:"tags"`
	Seeds             []StationSeed     `json:"seeds"`
}

func (c *Client) GetStation(ident StationIdent) (*Station, error) {
	args := &url.Values{}
	args.Set("stationIdent", string(ident))

	r, err := c.Call("GET", "radio/getStation", args, nil)
	if err != nil {
		return nil, err
	}

	var s Station
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
