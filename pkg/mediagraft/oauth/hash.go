package oauth

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

// Authorization generates the oauth Authorization header for a
// given request e.g.
//   Authorization: MAC token="IZAxYqW3gyxYMoXy7cAu33VH52slX6TfbxHEjajECUi6EOGH4dhN9Cy++tJ3iI\/WsqrSq04CM+S4Yu4R2QZBZQ==",timestamp="1312472895&",nonce="gfn2lfvn5asfo",signature="38kvZAJcf+Xq+W/Zs+7nG9ClZnI="
func (c *Credentials) Authorization(r *http.Request, t time.Time, nonce string) string {
	c.credLock.RLock()
	defer c.credLock.RUnlock()

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	fmt.Fprintf(w, "%s\n", c.AccessToken)
	fmt.Fprintf(w, "%s\n", strconv.FormatInt(t.Unix(), 10))
	fmt.Fprintf(w, "%s\n", nonce)
	fmt.Fprintf(w, "\n")             // Body hash, not using yet
	fmt.Fprintf(w, "%s\n", r.Method) // Body hash, not using yet

	h, p := requestedHostPort(r)
	fmt.Fprintf(w, "%s\n%s\n", h, p)

	fmt.Fprintf(w, "%s\n", r.URL.Path) // Body hash, not using yet

	// Must output the headers in sorted order
	var qs []string
	for k := range r.URL.Query() {
		key := url.QueryEscape(k)
		val := url.QueryEscape(r.URL.Query().Get(k))
		qs = append(qs, key+"="+val)
	}
	sort.Strings(qs)

	// To perform the opertion you want
	for _, k := range qs {
		fmt.Fprint(w, k+"\n")
	}

	w.Flush()
	//log.Println("blah", string(b.Bytes()))

	s := hmacSha1(&b, []byte(c.Secret))

	str := base64.StdEncoding.EncodeToString(s)

	return fmt.Sprintf(
		"%s token=\"%s\",timestamp=\"%d\",nonce=\"%s\",signature=\"%s\"",
		c.TokenType,
		c.AccessToken,
		t.Unix(),
		nonce,
		str,
	)
}

func hmacSha1(r io.Reader, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	io.Copy(mac, r)

	return mac.Sum(nil)
}
