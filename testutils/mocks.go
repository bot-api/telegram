package testutils

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/bot-api/telegram"
	"github.com/m0sth8/httpmock"
	"golang.org/x/net/context"
)

// UpdatesResponder provides a useful responder that returns updates from ch
func UpdatesResponder(ctx context.Context, ch <-chan []telegram.Update) httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case update := <-ch:
			data, err := json.Marshal(update)
			if err != nil {
				return httpmock.NewJsonResponse(200,
					telegram.APIResponse{
						Ok:          false,
						Description: err.Error(),
					},
				)
			}
			raw := json.RawMessage(data)
			return httpmock.NewJsonResponse(200,
				telegram.APIResponse{
					Ok:     true,
					Result: &raw,
				},
			)
		}
	}
}

// ReceivedMessage contains information about message that are sent by api
type ReceivedMessage struct {
	Values url.Values
	URL    url.URL
}

// SendResponder helps to test telegram api.
// Returns received message to out channel,
// Waits for result channel with object to return in API.
// If result channel is nil, then returns empty telegram.Message.
// If result channel has error, then returns APIResponse with Ok False and Description.
func SendResponder(ctx context.Context,
	result <-chan interface{}) (httpmock.Responder, <-chan ReceivedMessage) {
	out := make(chan ReceivedMessage)
	return func(req *http.Request) (*http.Response, error) {
		req.ParseForm()
		select {
		case out <- ReceivedMessage{
			Values: req.Form,
			URL:    *req.URL,
		}:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		select {
		case data := <-result:
			if err, casted := data.(error); casted {
				return httpmock.NewJsonResponse(200,
					telegram.APIResponse{
						Ok:          false,
						Description: err.Error(),
					},
				)
			}
			rawData, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			raw := json.RawMessage(rawData)
			return httpmock.NewJsonResponse(200,
				telegram.APIResponse{
					Ok:     true,
					Result: &raw,
				},
			)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}, out
}
