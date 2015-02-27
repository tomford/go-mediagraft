package oauth

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"testing"
	"time"
)

// Request holds the details of the
type Request struct {
	Token  string
	Secret string
	Time   time.Time
	Nonce  string
	Method string
	URL    string
	Host   string
}

var testReqs = []struct {
	r   Request
	out string
}{
	{
		Request{
			// From the docs
			//
			// http://api.we7.com/api/0.1/userPlaylistsInfo?appVersion=1&apiKey=myKey&format=xml&detail=full
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
			"IZAxYqW3gyxYMoXy7cAu33VH52slX6TfbxHEjajECUi6EOGH4dhN9Cy++tJ3iI\\/WsqrSq04CM+S4Yu4R2QZBZQ==",
			"5t4lGTb2",
			time.Unix(1312471030, 0),
			"gfn2lfvn5asfo",
			"GET",
			"http://127.0.0.1:80/api/0.1/userPlaylistsInfo?apiKey=myKey&appVersion=1&detail=full&format=xml",
			"api.we7.com",
		},
		`MAC token="IZAxYqW3gyxYMoXy7cAu33VH52slX6TfbxHEjajECUi6EOGH4dhN9Cy++tJ3iI\/WsqrSq04CM+S4Yu4R2QZBZQ==",timestamp="1312471030",nonce="gfn2lfvn5asfo",signature="MkvSv/FUo/3HQvTCzPQg2Vm/lUY="`,
	},
	{
		Request{
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
			"h480djs93hd8",
			"489dks293j39",
			time.Unix(137131200, 0),
			"dj83hs9s",
			"GET",
			"http://127.0.0.1:80/resource/1?a=2&b=1",
			"example.com",
		},
		`MAC token="h480djs93hd8",timestamp="137131200",nonce="dj83hs9s",signature="YTVjyNSujYs1WsDurFnvFi4JK6o="`,
	},
}

func TestHashClientReq(t *testing.T) {
	c := DefaultCredentials()
	c.TokenType = "MAC"

	for i, tt := range testReqs {
		c.AccessToken = tt.r.Token
		c.Secret = tt.r.Secret
		r, _ := http.NewRequest(tt.r.Method, tt.r.URL, nil)
		r.Host = tt.r.Host

		a := c.Authorization(r, tt.r.Time, tt.r.Nonce)
		if a != tt.out {
			t.Errorf("%d. failed: expected %s got %s\n", i, tt.out, a)
		}
	}
}

var testHmacSha1 = []struct {
	i string
	k string
	o []byte
}{
	{"", "", mustHexDecodeString("fbdb1d1b18aa6c08324b7d64b71fb76370690e1d")},                                                                                                                                                                                                                                //Wikipedia example empty hmac-sha1
	{"The quick brown fox jumps over the lazy dog", "key", mustHexDecodeString("de7c9b85b8b78aa6bc8a7a36f70a90701c9db4d9")},                                                                                                                                                                                  // Wikipedia example hmac-sha1
	{"IZAxYqW3gyxYMoXy7cAu33VH52slX6TfbxHEjajECUi6EOGH4dhN9Cy++tJ3iI\\/WsqrSq04CM+S4Yu4R2QZBZQ==\n1312471030\ngfn2lfvn5asfo\n\nGET\napi.we7.com\n80\n/api/0.1/userPlaylistsInfo\napiKey=myKey\nappVersion=1\ndetail=full\nformat=xml\n", "5t4lGTb2", mustBase64DecodeString("MkvSv/FUo/3HQvTCzPQg2Vm/lUY=")}, // mediagraft oauth
	{"h480djs93hd8\n137131200\ndj83hs9s\n\nGET\nexample.com\n80\n/resource/1\na=2\nb=1\n", "489dks293j39", mustBase64DecodeString("YTVjyNSujYs1WsDurFnvFi4JK6o=")},
}

func TestHmacSha1(t *testing.T) {
	for i, tt := range testHmacSha1 {
		r := bytes.NewBufferString(tt.i)
		out := hmacSha1(r, []byte(tt.k))
		if bytes.Compare(out, tt.o) != 0 {
			wanted := base64.StdEncoding.EncodeToString(tt.o)
			got := base64.StdEncoding.EncodeToString(out)
			t.Errorf("%d. failed: expected %s got %s\n", i, wanted, got)
		}
	}
}

func mustHexDecodeString(s string) []byte {
	v, err := hex.DecodeString(s)
	if err != nil {
		panic("failed to decode hex, " + err.Error())
	}
	return v
}

func mustBase64DecodeString(s string) []byte {
	v, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic("failed to deode base64, " + err.Error())
	}
	return v
}
