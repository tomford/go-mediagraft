package oauth

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
)

// Return the client and port we should use for the oauth hash
// this is the host and port as the endpoint would naturally see,
// vs the target IP and Port the client targets
func requestedHostPort(r clientReq) (h string, p string) {
	reqParts := strings.Split(r.Req.Host, ":")
	reqURLParts := strings.Split(r.Req.URL.Host, ":")
	reqURLScehem := r.Req.URL.Scheme

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

func hashClientReq(r clientReq) {
	log.Println(r)
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	fmt.Fprintf(w, "%s\n", r.Token)
	fmt.Fprintf(w, "%s\n", strconv.FormatInt(r.Time.Unix(), 10))
	fmt.Fprintf(w, "%s\n", r.Nonce)
	fmt.Fprintf(w, "\n")                 // Body hash, not using yet
	fmt.Fprintf(w, "%s\n", r.Req.Method) // Body hash, not using yet

	h, p := requestedHostPort(r)
	fmt.Fprintf(w, "%s\n%s\n", h, p)

	fmt.Fprintf(w, "%s\n", r.Req.URL.Path) // Body hash, not using yet

	// Must output the headers in sorted order
	var hks []string
	for k := range r.Req.Header {
		hks = append(hks, k)
	}
	sort.Strings(hks)

	// To perform the opertion you want
	for _, k := range hks {
		fmt.Fprintf(w, "%s=%s\n", k, r.Req.Header[k][0])
	}

	w.Flush()

	mac := hmac.New(sha1.New, []byte(r.Token))
	io.Copy(mac, w)

	str := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	log.Println(b.String())
	log.Println(str)
}
