package telegram

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
)

// Assert interfaces
var _ Filer = (*PhotoCfg)(nil)
var _ Filer = (*AudioCfg)(nil)
var _ Filer = (*VideoCfg)(nil)
var _ Filer = (*VoiceCfg)(nil)
var _ Filer = (*DocumentCfg)(nil)
var _ Filer = (*StickerCfg)(nil)
var _ Filer = (*WebhookCfg)(nil)

// BaseMessage is a base type for all message config types.
// Implements Messenger interface.
type BaseMessage struct {
	BaseChat
	// If the message is a reply, ID of the original message
	ReplyToMessageID int64
	// Additional interface options.
	// A JSON-serialized object for a custom reply keyboard,
	// instructions to hide keyboard or to force a reply from the user.
	ReplyMarkup ReplyMarkup
	// Sends the message silently.
	// iOS users will not receive a notification,
	// Android users will receive a notification with no sound.
	// Other apps coming soon.
	DisableNotification bool
}

// Values returns url.Values representation of BaseMessage
func (m BaseMessage) Values() (url.Values, error) {
	v, err := m.BaseChat.Values()
	if err != nil {
		return nil, err
	}

	if m.ReplyToMessageID != 0 {
		v.Add(
			"reply_to_message_id",
			strconv.FormatInt(m.ReplyToMessageID, 10),
		)
	}

	if m.ReplyMarkup != nil {
		data, err := json.Marshal(m.ReplyMarkup)
		if err != nil {
			return nil, err
		}
		v.Add("reply_markup", string(data))
	}
	if m.DisableNotification {
		v.Add(
			"disable_notification",
			strconv.FormatBool(m.DisableNotification),
		)
	}

	return v, nil
}

// Message returns instance of *Message type.
func (BaseMessage) Message() *Message {
	return &Message{}
}

// MessageCfg contains information about a SendMessage request.
// Use it to send text messages.
// Implements Messenger interface.
type MessageCfg struct {
	BaseMessage
	Text string
	// Send Markdown or HTML, if you want Telegram apps to show
	// bold, italic, fixed-width text or inline URLs in your bot's message.
	// Use one of constants: ModeHTML, ModeMarkdown.
	ParseMode string
	// Disables link previews for links in this message.
	DisableWebPagePreview bool
}

// Name returns method name
func (cfg MessageCfg) Name() string {
	return sendMessageMethod
}

// Values returns a url.Values representation of MessageCfg.
// Returns RequiredError if Text is empty.
func (cfg MessageCfg) Values() (url.Values, error) {
	v, err := cfg.BaseMessage.Values()
	if err != nil {
		return nil, err
	}
	if cfg.Text == "" {
		return nil, NewRequiredError("Text")
	}
	v.Add("text", cfg.Text)
	if cfg.DisableWebPagePreview {
		v.Add(
			"disable_web_page_preview",
			strconv.FormatBool(cfg.DisableWebPagePreview),
		)
	}
	if cfg.ParseMode != "" {
		v.Add("parse_mode", cfg.ParseMode)
	}

	return v, nil
}

// LocationCfg contains information about a SendLocation request.
type LocationCfg struct {
	BaseMessage
	Location
}

// Values returns a url.Values representation of LocationCfg.
func (cfg LocationCfg) Values() (url.Values, error) {
	v, err := cfg.BaseMessage.Values()
	if err != nil {
		return nil, err
	}
	updateValues(v, cfg.Location.Values())

	return v, nil
}

// Name method returns Telegram API method name for sending Location.
func (cfg LocationCfg) Name() string {
	return sendLocationMethod
}

// ContactCfg contains information about a SendContact request.
// Use it to send information about a venue
// Implements Messenger interface.
type ContactCfg struct {
	BaseMessage
	Contact
}

// Name returns method name
func (cfg ContactCfg) Name() string {
	return sendContactMethod
}

// Values returns a url.Values representation of ContactCfg.
// Returns RequiredError if Text is empty.
func (cfg ContactCfg) Values() (url.Values, error) {
	v, err := cfg.BaseMessage.Values()
	if err != nil {
		return nil, err
	}
	missed := []string{}
	if cfg.Contact.FirstName == "" {
		missed = append(missed, "FirstName")
	}
	if cfg.Contact.PhoneNumber == "" {
		missed = append(missed, "PhoneNumber")
	}
	if len(missed) > 0 {
		return nil, NewRequiredError(missed...)
	}
	updateValues(v, cfg.Contact.Values())
	return v, nil
}

// VenueCfg contains information about a SendVenue request.
// Use it to send information about a venue
// Implements Messenger interface.
type VenueCfg struct {
	BaseMessage
	Venue
}

// Name returns method name
func (cfg VenueCfg) Name() string {
	return sendVenueMethod
}

// Values returns a url.Values representation of VenueCfg.
// Returns RequiredError if Text is empty.
func (cfg VenueCfg) Values() (url.Values, error) {
	v, err := cfg.BaseMessage.Values()
	if err != nil {
		return nil, err
	}
	missed := []string{}
	if cfg.Venue.Title == "" {
		missed = append(missed, "Title")
	}
	if cfg.Venue.Address == "" {
		missed = append(missed, "Address")
	}
	if len(missed) > 0 {
		return nil, NewRequiredError(missed...)
	}
	updateValues(v, cfg.Venue.Values())
	return v, nil
}

// ForwardMessageCfg contains information about a ForwardMessage request.
// Use it to forward messages of any kind
// Implements Messenger interface.
type ForwardMessageCfg struct {
	BaseChat
	// Unique identifier for the chat where the original message was sent
	FromChat BaseChat
	// Unique message identifier
	MessageID int64
	// Sends the message silently.
	// iOS users will not receive a notification,
	// Android users will receive a notification with no sound.
	// Other apps coming soon.
	DisableNotification bool
}

// Message returns instance of *Message type.
func (ForwardMessageCfg) Message() *Message {
	return &Message{}
}

// Name returns method name
func (cfg ForwardMessageCfg) Name() string {
	return forwardMessageMethod
}

// Values returns a url.Values representation of MessageCfg.
// Returns RequiredError if Text is empty.
func (cfg ForwardMessageCfg) Values() (url.Values, error) {
	v, err := cfg.BaseChat.Values()
	if err != nil {
		return nil, err
	}
	from, err := cfg.FromChat.Values()
	if err != nil {
		return nil, err
	}
	updateValuesWithPrefix(v, from, "from_")
	if cfg.MessageID == 0 {
		return nil, NewRequiredError("MessageID")
	}

	v.Add("message_id", strconv.FormatInt(cfg.MessageID, 10))
	if cfg.DisableNotification {
		v.Add(
			"disable_notification",
			strconv.FormatBool(cfg.DisableNotification),
		)
	}

	return v, nil
}

type localFile struct {
	filename string
	reader   io.Reader
}

func (l localFile) Reader() io.Reader {
	return l.reader
}

func (l localFile) Name() string {
	return l.filename
}

// NewInputFile takes Reader object and returns InputFile.
func NewInputFile(filename string, r io.Reader) InputFile {
	return localFile{filename, r}
}

// NewBytesFile takes byte slice and returns InputFile.
func NewBytesFile(filename string, data []byte) InputFile {
	return NewInputFile(filename, bytes.NewBuffer(data))
}

// BaseFile describes file settings. It's an abstract type.
type BaseFile struct {
	BaseMessage
	FileID    string
	MimeType  string
	InputFile InputFile
}

// GetFileID returns fileID if it's exist
func (b BaseFile) GetFileID() string {
	return b.FileID
}

// Exist returns true if file exists on telegram servers
func (b BaseFile) Exist() bool {
	return b.FileID != ""
}

// File returns InputFile object that are used to create request
func (b BaseFile) File() InputFile {
	return b.InputFile
}

// Values returns a url.Values representation of BaseFile.
func (b BaseFile) Values() (url.Values, error) {
	v, err := b.BaseMessage.Values()
	if err != nil {
		return nil, err
	}
	if b.MimeType != "" {
		v.Add("mime_type", b.MimeType)
	}
	return v, nil
}

// Reset method removes FileID and sets new InputFile
func (b *BaseFile) Reset(i InputFile) {
	b.FileID = ""
	b.InputFile = i
}

// PhotoCfg contains information about a SendPhoto request.
// Use it to send information about a venue
// Implements Filer and Messenger interfaces.
type PhotoCfg struct {
	BaseFile
	Caption string
}

// Name returns method name
func (cfg PhotoCfg) Name() string {
	return sendPhotoMethod
}

// Values returns a url.Values representation of PhotoCfg.
func (cfg PhotoCfg) Values() (url.Values, error) {
	v, err := cfg.BaseFile.Values()
	if err != nil {
		return nil, err
	}
	if cfg.BaseFile.FileID != "" {
		v.Add(cfg.Field(), cfg.BaseFile.FileID)
	}
	if cfg.Caption != "" {
		v.Add("caption", cfg.Caption)
	}
	return v, nil
}

// Field returns name for photo file data
func (cfg PhotoCfg) Field() string {
	return photoField
}

// AudioCfg contains information about a SendAudio request.
// Use it to send information about an audio
// Implements Filer and Messenger interfaces.
type AudioCfg struct {
	BaseFile
	Duration  int
	Performer string
	Title     string
}

// Name returns method name
func (cfg AudioCfg) Name() string {
	return sendAudioMethod
}

// Values returns a url.Values representation of AudioCfg.
func (cfg AudioCfg) Values() (url.Values, error) {
	v, err := cfg.BaseFile.Values()
	if err != nil {
		return nil, err
	}
	if cfg.BaseFile.FileID != "" {
		v.Add(cfg.Field(), cfg.BaseFile.FileID)
	}
	if cfg.Duration != 0 {
		v.Add("duration", strconv.Itoa(cfg.Duration))
	}

	if cfg.Performer != "" {
		v.Add("performer", cfg.Performer)
	}
	if cfg.Title != "" {
		v.Add("title", cfg.Title)
	}
	return v, nil
}

// Field returns name for audio file data
func (cfg AudioCfg) Field() string {
	return audioField
}

// VideoCfg contains information about a SendVideo request.
// Use it to send information about a video
// Implements Filer and Messenger interfaces.
type VideoCfg struct {
	BaseFile
	Duration int
	Caption  string
}

// Name returns method name
func (cfg VideoCfg) Name() string {
	return sendVideoMethod
}

// Values returns a url.Values representation of VideoCfg.
func (cfg VideoCfg) Values() (url.Values, error) {
	v, err := cfg.BaseFile.Values()
	if err != nil {
		return nil, err
	}
	if cfg.BaseFile.FileID != "" {
		v.Add(cfg.Field(), cfg.BaseFile.FileID)
	}
	if cfg.Duration != 0 {
		v.Add("duration", strconv.Itoa(cfg.Duration))
	}

	if cfg.Caption != "" {
		v.Add("caption", cfg.Caption)
	}
	return v, nil
}

// Field returns name for video file data
func (cfg VideoCfg) Field() string {
	return videoField
}

// VoiceCfg contains information about a SendVoice request.
// Use it to send information about a venue
// Implements Filer and Messenger interfaces.
type VoiceCfg struct {
	BaseFile
	Duration int
}

// Name returns method name
func (cfg VoiceCfg) Name() string {
	return sendVoiceMethod
}

// Values returns a url.Values representation of VoiceCfg.
func (cfg VoiceCfg) Values() (url.Values, error) {
	v, err := cfg.BaseFile.Values()
	if err != nil {
		return nil, err
	}
	if cfg.BaseFile.FileID != "" {
		v.Add(cfg.Field(), cfg.BaseFile.FileID)
	}
	if cfg.BaseFile.FileID != "" {
		v.Add(cfg.Field(), cfg.BaseFile.FileID)
	}
	if cfg.Duration != 0 {
		v.Add("duration", strconv.Itoa(cfg.Duration))
	}

	return v, nil
}

// Field returns name for voice file data
func (cfg VoiceCfg) Field() string {
	return voiceField
}

// DocumentCfg contains information about a SendDocument request.
// Use it to send information about a documents
// Implements Filer and Messenger interfaces.
type DocumentCfg struct {
	BaseFile
}

// Name returns method name
func (cfg DocumentCfg) Name() string {
	return sendDocumentMethod
}

// Values returns a url.Values representation of DocumentCfg.
func (cfg DocumentCfg) Values() (url.Values, error) {
	v, err := cfg.BaseFile.Values()
	if err != nil {
		return nil, err
	}
	if cfg.BaseFile.FileID != "" {
		v.Add(cfg.Field(), cfg.BaseFile.FileID)
	}

	return v, nil
}

// Field returns name for document file data
func (cfg DocumentCfg) Field() string {
	return documentField
}

// StickerCfg contains information about a SendSticker request.
// Implements Filer and Messenger interfaces.
type StickerCfg struct {
	BaseFile
}

// Values returns a url.Values representation of StickerCfg.
func (cfg StickerCfg) Values() (url.Values, error) {
	v, err := cfg.BaseFile.Values()
	if err != nil {
		return nil, err
	}
	if cfg.BaseFile.FileID != "" {
		v.Add(cfg.Field(), cfg.BaseFile.FileID)
	}

	return v, nil
}

// Name method returns Telegram API method name for sending Sticker.
func (cfg StickerCfg) Name() string {
	return sendStickerMethod
}

// Field returns name for sticker file data
func (cfg StickerCfg) Field() string {
	return stickerField
}
