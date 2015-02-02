package oauth

import (
	"net/http"
	"time"
)

type Request struct {
	Token  string
	Secret string
	Time   time.Time
	Nonce  string
	Req    *http.Request
}
