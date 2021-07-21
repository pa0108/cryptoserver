package crypto

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client ... struct
type Client struct {
	apiKey      string
	apiSecret   string
	httpClient  *http.Client
	httpTimeout time.Duration
}

// NewClient ... return a new HTTP client
func NewClient(apiKey, apiSecret string) (c *Client) {
	return &Client{apiKey, apiSecret, &http.Client{}, 10 * time.Second}
}

// HttpRequest ... do a HTTP request with timeout
func (c *Client) HttpRequest(req *http.Request) (*http.Response, error) {
	timer := time.NewTimer(c.httpTimeout)
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout error")
	}
}

// do ... http request to crypto server
func (c *Client) do(method string, resource string, payload map[string]string, authNeeded bool) (response []byte, err error) {
	var requesturl string
	if strings.HasPrefix(resource, "http") {
		requesturl = resource
	} else {
		requesturl = fmt.Sprintf("%s/%s", API_BASE, resource)
	}
	var formData string
	if method == "GET" {
		var URL *url.URL
		URL, err = url.Parse(requesturl)
		if err != nil {
			return
		}
		q := URL.Query()
		for key, value := range payload {
			q.Set(key, value)
		}
		formData = q.Encode()
		URL.RawQuery = formData
		requesturl = URL.String()
	} else {
		formValues := url.Values{}
		for key, value := range payload {
			formValues.Set(key, value)
		}
		formData = formValues.Encode()
	}
	req, err := http.NewRequest(method, requesturl, strings.NewReader(formData))
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")

	// Auth
	if authNeeded {
		if len(c.apiKey) == 0 || len(c.apiSecret) == 0 {
			err = errors.New("Authentication failed")
			return
		}
		req.SetBasicAuth(c.apiKey, c.apiSecret)
	}

	resp, err := c.HttpRequest(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 401 {
		err = errors.New(resp.Status)
	}
	return response, err
}
