// +build go1.7

package telegram

import (
	"net/http"
	"net/url"

	"golang.org/x/net/context"
)

// these errors are from net/http for go 1.7
const (
	errRequestCanceled     = "net/http: request canceled"
	errRequestCanceledConn = "net/http: request canceled while waiting for connection"
)

func makeRequest(ctx context.Context, client HTTPDoer, req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)
	if httpClient, ok := client.(*http.Client); ok {
		resp, err = httpClient.Do(req.WithContext(ctx))
	} else {
		// TODO: implement cancel logic for non http.Client
		resp, err = client.Do(req)
	}
	if err != nil {
		if urlErr, casted := err.(*url.Error); casted {
			if urlErr.Err == context.Canceled {
				return resp, context.Canceled
			}
			errMsg := urlErr.Err.Error()
			if errMsg == errRequestCanceled ||
				errMsg == errRequestCanceledConn {
				return resp, context.Canceled
			}
		}
	}
	return resp, err
}
