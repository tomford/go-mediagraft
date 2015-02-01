package oauth

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

// From the docs
//
// IZAxYqW3gyxYMoXy7cAu33VH52slX6TfbxHEjajECUi6EOGH4dhN9Cy++tJ3iI\/WsqrSq04CM+S4Yu4R2QZBZQ==
// 1312471030
// gfn2lfvn5asfo
//
// GET
// api.we7.com
// 80
// /api/0.1/userPlaylistsInfo
// apiKey=myKey
// appVersion=1
// detail=full
// format=xml
//
//
// Gives
//
// Authorization: MAC token="IZAxYqW3gyxYMoXy7cAu33VH52slX6TfbxHEjajECUi6EOGH4dhN9Cy++tJ3iI\/WsqrSq04CM+S4Yu4R2QZBZQ==",timestamp="1312472895&",nonce="gfn2lfvn5asfo",signature="38kvZAJcf+Xq+W/Zs+7nG9ClZnI="
//
// http://api.we7.com/api/0.1/userPlaylistsInfo?appVersion=1&apiKey=myKey&format=xml&detail=full
var testReq = struct {
	Token string
	Time  time.Time
	Nonce string
	Req   *http.Request
}{
	"IZAxYqW3gyxYMoXy7cAu33VH52slX6TfbxHEjajECUi6EOGH4dhN9Cy++tJ3iI\\/WsqrSq04CM+S4Yu4R2QZBZQ==",
	time.Unix(1312471030, 0),
	"gfn2lfvn5asfo",
	&http.Request{
		Method: "GET",
		Host:   "api.we7.com",
		URL: &url.URL{
			Scheme: "http",
			Host:   "127.0.0.1:80",
			Path:   "/api/0.1/userPlaylistsInfo",
		},
		Header: map[string][]string{
			"apiKey":     {"myKey"},
			"appVersion": {"1"},
			"detail":     {"full"},
			"format":     {"xml"},
		},
	},
}

func TestGet(t *testing.T) {
	hashClientReq(testReq)
}
