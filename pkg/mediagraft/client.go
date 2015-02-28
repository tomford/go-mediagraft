// Package provides access to We7's mediagraft API
package mediagraft

import (
	"fmt"
	"io"
	"net/http"

	"github.com/we7/go-mediagraft/pkg/mediagraft/oauth"
)

// Client implements an oauth client that transparently handles
// token acquisition refresh
type Client struct {
	Proto    string //http or https, defaults to https
	Host     string //The host:port to connect to, default to the domain
	HostName string //
	Port     string //

	ApiBase    string
	ApiKey     string
	ApiVersion string
	AppVersion string

	verbosity   int
	oauthClient *oauth.Client
}

var DefaultClient = &Client{
	Proto:    "http",
	Host:     "api.we7.com",
	HostName: "",
	Port:     "80",

	ApiBase:    "/api",
	ApiKey:     "gomg",
	ApiVersion: "0.1",
	AppVersion: "1",

	oauthClient: oauth.DefaultClient,
	verbosity:   0,
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

func (c *Client) Verbosity() int {
	return c.verbosity
}

// OAuthClient sets the underlying oauth.Clinet we'll be using
func OAuthClient(o *oauth.Client) option {
	return func(c *Client) option {
		previous := c.oauthClient
		c.oauthClient = o
		return OAuthClient(previous)
	}
}

func (c *Client) OAuthClient() *oauth.Client {
	return c.oauthClient
}

func (c *Client) Call(httpmethod string, method, qs string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s://%s/%s/%s/%s?apiKey=%s&appVersion=%s",
		c.Proto,
		c.Host,
		c.ApiBase,
		c.ApiVersion,
		method,
		c.ApiKey,
		c.AppVersion,
	)

	if qs != "" {
		url += "&" + qs
	}

	r, err := http.NewRequest(httpmethod, url, body)
	if err != nil {
		return nil, err
	}

	if c.HostName != "" {
		r.Host = c.HostName
	}

	return c.OAuthClient().Do(r)
}
