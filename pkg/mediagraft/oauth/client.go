package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Package oauth implements OAuth 2.0 draft 15 specification, as used
// by Mediagraft

// Client implements an oauth client that transparently handles
// token acquisition refresh
type Client struct {
	verbosity   int
	httpClient  *http.Client
	credentials *credentialMap
}

var DefaultClient = &Client{
	httpClient:  http.DefaultClient,
	verbosity:   0,
	credentials: &credentialMap{},
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

func (c *Client) Verbosity() int {
	return c.verbosity
}

// HTTPClient sets the underlying http.Clinet we'll be using
func HTTPClient(h *http.Client) option {
	return func(c *Client) option {
		previous := c.httpClient
		c.httpClient = h
		return HTTPClient(previous)
	}
}

func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// Credentials is the full set of oauth credentials required to
// log into an oauth service
type Credentials struct {
	Proto        string //http or https, defaults to https
	Host         string //The host:port to connect to, default to the domain
	HostName     string //If set, this is used as the host header, default to the domain
	TokenPath    string
	AuthPath     string
	ClientID     string
	ClientSecret string
	ApiKey       string
	CheckEnabled bool
	Username     string
	Password     string
}

func DefaultCredentials() Credentials {
	return Credentials{
		Proto:        "https",
		TokenPath:    "/oauth/2/token",
		AuthPath:     "/oauth/2/authorize",
		ApiKey:       "test",
		CheckEnabled: false,
	}
}

// CredentialMap maps oauth credentials to the domain they are used within
type credentialMap struct {
	credsLock sync.RWMutex
	creds     map[string]*Credentials
}

func (c *Client) AddDomain(domain string, creds Credentials) {
	c.credentials.credsLock.Lock()
	defer c.credentials.credsLock.Unlock()
	if c.credentials.creds == nil {
		c.credentials.creds = make(map[string]*Credentials)
	}
	c.credentials.creds[domain] = &creds
}

func (c *Client) getDomains(domain string) (creds *Credentials, ok bool) {
	c.credentials.credsLock.RLock()
	defer c.credentials.credsLock.RUnlock()
	creds, ok = c.credentials.creds[domain]
	return creds, ok
}

// Do is the http.Do implementation that hides oauth
func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	h, _ := requestedHostPort(req)
	creds, ok := c.getDomains(h)

	if !ok {
		// We have no oauth creds for this domain, pass it directly
		// to the http.Client
		return c.httpClient.Do(req)
	}

	err = c.getToken(h, creds)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Get is the http.Get implementation that hides oauth
func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Get uses the DefaultClient to perform a GET request
func Get(url string) (resp *http.Response, err error) {
	return DefaultClient.Get(url)
}

func (c *Client) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func Head(url string) (resp *http.Response, err error) {
	return DefaultClient.Head(url)
}

func (c *Client) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return c.Do(req)
}

func Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	return DefaultClient.Post(url, bodyType, body)
}

func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return DefaultClient.PostForm(url, data)
}

// oauthJSONResp maps to the json data returned from the
// mediagraft oauth implementation
type oauthJSONResp struct {
	TokenType        string `json:"token_type"`
	Algorithm        string `json:"algorithm"`
	Secret           string `json:"secret"`
	AccessToken      string `json:"access_token"`
	ExpiresIn        string `json:"expires_in"`
	RefreshToken     string `json:"efresh_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (c *Client) getToken(domain string, creds *Credentials) error {
	h := domain
	if creds.Host != "" {
		h = creds.Host
	}

	hh := domain
	if creds.HostName != "" {
		hh = creds.HostName
	}

	urlArgs := "grant_type=password"
	urlArgs += fmt.Sprintf("&client_id=%s", url.QueryEscape(creds.ClientID))
	urlArgs += fmt.Sprintf("&client_secret=%s", url.QueryEscape(creds.ClientSecret))
	urlArgs += fmt.Sprintf("&username=%s", url.QueryEscape(creds.Username))
	urlArgs += fmt.Sprintf("&password=%s", url.QueryEscape(creds.Password))
	url := fmt.Sprintf("%s://%s/%s?%s", creds.Proto, h, creds.TokenPath, urlArgs)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Host = hh

	resp, err := c.httpClient.Do(req)

	log.Println(url, resp, err)

	var oresp oauthJSONResp
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&oresp); err == io.EOF {
	} else if err != nil {
		return err
	}

	if err = oresp.Err(); err != nil {
		return err
	}

	log.Println(oresp)

	return nil
}

func (r *oauthJSONResp) Err() error {
	switch r.Error {
	case "":
		{
			return nil
		}
	default:
		{
			return fmt.Errorf("%s: %s", r.Error, r.ErrorDescription)
		}
	}
}
