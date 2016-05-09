package telebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bot-api/telegram"
	"github.com/jarcoal/httpmock"
	"golang.org/x/net/context"
	"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/require"
)

func TestNewWithApi(t *testing.T) {
	api := &telegram.API{}
	b := NewWithAPI(api)
	assert.Equal(t, api, b.api)
	assert.Equal(t, 0, len(b.middleware))
	assert.NotNil(t, b.errFunc)
	assert.Nil(t, b.handler)
}

func TestBot_Use(t *testing.T) {
	b := New("")
	b.Use(nil, nil)
	assert.Equal(t, 2, len(b.middleware))
}

func TestBot_Handle(t *testing.T) {
	b := New("")
	h := HandlerFunc(func(context.Context) error { return nil })
	b.Handle(h)
	assert.NotNil(t, b.handler)
	assert.Equal(t, fmt.Sprintf("%#v", h), fmt.Sprintf("%#v", b.handler))
}

func TestBot_HandleFunc(t *testing.T) {
	b := New("")
	h := HandlerFunc(func(context.Context) error { return nil })
	b.HandleFunc(h)
	assert.NotNil(t, b.handler)
	assert.Equal(t, fmt.Sprintf("%#v", h), fmt.Sprintf("%#v", b.handler))
}

func TestBot_ErrorFunc(t *testing.T) {
	b := New("")
	h := ErrorFunc(func(context.Context, error) {})
	b.ErrorFunc(h)
	assert.NotNil(t, b.errFunc)
	assert.Equal(t, fmt.Sprintf("%#v", h), fmt.Sprintf("%#v", b.errFunc))
}

func TestBot_handleUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	expErr := fmt.Errorf("expected error")

	api := telegram.New("token")
	b := NewWithAPI(api)
	u := &telegram.Update{
		UpdateID: 10,
		Message: &telegram.Message{
			Text: "message",
		},
	}
	handlerInvoked := false
	errFuncInvoked := false
	middlewareInvoked := false
	b.ErrorFunc(func(_ context.Context, err error) {
		errFuncInvoked = true
		assert.Equal(t, expErr, err)
	})
	b.Use(func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context) error {
			middlewareInvoked = true
			return next.Handle(ctx)
		})
	})
	b.HandleFunc(func(ctx context.Context) error {
		handlerInvoked = true
		assert.Equal(t, api, GetAPI(ctx))
		assert.Equal(t, u, GetUpdate(ctx))
		assert.Equal(t, u.UpdateID, ctx.Value("update.id"))
		return expErr
	})
	b.handleUpdate(ctx, u)

	assert.True(t, handlerInvoked, "handler wasn't invoked")
	assert.True(t, errFuncInvoked, "errFunc wasn't invoked")
	assert.True(t, middlewareInvoked, "middleware wasn't invoked")

	// test EmptyHandler
	b = NewWithAPI(api)
	b.Use(func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context) error {
			assert.Equal(t,
				fmt.Sprintf("%#v", EmptyHandler()),
				fmt.Sprintf("%#v", next))
			return next.Handle(ctx)
		})
	})
	b.handleUpdate(ctx, u)

}

func TestBot_getWebhookHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	api := telegram.New("token")
	b := NewWithAPI(api)
	expUpd := telegram.Update{
		UpdateID: 10,
		Message: &telegram.Message{
			Text: "message",
		},
	}
	ch := make(chan telegram.Update, 1)
	whHandler := b.getWebhookHandler(ctx, ch)

	{
		w := httptest.NewRecorder()
		// prepare request
		buf := &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(expUpd)
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "", buf)
		require.NoError(t, err)

		whHandler(w, req)

		select {
		case <-ctx.Done():
			require.NoError(t, ctx.Err())
		case upd := <-ch:
			assert.EqualValues(t, expUpd, upd)
		}
	}

	// TODO (m0sth8): what to do with errors during webhook handling??
	// right now it's just logging
	//{
	//	w := httptest.NewRecorder()
	//	// prepare request
	//	req, err := http.NewRequest("POST", "",
	//		bytes.NewBufferString("bad json"))
	//	require.NoError(t, err)
	//
	//	whHandler(w, req)
	//
	//	select {
	//	case <- ctx.Done():
	//		require.NoError(t, ctx.Err())
	//	case upd := <- ch:
	//		assert.EqualValues(t, expUpd, upd)
	//	}
	//}

}

func NewAPIResponder(status int, result interface{}) httpmock.Responder {
	data, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	raw := json.RawMessage(data)
	apiResponse := telegram.APIResponse{
		Ok:     true,
		Result: &raw,
	}
	responder, err := httpmock.NewJsonResponder(status, apiResponse)
	if err != nil {
		panic(err)
	}
	return responder
}

func TestBot_updateMe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	client := &http.Client{}
	api := telegram.NewWithClient("_token", client)
	b := NewWithAPI(api)
	expMe := telegram.User{
		ID:       10,
		Username: "test_bot",
	}

	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/bot_token/getMe",
		NewAPIResponder(200, expMe),
	)
	err := b.updateMe(ctx)
	require.NoError(t, err)
	assert.Equal(t, expMe, *b.me)
}

func TestBot_Serve(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	client := &http.Client{}
	api := telegram.NewWithClient("_token", client)

	expUpd := []telegram.Update{
		{
			UpdateID: 10,
			Message: &telegram.Message{
				Text: "message",
			},
		},
		{
			UpdateID: 11,
			Message: &telegram.Message{
				Text: "message",
			},
		},
	}
	expMe := telegram.User{
		ID:       10,
		Username: "test_bot",
	}
	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/bot_token/getMe",
		NewAPIResponder(200, expMe),
	)
	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/bot_token/getUpdates",
		NewAPIResponder(200, expUpd),
	)

	b := NewWithAPI(api)
	handleCh := make(chan context.Context, 1)
	b.HandleFunc(func(ctx context.Context) error {
		select {
		case handleCh <- ctx:
		case <-ctx.Done():
		}
		return nil
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- b.Serve(ctx)
	}()

	// got first update
	var update1 *telegram.Update
	select {
	case err := <-errCh:
		require.NoError(t, err)
	case handleCtx1 := <-handleCh:
		update1 = GetUpdate(handleCtx1)
	}

	assert.Equal(t, expMe, *b.me)
	assert.Equal(t, expUpd[0], *update1)

	// got second update
	var update2 *telegram.Update
	select {
	case err := <-errCh:
		require.NoError(t, err)
	case handleCtx1 := <-handleCh:
		update2 = GetUpdate(handleCtx1)
	}

	assert.Equal(t, expUpd[1], *update2)

	cancel()
	select {
	case err := <-errCh:
		require.Equal(t, err, context.Canceled)
	}
}

func TestBot_ServeByWebhook(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	api := telegram.New("_token")

	expUpd := []telegram.Update{
		{
			UpdateID: 10,
			Message: &telegram.Message{
				Text: "message",
			},
		},
		{
			UpdateID: 11,
			Message: &telegram.Message{
				Text: "message",
			},
		},
	}
	expMe := telegram.User{
		ID:       10,
		Username: "test_bot",
	}

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/bot_token/getMe",
		NewAPIResponder(200, expMe),
	)

	b := NewWithAPI(api)

	handleCh := make(chan context.Context, 1)
	b.HandleFunc(func(ctx context.Context) error {
		select {
		case handleCh <- ctx:
		case <-ctx.Done():
		}
		return nil
	})

	whHandler, err := b.ServeByWebhook(ctx)
	require.NoError(t, err)

	{
		w := httptest.NewRecorder()
		// prepare request
		buf := &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(expUpd[0])
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "", buf)
		require.NoError(t, err)

		go whHandler(w, req)

	}

	// got first update
	var update1 *telegram.Update
	select {
	case handleCtx1 := <-handleCh:
		update1 = GetUpdate(handleCtx1)
	case <-ctx.Done():
		require.NoError(t, ctx.Err())
	}

	assert.Equal(t, expMe, *b.me)
	assert.Equal(t, expUpd[0], *update1)

	{
		w := httptest.NewRecorder()
		// prepare request
		buf := &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(expUpd[1])
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "", buf)
		require.NoError(t, err)

		go whHandler(w, req)

	}

	// got second update
	var update2 *telegram.Update
	select {
	case handleCtx1 := <-handleCh:
		update2 = GetUpdate(handleCtx1)
	case <-ctx.Done():
		require.NoError(t, ctx.Err())
	}

	assert.Equal(t, expUpd[1], *update2)

}
