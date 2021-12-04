package bapi

import (
	"fmt"
	"github.com/1makarov/binance-nft-buy/internal/domain/account"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type Api struct {
	request *fasthttp.Request
	http    *fasthttp.Client
}

func New(setting account.Setting) (*Api, error) {
	http := initHttp(setting.Proxy)
	request, err := initHeaders(setting.BAuth)
	if err != nil {
		return nil, err
	}
	return &Api{http: http, request: request}, nil
}

func initHttp(proxy string) *fasthttp.Client {
	c := &fasthttp.Client{}
	if proxy != "" {
		c.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
		return c
	}
	return c
}

func initHeaders(bAuth *account.BAuth) (*fasthttp.Request, error) {
	if bAuth.Cookie == "" || bAuth.Csrf == "" {
		return nil, fmt.Errorf("empty field .env")
	}
	r := &fasthttp.Request{}
	r.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36")
	r.Header.Set("clienttype", "web")
	r.Header.Set("cookie", bAuth.Cookie)
	r.Header.Set("csrftoken", bAuth.Csrf)
	r.Header.Set("device-info","eyJzY3JlZW5fcmVzb2x1dGlvbiI6IjE5MjAsMTA4MCIsImF2YWlsYWJsZV9zY3JlZW5fcmVzb2x1dGlvbiI6IjE5MjAsMTA0MCIsInN5c3RlbV92ZXJzaW9uIjoiV2luZG93cyAxMCIsImJyYW5kX21vZGVsIjoidW5rbm93biIsInN5c3RlbV9sYW5nIjoibmwtTkwiLCJ0aW1lem9uZSI6IkdNVCsxIiwidGltZXpvbmVPZmZzZXQiOi02MCwidXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQpIEFwcGxlV2ViS2l0LzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZS85My4wLjQ1NzcuODIgU2FmYXJpLzUzNy4zNiIsImxpc3RfcGx1Z2luIjoiQ2hyb21lIFBERiBQbHVnaW4sQ2hyb21lIFBERiBWaWV3ZXIsTmF0aXZlIENsaWVudCIsImNhbnZhc19jb2RlIjoiY2MzNDBjNDkiLCJ3ZWJnbF92ZW5kb3IiOiJHb29nbGUgSW5jLiIsIndlYmdsX3JlbmRlcmVyIjoiQU5HTEUgKE5WSURJQSBHZUZvcmNlIEdUWCAxMDYwIDZHQiBEaXJlY3QzRDExIHZzXzVfMCBwc181XzApXHIiLCJhdWRpbyI6IjEyNC4wNDM0NzQ2ODY1NjQ4OCIsInBsYXRmb3JtIjoiV2luMzIiLCJ3ZWJfdGltZXpvbmUiOiJFdXJvcGUvQW1zdGVyZGFtIiwiZGV2aWNlX25hbWUiOiJDaHJvbWUgVjkzLjAuNDU3Ny44MiAoV2luZG93cykiLCJmaW5nZXJwcmludCI6IjQ2MDFjZDUxYzkyZTZlNTZkYzc3ODIyNTBjYTVkYmFkIiwiZGV2aWNlX2lkIjoiIiwicmVsYXRlZF9kZXZpY2VfaWRzIjoiMTYzODU0MDQ2NTk5NG9vMWFrMHdxUThmRmRsdzlWeTYifQ==")
	return r, nil
}
