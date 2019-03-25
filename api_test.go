// telegram_test package tests only public interface
package telegram_test

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/bot-api/telegram"
	"github.com/jarcoal/httpmock"
	"golang.org/x/net/context"
	"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/require"
)

var apiToken = "token"

var forbiddenResponder = httpmock.NewStringResponder(
	http.StatusForbidden, "forbidden")

func TestAPI_GetMe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	api := telegram.New(apiToken)

	testTable := []struct {
		resp httpmock.Responder

		expErr  string
		expUser *telegram.User
		expURL  string
	}{
		{
			resp: httpmock.NewStringResponder(200, `
			{
			    "ok": true,
			    "result": {
				"id": 100,
				"first_name": "first_name",
				"last_name": "last_name",
				"username": "username"
			    }
			}`),
			expUser: &telegram.User{
				ID:        100,
				FirstName: "first_name",
				LastName:  "last_name",
				Username:  "username",
			},
			expURL: "/bottoken/getMe",
		},
		{
			resp:   forbiddenResponder,
			expErr: "forbidden",
			expURL: "/bottoken/getMe",
		},
	}
	for i, tt := range testTable {
		t.Logf("Experiment %d", i)
		httpmock.RegisterResponder(
			"POST",
			"https://api.telegram.org"+tt.expURL,
			tt.resp,
		)

		user, err := api.GetMe(ctx)
		if tt.expErr != "" {
			assert.EqualError(t, err, tt.expErr)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tt.expUser, user)

		httpmock.Reset()

	}
}

func TestAPI_GetUpdates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	api := telegram.New(apiToken)

	testTable := []struct {
		cfg  telegram.UpdateCfg
		resp httpmock.Responder

		expErr    string
		expResult []telegram.Update
		expURL    string
	}{
		{
			cfg: telegram.UpdateCfg{
				Offset:  10,
				Limit:   100,
				Timeout: 1000,
			},
			resp: httpmock.NewStringResponder(200, `
			{
			    "ok": true,
			    "result": [
			    {
			    	"update_id": 100,
			    	"message": {
			    		"message_id": 135
			    	}
			    }
			    ]
			}`),
			expResult: []telegram.Update{
				{
					UpdateID: 100,
					Message: &telegram.Message{
						MessageID: 135,
					},
				},
			},
			expURL: "/bottoken/getUpdates",
		},
		{
			cfg: telegram.UpdateCfg{
				Offset:  10,
				Limit:   100,
				Timeout: 1000,
			},
			resp: httpmock.NewStringResponder(200, `
			{
			    "ok": true,
			    "result": []
			}`),
			expResult: []telegram.Update{},
			expURL:    "/bottoken/getUpdates",
		},
		{
			cfg:    telegram.UpdateCfg{},
			resp:   forbiddenResponder,
			expErr: "forbidden",
			expURL: "/bottoken/getUpdates",
		},
	}
	for i, tt := range testTable {
		t.Logf("Experiment %d", i)

		httpmock.RegisterResponder(
			"POST",
			"https://api.telegram.org"+tt.expURL,
			tt.resp,
		)

		result, err := api.GetUpdates(ctx, tt.cfg)
		if tt.expErr != "" {
			assert.EqualError(t, err, tt.expErr)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, result, tt.expResult)

		httpmock.Reset()

	}
}

func TestAPI_GetUserProfilePhotos(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	api := telegram.New(apiToken)

	testTable := []struct {
		cfg  telegram.UserProfilePhotosCfg
		resp httpmock.Responder

		expErr    string
		expResult *telegram.UserProfilePhotos
		expURL    string
	}{
		{
			cfg: telegram.UserProfilePhotosCfg{
				UserID: 100,
			},
			resp: httpmock.NewStringResponder(200, `
			{
			    "ok": true,
			    "result":
			    {
			    	"total_count": 2,
			    	"photos": [
			    		[
					  {
					    "file_id": "value1",
					    "file_size": 345,
					    "width": 100,
					    "height": 200
					  }
			    		],
			    		[
					  {
					    "file_id": "value2",
					    "file_size": 500,
					    "width": 200,
					    "height": 100
					  }
			    		]
			    	]
			    }
			}`),
			expResult: &telegram.UserProfilePhotos{
				TotalCount: 2,
				Photos: [][]telegram.PhotoSize{
					{
						{
							MetaFile: telegram.MetaFile{
								FileID:   "value1",
								FileSize: 345,
							},
							Size: telegram.Size{
								Width:  100,
								Height: 200,
							},
						},
					},
					{
						{
							MetaFile: telegram.MetaFile{
								FileID:   "value2",
								FileSize: 500,
							},
							Size: telegram.Size{
								Width:  200,
								Height: 100,
							},
						},
					},
				},
			},
			expURL: "/bottoken/getUserProfilePhotos",
		},
		{
			cfg:    telegram.UserProfilePhotosCfg{},
			resp:   nil,
			expErr: "UserID required",
			expURL: "/bottoken/getUserProfilePhotos",
		},
		{
			cfg: telegram.UserProfilePhotosCfg{
				UserID: 10,
				Limit:  1000,
			},
			resp:   nil,
			expErr: "field Limit is invalid: should be between 1 and 100",
			expURL: "/bottoken/getUserProfilePhotos",
		},
		{
			cfg: telegram.UserProfilePhotosCfg{
				UserID: 10,
			},
			resp:   forbiddenResponder,
			expErr: "forbidden",
			expURL: "/bottoken/getUserProfilePhotos",
		},
	}
	for i, tt := range testTable {
		t.Logf("Experiment %d", i)

		httpmock.RegisterResponder(
			"POST",
			"https://api.telegram.org"+tt.expURL,
			tt.resp,
		)

		result, err := api.GetUserProfilePhotos(ctx, tt.cfg)
		if tt.expErr != "" {
			assert.EqualError(t, err, tt.expErr)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, result, tt.expResult)

		httpmock.Reset()

	}
}

func TestAPI_DownloadFile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	api := telegram.New(apiToken)

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/bottoken/getFile",
		httpmock.NewStringResponder(200, `
			{
			    "ok": true,
			    "result":
			    {
			    	"file_path": "path_to_file",
                    		"file_id": "file_id"
			    }
			}`,
		),
	)
	httpmock.RegisterResponder(
		"GET",
		"https://api.telegram.org/file/bottoken/path_to_file",
		httpmock.NewStringResponder(200, "FILE DATA"),
	)
	f, err := api.GetFile(ctx, telegram.FileCfg{FileID: "file_id"})
	require.NoError(t, err)
	assert.Equal(t, "file_id", f.FileID)
	assert.Equal(t, "https://api.telegram.org/file/bottoken/path_to_file", f.Link)

	buf := bytes.NewBuffer(nil)
	err = api.DownloadFile(ctx, telegram.FileCfg{FileID: "file_id"}, buf)
	require.NoError(t, err)
	assert.Equal(t, "FILE DATA", buf.String())

}

func TestSplittingBigMessages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	var buffer bytes.Buffer
	for i := 0; i < telegram.MaximumMessageLength+1; i++ {
		buffer.WriteString("a")
	}
	message := buffer.String()
	msg := telegram.NewMessage(123, message)
	defer cancel()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	api := telegram.New(apiToken)

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/bottoken/sendMessage",
		httpmock.NewStringResponder(200, `
		{
			"result": {
				"message": {
				  "text": "sdfsdfs",
				  "date": 1552220293,
				  "chat": {
					"type": "private",
					"username": "YappoFTW",
					"first_name": "Alexander",
					"id": 53853963
				  },
				  "from": {
					"language_code": "en",
					"username": "YappoFTW",
					"first_name": "Alexander",
					"is_bot": false,
					"id": 53853963
				  },
				  "message_id": 2
				},
				"update_id": 576705363
			},
			"ok": true
		}`,
		),
	)

	_, err := api.SplitAndSendMessage(ctx, msg)
	require.NoError(t, err)
	assert.Equal(t, 2, httpmock.GetTotalCallCount())
}
