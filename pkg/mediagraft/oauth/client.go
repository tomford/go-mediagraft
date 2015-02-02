package oauth

import (
	"errors"
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

var DefaultClient = &Client{
	httpClient: http.DefaultClient,
	verbosity:  0,
}

var (
	// OAuth specific Errors
	ErrUnactivatedUser    = errors.New("The user account is not activated. The user must respond to the validation email, or else clients may make the token request again with checkEnabled=false to allow a short grace-period decided by the client")
	ErrBadUserCredentials = errors.New("Incorrect username or password")
	ErrGrantTypeMismatch  = errors.New("The grant type was not set or set to an invalid value")
	ErrUnknownClientID    = errors.New("Unknown client ID")
	ErrBadClientSecret    = errors.New("The client secret supplied does not match the client ID")
)

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

type Credentials struct {
	ClientID     string
	ClientSecret string
	ApiKey       string
	CheckEnabled bool
	Username     string
	Password     string
}

// Do is the http.Do implementation that hides oauth
func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	return nil, nil
}

// Do is the http.Get implementation that hides oauth
func (c *Client) Get(url string) (resp *http.Response, err error) {
	return nil, nil
}

func Get(url string) (resp *http.Response, err error) {
	return DefaultClient.Get(url)
}

func (c *Client) Head(url string) (resp *http.Response, err error) {
	return nil, nil
}

func Head(url string) (resp *http.Response, err error) {
	return DefaultClient.Head(url)
}

func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	return nil, nil
}

func Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	return DefaultClient.Post(url, bodyType, body)
}

func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return nil, nil
}

func PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return DefaultClient.PostForm(url, data)
}

// oauthJSONResp maps to the json data returned from the
// mediagraft oauth implementation
type oauthJSONResp struct {
	TokenType    string `json:token_type`
	Algorithm    string `json:algorithm`
	Secret       string `json:secret`
	AccessToken  string `json:access_token`
	ExpiresIn    string `json:expires_in`
	RefreshToken string `rjson:efresh_token`
}
