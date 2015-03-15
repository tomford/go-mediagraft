package mediagraft

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) AlbumsInfo(albumId ...int32) ([]Album, error) {
	args := &url.Values{}

	var strids []string
	for _, t := range albumId {
		strids = append(strids, strconv.Itoa(int(t)))
	}
	args.Set("ids", strings.Join(strids, ","))
	args.Set("detail", "full")

	r, err := c.Call("GET", "albumsInfo", args, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var a []Album
	//io.Copy(os.Stdout, r.Body)
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&a)
	if err != nil {
		return nil, err
	}

	return a, nil
}
