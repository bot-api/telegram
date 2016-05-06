// +build integration

package telegram_test

import (
	"os"
	"testing"
	"time"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bot-api/telegram"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

var (
	apiBotToken string
	botUserID   int64 = 201910478
)

func init() {
	apiBotToken = os.Getenv("API_BOT_TOKEN")
}

func TestI_Api_GetMe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	{
		api := telegram.New(apiBotToken)
		user, err := api.GetMe(ctx)

		require.NoError(t, err)
		assert.Equal(t, &telegram.User{
			ID:        201910478,
			FirstName: "Chatter",
			LastName:  "",
			Username:  "PoboltaemBot",
		}, user)
	}
	{
		// send bad token
		api := telegram.New("3" + apiBotToken[1:])
		api.Debug(true)
		_, err := api.GetMe(ctx)

		require.EqualError(t, err, "unauthorized")
	}
}

func TestI_Api_GetUserProfilePhotos(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	{
		api := telegram.New(apiBotToken)
		photos, err := api.GetUserProfilePhotos(
			ctx,
			telegram.NewUserProfilePhotos(botUserID),
		)

		require.NoError(t, err)
		expected := &telegram.UserProfilePhotos{}
		err = parseTestData(
			"integration_user_profile_photos.json",
			expected)
		require.NoError(t, err)
		assert.Equal(t, expected, photos)
	}
}

func TestI_Api_GetFile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	api := telegram.New(apiBotToken)
	api.Debug(true)

	photos := &telegram.UserProfilePhotos{}
	err := parseTestData(
		"integration_user_profile_photos.json",
		photos)
	require.NoError(t, err)

	for _, photo := range photos.Photos[0] {
		file, err := api.GetFile(ctx, telegram.FileCfg{
			FileID: photo.FileID,
		})
		require.NoError(t, err)
		assert.Equal(t, photo.FileSize, file.FileSize)
		assert.Equal(t, photo.FileID, file.FileID)
		assert.NotEmpty(t, file.FilePath)
		assert.NotEmpty(t, file.Link)
		resp, err := ctxhttp.Get(ctx, http.DefaultClient, file.Link)
		require.NoError(t, err)
		defer resp.Body.Close()
		actualData, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Len(t, actualData, file.FileSize)
		expectedData, err := ioutil.ReadFile(
			fmt.Sprintf("./testdata/files/%s.jpg", file.FileID))
		require.NoError(t, err)
		assert.Equal(t, expectedData, actualData)
		//ioutil.WriteFile(
		//	fmt.Sprintf("./testdata/files/%s.jpg", file.FileID),
		//	fileData, 0666,
		//)
	}

}

func parseTestData(filename string, dst interface{}) error {
	f, err := os.Open(fmt.Sprintf("./testdata/%s", filename))
	if err != nil {
		return err
	}
	return json.NewDecoder(f).Decode(dst)
}
