package httpClient

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)


var maxDuration = time.Second * 30
var Client *HttpClient

type OptionHttpClientFunc func(*HttpClient)

type HttpClient struct {
	Client *http.Client
}

func init() {
	Client = InitClient(WithTimeOut(maxDuration))
}

func InitClient(opts ...OptionHttpClientFunc) (hc *HttpClient) {
	hc = &HttpClient{Client: http.DefaultClient}
	for _, opt := range opts {
		opt(hc)
	}
	return
}

func WithTimeOut(duration time.Duration) OptionHttpClientFunc {
	return func(hc *HttpClient) {
		hc.Client.Timeout = time.Duration(duration)
	}
}

func WithProxy(proxyAddr string) OptionHttpClientFunc {
	return func(hc *HttpClient) {
		httpTransport := http.DefaultTransport.(*http.Transport).Clone()
		proxyURL, _ := url.Parse(proxyAddr)
		httpTransport.Proxy = http.ProxyURL(proxyURL)
		httpTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		hc.Client.Transport = httpTransport
	}
}
