package telegram

import (
	"net/url"
	"strconv"
)

// BaseChat describes chat settings. It's an abstract type.
//
// You must set ID or ChannelUsername.
type BaseChat struct {
	// Unique identifier for the target chat
	ID int64
	// Username of the target channel (in the format @channelusername)
	ChannelUsername string
}

// Values returns RequiredError if neither ID or ChannelUsername are empty.
// Prefers ChannelUsername if both ID and ChannelUsername are not empty.
func (c BaseChat) Values() (url.Values, error) {
	v := url.Values{}
	if c.ChannelUsername != "" {
		v.Add("chat_id", c.ChannelUsername)
		return v, nil
	}
	if c.ID != 0 {
		v.Add("chat_id", strconv.FormatInt(c.ID, 10))
		return v, nil
	}
	return nil, NewRequiredError("ID", "ChannelUsername")
}

// MeCfg contains information about a getMe request.
type MeCfg struct{}

// Name returns method name
func (cfg MeCfg) Name() string {
	return getMeMethod
}

// Values for getMe is empty
func (cfg MeCfg) Values() (url.Values, error) {
	return nil, nil
}

// UpdateCfg contains information about a getUpdates request.
type UpdateCfg struct {
	Offset uint64
	// Limits the number of updates to be retrieved.
	// Values between 1—100 are accepted. Defaults to 100.
	Limit int
	// Timeout in seconds for long polling.
	// Defaults to 0, i.e. usual short polling
	Timeout int
}

// Name returns method name
func (cfg UpdateCfg) Name() string {
	return getUpdatesMethod
}

// Values returns getUpdate params.
// It returns error if Limit is not between 0 and 100.
// Zero params are not included to request.
func (cfg UpdateCfg) Values() (url.Values, error) {
	if cfg.Limit < 0 || cfg.Limit > 100 {
		return nil, NewValidationError(
			"Limit",
			"should be between 1 and 100",
		)
	}
	v := url.Values{}
	if cfg.Offset > 0 {
		v.Add("offset", strconv.FormatUint(cfg.Offset, 10))
	}
	if cfg.Limit > 0 {
		v.Add("limit", strconv.Itoa(cfg.Limit))
	}
	if cfg.Timeout > 0 {
		v.Add("timeout", strconv.Itoa(cfg.Timeout))
	}
	return v, nil
}

// ChatActionCfg contains information about a SendChatAction request.
// Action field is required.
type ChatActionCfg struct {
	BaseChat
	// Type of action to broadcast.
	// Choose one, depending on what the user is about to receive:
	// typing for text messages, upload_photo for photos,
	// record_video or upload_video for videos,
	// record_audio or upload_audio for audio files,
	// upload_document for general files,
	// find_location for location data.
	// Use one of constants: ActionTyping, ActionFindLocation, etc
	Action string
}

// Name returns method name
func (cfg ChatActionCfg) Name() string {
	return sendChatActionMethod
}

// Values returns a url.Values representation of ChatActionCfg.
// Returns a RequiredError if Action is empty.
func (cfg ChatActionCfg) Values() (url.Values, error) {
	v, err := cfg.BaseChat.Values()
	if err != nil {
		return nil, err
	}
	if cfg.Action == "" {
		return nil, NewRequiredError("Action")
	}
	v.Add("action", cfg.Action)
	return v, nil
}

// UserProfilePhotosCfg contains information about a
// GetUserProfilePhotos request.
type UserProfilePhotosCfg struct {
	UserID int64
	// Sequential number of the first photo to be returned.
	// By default, all photos are returned.
	Offset int
	// Limits the number of photos to be retrieved.
	// Values between 1—100 are accepted. Defaults to 100.
	Limit int
}

// Name returns method name
func (cfg UserProfilePhotosCfg) Name() string {
	return getUserProfilePhotosMethod
}

// Values returns a url.Values representation of UserProfilePhotosCfg.
// Returns RequiredError if UserID is empty.
func (cfg UserProfilePhotosCfg) Values() (url.Values, error) {
	if cfg.Limit < 0 || cfg.Limit > 100 {
		return nil, NewValidationError(
			"Limit",
			"should be between 1 and 100",
		)
	}
	v := url.Values{}
	if cfg.UserID == 0 {
		return nil, NewRequiredError("UserID")
	}
	v.Add("user_id", strconv.FormatInt(cfg.UserID, 10))
	if cfg.Offset != 0 {
		v.Add("offset", strconv.Itoa(cfg.Offset))
	}
	if cfg.Limit != 0 {
		v.Add("limit", strconv.Itoa(cfg.Limit))
	}
	return v, nil
}

// FileCfg has information about a file hosted on Telegram.
type FileCfg struct {
	FileID string
}

// Name returns method name
func (cfg FileCfg) Name() string {
	return getFileMethod
}

// Values returns a url.Values representation of FileCfg.
func (cfg FileCfg) Values() (url.Values, error) {
	v := url.Values{}
	v.Add("file_id", cfg.FileID)
	return v, nil
}

// WebhookCfg contains information about a SetWebhook request.
type WebhookCfg struct {
	URL string
	// self generated TLS certificate
	Certificate InputFile
}

// Name method returns Telegram API method name for sending Location.
func (cfg WebhookCfg) Name() string {
	return setWebhookMethod
}

// Values returns a url.Values representation of Webhook config.
func (cfg WebhookCfg) Values() (url.Values, error) {
	v := url.Values{}
	v.Add("url", cfg.URL)
	return v, nil
}

func (cfg WebhookCfg) Field() string {
	return "certificate"
}

func (cfg WebhookCfg) File() InputFile {
	return cfg.Certificate
}

// Exist is true if we don't have a certificate to upload.
// It's kind of confusing.
func (cfg WebhookCfg) Exist() bool {
	return cfg.Certificate == nil
}

// AnswerCallbackCfg contains information on making a anserCallbackQuery response.
type AnswerCallbackCfg struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert"`
}

// Name returns method name
func (cfg AnswerCallbackCfg) Name() string {
	return answerCallbackQueryMethod
}

// Values returns a url.Values representation of AnswerCallbackCfg.
// Returns a RequiredError if Action is empty.
func (cfg AnswerCallbackCfg) Values() (url.Values, error) {
	v := url.Values{}
	v.Add("callback_query_id", cfg.CallbackQueryID)
	v.Add("text", cfg.Text)
	v.Add("show_alert", strconv.FormatBool(cfg.ShowAlert))
	return v, nil
}
