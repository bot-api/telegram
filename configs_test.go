// telegram_test package tests only public interface
package telegram_test

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"

	"github.com/bot-api/telegram"
	"gopkg.in/stretchr/testify.v1/assert"
)

func TestMeCfg(t *testing.T) {
	name := "getMe"
	c := telegram.MeCfg{}
	assert.Equal(t, name, c.Name(), "method Name() has wrong value")
	values, err := c.Values()
	assert.Nil(t, values)
	assert.NoError(t, err)
}

func TestUpdateCfg_Name(t *testing.T) {
	name := "getUpdates"
	c := telegram.UpdateCfg{}
	assert.Equal(t, name, c.Name(), "method Name() has wrong value")
}

func TestUpdateCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"offset":  {"100"},
				"limit":   {"10"},
				"timeout": {"30"},
			},
			cfg: telegram.UpdateCfg{
				Offset:  100,
				Limit:   10,
				Timeout: 30,
			},
		},
		{
			exp: url.Values{},
			cfg: telegram.UpdateCfg{},
		},
		{
			exp: nil,
			expErr: telegram.NewValidationError(
				"Limit",
				"should be between 1 and 100",
			),
			cfg: telegram.UpdateCfg{
				Limit: -10,
			},
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestChatAction_Name(t *testing.T) {
	name := "sendChatAction"
	c := telegram.ChatActionCfg{}
	if c.Name() != name {
		t.Errorf("Expected Name() to be %s, actual %s", name, c.Name())
	}
}

func TestChatActionCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id": {"10"},
				"action":  {"typing"},
			},
			cfg: telegram.ChatActionCfg{
				BaseChat: telegram.BaseChat{ID: 10},
				Action:   telegram.ActionTyping,
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"Action",
			),
			cfg: telegram.ChatActionCfg{
				BaseChat: telegram.BaseChat{ID: 10},
			},
		},
		{
			exp: nil,
			cfg: telegram.ChatActionCfg{},
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestUserProfilePhotosCfg_Name(t *testing.T) {
	name := "getUserProfilePhotos"
	c := telegram.UserProfilePhotosCfg{}
	if c.Name() != name {
		t.Errorf("Expected Name() to be %s, actual %s", name, c.Name())
	}
}

func TestUserProfilePhotosCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"user_id": {"10"},
			},
			cfg: telegram.UserProfilePhotosCfg{
				UserID: 10,
			},
		},
		{
			exp: url.Values{
				"user_id": {"10"},
				"offset":  {"100"},
				"limit":   {"5"},
			},
			cfg: telegram.UserProfilePhotosCfg{
				UserID: 10,
				Offset: 100,
				Limit:  5,
			},
		},
		{
			expErr: telegram.NewValidationError(
				"Limit",
				"should be between 1 and 100",
			),
			cfg: telegram.UserProfilePhotosCfg{
				UserID: 10,
				Limit:  1000,
			},
		},
		{
			cfg: telegram.UserProfilePhotosCfg{},
			expErr: telegram.NewRequiredError(
				"UserID",
			),
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestAnswerInlineQueryCfg_Name(t *testing.T) {
	name := "answerInlineQuery"
	c := telegram.AnswerInlineQueryCfg{}
	if c.Name() != name {
		t.Errorf("Expected Name() to be %s, actual %s", name, c.Name())
	}
}

func TestAnswerInlineQueryCfg_Values(t *testing.T) {
	results := []telegram.InlineQueryResult{
		&telegram.InlineQueryResultArticle{
			BaseInlineQueryResult: telegram.BaseInlineQueryResult{
				ID:   "id",
				Type: "type",
			},
			Title: "title",
		},
	}
	resultsEncoded := "[{\"type\":\"type\",\"id\":\"id\",\"title\":\"title\"}]"
	testTable := []cfgTT{
		{
			exp: url.Values{
				"results":         {resultsEncoded},
				"inline_query_id": {"10"},
			},
			cfg: telegram.AnswerInlineQueryCfg{
				InlineQueryID: "10",
				Results:       results,
			},
		},
		{
			exp: url.Values{
				"results":             {resultsEncoded},
				"inline_query_id":     {"10"},
				"cache_time":          {"60"},
				"is_personal":         {"true"},
				"next_offset":         {"offset"},
				"switch_pm_text":      {"switch_pm_text"},
				"switch_pm_parameter": {"switch_pm_parameter"},
			},
			cfg: telegram.AnswerInlineQueryCfg{
				InlineQueryID:     "10",
				Results:           results,
				CacheTime:         60,
				IsPersonal:        true,
				NextOffset:        "offset",
				SwitchPMText:      "switch_pm_text",
				SwitchPMParameter: "switch_pm_parameter",
			},
		},
		{
			cfg: telegram.AnswerInlineQueryCfg{
				InlineQueryID: "10",
			},
			expErr: telegram.NewRequiredError(
				"Results",
			),
		},
		{
			cfg: telegram.AnswerInlineQueryCfg{
				InlineQueryID: "10",
				Results:       []telegram.InlineQueryResult{},
			},
			expErr: telegram.NewRequiredError(
				"Results",
			),
		},
		{
			cfg: telegram.AnswerInlineQueryCfg{
				InlineQueryID: "10",
				Results: []telegram.InlineQueryResult{
					badInlineQueryResult{},
				},
			},
			expErr: &json.MarshalerError{
				Type: reflect.TypeOf(badInlineQueryResult{}),
				Err:  marshalError,
			},
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestGetChat_Name(t *testing.T) {
	name := "getChat"
	c := telegram.GetChatCfg{}
	if c.Name() != name {
		t.Errorf("Expected Name() to be %s, actual %s", name, c.Name())
	}
}

func TestGetChatCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id": {"10"},
			},
			cfg: telegram.GetChatCfg{
				BaseChat: telegram.BaseChat{ID: 10},
			},
		},
		{
			exp: nil,
			cfg: telegram.GetChatCfg{},
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestGetChatAdministrators_Name(t *testing.T) {
	name := "getChatAdministrators"
	c := telegram.GetChatAdministratorsCfg{}
	if c.Name() != name {
		t.Errorf("Expected Name() to be %s, actual %s", name, c.Name())
	}
}

func TestGetChatAdministratorsCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id": {"10"},
			},
			cfg: telegram.GetChatAdministratorsCfg{
				BaseChat: telegram.BaseChat{ID: 10},
			},
		},
		{
			exp: nil,
			cfg: telegram.GetChatAdministratorsCfg{},
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestGetChatMembersCount_Name(t *testing.T) {
	name := "getChatMembersCount"
	c := telegram.GetChatMembersCountCfg{}
	if c.Name() != name {
		t.Errorf("Expected Name() to be %s, actual %s", name, c.Name())
	}
}

func TestGetChatMembersCountCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id": {"10"},
			},
			cfg: telegram.GetChatMembersCountCfg{
				BaseChat: telegram.BaseChat{ID: 10},
			},
		},
		{
			exp: nil,
			cfg: telegram.GetChatMembersCountCfg{},
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestGetChatMember_Name(t *testing.T) {
	name := "getChatMember"
	c := telegram.GetChatMemberCfg{}
	if c.Name() != name {
		t.Errorf("Expected Name() to be %s, actual %s", name, c.Name())
	}
}

func TestGetChatMemberCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id": {"10"},
				"user_id": {"11"},
			},
			cfg: telegram.GetChatMemberCfg{
				BaseChat: telegram.BaseChat{ID: 10},
				UserID:   11,
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"UserID",
			),
			cfg: telegram.GetChatMemberCfg{
				BaseChat: telegram.BaseChat{ID: 10},
			},
		},
		{
			exp: nil,
			cfg: telegram.GetChatMemberCfg{},
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestKickChatMember_Name(t *testing.T) {
	name := "kickChatMember"
	c := telegram.KickChatMemberCfg{}
	if c.Name() != name {
		t.Errorf("Expected Name() to be %s, actual %s", name, c.Name())
	}
}

func TestKickChatMemberCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id": {"10"},
				"user_id": {"11"},
			},
			cfg: telegram.KickChatMemberCfg{
				BaseChat: telegram.BaseChat{ID: 10},
				UserID:   11,
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"UserID",
			),
			cfg: telegram.KickChatMemberCfg{
				BaseChat: telegram.BaseChat{ID: 10},
			},
		},
		{
			exp: nil,
			cfg: telegram.KickChatMemberCfg{},
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestUnbanChatMember_Name(t *testing.T) {
	name := "unbanChatMember"
	c := telegram.UnbanChatMemberCfg{}
	if c.Name() != name {
		t.Errorf("Expected Name() to be %s, actual %s", name, c.Name())
	}
}

func TestUnbanChatMemberCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id": {"10"},
				"user_id": {"11"},
			},
			cfg: telegram.UnbanChatMemberCfg{
				BaseChat: telegram.BaseChat{ID: 10},
				UserID:   11,
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"UserID",
			),
			cfg: telegram.UnbanChatMemberCfg{
				BaseChat: telegram.BaseChat{ID: 10},
			},
		},
		{
			exp: nil,
			cfg: telegram.UnbanChatMemberCfg{},
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestLeaveChat_Name(t *testing.T) {
	name := "leaveChat"
	c := telegram.LeaveChatCfg{}
	if c.Name() != name {
		t.Errorf("Expected Name() to be %s, actual %s", name, c.Name())
	}
}

func TestLeaveChatCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id": {"10"},
			},
			cfg: telegram.LeaveChatCfg{
				BaseChat: telegram.BaseChat{ID: 10},
			},
		},
		{
			exp: nil,
			cfg: telegram.LeaveChatCfg{},
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}
