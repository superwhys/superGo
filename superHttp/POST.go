package superHttp

import (
	"context"
	"github.com/superwhys/superGo/superHttp/httpClient"
	"github.com/superwhys/superGo/superHttp/httpRequests"
	"io"
)

func POST(ctx context.Context, url string, body io.Reader, opts ...httpRequests.OptionHttpRequestsFunc) []byte {
	var req *httpRequests.HttpRequests
	req = httpRequests.InitRequests("POST", url, ctx, body, opts...)

	readResp, err := httpRequests.ReadResp(httpRequests.SuperRequests(httpClient.Client, req))

	if err != nil {
		return nil
	}
	return readResp
}
