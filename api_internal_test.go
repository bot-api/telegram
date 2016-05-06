package telegram

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"mime/multipart"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestNew(t *testing.T) {
	api := New("token123")
	assert.Equal(t, "token123", api.token)
	assert.Equal(t,
		http.DefaultClient,
		api.client,
		"Expected client to be http.DefaultClient")
	assert.Equal(t, APIEndpoint, api.apiEndpoint)
	assert.Equal(t, FileEndpoint, api.fileEndpoint)
}

func TestNewWithClient(t *testing.T) {
	customClient := &http.Client{}
	apiActual := NewWithClient("token123", customClient)
	assert.Equal(t, apiActual.client, customClient)
}

func TestApi_Debug(t *testing.T) {
	api := New("token123")
	assert.False(t, api.debug, "Expected debug to be false by default")
	api.Debug(true)
	assert.True(t, api.debug, "Expected debug to be true")
	api.Debug(false)
	assert.False(t, api.debug, "Expected debug to be false")

}

func TestApi_makeRequest_testContextCancel(t *testing.T) {
	// Use real http.Client for this test to test ctxhttp
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	handlerCh := make(chan bool, 1)
	testHandle := func(res http.ResponseWriter, req *http.Request) {
		// handler received request
		handlerCh <- true
		// handler waits till context is done
		select {
		case <-ctx.Done():
			return
		}
	}
	testServ := httptest.NewServer(http.HandlerFunc(testHandle))
	api := New("token")
	api.apiEndpoint = fmt.Sprintf("%s/bot%%s/%%s", testServ.URL)

	reqCtx, cancelReq := context.WithTimeout(
		context.Background(),
		time.Millisecond*500)
	errCh := make(chan error, 1)
	go func() {
		req, err := http.NewRequest(
			"GET",
			testServ.URL,
			nil,
		)
		if err != nil {
			errCh <- err
			return
		}
		err = api.makeRequest(reqCtx, req, nil)
		errCh <- err
	}()
	// wait till http handler receives request
	<-handlerCh
	cancelReq()
	// receive error from makeRequest
	err := <-errCh
	assert.Equal(t, context.Canceled, err)

}

type fakeClient struct {
	res *http.Response
	err error
	req *http.Request
	ctx context.Context
}

func (c *fakeClient) Do(req *http.Request) (*http.Response, error) {
	c.req = req
	return c.res, c.err
}

func newFakeClient(ctx context.Context, res *http.Response, err error) *fakeClient {
	return &fakeClient{
		res: res,
		err: err,
		ctx: ctx,
	}
}

func TestApi_makeRequest(t *testing.T) {
	testTable := []struct {
		method string
		params url.Values
		dst    interface{}

		resp     string
		respCode int

		expErr  string
		expForm url.Values
		expURL  string
		expDst  interface{}
	}{
		{
			params: url.Values{
				"key1": {"value1"},
			},
			method: "search_method",
			resp:   "{\"ok\": true, \"result\": \"data\"}",
			dst:    stringP(""),
			expDst: stringP("data"),
			expURL: "/bottoken/search_method",
		},
		{
			resp:   "not a json string",
			expErr: "invalid character 'o' in literal null (expecting 'u')",
		},
		{
			resp:   "{\"ok\": false, \"description\": \"someError\"}",
			expErr: "apiError: someError",
		},
		{
			resp:   "{\"ok\": true, \"result\": \"wrong json\"}",
			dst:    intP(0),
			expErr: "json: cannot unmarshal string into Go value of type int",
		},
		{
			respCode: http.StatusForbidden,
			expErr:   "forbidden",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	api := New("token")
	api.debug = true

	for i, tt := range testTable {
		t.Logf("Experiment %d", i)
		var dst interface{}
		if tt.dst != nil {
			dst = tt.dst
		} else {
			dst = nil
		}
		fResBody := newNopCloser(bytes.NewBufferString(tt.resp))
		fRes := &http.Response{
			Body:       fResBody,
			StatusCode: http.StatusOK,
		}
		if tt.respCode != 0 {
			fRes.StatusCode = tt.respCode
		}
		fc := newFakeClient(ctx, fRes, nil)
		api.client = fc
		req, err := api.getFormRequest(tt.method, tt.params)
		require.NoError(t, err)
		err = api.makeRequest(ctx, req, dst)
		// check error from makeRequest
		if tt.expErr != "" {
			assert.EqualError(t, err, tt.expErr)
			continue
		} else {
			if !assert.NoError(t, err) {
				continue
			}
		}

		if tt.expURL != "" {
			assert.Equal(t, tt.expURL, fc.req.URL.Path)
		}
		assert.Equal(t,
			fc.req.Header.Get("Content-Type"),
			"application/x-www-form-urlencoded",
		)
		reqBody, err := ioutil.ReadAll(fc.req.Body)
		if err != nil {
			t.Fatal(err)
		}
		// check request body
		if tt.expForm != nil {
			assert.Equal(t, tt.expForm.Encode(), string(reqBody))
		} else {
			assert.Equal(t, tt.params.Encode(), string(reqBody))
		}

		// check dst
		assert.Equal(t, tt.expDst, tt.dst)

		// body should be closed
		assert.True(t,
			fResBody.closed,
			"response body should be closed")
		// body should be empty
		data := make([]byte, 512)
		n, err := fRes.Body.Read(data)
		if !assert.Equal(t, 0, n, "response body should be empty") {
			t.Logf("Body has next data %s", string(data[:n]))
		}
		assert.Equal(t, io.EOF, err)

	}

}

func TestApi_getFormRequest(t *testing.T) {
	api := New("token")
	api.debug = true
	req, err := api.getFormRequest("method", url.Values{
		"key1": {"value1"},
		"key2": {"value2"},
	})
	require.NoError(t, err)
	assert.Equal(t,
		"https://api.telegram.org/bottoken/method",
		req.URL.String())
	data, err := ioutil.ReadAll(req.Body)
	require.NoError(t, err)
	assert.Equal(t, "key1=value1&key2=value2", string(data))
	assert.Equal(t,
		"application/x-www-form-urlencoded",
		req.Header.Get("Content-Type"))

	// check errors
	api.apiEndpoint = "not a url"
	_, err = api.getFormRequest("method", url.Values{})
	require.Error(t, err)
	assert.EqualError(t,
		err, "parse not a url%!(EXTRA string=token,"+
			" string=method): invalid URL escape \"%!(\"")
}

func TestApi_getUploadRequest(t *testing.T) {
	api := New("token")
	api.debug = true
	buf := bytes.NewBufferString("file content")
	req, err := api.getUploadRequest(
		"method",
		url.Values{
			"key1": {"value1"},
			"key2": {"value2"},
		},
		"field_name",
		NewInputFile("filename", buf),
	)
	require.NoError(t, err)
	assert.Equal(t,
		"https://api.telegram.org/bottoken/method",
		req.URL.String())
	contentType := req.Header.Get("Content-Type")
	require.True(t, strings.HasPrefix(contentType,
		"multipart/form-data; boundary="))
	boundary := contentType[30:]

	r := multipart.NewReader(req.Body, boundary)
	for i := 0; i < 3; i++ {
		part, err := r.NextPart()
		require.NoError(t, err)
		data, err := ioutil.ReadAll(part)
		require.NoError(t, err)
		if part.FormName() == "key1" {
			assert.Equal(t, "value1", string(data))
		}
		if part.FormName() == "key2" {
			assert.Equal(t, "value2", string(data))
		}
		if part.FormName() == "field_name" {
			assert.Equal(t, "file content", string(data))
			assert.Equal(t,
				"application/octet-stream",
				part.Header.Get("Content-Type"),
			)
		}
	}

	// check errors
	api.apiEndpoint = "not a url"
	buf = bytes.NewBufferString("file content")
	req, err = api.getUploadRequest(
		"method",
		url.Values{
			"key1": {"value1"},
			"key2": {"value2"},
		},
		"field_name",
		NewInputFile("filename", buf),
	)
	require.Error(t, err)
	assert.EqualError(t,
		err, "parse not a url%!(EXTRA string=token,"+
			" string=method): invalid URL escape \"%!(\"")
}

// helpers

func stringP(s string) *string {
	return &s
}

func intP(s int) *int {
	return &s
}

type nopCloser struct {
	io.Reader
	closed bool
}

func (c *nopCloser) Close() error {
	c.closed = true
	return nil
}

// NopCloser returns a ReadCloser with a no-op Close method wrapping
// the provided Reader r.
func newNopCloser(r io.Reader) *nopCloser {
	return &nopCloser{r, false}
}
