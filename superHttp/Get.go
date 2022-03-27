package superHttp

import (
	"context"
	"github.com/superwhys/superGo/superHttp/httpClient"
	"github.com/superwhys/superGo/superHttp/httpRequests"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3100.0 Safari/537.36"

func Get(ctx context.Context, url string, opts ...httpRequests.OptionHttpRequestsFunc) []byte {
	var req *httpRequests.HttpRequests

	req = httpRequests.InitRequests("GET", url, ctx, nil, opts...)

	readResp, err := httpRequests.ReadResp(httpRequests.SuperRequests(httpClient.Client, req))
	if err != nil {
		return nil
	}
	return readResp
}
