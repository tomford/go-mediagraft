package mediagraft

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) TracksInfo(trackId ...int32) ([]Track, error) {
	args := &url.Values{}

	var strids []string
	for _, t := range trackId {
		strids = append(strids, strconv.Itoa(int(t)))
	}
	args.Set("ids", strings.Join(strids, ","))
	args.Set("detail", "full")

	r, err := c.Call("GET", "tracksInfo", args, nil)
	if err != nil {
		return nil, err
	}

	var t []Track
	//io.Copy(os.Stdout, r.Body)
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (c *Client) TrackVersionsInfo(trackId ...int32) ([]Track, error) {
	args := &url.Values{}

	var strids []string
	for _, t := range trackId {
		strids = append(strids, strconv.Itoa(int(t)))
	}

	args.Set("versionIds", strings.Join(strids, ","))
	args.Set("detail", "full")

	r, err := c.Call("GET", "tracksInfo", args, nil)
	if err != nil {
		return nil, err
	}

	var t []Track
	//io.Copy(os.Stdout, r.Body)
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
