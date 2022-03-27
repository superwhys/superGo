package superHttp

import (
	"context"
	"fmt"
	"github.com/superwhys/superGo/superHttp/httpClient"
	"github.com/superwhys/superGo/superHttp/httpRequests"
)

func main() {
	ctx := context.Background()

	var oSlice []httpRequests.OptionHttpRequestsFunc
	oSlice = append(oSlice, httpRequests.AddParams("name", "why"))
	oSlice = append(oSlice, httpRequests.AddUserAgent(UserAgent))

	req := httpRequests.InitRequests("GET", "http://127.0.0.1:8000/get_test/1234", ctx, nil, oSlice...)

	readResp, err := httpRequests.ReadResp(httpRequests.SuperRequests(httpClient.Client, req))
	if err != nil {
		return
	}

	fmt.Println(string(readResp))

	readResp = Get(ctx, "http://127.0.0.1:8000/get_test/1234",
		httpRequests.AddParams("name", "why"))
	fmt.Println(string(readResp))
}
