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
	Token  string
	Secret string
	Time   time.Time
	Nonce  string
	Req    *http.Request
}{
	"IZAxYqW3gyxYMoXy7cAu33VH52slX6TfbxHEjajECUi6EOGH4dhN9Cy++tJ3iI\\/WsqrSq04CM+S4Yu4R2QZBZQ==",
	"5t4lGTb2",
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

// From RFC
//     h480djs93hd8\n
//     137131200\n
//     dj83hs9s\n
//     \n
//     GET\n
//     example.com\n
//     80\n
//     /resource/1\n
//     a=2\n
//     b=1\n
var testReq2 = struct {
	Token  string
	Secret string
	Time   time.Time
	Nonce  string
	Req    *http.Request
}{
	"h480djs93hd8",
	"dj83hs9s",
	time.Unix(137131200, 0),
	"dj83hs9s",
	&http.Request{
		Method: "GET",
		Host:   "example.com",
		URL: &url.URL{
			Scheme: "http",
			Host:   "127.0.0.1:80",
			Path:   "/resource/1",
		},
		Header: map[string][]string{
			"a": {"2"},
			"b": {"1"},
		},
	},
}

func TestGet(t *testing.T) {
	hashClientReq(testReq)
	hashClientReq(testReq2)
}
