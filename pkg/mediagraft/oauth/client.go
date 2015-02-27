package oauth

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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
	ErrBadExpiresAt       = errors.New("The expires_at value was unparsable")
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
	RedirectURI  string

	credLock          *sync.RWMutex //http or https, defaults to https
	TokenType         string
	Algorithm         string
	Secret            string
	ExpiresAt         time.Time
	AccessToken       string
	RefreshToken      string
	AuthorizationCode string
}

func DefaultCredentials() Credentials {
	return Credentials{
		Proto:        "https",
		TokenPath:    "/oauth/2/token",
		AuthPath:     "/oauth/2/authorize",
		ApiKey:       "test",
		CheckEnabled: false,
		credLock:     &sync.RWMutex{},
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
func (c *Client) Do(r *http.Request) (resp *http.Response, err error) {
	h, _ := requestedHostPort(r)
	creds, ok := c.getDomains(h)

	if !ok {
		// We have no oauth creds for this domain, pass it directly
		// to the http.Client
		return c.httpClient.Do(r)
	}

	// If we have no token, get one: grant_type=passord
	// If we have check the expiry, if < 1 min (or configurable), do a refresh
	//
	// if we think we have a valid token, make the call.
	//
	// if we get a 401 back,
	//   if our token is still valid now return 401 to the user
	//   if our token is invalid now, refresh the token
	//
	// Make the call again with the new token
	//   if we get another 401 back, assume either our auth is failing, or we
	//   just aren't allowed to call that endpoint

	err = creds.updateCreds(h, c.httpClient)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorization", creds.Authorization(r, time.Now(), nonce()))

	return c.httpClient.Do(r)
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
	ExpiresIn        string `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Reason           string `json:"reason"`
	// These fields are optionsal
	AccessToken       *string `json:"access_token"`
	RefreshToken      *string `json:"refresh_token"`
	AuthorizationCode *string `json:"authorizationCode"`
}

func (c *Credentials) updateCreds(domain string, cl *http.Client) error {
	c.credLock.Lock()
	defer c.credLock.Unlock()

	var err error
	var oresp *oauthJSONResp
	switch {
	case c.AccessToken == "":
		oresp, err = c.getNewToken(domain, "password", cl)
	case time.Now().After(c.ExpiresAt):
		oresp, err = c.getNewToken(domain, "refresh_token", cl)
	default:
		// We have a token, and we think it is valid
		return nil
	}
	if err != nil {
		return err
	}

	c.Algorithm = oresp.Algorithm
	c.TokenType = oresp.TokenType
	c.Secret = oresp.Secret
	if oresp.AccessToken != nil {
		c.AccessToken = *oresp.AccessToken
	}
	if oresp.RefreshToken != nil {
		c.RefreshToken = *oresp.RefreshToken
	}
	if oresp.AuthorizationCode != nil {
		c.AuthorizationCode = *oresp.AuthorizationCode
	}

	expAt, err := strconv.Atoi(oresp.ExpiresIn)
	if err != nil {
		return ErrBadExpiresAt
	}

	c.ExpiresAt = time.Now().Add(time.Second * time.Duration(expAt))

	return nil
}

func (c *Credentials) getNewToken(domain string, grantType string, cl *http.Client) (*oauthJSONResp, error) {
	h := domain
	if c.Host != "" {
		h = c.Host
	}

	hh := domain
	if c.HostName != "" {
		hh = c.HostName
	}

	urlArgs := "grant_type=" + grantType
	urlArgs += fmt.Sprintf("&client_id=%s", url.QueryEscape(c.ClientID))
	urlArgs += fmt.Sprintf("&client_secret=%s", url.QueryEscape(c.ClientSecret))
	switch grantType {
	case "password":
		urlArgs += fmt.Sprintf("&username=%s", url.QueryEscape(c.Username))
		urlArgs += fmt.Sprintf("&password=%s", url.QueryEscape(c.Password))
	case "refresh_token":
		urlArgs += fmt.Sprintf("&refresh_token=%s", url.QueryEscape(c.RefreshToken))
	case "authorization_code":
		urlArgs += fmt.Sprintf("&redirect_uri=%s", url.QueryEscape(c.RedirectURI))
	default:
		return nil, ErrGrantTypeMismatch
	}
	url := fmt.Sprintf("%s://%s/%s?%s", c.Proto, h, c.TokenPath, urlArgs)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Host = hh

	resp, err := cl.Do(req)

	var oresp oauthJSONResp
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&oresp); err == io.EOF {
	} else if err != nil {
		return nil, err
	}

	if err = oresp.Err(); err != nil {
		return nil, err
	}

	return &oresp, err
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

var nonceCounter uint64

// nonce returns a unique string, stolen from another oauth library
func nonce() string {
	n := atomic.AddUint64(&nonceCounter, 1)
	if n == 1 {
		binary.Read(rand.Reader, binary.BigEndian, &n)
		n ^= uint64(time.Now().UnixNano())
		atomic.CompareAndSwapUint64(&nonceCounter, 1, n)
	}
	return strconv.FormatUint(n, 16)
}

// Return the client and port we should use for the oauth hash
// this is the host and port as the endpoint would naturally see,
// vs the target IP and Port the client targets
func requestedHostPort(r *http.Request) (h string, p string) {
	reqParts := strings.Split(r.Host, ":")
	reqURLParts := strings.Split(r.URL.Host, ":")
	reqURLScehem := r.URL.Scheme

	switch {
	case reqParts[0] != "": // A hostname was explicitly given in the request
		h = reqParts[0]
	case reqURLParts[0] != "": // A hostname was explicitly given in the URL
		h = reqURLParts[0]
	default:
		panic("Could not determine requested host") // Really shouldn't get here
	}

	switch {
	case len(reqParts) == 2: //A port was explicitly given in the request
		p = reqParts[1]
	case len(reqURLParts) == 2: //A port was explicitly given in URL
		p = reqURLParts[1]
	case reqURLScehem == "http": //default to 80 for http
		p = "80"
	case reqURLScehem == "https": //default to 443 for https
		p = "443"
	default:
		panic("Could not determine requested port")
	}

	return h, p
}
