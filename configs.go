package telegram

import (
	"encoding/json"
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

// SetChatID sets new chat id
func (c *BaseChat) SetChatID(id int64) {
	c.ID = id
}

// GetChatCfg contains information about a getChat request.
type GetChatCfg struct {
	BaseChat
}

// Name returns method name
func (cfg GetChatCfg) Name() string {
	return getChatMethod
}

// Values returns a url.Values representation of GetChatCfg.
// Returns RequiredError if Chat is not set.
func (cfg GetChatCfg) Values() (url.Values, error) {
	return cfg.BaseChat.Values()
}

// GetChatAdministratorsCfg contains information about a getChat request.
type GetChatAdministratorsCfg struct {
	BaseChat
}

// Name returns method name
func (cfg GetChatAdministratorsCfg) Name() string {
	return getChatAdministratorsMethod
}

// Values returns a url.Values representation of GetChatCfg.
// Returns RequiredError if Chat is not set.
func (cfg GetChatAdministratorsCfg) Values() (url.Values, error) {
	return cfg.BaseChat.Values()
}

// GetChatMembersCountCfg contains information about a getChatMemberCount request.
type GetChatMembersCountCfg struct {
	BaseChat
}

// Name returns method name
func (cfg GetChatMembersCountCfg) Name() string {
	return getChatMembersCountMethod
}

// Values returns a url.Values representation of GetChatMembersCountCfg.
// Returns RequiredError if Chat is not set.
func (cfg GetChatMembersCountCfg) Values() (url.Values, error) {
	return cfg.BaseChat.Values()
}

// GetChatMemberCfg contains information about a getChatMember request.
type GetChatMemberCfg struct {
	BaseChat
	UserID int64 `json:"user_id"`
}

// Name returns method name
func (cfg GetChatMemberCfg) Name() string {
	return getChatMemberMethod
}

// Values returns a url.Values representation of GetChatMemberCfg.
// Returns RequiredError if Chat or UserID are not set.
func (cfg GetChatMemberCfg) Values() (url.Values, error) {
	v, err := cfg.BaseChat.Values()
	if err != nil {
		return nil, err
	}
	if cfg.UserID == 0 {
		return nil, NewRequiredError("UserID")
	}
	v.Add("user_id", strconv.FormatInt(cfg.UserID, 10))
	return v, nil
}

// KickChatMemberCfg contains information about a kickChatMember request.
type KickChatMemberCfg struct {
	BaseChat
	UserID int64 `json:"user_id"`
}

// Name returns method name
func (cfg KickChatMemberCfg) Name() string {
	return kickChatMemberMethod
}

// Values returns a url.Values representation of KickChatMemberCfg.
// Returns RequiredError if Chat or UserID are not set.
func (cfg KickChatMemberCfg) Values() (url.Values, error) {
	v, err := cfg.BaseChat.Values()
	if err != nil {
		return nil, err
	}
	if cfg.UserID == 0 {
		return nil, NewRequiredError("UserID")
	}
	v.Add("user_id", strconv.FormatInt(cfg.UserID, 10))
	return v, nil
}

// UnbanChatMemberCfg contains information about a unbanChatMember request.
type UnbanChatMemberCfg struct {
	BaseChat
	UserID int64 `json:"user_id"`
}

// Name returns method name
func (cfg UnbanChatMemberCfg) Name() string {
	return unbanChatMemberMethod
}

// Values returns a url.Values representation of UnbanChatMemberCfg.
// Returns RequiredError if Chat or UserID are not set.
func (cfg UnbanChatMemberCfg) Values() (url.Values, error) {
	v, err := cfg.BaseChat.Values()
	if err != nil {
		return nil, err
	}
	if cfg.UserID == 0 {
		return nil, NewRequiredError("UserID")
	}
	v.Add("user_id", strconv.FormatInt(cfg.UserID, 10))
	return v, nil
}

// LeaveChatCfg contains information about a leaveChat request.
type LeaveChatCfg struct {
	BaseChat
}

// Name returns method name
func (cfg LeaveChatCfg) Name() string {
	return leaveChatMethod
}

// Values returns a url.Values representation of LeaveChatCfg.
// Returns RequiredError if Chat is not set.
func (cfg LeaveChatCfg) Values() (url.Values, error) {
	return cfg.BaseChat.Values()
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
	// Identifier of the first update to be returned.
	// Must be greater by one than the highest
	// among the identifiers of previously received updates.
	// By default, updates starting with the earliest
	// unconfirmed update are returned. An update is considered confirmed
	// as soon as getUpdates is called with an offset
	// higher than its update_id. The negative offset
	// can be specified to retrieve updates starting
	// from -offset update from the end of the updates queue.
	// All previous updates will forgotten.
	Offset int64
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
	if cfg.Offset != 0 {
		v.Add("offset", strconv.FormatInt(cfg.Offset, 10))
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
// Implements Method and Filer interface
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

// Field returns name for webhook file data
func (cfg WebhookCfg) Field() string {
	return "certificate"
}

// File returns certificate data
func (cfg WebhookCfg) File() InputFile {
	return cfg.Certificate
}

// Exist is true if we don't have a certificate to upload.
// It's kind of confusing.
func (cfg WebhookCfg) Exist() bool {
	return cfg.Certificate == nil
}

// Reset method sets new Certificate
func (cfg *WebhookCfg) Reset(i InputFile) {
	cfg.Certificate = i
}

// GetFileID for webhook is always empty
func (cfg *WebhookCfg) GetFileID() string {
	return ""
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

// AnswerInlineQueryCfg contains information on making an InlineQuery response.
type AnswerInlineQueryCfg struct {
	// Unique identifier for the answered query
	InlineQueryID string              `json:"inline_query_id"`
	Results       []InlineQueryResult `json:"results"`
	// The maximum amount of time in seconds
	// that the result of the inline query may be cached on the server.
	// Defaults to 300.
	CacheTime int `json:"cache_time,omitempty"`
	// Pass True, if results may be cached on the server side
	// only for the user that sent the query.
	// By default, results may be returned to any user
	// who sends the same query
	IsPersonal bool `json:"is_personal,omitempty"`
	// Pass the offset that a client should send in the next query
	// with the same text to receive more results.
	// Pass an empty string if there are no more results
	// or if you don‘t support pagination.
	// Offset length can’t exceed 64 bytes.
	NextOffset string `json:"next_offset,omitempty"`
	// If passed, clients will display a button with specified text
	// that switches the user to a private chat with the bot and
	// sends the bot a start message with the parameter switch_pm_parameter
	SwitchPMText string `json:"switch_pm_text,omitempty"`
	// Parameter for the start message sent to the bot
	// when user presses the switch button
	SwitchPMParameter string `json:"switch_pm_parameter"`
}

// Name returns method name
func (cfg AnswerInlineQueryCfg) Name() string {
	return answerInlineQueryMethod
}

// Values returns a url.Values representation of AnswerInlineQueryCfg.
// Returns a RequiredError if Action is empty.
func (cfg AnswerInlineQueryCfg) Values() (url.Values, error) {
	v := url.Values{}
	if cfg.Results == nil || len(cfg.Results) == 0 {
		return nil, NewRequiredError("Results")
	}
	data, err := json.Marshal(cfg.Results)
	if err != nil {
		return nil, err
	}
	v.Add("results", string(data))
	v.Add("inline_query_id", cfg.InlineQueryID)
	if cfg.CacheTime > 0 {
		v.Add("cache_time", strconv.Itoa(cfg.CacheTime))
	}
	if cfg.IsPersonal {
		v.Add("is_personal", strconv.FormatBool(cfg.IsPersonal))
	}
	if cfg.NextOffset != "" {
		v.Add("next_offset", cfg.NextOffset)
	}
	if cfg.SwitchPMText != "" {
		v.Add("switch_pm_text", cfg.SwitchPMText)
	}
	if cfg.SwitchPMParameter != "" {
		v.Add("switch_pm_parameter", cfg.SwitchPMParameter)
	}
	return v, nil
}
