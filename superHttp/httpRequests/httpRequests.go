package httpRequests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/superwhys/superGo/superHttp/httpClient"
	"github.com/superwhys/superGo/superLog"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type OptionHttpRequestsFunc func(requests *HttpRequests)
type HttpRequests struct {
	Requests *http.Request
}

func InitRequests(method, url string, ctx context.Context, rb io.Reader, opts ...OptionHttpRequestsFunc) (hr *HttpRequests) {

	req, _ := http.NewRequestWithContext(ctx, method, url, rb)
	hr = &HttpRequests{Requests: req}
	for _, opt := range opts {
		opt(hr)
	}
	return
}

func AddHeader(key, val string) OptionHttpRequestsFunc {
	return func(hr *HttpRequests) {
		hr.Requests.Header.Add(key, val)
	}
}

func AddUserAgent(userAgent string) OptionHttpRequestsFunc {
	return func(hr *HttpRequests) {
		hr.Requests.Header.Add("User-Agent", userAgent)
	}
}

func AddParams(key, val string) OptionHttpRequestsFunc {
	return func(hr *HttpRequests) {
		data := url.Values{}
		data.Set(key, val)
		if hr.Requests.URL.RawQuery == "" {
			hr.Requests.URL.RawQuery += data.Encode()
		} else {
			hr.Requests.URL.RawQuery += fmt.Sprintf("&%s", data.Encode())
		}
	}
}

func ReadResp(resp *http.Response, err error) ([]byte, error) {
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read response")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("%s[%d]:%s", resp.Status, resp.StatusCode, string(data))
	}
	return data, nil
}

func ReadJson(resp *http.Response, err error, i interface{}) error {
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("%s[%d]", resp.Status, resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(i); err != nil {
		return errors.Wrap(err, "Decode response")
	}

	return nil
}

func SuperRequests(hc *httpClient.HttpClient, hq *HttpRequests) (*http.Response, error) {
	res, err := hc.Client.Do(hq.Requests)
	if err != nil {
		superLog.PanicError(err)
	}
	return res, nil
}
