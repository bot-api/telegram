// +build !go1.7

package telegram

import (
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

func makeRequest(ctx context.Context, client HTTPDoer, req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)
	if httpClient, ok := client.(*http.Client); ok {
		resp, err = ctxhttp.Do(ctx, httpClient, req)
	} else {
		// TODO: implement cancel logic for non http.Client
		resp, err = client.Do(req)
	}
	return resp, err
}
