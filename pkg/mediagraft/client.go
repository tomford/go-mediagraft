// Package provides access to We7's mediagraft API
package mediagraft

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

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

// New returns a new oauth client with the default settings
func New() *Client {
	c := *DefaultClient
	return &c
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

func (c *Client) Call(httpmethod string, method string, vs *url.Values, body io.Reader) (*http.Response, error) {
	u, err := url.Parse(fmt.Sprintf("%s://%s/%s/%s/%s",
		c.Proto,
		c.Host,
		c.ApiBase,
		c.ApiVersion,
		method,
	))
	if err != nil {
		return nil, err
	}
	if vs == nil {
		vs = &url.Values{}
	}

	vs.Set("apiKey", c.ApiKey)
	vs.Set("appVersion", c.AppVersion)
	vs.Set("format", "json")

	u.RawQuery = vs.Encode()

	r, err := http.NewRequest(httpmethod, u.String(), body)
	if err != nil {
		return nil, err
	}

	if c.HostName != "" {
		r.Host = c.HostName
	}

	return c.OAuthClient().Do(r)
}
