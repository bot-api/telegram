package telegram_test

import (
	"encoding/json"
	"net/url"
	"reflect"
	"testing"

	"github.com/bot-api/telegram"
	"gopkg.in/stretchr/testify.v1/assert"
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
						{
							{
								Text: "1",
							},
							{
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
			expErr: &json.MarshalerError{
				Type: reflect.TypeOf(replyBadMarkup{}),
				Err:  marshalError,
			},
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

func TestForwardMessageCfg_Message(t *testing.T) {
	c := telegram.ForwardMessageCfg{}
	assert.Equal(t, c.Message(), &telegram.Message{})
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
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
			cfg: telegram.ForwardMessageCfg{},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
			cfg: telegram.ForwardMessageCfg{
				BaseChat: telegram.BaseChat{ID: 10},
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

func TestLocationCfg_Name(t *testing.T) {
	name := "sendLocation"
	c := telegram.LocationCfg{}
	assert.Equal(t, name, c.Name())
}

func TestLocationCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id":   {"10"},
				"longitude": {"20"},
				"latitude":  {"30"},
			},
			cfg: telegram.LocationCfg{
				BaseMessage: telegram.BaseMessage{
					BaseChat: telegram.BaseChat{ID: 10},
				},
				Location: telegram.Location{
					Longitude: 20,
					Latitude:  30,
				},
			},
		},
		{
			exp: url.Values{
				"chat_id":              {"10"},
				"longitude":            {"0"},
				"latitude":             {"0"},
				"reply_to_message_id":  {"20"},
				"disable_notification": {"true"},
			},
			cfg: telegram.LocationCfg{
				BaseMessage: telegram.BaseMessage{
					BaseChat:            telegram.BaseChat{ID: 10},
					ReplyToMessageID:    20,
					DisableNotification: true,
				},
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
			cfg: telegram.LocationCfg{},
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestContactCfg_Name(t *testing.T) {
	name := "sendContact"
	c := telegram.ContactCfg{}
	assert.Equal(t, name, c.Name())
}

func TestContactCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id":      {"10"},
				"user_id":      {"30"},
				"phone_number": {"phone_number_value"},
				"first_name":   {"first_name_value"},
				"last_name":    {"last_name_value"},
			},
			cfg: telegram.ContactCfg{
				BaseMessage: telegram.BaseMessage{
					BaseChat: telegram.BaseChat{ID: 10},
				},
				Contact: telegram.Contact{
					FirstName:   "first_name_value",
					LastName:    "last_name_value",
					PhoneNumber: "phone_number_value",
					UserID:      30,
				},
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"FirstName", "PhoneNumber",
			),
			cfg: telegram.ContactCfg{
				BaseMessage: telegram.BaseMessage{
					BaseChat: telegram.BaseChat{ID: 10},
				},
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
			cfg: telegram.ContactCfg{},
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}

func TestVenueCfg_Name(t *testing.T) {
	name := "sendVenue"
	c := telegram.VenueCfg{}
	assert.Equal(t, name, c.Name())
}

func TestVenueCfg_Values(t *testing.T) {
	testTable := []cfgTT{
		{
			exp: url.Values{
				"chat_id":       {"10"},
				"longitude":     {"20"},
				"latitude":      {"30"},
				"title":         {"venue title"},
				"address":       {"venue address"},
				"foursquare_id": {"foursquare-id"},
			},
			cfg: telegram.VenueCfg{
				BaseMessage: telegram.BaseMessage{
					BaseChat: telegram.BaseChat{ID: 10},
				},
				Venue: telegram.Venue{
					Location: telegram.Location{
						Longitude: 20,
						Latitude:  30,
					},
					Title:        "venue title",
					Address:      "venue address",
					FoursquareID: "foursquare-id",
				},
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"Title", "Address",
			),
			cfg: telegram.VenueCfg{
				BaseMessage: telegram.BaseMessage{
					BaseChat: telegram.BaseChat{ID: 10},
				},
			},
		},
		{
			exp: nil,
			expErr: telegram.NewRequiredError(
				"ID", "ChannelUsername",
			),
			cfg: telegram.VenueCfg{},
		},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		values, err := tt.cfg.Values()
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, tt.exp, values)
	}
}
