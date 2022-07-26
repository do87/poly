package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// Client is the HTTP client service for poly mesh API
type Client struct {
	client  *http.Client
	baseURL string
	token   string
}

// New creates a new http client wrapper
func New(baseURL string) *Client {
	return &Client{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

// SetToken sets an authorization bearer token
func (c *Client) SetToken(bearer string) {
	c.token = bearer
}

// Do prepares the body and runs the request
func (c *Client) Do(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	b, err := c.prepareBody(body)
	if err != nil {
		return nil, err
	}
	req, err := c.getRequest(ctx, method, path, b)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %v", err)
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func (c *Client) prepareBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

func (c *Client) getRequest(ctx context.Context, method, relPath string, body io.Reader) (req *http.Request, err error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return
	}
	u.Path = path.Join(u.Path, relPath)
	req, err = http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return
	}
	if c.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	return
}
