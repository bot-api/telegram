package telegram

import (
	"fmt"
	"net/url"

	"golang.org/x/net/context"
)

func newBM(chatID int64) BaseMessage {
	return BaseMessage{
		BaseChat: BaseChat{
			ID: chatID,
		},
	}
}

// NewMessage creates a new Message.
//
// chatID is where to send it, text is the message text.
func NewMessage(chatID int64, text string) MessageCfg {
	return MessageCfg{
		BaseMessage: newBM(chatID),
		Text:        text,
		DisableWebPagePreview: false,
	}
}

// NewMessagef creates a new Message with formatting.
//
// chatID is where to send it, text is the message text
func NewMessagef(chatID int64, text string, args ...interface{}) MessageCfg {
	return MessageCfg{
		BaseMessage: newBM(chatID),
		Text:        fmt.Sprintf(text, args...),
		DisableWebPagePreview: false,
	}
}

// NewKeyboard creates keyboard by matrix i*j.
func NewKeyboard(buttons [][]string) [][]KeyboardButton {
	rows := make([][]KeyboardButton, len(buttons))
	for i, colButtons := range buttons {
		cols := make([]KeyboardButton, len(colButtons))
		for j, button := range colButtons {
			cols[j].Text = button
		}
		rows[i] = cols
	}
	return rows
}

// NewHKeyboard creates keyboard with horizontal buttons only.
// [ first ] [ second ] [ third ]
func NewHKeyboard(buttons ...string) [][]KeyboardButton {
	row := make([]KeyboardButton, len(buttons))
	for i, button := range buttons {
		row[i].Text = button
	}
	return [][]KeyboardButton{row}
}

// NewVKeyboard creates keyboard with vertical buttons only
// [ first  ]
// [ second ]
// [ third  ]
func NewVKeyboard(buttons ...string) [][]KeyboardButton {
	r := make([][]KeyboardButton, len(buttons))
	for i, button := range buttons {
		r[i] = []KeyboardButton{{Text: button}}
	}
	return r
}

// NewHInlineKeyboard creates inline keyboard with horizontal buttons only.
// [ first ] [ second ] [ third ]
func NewHInlineKeyboard(prefix string, text []string, data []string) [][]InlineKeyboardButton {
	row := make([]InlineKeyboardButton, len(text))

	for i, t := range text {
		row[i].Text = t
		row[i].CallbackData = prefix + data[i]
	}
	return [][]InlineKeyboardButton{row}
}

// NewVInlineKeyboard creates inline keyboard with vertical buttons only
// [ first  ]
// [ second ]
// [ third  ]
func NewVInlineKeyboard(prefix string, text []string, data []string) [][]InlineKeyboardButton {
	r := make([][]InlineKeyboardButton, len(text))
	for i, button := range text {
		r[i] = []InlineKeyboardButton{
			{Text: button, CallbackData: prefix + data[i]},
		}
	}
	return r
}

// NewForwardMessage creates a new Message.
//
// chatID is where to send it, text is the message text.
func NewForwardMessage(chatID, fromChatID, messageID int64) ForwardMessageCfg {
	return ForwardMessageCfg{
		BaseChat: BaseChat{
			ID: chatID,
		},
		FromChat: BaseChat{
			ID: fromChatID,
		},
		MessageID: messageID,
	}
}

// NewUserProfilePhotos gets user profile photos.
//
// userID is the ID of the user you wish to get profile photos from.
func NewUserProfilePhotos(userID int64) UserProfilePhotosCfg {
	return UserProfilePhotosCfg{
		UserID: userID,
		Offset: 0,
		Limit:  0,
	}
}

// NewUpdate gets updates since the last Offset with Timeout 30 seconds
//
// offset is the last Update ID to include.
// You likely want to set this to the last Update ID plus 1.
// The negative offset can be specified to retrieve updates starting
// from -offset update from the end of the updates queue.
// All previous updates will forgotten.
func NewUpdate(offset int64) UpdateCfg {
	return UpdateCfg{
		Offset:  offset,
		Limit:   0,
		Timeout: 30,
	}
}

// NewChatAction sets a chat action.
// Actions last for 5 seconds, or until your next action.
//
// chatID is where to send it, action should be set via Action constants.
func NewChatAction(chatID int64, action string) ChatActionCfg {
	return ChatActionCfg{
		BaseChat: BaseChat{ID: chatID},
		Action:   action,
	}
}

// NewLocation shares your location.
//
// chatID is where to send it, latitude and longitude are coordinates.
func NewLocation(chatID int64, lat float64, lon float64) LocationCfg {
	return LocationCfg{
		BaseMessage: newBM(chatID),
		Location: Location{
			Latitude:  lat,
			Longitude: lon,
		},
	}
}

// NewPhotoUpload creates a new photo uploader.
//
// chatID is where to send it, inputFile is a file representation.
func NewPhotoUpload(chatID int64, inputFile InputFile) PhotoCfg {
	return PhotoCfg{
		BaseFile: BaseFile{
			BaseMessage: newBM(chatID),
			InputFile:   inputFile,
		},
	}
}

// NewPhotoShare creates a new photo uploader.
//
// chatID is where to send it
func NewPhotoShare(chatID int64, fileID string) PhotoCfg {
	return PhotoCfg{
		BaseFile: BaseFile{
			BaseMessage: newBM(chatID),
			FileID:      fileID,
		},
	}
}

// NewAnswerCallback creates a new callback message.
func NewAnswerCallback(id, text string) AnswerCallbackCfg {
	return AnswerCallbackCfg{
		CallbackQueryID: id,
		Text:            text,
		ShowAlert:       false,
	}
}

// NewAnswerCallbackWithAlert creates a new callback message that alerts
// the user.
func NewAnswerCallbackWithAlert(id, text string) AnswerCallbackCfg {
	return AnswerCallbackCfg{
		CallbackQueryID: id,
		Text:            text,
		ShowAlert:       true,
	}
}

// NewEditMessageText allows you to edit the text of a message.
func NewEditMessageText(chatID, messageID int64, text string) EditMessageTextCfg {
	return EditMessageTextCfg{
		BaseEdit: BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Text: text,
	}
}

// NewEditMessageCaption allows you to edit the caption of a message.
func NewEditMessageCaption(chatID, messageID int64, caption string) EditMessageCaptionCfg {
	return EditMessageCaptionCfg{
		BaseEdit: BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Caption: caption,
	}
}

// NewEditMessageReplyMarkup allows you to edit the inline
// keyboard markup.
func NewEditMessageReplyMarkup(chatID, messageID int64, replyMarkup *InlineKeyboardMarkup) EditMessageReplyMarkupCfg {
	return EditMessageReplyMarkupCfg{
		BaseEdit: BaseEdit{
			ChatID:      chatID,
			MessageID:   messageID,
			ReplyMarkup: replyMarkup,
		},
	}
}

// NewWebhook creates a new webhook.
//
// link is the url parsable link you wish to get the updates.
func NewWebhook(link string) WebhookCfg {
	u, _ := url.Parse(link)

	return WebhookCfg{
		URL: u.String(),
	}
}

// NewWebhookWithCert creates a new webhook with a certificate.
//
// link is the url you wish to get webhooks,
// file contains a string to a file, FileReader, or FileBytes.
func NewWebhookWithCert(link string, file InputFile) WebhookCfg {
	u, _ := url.Parse(link)

	return WebhookCfg{
		URL:         u.String(),
		Certificate: file,
	}
}

// CloneMessage convert message to Messenger type to send it to another chat.
// It supports only data message: Text, Sticker, Audio, Photo, Location,
// Contact, Audio, Voice, Document.
func CloneMessage(msg *Message, baseMessage *BaseMessage) Messenger {
	var base BaseMessage
	if baseMessage == nil {
		base = newBM(msg.Chat.ID)

	} else {
		base = *baseMessage
	}
	if msg.Text != "" {
		return &MessageCfg{
			BaseMessage: base,
			Text:        msg.Text,
		}
	}
	if msg.Sticker != nil {
		return &StickerCfg{
			BaseFile: BaseFile{
				BaseMessage: base,
				FileID:      msg.Sticker.FileID,
			},
		}
	}
	if msg.Photo != nil && len(msg.Photo) > 0 {
		return &PhotoCfg{
			BaseFile: BaseFile{
				BaseMessage: base,
				FileID:      msg.Photo[len(msg.Photo)-1].FileID,
			},
			Caption: msg.Caption,
		}
	}
	if msg.Location != nil {
		return &LocationCfg{
			BaseMessage: base,
			Location:    *msg.Location,
		}
	}
	if msg.Contact != nil {
		return &ContactCfg{
			BaseMessage: base,
			Contact:     *msg.Contact,
		}
	}
	if msg.Audio != nil {
		return &AudioCfg{
			BaseFile: BaseFile{
				BaseMessage: base,
				FileID:      msg.Audio.FileID,
			},
			Duration:  msg.Audio.Duration,
			Performer: msg.Audio.Performer,
			Title:     msg.Audio.Title,
		}
	}
	if msg.Voice != nil {
		return &VoiceCfg{
			BaseFile: BaseFile{
				BaseMessage: base,
				FileID:      msg.Voice.FileID,
			},
			Duration: msg.Voice.Duration,
		}
	}
	if msg.Document != nil {
		return &DocumentCfg{
			BaseFile: BaseFile{
				BaseMessage: base,
				FileID:      msg.Document.FileID,
			},
		}
	}
	return nil
}

// GetUpdates runs loop and requests updates from telegram.
// It breaks loop, close out channel and returns error
// if something happened during update cycle.
func GetUpdates(
	ctx context.Context,
	api *API,
	cfg UpdateCfg,
	out chan<- Update) error {

	var rErr error
	defer close(out)

loop:
	for {
		updates, err := api.GetUpdates(
			ctx,
			cfg,
		)
		if err != nil {
			rErr = err
			break loop
		}
		for _, update := range updates {
			if update.UpdateID >= cfg.Offset {
				cfg.Offset = update.UpdateID + 1
				select {
				case <-ctx.Done():
					rErr = ctx.Err()
					break loop
				case out <- update:
				}
			}
		}
	}
	return rErr
}

// InlineQuery helpers

// NewInlineQueryResultArticle creates a new inline query article.
func NewInlineQueryResultArticle(id, title, messageText string) *InlineQueryResultArticle {
	return &InlineQueryResultArticle{
		BaseInlineQueryResult: BaseInlineQueryResult{
			Type: "article",
			ID:   id,
			InputMessageContent: InputTextMessageContent{
				MessageText: messageText,
			},
		},
		Title: title,
	}
}
