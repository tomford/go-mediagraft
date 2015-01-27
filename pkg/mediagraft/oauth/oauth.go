package oauth

import (
	"io"
	"net/http"
	"net/url"
)

// Package oauth implements OAuth 2.0 draft 15 specification, as used
// by Mediagraft

// Client implements an oauth client that transparently handles
// token acquisition refresh
type Client struct {
	verbosity  int
	httpClient *http.Client
}

type option func(c *Client) option

// Option sets the options specified.
// It returns an option to restore the last arg's previous value.
func (c *Client) Option(opts ...option) (previous option) {
	for _, opt := range opts {
		previous = opt(c)
	}
	return previous
}

// Verbosity sets the oauth client's log level
func Verbosity(v int) option {
	return func(c *Client) option {
		previous := c.verbosity
		c.verbosity = v
		return Verbosity(previous)
	}
}

// HTTPClient sets the underlying http.Clinet we'll be using
func HTTPClient(h *http.Client) option {
	return func(c *Client) option {
		previous := c.httpClient
		c.httpClient = h
		return HTTPClient(previous)
	}
}

// New createst a new oauth client
func New() *Client {
	c := &Client{}
	c.Option(
		Verbosity(0),
		HTTPClient(http.DefaultClient),
	)
	return c
}

// Do is the http.Do implementation that hides oauth
func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	return nil, nil
}

// Do is the http.Get implementation that hides oauth
func (c *Client) Get(url string) (resp *http.Response, err error) {
	return nil, nil
}

// Do is the http.Head implementation that hides oauth
func (c *Client) Head(url string) (resp *http.Response, err error) {
	return nil, nil
}

// Do is the http.Post implementation that hides oauth
func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	return nil, nil
}

// Do is the http.PostForm implementation that hides oauth
func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return nil, nil
}
