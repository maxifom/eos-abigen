package client

import (
	"net/http"
	"net/url"
)

type RPCClient struct {
	baseURL *url.URL
	client  *http.Client
}

func NewRPCClient(baseURL string, client *http.Client) (*RPCClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &RPCClient{baseURL: u, client: client}, nil
}

func (c *RPCClient) cloneURL() *url.URL {
	u, _ := url.Parse(c.baseURL.String())
	return u
}
