package telegram_test

import (
	"net/url"
	"testing"

	"github.com/bot-api/telegram"
	"github.com/stretchr/testify/assert"
)

func TestBaseChat_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id": {"100"},
			},
			cfg: telegram.BaseChat{
				ID: 100,
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
			cfg: telegram.BaseChat{},
		},
		{
			exp: url.Values{
				"chat_id": {"username"},
			},
			cfg: telegram.BaseChat{
				ID:              10,
				ChannelUsername: "username",
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

func TestBaseMessage_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id": {"10"},
			},
			cfg: telegram.BaseMessage{
				BaseChat: telegram.BaseChat{ID: 10},
			},
		},
		{
			exp: url.Values{
				"chat_id":              {"100"},
				"reply_to_message_id":  {"100"},
				"disable_notification": {"true"},
				"reply_markup": {
					"{\"keyboard\":[[{\"text\":\"1\"}" +
						",{\"text\":\"2\"}]]}",
				},
			},
			cfg: telegram.BaseMessage{
				BaseChat: telegram.BaseChat{
					ID: 100,
				},
				ReplyToMessageID:    100,
				DisableNotification: true,
				ReplyMarkup: &telegram.ReplyKeyboardMarkup{
					Keyboard: [][]telegram.KeyboardButton{
						[]telegram.KeyboardButton{
							telegram.KeyboardButton{
								Text: "1",
							},
							telegram.KeyboardButton{
								Text: "2",
							},
						},
					},
				},
			},
		},
		{
			exp: nil,
			cfg: telegram.BaseMessage{
				BaseChat: telegram.BaseChat{
					ID: 100,
				},
				ReplyMarkup: replyBadMarkup{},
			},
			expErr: marshalError,
		},
		{
			exp: nil,
			cfg: telegram.BaseMessage{},
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

func TestBaseMessage_Message(t *testing.T) {
	m := telegram.BaseMessage{}
	msg := m.Message()
	assert.NotNil(t, msg, "Message shouln't be nil")
}

func TestMessageCfg_Name(t *testing.T) {
	name := "sendMessage"
	c := telegram.MessageCfg{}
	assert.Equal(t, name, c.Name())
}

func TestMessageCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id": {"10"},
				"text":    {"text"},
			},
			cfg: telegram.MessageCfg{
				BaseMessage: telegram.BaseMessage{
					BaseChat: telegram.BaseChat{ID: 10},
				},
				Text: "text",
			},
		},
		{
			exp: url.Values{
				"chat_id": {"10"},
				"text":    {"text2"},
				"disable_web_page_preview": {"true"},
				"parse_mode":               {"HTML"},
			},
			cfg: telegram.MessageCfg{
				BaseMessage: telegram.BaseMessage{
					BaseChat: telegram.BaseChat{ID: 10},
				},
				Text: "text2",
				DisableWebPagePreview: true,
				ParseMode:             telegram.HTMLMode,
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
			cfg: telegram.MessageCfg{},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"Text",
			),
			cfg: telegram.MessageCfg{
				BaseMessage: telegram.BaseMessage{
					BaseChat: telegram.BaseChat{ID: 10},
				},
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

func TestForwardMessageCfg_Name(t *testing.T) {
	name := "forwardMessage"
	c := telegram.ForwardMessageCfg{}
	assert.Equal(t, name, c.Name())
}

func TestForwardMessageCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id":      {"10"},
				"from_chat_id": {"20"},
				"message_id":   {"30"},
			},
			cfg: telegram.ForwardMessageCfg{
				BaseChat:  telegram.BaseChat{ID: 10},
				FromChat:  telegram.BaseChat{ID: 20},
				MessageID: 30,
			},
		},
		{
			exp: url.Values{
				"chat_id":              {"10"},
				"from_chat_id":         {"20"},
				"message_id":           {"30"},
				"disable_notification": {"true"},
			},
			cfg: telegram.ForwardMessageCfg{
				BaseChat:            telegram.BaseChat{ID: 10},
				FromChat:            telegram.BaseChat{ID: 20},
				MessageID:           30,
				DisableNotification: true,
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"MessageID",
			),
			cfg: telegram.ForwardMessageCfg{
				BaseChat: telegram.BaseChat{ID: 10},
				FromChat: telegram.BaseChat{ID: 20},
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
