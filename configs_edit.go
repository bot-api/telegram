package telegram

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// BaseEdit is base type of all chat edits.
type BaseEdit struct {
	// Required if inline_message_id is not specified.
	// Unique identifier for the target chat or
	// username of the target channel (in the format @channelusername)
	ChatID          int64
	ChannelUsername string
	// Required if inline_message_id is not specified.
	// Unique identifier of the sent message
	MessageID int64
	// Required if chat_id and message_id are not specified.
	// Identifier of the inline message
	InlineMessageID string
	// Only InlineKeyboardMarkup supported right now.
	ReplyMarkup ReplyMarkup
}

// Values returns a url.Values representation of BaseEdit.
func (m BaseEdit) Values() (url.Values, error) {
	v := url.Values{}

	if m.ChannelUsername != "" {
		v.Add("chat_id", m.ChannelUsername)
	} else {
		v.Add("chat_id", strconv.FormatInt(m.ChatID, 10))
	}
	if m.MessageID != 0 {
		v.Add("message_id", strconv.FormatInt(m.MessageID, 10))
	}
	if m.InlineMessageID != "" {
		v.Add("inline_message_id", m.InlineMessageID)
	}

	if m.ReplyMarkup != nil {
		data, err := json.Marshal(m.ReplyMarkup)
		if err != nil {
			return nil, err
		}
		v.Add("reply_markup", string(data))
	}

	return v, nil
}

// EditMessageTextCfg allows you to modify the text in a message.
type EditMessageTextCfg struct {
	BaseEdit
	// New text of the message
	Text string
	// Send Markdown or HTML, if you want Telegram apps
	// to show bold, italic, fixed-width text
	// or inline URLs in your bot's message. Optional.
	ParseMode string
	// Disables link previews for links in this message. Optional.
	DisableWebPagePreview bool
}

// Values returns a url.Values representation of EditMessageTextCfg.
func (cfg EditMessageTextCfg) Values() (url.Values, error) {
	v, err := cfg.BaseEdit.Values()
	if err != nil {
		return nil, err
	}

	v.Add("text", cfg.Text)
	if cfg.ParseMode != "" {
		v.Add("parse_mode", cfg.ParseMode)
	}
	if cfg.DisableWebPagePreview {
		v.Add("disable_web_page_preview", "true")
	}

	return v, nil
}

// Name returns method name
func (EditMessageTextCfg) Name() string {
	return editMessageTextMethod
}

// EditMessageCaptionCfg allows you to modify the caption of a message.
type EditMessageCaptionCfg struct {
	BaseEdit
	// New caption of the message
	Caption string
}

// Values returns a url.Values representation of EditMessageCaptionCfg.
func (cfg EditMessageCaptionCfg) Values() (url.Values, error) {
	v, err := cfg.BaseEdit.Values()
	if err != nil {
		return nil, err
	}
	v.Add("text", cfg.Caption)

	return v, nil
}

// Name returns method name
func (EditMessageCaptionCfg) Name() string {
	return editMessageCaptionMethod
}

// EditMessageReplyMarkupCfg allows you to modify the reply markup of a message.
type EditMessageReplyMarkupCfg struct {
	BaseEdit
}

// Values returns a url.Values representation of EditMessageReplyMarkupCfg.
func (cfg EditMessageReplyMarkupCfg) Values() (url.Values, error) {
	return cfg.BaseEdit.Values()
}

// Name returns method name
func (EditMessageReplyMarkupCfg) Name() string {
	return editMessageReplyMarkupMethod
}
