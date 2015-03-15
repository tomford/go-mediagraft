package mediagraft

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"strings"
)

var spaceReplacer = strings.NewReplacer(" ", "+")

type SearchResult struct {
	Artists       []Artist
	Albums        []Album
	Tracks        []Track
	TrackVersions []TrackVersion
	Genres        []Genre
	RadioStations []RadioStation
	Playlists     []Playlist

	SearchResultInfo
}

type SearchResultInfo struct {
	IsSearchFuzzy bool `json:",string"`
	DidYouMean    string
}

type SearchResultsWithInfo struct {
	Data struct {
		SearchResults     SearchResult
		SearchResultsInfo struct {
			TotalNumberOfResults int  `json:",string"`
			ResultsExhaused      bool `json:",string"`
			SearchResultInfo
		}
	}
	Status string
}

type searchOpt func(c *Search) searchOpt
type Search struct {
	order                  *string
	orderDirection         *string
	limitBegin             *int
	limitEnd               *int
	exact                  *bool
	artistIDs              []int
	restrictedToStreamable *bool
	allowExplicit          *bool
	useSpellCheck          *bool
}

func (s *Search) Option(opts ...searchOpt) (previous searchOpt) {
	for _, opt := range opts {
		previous = opt(s)
	}
	return previous
}

// Verbosity sets the oauth client's log level
func Order(o *string) searchOpt {
	return func(s *Search) searchOpt {
		previous := s.order
		s.order = o
		return Order(previous)
	}
}

func (s *Search) Order() *string {
	return s.order
}

// Verbosity sets the oauth client's log level
func OrderDirection(o *string) searchOpt {
	return func(s *Search) searchOpt {
		previous := s.orderDirection
		s.orderDirection = o
		return OrderDirection(previous)
	}
}

func (s *Search) OrderDirection() *string {
	return s.orderDirection
}

// Verbosity sets the oauth client's log level
func LimitBegin(l *int) searchOpt {
	return func(s *Search) searchOpt {
		previous := s.limitBegin
		s.limitBegin = l
		return LimitBegin(previous)
	}
}

func (s *Search) LimitBegin() *int {
	return s.limitBegin
}

// Verbosity sets the oauth client's log level
func LimitEnd(l *int) searchOpt {
	return func(s *Search) searchOpt {
		previous := s.limitEnd
		s.limitEnd = l
		return LimitEnd(previous)
	}
}

func (s *Search) LimitEnd() *int {
	return s.limitEnd
}

func Exact(b *bool) searchOpt {
	return func(s *Search) searchOpt {
		previous := s.exact
		s.exact = b
		return Exact(previous)
	}
}

func (s *Search) Exact() *bool {
	return s.exact
}

func ArtistIDs(l []int) searchOpt {
	return func(s *Search) searchOpt {
		previous := s.artistIDs
		s.artistIDs = l
		return ArtistIDs(previous)
	}
}

func (s *Search) ArtistIDs() []int {
	return s.artistIDs
}

func RestrictedToStreamable(r *bool) searchOpt {
	return func(s *Search) searchOpt {
		previous := s.restrictedToStreamable
		s.restrictedToStreamable = r
		return RestrictedToStreamable(previous)
	}
}

func (s *Search) RestrictedToStreamable() *bool {
	return s.restrictedToStreamable
}

func AllowExplicit(e *bool) searchOpt {
	return func(s *Search) searchOpt {
		previous := s.allowExplicit
		s.allowExplicit = e
		return AllowExplicit(previous)
	}
}

func (s *Search) AllowExplicit() *bool {
	return s.allowExplicit
}

func UseSpellCheck(c *bool) searchOpt {
	return func(s *Search) searchOpt {
		previous := s.useSpellCheck
		s.useSpellCheck = c
		return UseSpellCheck(previous)
	}
}

func (s *Search) UseSpellCheck() *bool {
	return s.useSpellCheck
}

func (s *Search) args() *url.Values {
	a := &url.Values{}

	if v := s.Order(); v != nil {
		a.Add("order", *v)
	}

	if v := s.OrderDirection(); v != nil {
		a.Add("orderDirection", *v)
	}

	if v := s.AllowExplicit(); v != nil {
		a.Add("allowExplicit", strconv.FormatBool(*v))
	}

	if vs := s.ArtistIDs(); len(vs) != 0 {
		var ss []string
		for _, v := range vs {
			ss = append(ss, strconv.Itoa(v))
		}
		a.Add("artistIds", strings.Join(ss, ","))
	}

	if v := s.Exact(); v != nil {
		a.Add("exact", strconv.FormatBool(*v))
	}

	if v := s.RestrictedToStreamable(); v != nil {
		a.Add("restrictedToStreamable", strconv.FormatBool(*v))
	}

	if v := s.UseSpellCheck(); v != nil {
		a.Add("useSpellCheck", strconv.FormatBool(*v))
	}

	return a
}

func (c *Client) SimpleSearch(q string, types []string, opts ...searchOpt) (*SearchResult, error) {
	r, err := c.doSearch("simpleSearch", q, types, opts...)
	if err != nil {
		return nil, err
	}

	var sr SearchResult
	dec := json.NewDecoder(r)
	err = dec.Decode(&sr)
	if err != nil {
		return nil, err
	}

	return &sr, nil
}

func (c *Client) SimpleSearchWithInfo(q string, types []string, opts ...searchOpt) (*SearchResultsWithInfo, error) {
	r, err := c.doSearch("simpleSearchWithInfo", q, types, opts...)
	if err != nil {
		return nil, err
	}

	var sr SearchResultsWithInfo

	dec := json.NewDecoder(r)
	err = dec.Decode(&sr)
	if err != nil {
		return nil, err
	}

	return &sr, nil
}

func (c *Client) doSearch(method string, q string, types []string, opts ...searchOpt) (io.Reader, error) {
	s := Search{}
	s.Option(opts...)

	args := s.args()
	args.Add("query", spaceReplacer.Replace(q))
	args.Add("type", strings.Join(types, ","))

	r, err := c.Call("GET", method, args, nil)
	if err != nil {
		return nil, err
	}
	return r.Body, nil

}

func (c *Client) FindMatch(title string, artistname string, types []string) (*SearchResult, error) {
	args := &url.Values{}
	args.Add("title", spaceReplacer.Replace(title))
	args.Add("artistName", spaceReplacer.Replace(artistname))
	args.Add("type", strings.Join(types, ","))

	r, err := c.Call("GET", "findMatch", args, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var sr SearchResult
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&sr)
	if err != nil {
		return nil, err
	}

	return &sr, nil
}
