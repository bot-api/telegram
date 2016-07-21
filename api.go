package telegram

// Package telegram provides implementation for Telegram Bot API

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"
)

const (
	// APIEndpoint is the endpoint for all API methods,
	// with formatting for Sprintf.
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram.
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"
)

// HTTPDoer interface helps to test api
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// DebugFunc describes function for debugging.
type DebugFunc func(msg string, fields map[string]interface{})

// DefaultDebugFunc prints debug message to default logger
var DefaultDebugFunc = func(msg string, fields map[string]interface{}) {
	log.Printf("%s %v", msg, fields)
}

// API implements Telegram bot API
// described on https://core.telegram.org/bots/api
type API struct {
	// token is a unique authentication string,
	// obtained by each bot when it is created.
	// The token looks something like
	// 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
	token        string
	client       HTTPDoer
	apiEndpoint  string
	fileEndpoint string
	debug        bool
	debugFunc    DebugFunc
}

// New returns API instance with default http client
func New(token string) *API {
	return NewWithClient(token, http.DefaultClient)
}

// NewWithClient returns API instance with custom http client
func NewWithClient(token string, client HTTPDoer) *API {
	return &API{
		token:        token,
		client:       client,
		apiEndpoint:  APIEndpoint,
		fileEndpoint: FileEndpoint,
		debugFunc:    DefaultDebugFunc,
	}
}

// Invoke is a generic method that helps to make request to Telegram Api.
// Use particular methods instead (e.x. GetMe, GetUpdates etc).
// The only case when this method seems useful is
// when Telegram Api has method
// that still doesn't exist in this implementation.
func (c *API) Invoke(ctx context.Context, m Method, dst interface{}) error {
	params, err := m.Values()
	if err != nil {
		return err
	}
	var req *http.Request
	if mf, casted := m.(Filer); casted && !mf.Exist() {
		// upload a file, if FileID doesn't exist
		req, err = c.getUploadRequest(
			m.Name(),
			params,
			mf.Field(),
			mf.File(),
		)
	} else {
		req, err = c.getFormRequest(m.Name(), params)
	}
	if err != nil {
		return err
	}
	return c.makeRequest(ctx, req, dst)
}

// Debug enables sending debug messages to default log
func (c *API) Debug(val bool) {
	c.debug = val
}

// DebugFunc replaces default debug function
func (c *API) DebugFunc(f DebugFunc) {
	c.debugFunc = f
}

// Telegram Bot API methods

// GetMe returns basic information about the bot in form of a User object
func (c *API) GetMe(ctx context.Context) (*User, error) {
	u := &User{}
	if err := c.Invoke(ctx, MeCfg{}, u); err != nil {
		return nil, err
	}
	return u, nil
}

// GetChat returns up to date information about the chat
// (current name of the user for one-on-one conversations,
// current username of a user, group or channel, etc.).
// Returns a Chat object on success.
func (c *API) GetChat(ctx context.Context, cfg GetChatCfg) (*Chat, error) {
	chat := &Chat{}
	if err := c.Invoke(ctx, cfg, chat); err != nil {
		return nil, err
	}
	return chat, nil
}

// GetChatAdministrators returns a list of administrators in a chat.
// On success, returns an Array of ChatMember objects
// that contains information about all chat administrators
// except other bots. If the chat is a group or a supergroup
// and no administrators were appointed, only the creator will be returned.
func (c *API) GetChatAdministrators(
	ctx context.Context,
	cfg GetChatAdministratorsCfg) ([]ChatMember, error) {

	chatMembers := []ChatMember{}
	if err := c.Invoke(ctx, cfg, &chatMembers); err != nil {
		return nil, err
	}
	return chatMembers, nil
}

// GetChatMembersCount returns the number of members in a chat.
func (c *API) GetChatMembersCount(
	ctx context.Context,
	cfg GetChatMembersCountCfg) (int, error) {

	var count int
	if err := c.Invoke(ctx, cfg, &count); err != nil {
		return count, err
	}
	return count, nil
}

// GetChatMember returns information about a member of a chat.
func (c *API) GetChatMember(
	ctx context.Context,
	cfg GetChatMemberCfg) (*ChatMember, error) {

	member := &ChatMember{}
	if err := c.Invoke(ctx, cfg, member); err != nil {
		return nil, err
	}
	return member, nil
}

// KickChatMember kicks a user from a group or a supergroup.
// In the case of supergroups, the user will not be able to return
// to the group on their own using invite links, etc., unless unbanned first.
// The bot must be an administrator in the group for this to work.
// Returns True on success.
func (c *API) KickChatMember(
	ctx context.Context,
	cfg KickChatMemberCfg) (bool, error) {

	var result bool
	if err := c.Invoke(ctx, cfg, &result); err != nil {
		return result, err
	}
	return result, nil
}

// UnbanChatMember unbans a previously kicked user in a supergroup.
// The user will not return to the group automatically,
// but will be able to join via link, etc.
// The bot must be an administrator in the group for this to work.
// Returns True on success.
func (c *API) UnbanChatMember(
	ctx context.Context,
	cfg UnbanChatMemberCfg) (bool, error) {

	var result bool
	if err := c.Invoke(ctx, cfg, &result); err != nil {
		return result, err
	}
	return result, nil
}

// LeaveChat method helps your bot to leave a group, supergroup or channel.
// Returns True on success.
func (c *API) LeaveChat(
	ctx context.Context,
	cfg LeaveChatCfg) (bool, error) {

	var result bool
	if err := c.Invoke(ctx, cfg, &result); err != nil {
		return result, err
	}
	return result, nil
}

// GetUpdates requests incoming updates using long polling.
// This method will not work if an outgoing webhook is set up.
// In order to avoid getting duplicate updates,
// recalculate offset after each server response.
func (c *API) GetUpdates(
	ctx context.Context,
	cfg UpdateCfg) ([]Update, error) {

	updates := []Update{}
	if err := c.Invoke(ctx, cfg, &updates); err != nil {
		return nil, err
	}
	return updates, nil
}

// GetUserProfilePhotos requests a list of profile pictures for a user.
func (c *API) GetUserProfilePhotos(
	ctx context.Context,
	cfg UserProfilePhotosCfg) (*UserProfilePhotos, error) {

	photos := &UserProfilePhotos{}
	if err := c.Invoke(ctx, cfg, photos); err != nil {
		return nil, err
	}
	return photos, nil
}

// SendChatAction tells the user that something is happening
// on the bot's side. The status is set for 5 seconds or less
// (when a message arrives from your bot,
// Telegram clients clear its typing status).
func (c *API) SendChatAction(ctx context.Context, cfg ChatActionCfg) error {
	return c.Invoke(ctx, cfg, nil)
}

// GetFile returns a File which can download a file from Telegram.
//
// Requires FileID.
func (c *API) GetFile(ctx context.Context, cfg FileCfg) (*File, error) {
	var file File
	err := c.Invoke(ctx, cfg, &file)
	if err != nil {
		return nil, err
	}
	file.Link = fmt.Sprintf(c.fileEndpoint, c.token, file.FilePath)
	return &file, nil
}

// DownloadFile downloads file from telegram servers to w
//
// Requires FileID
func (c *API) DownloadFile(ctx context.Context, cfg FileCfg, w io.Writer) error {
	f, err := c.GetFile(ctx, cfg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", f.Link, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			if c.debug {
				c.print("body close error", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// AnswerCallbackQuery sends a response to an inline query callback.
func (c *API) AnswerCallbackQuery(
	ctx context.Context,
	cfg AnswerCallbackCfg) (bool, error) {

	var result bool
	return result, c.Invoke(ctx, cfg, &result)
}

// Edit method allows you to change an existing message in the message history
// instead of sending a new one with a result of an action.
// This is most useful for messages with inline keyboards using callback queries,
// but can also help reduce clutter in conversations with regular chat bots.
// Please note, that it is currently only possible to edit messages without
// reply_markup or with inline keyboards.
//
// You can use this method directly or one of:
// EditMessageText, EditMessageCaption, EditMessageReplyMarkup,
func (c *API) Edit(ctx context.Context, cfg Method) (*EditResult, error) {
	er := &EditResult{}
	return er, c.Invoke(ctx, cfg, er)
}

// Send method sends message.
//
// TODO m0sth8: rewrite this doc
func (c *API) Send(ctx context.Context, cfg Messenger) (*Message, error) {
	msg := cfg.Message()
	return msg, c.Invoke(ctx, cfg, msg)
}

// === Methods based on Send method

// SendMessage sends text message.
func (c *API) SendMessage(
	ctx context.Context,
	cfg MessageCfg) (*Message, error) {

	return c.Send(ctx, cfg)
}

// SendSticker sends message with sticker.
func (c *API) SendSticker(
	ctx context.Context,
	cfg StickerCfg) (*Message, error) {

	return c.Send(ctx, cfg)
}

// SendVenue sends venue message.
func (c *API) SendVenue(
	ctx context.Context,
	cfg VenueCfg) (*Message, error) {

	return c.Send(ctx, cfg)
}

// SendContact sends phone contact message.
func (c *API) SendContact(
	ctx context.Context,
	cfg ContactCfg) (*Message, error) {

	return c.Send(ctx, cfg)
}

// SendPhoto sends photo.
func (c *API) SendPhoto(
	ctx context.Context,
	cfg PhotoCfg) (*Message, error) {

	return c.Send(ctx, cfg)
}

// SendAudio sends Audio.
func (c *API) SendAudio(
	ctx context.Context,
	cfg AudioCfg) (*Message, error) {

	return c.Send(ctx, cfg)
}

// SendVideo sends Video.
func (c *API) SendVideo(
	ctx context.Context,
	cfg VideoCfg) (*Message, error) {

	return c.Send(ctx, cfg)
}

// SendVoice sends Voice.
func (c *API) SendVoice(
	ctx context.Context,
	cfg VoiceCfg) (*Message, error) {

	return c.Send(ctx, cfg)
}

// SendDocument sends Document.
func (c *API) SendDocument(
	ctx context.Context,
	cfg DocumentCfg) (*Message, error) {

	return c.Send(ctx, cfg)
}

// ForwardMessage forwards messages of any kind.
func (c *API) ForwardMessage(
	ctx context.Context,
	cfg ForwardMessageCfg) (*Message, error) {

	return c.Send(ctx, cfg)
}

// === Methods based on Edit method

// EditMessageText modifies the text of message.
// Use this method to edit only the text of messages
// sent by the bot or via the bot (for inline bots).
// On success, if edited message is sent by the bot,
// the edited Message is returned, otherwise True is returned.
func (c *API) EditMessageText(
	ctx context.Context,
	cfg EditMessageTextCfg) (*EditResult, error) {

	return c.Edit(ctx, cfg)
}

// EditMessageCaption modifies the caption of message.
// Use this method to edit only the caption of messages
// sent by the bot or via the bot (for inline bots).
// On success, if edited message is sent by the bot,
// the edited Message is returned, otherwise True is returned.
func (c *API) EditMessageCaption(
	ctx context.Context,
	cfg EditMessageCaptionCfg) (*EditResult, error) {

	return c.Edit(ctx, cfg)
}

// EditMessageReplyMarkup modifies the reply markup of message.
// Use this method to edit only the reply markup of messages
// sent by the bot or via the bot (for inline bots).
// On success, if edited message is sent by the bot,
// the edited Message is returned, otherwise True is returned.
func (c *API) EditMessageReplyMarkup(
	ctx context.Context,
	cfg EditMessageReplyMarkupCfg) (*EditResult, error) {

	return c.Edit(ctx, cfg)
}

// SetWebhook sets a webhook.
// Use this method to specify a url and receive incoming updates
// via an outgoing webhook. Whenever there is an update for the bot,
// we will send an HTTPS POST request to the specified url,
// containing a JSON‐serialized Update.
// In case of an unsuccessful request,
// we will give up after a reasonable amount of attempts.
//
// If this is set, GetUpdates will not get any data!
//
// If you do not have a legitimate TLS certificate,
// you need to include your self signed certificate with the config.
func (c *API) SetWebhook(ctx context.Context, cfg WebhookCfg) error {
	return c.Invoke(ctx, cfg, nil)
}

// AnswerInlineQuery sends answers to an inline query.
// On success, True is returned. No more than 50 results per query are allowed.
func (c *API) AnswerInlineQuery(ctx context.Context, cfg AnswerInlineQueryCfg) (bool, error) {
	var result bool
	return result, c.Invoke(ctx, cfg, &result)
}

// Internal methods

func (c *API) print(msg string, fields map[string]interface{}) {
	if c.debugFunc != nil {
		c.debugFunc(msg, fields)
	}
}

func (c *API) getFormRequest(
	method string,
	params url.Values) (*http.Request, error) {

	urlStr := fmt.Sprintf(c.apiEndpoint, c.token, method)
	body := params.Encode()
	if c.debug {
		c.print("request", map[string]interface{}{
			"url":  urlStr,
			"data": body,
		})
	}

	req, err := http.NewRequest(
		"POST",
		urlStr,
		strings.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (c *API) getUploadRequest(
	method string,
	params url.Values,
	field string,
	file InputFile) (*http.Request, error) {

	urlStr := fmt.Sprintf(c.apiEndpoint, c.token, method)

	if c.debug {
		c.print("file request", map[string]interface{}{
			"url":        urlStr,
			"data":       params.Encode(),
			"file_field": field,
			"file_name":  file.Name(),
		})
	}

	buf := &bytes.Buffer{}

	w := multipart.NewWriter(buf)

	for key, values := range params {
		for _, value := range values {
			err := w.WriteField(key, value)
			if err != nil {
				return nil, fmt.Errorf(
					"can't write field %s, cause %s",
					key, err.Error(),
				)
			}
		}
	}
	fw, err := w.CreateFormFile(field, file.Name())
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, file.Reader()); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		urlStr,
		buf,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return req, nil
}

func (c *API) makeRequest(
	ctx context.Context,
	req *http.Request,
	dst interface{}) error {

	var err error
	var resp *http.Response

	resp, err = makeRequest(ctx, c.client, req)

	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			if c.debug {
				c.print("body close error", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
	if c.debug {
		c.print("response", map[string]interface{}{
			"status_code": resp.StatusCode,
		})
	}
	if resp.StatusCode == http.StatusForbidden {
		// read all from body to save keep-alive connection.
		if _, err = io.Copy(ioutil.Discard, resp.Body); err != nil {
			if c.debug {
				c.print("discard error", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
		return errForbidden
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if c.debug {
		c.print("response", map[string]interface{}{
			"data": string(data),
		})
	}

	apiResponse := APIResponse{}
	err = json.Unmarshal(data, &apiResponse)
	if err != nil {
		return err
	}
	if !apiResponse.Ok {
		if apiResponse.ErrorCode == 401 {
			return errUnauthorized
		}
		return &APIError{
			Description: apiResponse.Description,
			ErrorCode:   apiResponse.ErrorCode,
		}
	}
	if dst != nil && apiResponse.Result != nil {
		err = json.Unmarshal(*apiResponse.Result, dst)
	}
	return err
}
