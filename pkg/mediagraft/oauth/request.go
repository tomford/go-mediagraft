package oauth

import (
	"net/http"
	"strings"
	"time"
)

// Request holds the details of the
type Request struct {
	Token  string
	Secret string
	Time   time.Time
	Nonce  string
	Req    *http.Request
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
