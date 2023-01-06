package httphandler

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type HTTPClient struct {
	Url    string
	Client *http.Client
}

func NewClient(Url string, Client *http.Client) *HTTPClient {
	return &HTTPClient{
		Url:    Url,
		Client: Client,
	}
}

func InitHTTPClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
	}
	return client
}
