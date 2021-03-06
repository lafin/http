// Package http handle work with http
package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"
)

var once sync.Once
var client *http.Client

// Get - wrapper to execute http GET request
func Get(url string, headers map[string]string) ([]byte, error) {
	client = Client()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%d", res.StatusCode)
	}
	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Post - wrapper to execute http POST request
func Post(url string, body io.Reader, headers map[string]string) ([]byte, error) {
	client = Client()
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%d", res.StatusCode)
	}
	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Client - get instance of http client
func Client() *http.Client {
	once.Do(func() {
		transport := &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Second * 10,
			}).Dial,
			MaxIdleConns:        10,
			IdleConnTimeout:     10 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
			DisableKeepAlives:   true,
		}
		cookieJar, err := cookiejar.New(nil)
		if err != nil {
			panic(err)
		}
		client = &http.Client{
			Timeout:   300 * time.Second,
			Transport: transport,
			Jar:       cookieJar,
		}
	})

	return client
}
