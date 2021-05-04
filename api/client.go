package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/net/publicsuffix"
)

type Client struct {
	BaseURL     *url.URL
	Credentials Login

	HTTP      *http.Client
	WebSocket *websocket.Dialer

	System    *System
	User      *User
	csrfToken string
}

func New(controllerURL string, login Login) (*Client, error) {
	u, err := url.Parse(controllerURL)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}

	// TODO: add explicit certificate option
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	client := &Client{
		BaseURL:     u,
		Credentials: login,
		HTTP: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
		WebSocket: &websocket.Dialer{
			Jar:             jar,
			TLSClientConfig: tlsConfig,
		},
	}

	return client, nil
}

func (c *Client) Do(method, path string, reqBody interface{}, resBody interface{}) (*http.Response, error) {
	var bodyBuffer io.Reader
	if reqBody != nil {
		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}

		bodyBuffer = bytes.NewBuffer(body)
	}

	// TODO: context
	req, err := http.NewRequest(method, c.URL(path).String(), bodyBuffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.csrfToken != "" {
		req.Header.Set("X-Csrf-Token", c.csrfToken)
	}

	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}

	log.Printf("%s %s -> %s\n", method, path, res.Status)
	if res.StatusCode >= 400 {
		return res, fmt.Errorf("got %s response", res.Status)
	}

	defer res.Body.Close()

	if resBody != nil {
		contentType := res.Header.Get("Content-Type")

		if !strings.HasPrefix(contentType, "application/json") {
			if strings.HasPrefix(contentType, "text/plain") {
				var buf bytes.Buffer
				if _, err := io.Copy(&buf, res.Body); err != nil {
					return res, fmt.Errorf("got %#v response but failed to read it: %w", contentType, err)
				}

				return res, fmt.Errorf("got error response: %s", buf.String())
			}

			return res, fmt.Errorf("got %#v response instead of JSON", contentType)
		}

		dec := json.NewDecoder(res.Body)
		if err := dec.Decode(resBody); err != nil {
			return res, fmt.Errorf("failed to JSON-decode %s: %w", path, err)
		}
	}
	return res, nil
}

func (c *Client) URL(path string) *url.URL {
	pu, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	return c.BaseURL.ResolveReference(pu)
}
