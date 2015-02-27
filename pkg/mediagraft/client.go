// Package provides access to We7's mediagraft API
package mediagraft

import "github.com/we7/go-mediagraft/pkg/mediagraft/oauth"

// Client implements an oauth client that transparently handles
// token acquisition refresh
type Client struct {
	verbosity   int
	oauthClient *oauth.Client
}

var DefaultClient = &Client{
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
