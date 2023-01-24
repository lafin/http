// Package http handle work with http
package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"
)

var once sync.Once
var client *http.Client

// Get - wrapper to execute http GET request
func Get(url string, headers map[string]string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, nil, err
	}
	return doRequest(req, headers)
}

// Post - wrapper to execute http POST request
func Post(url string, data io.Reader, headers map[string]string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return nil, nil, err
	}
	return doRequest(req, headers)
}

func doRequest(req *http.Request, headers map[string]string) ([]byte, *http.Response, error) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	client = Client(Params{})
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode >= 400 {
		return nil, nil, fmt.Errorf("%d", res.StatusCode)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	return body, res, nil
}

// Params - http client settings
type Params struct {
	MaxIdleConns        int
	IdleConnTimeout     time.Duration
	TLSHandshakeTimeout time.Duration
	DisableKeepAlives   string
	Timeout             time.Duration
}

// Client - get instance of http client
func Client(p Params) *http.Client {
	if p.MaxIdleConns == 0 {
		p.MaxIdleConns = 10
	}
	if p.IdleConnTimeout == 0 {
		p.IdleConnTimeout = 10 * time.Second
	}
	if p.TLSHandshakeTimeout == 0 {
		p.TLSHandshakeTimeout = 5 * time.Second
	}
	if p.DisableKeepAlives != "no" {
		p.DisableKeepAlives = "yes"
	}
	if p.Timeout == 0 {
		p.Timeout = 300 * time.Second
	}
	once.Do(func() {
		transport := &http.Transport{
			MaxIdleConns:        p.MaxIdleConns,
			IdleConnTimeout:     p.IdleConnTimeout,
			TLSHandshakeTimeout: p.TLSHandshakeTimeout,
			DisableKeepAlives:   p.DisableKeepAlives == "yes",
		}
		cookieJar, err := cookiejar.New(nil)
		if err != nil {
			panic(err)
		}
		client = &http.Client{
			Timeout:   p.Timeout,
			Transport: transport,
			Jar:       cookieJar,
		}
	})
	return client
}
