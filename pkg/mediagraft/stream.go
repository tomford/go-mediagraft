package mediagraft

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (c *Client) StreamInfo(trackId int, playSource string, playlistId int, musicFormats []string) (*Stream, error) {
	args := &url.Values{}
	args.Set("trackId", strconv.Itoa(trackId))
	args.Set("playSource", playSource)
	if playSource == "PLAYLIST" {
		args.Set("playlistId", strconv.Itoa(playlistId))
	}
	args.Set("musicFormats", strings.Join(musicFormats, ","))

	r, err := c.Call("GET", "streaming/streamInfoWithOAuth", args, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var s Stream
	//io.Copy(os.Stdout, r.Body)
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (c *Client) StreamEnd(u StreamUnique, played, paused time.Duration) error {
	args := &url.Values{}
	args.Set("streamUnique", string(u))
	args.Set("playedTime", strconv.Itoa(int(played/time.Millisecond)))
	args.Set("pausedTime", strconv.Itoa(int(paused/time.Millisecond)))

	_, err := c.Call("POST", "streamEnd", args, nil)
	if err != nil {
		return err
	}

	return nil
}
