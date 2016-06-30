package telegram

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

// APIResponse is a response from the Telegram API with the result
// stored raw.
type APIResponse struct {
	Ok          bool             `json:"ok"`
	Result      *json.RawMessage `json:"result"`
	ErrorCode   int              `json:"error_code,omitempty"`
	Description string           `json:"description,omitempty"`
}

// Update object represents an incoming update.
// Only one of the optional parameters can be present in any given update
type Update struct {
	// UpdateID is the update‘s unique identifier.
	// Update identifiers start from a certain positive number
	// and increase sequentially
	UpdateID int64 `json:"update_id"`
	// Message is a new incoming message of any kind:
	// text, photo, sticker, etc. Optional.
	Message *Message `json:"message,omitempty"`
	// New version of a message that is known to the bot and was edited.
	// Optional.
	EditedMessage *Message `json:"edited_message,omitempty"`
	// InlineQuery is a new incoming inline query. Optional.
	InlineQuery *InlineQuery `json:"inline_query,omitempty"`
	// ChosenInlineResult is a result of an inline query
	// that was chosen by a user and sent to their chat partner. Optional.
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	// CallbackQuery is a new incoming callback query. Optional.
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

// HasMessage returns true if update object contains Message field
func (u Update) HasMessage() bool {
	return u.Message != nil
}

// IsEdited returns true if update object contains EditedMessage field
func (u Update) IsEdited() bool {
	return u.EditedMessage != nil
}

// From takes User from Message, CallbackQuery, InlineQuery or ChosenInlineResult
func (u Update) From() (from *User) {
	switch {
	case u.Message != nil:
		from = u.Message.From
	case u.EditedMessage != nil:
		from = u.Message.From
	case u.CallbackQuery != nil:
		from = u.CallbackQuery.From
	case u.InlineQuery != nil:
		from = &u.InlineQuery.From
	case u.ChosenInlineResult != nil:
		from = &u.ChosenInlineResult.From
	}
	return from
}

// Chat takes chat from Message and CallbackQuery
func (u Update) Chat() (chat *Chat) {
	switch {
	case u.Message != nil:
		chat = &u.Message.Chat
	case u.EditedMessage != nil:
		chat = &u.EditedMessage.Chat
	case u.CallbackQuery != nil && u.CallbackQuery.Message != nil:
		chat = &u.CallbackQuery.Message.Chat
	}
	return chat
}

// Message object represents a message.
type Message struct {
	// MessageID is a unique message identifier.
	MessageID int64 `json:"message_id"`

	// From is a sender, can be empty for messages sent to channels.
	// Optional.
	From *User `json:"from,omitempty"`
	// Date  the message was sent in Unix time.
	Date int `json:"date"`
	// Chat is a conversation the message belongs to.
	Chat Chat `json:"chat"`

	// ForwardFrom is a sender of the original message
	// for forwarded messages. Optional.
	ForwardFrom *User `json:"forward_from,omitempty"`

	// For messages forwarded from a channel,
	// information about the original channel. Optional.
	ForwardFromChat *Chat `json:"forward_from_chat,omitempty"`

	// ForwardDate is a unixtime of the original message
	// for forwarded messages. Optional.
	ForwardDate int `json:"forward_date"`

	// ReplyToMessage is an original message for replies.
	// Note that the Message object in this field will not
	// contain further ReplyToMessage fields even if it
	// itself is a reply. Optional.
	ReplyToMessage *Message `json:"reply_to_message,omitempty"`

	// Date the message was last edited in Unix time
	// Zero value means object wasn't edited.
	// Optional.
	EditDate int `json:"edit_date,omitempty"`

	// For text messages, special entities like usernames,
	// URLs, bot commands, etc. that appear in the text. Optional
	Entities []MessageEntity `json:"entities"`

	// Text is an actual UTF-8 text of the message for a text message,
	// 0-4096 characters. Optional.
	Text string `json:"text,omitempty"`

	// Audio has information about the audio file. Optional.
	Audio *Audio `json:"audio,omitempty"`

	// Document has information about a general file. Optional.
	Document *Document `json:"document,omitempty"`

	// Photo has a slice of available sizes of photo. Optional.
	Photo []PhotoSize `json:"photo,omitempty"`

	// Sticker has information about the sticker. Optional.
	Sticker *Sticker `json:"sticker,omitempty"`

	// For a video, information about it.
	Video *Video `json:"video,omitempty"`

	// Message is a voice message, information about the file
	Voice *Voice `json:"voice,omitempty"`

	// Caption for the document, photo or video, 0‐200 characters
	Caption string `json:"caption,omitempty"`

	// For a contact, contact information itself.
	Contact *Contact `json:"contact,omitempty"`

	// For a location, its longitude and latitude.
	Location *Location `json:"location,omitempty"`

	// Message is a venue, information about the venue
	Venue *Venue `json:"venue,omitempty"`

	// NewChatMember has an information about a new member
	// that was added to the group
	// (this member may be the bot itself). Optional.
	NewChatMember *User `json:"new_chat_member,omitempty"`

	// LeftChatMember has an information about a member
	// that was removed from the group
	// (this member may be the bot itself). Optional.
	LeftChatMember *User `json:"left_chat_member,omitempty"`

	// For a service message, represents a new title
	// for chat this message came from.
	//
	// Sender would lead to a User, capable of change.
	NewChatTitle string `json:"new_chat_title,omitempty"`

	// For a service message, represents all available
	// thumbnails of new chat photo.
	//
	// Sender would lead to a User, capable of change.
	NewChatPhoto []PhotoSize `json:"new_chat_photo,omitempty"`

	// For a service message, true if chat photo just
	// got removed.
	//
	// Sender would lead to a User, capable of change.
	DeleteChatPhoto bool `json:"delete_chat_photo,omitempty"`

	// For a service message, true if group has been created.
	//
	// You would receive such a message if you are one of
	// initial group chat members.
	//
	// Sender would lead to creator of the chat.
	GroupChatCreated bool `json:"group_chat_created,omitempty"`

	// For a service message, true if super group has been created.
	//
	// You would receive such a message if you are one of
	// initial group chat members.
	//
	// Sender would lead to creator of the chat.
	SuperGroupChatCreated bool `json:"supergroup_chat_created,omitempty"`

	// For a service message, true if channel has been created.
	//
	// You would receive such a message if you are one of
	// initial channel administrators.
	//
	// Sender would lead to creator of the chat.
	ChannelChatCreated bool `json:"channel_chat_created,omitempty"`

	// For a service message, the destination (super group) you
	// migrated to.
	//
	// You would receive such a message when your chat has migrated
	// to a super group.
	//
	// Sender would lead to creator of the migration.
	MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"`

	// For a service message, the Origin (normal group) you migrated
	// from.
	//
	// You would receive such a message when your chat has migrated
	// to a super group.
	//
	// Sender would lead to creator of the migration.
	MigrateFromChatID int64 `json:"migrate_from_chat_id,omitempty"`
	// Specified message was pinned.
	// Note that the Message object in this field
	// will not contain further reply_to_message fields
	// even if it is itself a reply.
	PinnedMessage *Message `json:"pinned_message,omitempty"`
}

// IsCommand returns true if message starts with '/'.
func (m *Message) IsCommand() bool {
	return m.Text != "" && m.Text[0] == '/'
}

// Command checks if the message was a command and if it was,
// returns the command and agrument.
// If the Message was not a command, it returns an empty strings.
//
// If the command contains the at bot syntax, it removes the bot name.
func (m *Message) Command() (string, string) {
	arg := ""
	if !m.IsCommand() {
		return "", arg
	}
	splits := strings.SplitN(m.Text, " ", 2)
	command := splits[0][1:]
	if len(splits) == 2 {
		arg = splits[1]
	}

	if i := strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}

	return command, arg
}

// EditResult is an option type, because telegram may return bool or Message
type EditResult struct {
	Message *Message
	Ok      bool
}

// UnmarshalJSON helps to parse EditResult.
// On success, if edited message is sent by the bot,
// the edited Message is returned, otherwise True is returned.
func (e *EditResult) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &e.Ok)
	if err == nil {
		return nil
	}
	msg := Message{}
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}
	e.Message = &msg
	return nil
}

// MessageEntity object represents one special entity in a text message.
// For example, hashtags, usernames, URLs, etc.
type MessageEntity struct {
	// Type of the entity. One of mention ( @username ), hashtag,
	// bot_command, url, email, bold (bold text),
	// italic (italic text), code (monowidth string),
	// pre (monowidth block), text_link (for clickable text URLs),
	// text_mention (for users without usernames)
	// Use constants SomethingEntityType instead of string.
	Type string `json:"type"`
	// Offset in UTF‐16 code units to the start of the entity
	Offset int `json:"offset"`
	// Length of the entity in UTF‐16 code units
	Length int `json:"length"`
	// For “text_link” only, url that will be opened
	// after user taps on the text. Optional
	URL string `json:"url,omitempty"`
	// For “text_mention” only, the mentioned user. Optional.
	User *User `json:"user,omitempty"`
}

// User object represents a Telegram user or bot.
//
// object represents a group chat if Title is empty.
type User struct {
	// ID is a unique identifier for this user or bot.
	ID int64 `json:"id"`
	// FirstName is a user‘s or bot’s first name
	FirstName string `json:"first_name"`
	// LastName is a user‘s or bot’s last name. Optional.
	LastName string `json:"last_name,omitempty"`
	// Username is a user‘s or bot’s username. Optional.
	Username string `json:"username,omitempty"`
}

// Chat object represents a Telegram user, bot or group chat.
// Title for channels and group chats
type Chat struct {
	// ID is a Unique identifier for this chat, not exceeding 1e13 by absolute value.
	ID int64 `json:"id"`
	// Type of chat, can be either “private”, “group”, "supergroup" or “channel”
	Type string `json:"type"`

	Title     string `json:"title,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

// ChatMember object contains information about one member of the chat.
type ChatMember struct {
	// Information about the user.
	User User `json:"user"`
	// The member's status in the chat.
	// One of MemberStatus constants.
	Status string `json:"status"`
}

// MetaFile represents meta information about file.
type MetaFile struct {
	// FileID is a Unique identifier for this file.
	FileID string `json:"file_id"`
	// FileSize is a size of file if known. Optional.
	FileSize int `json:"file_size,omitempty"`
}

// File object represents any sort of file.
// The file can be downloaded via the Link.
// It is guaranteed that the link will be valid for at least 1 hour.
// When the link expires, a new one can be requested by calling GetFile.
// Maximum file size to download is 20 MB.
type File struct {
	MetaFile
	// FilePath is a relative path to file.
	// Use https://api.telegram.org/file/bot<token>/<file_path>
	// to get the file.
	FilePath string `json:"file_path,omitempty"`

	// Link is inserted by Api client after GetFile request
	Link string `json:"link"`
}

// Size object represent size information.
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// PhotoSize object represents one size of a photo
// or a file / sticker thumbnail.
type PhotoSize struct {
	MetaFile
	Size
}

// Audio object represents an audio file
// to be treated as music by the Telegram clients.
type Audio struct {
	MetaFile

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`
	//  Performer of the audio as defined by sender
	// or by audio tags. Optional.
	Performer string `json:"performer,omitempty"`
	// Title of the audio as defined by sender
	// or by audio tags. Optional.
	Title string `json:"title,omitempty"`
	// MIMEType of the file as defined by sender. Optional.
	MIMEType string `json:"mime_type,omitempty"`
}

// Document object represents a general file (as opposed to Photo or Audio).
// Telegram users can send files of any type of up to 1.5 GB in size.
type Document struct {
	MetaFile

	// Document thumbnail as defined by sender. Optional.
	Thumb *PhotoSize `json:"thumb,omitempty"`

	// Original filename as defined by sender. Optional.
	FileName string `json:"file_name,omitempty"`

	// MIMEType of the file as defined by sender. Optional.
	MIMEType string `json:"mime_type,omitempty"`
}

// Sticker object represents a WebP image, so-called sticker.
type Sticker struct {
	MetaFile
	Size // Sticker width and height

	// Sticker thumbnail in .webp or .jpg format. Optional.
	Thumb *PhotoSize `json:"thumb,omitempty"`
	// Emoji associated with the sticker. Optional.
	Emoji string `json:"emoji,omitempty"`
}

// Video object represents an MP4-encoded video.
type Video struct {
	MetaFile
	Size

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`
	// MIMEType of the file as defined by sender. Optional.
	MIMEType string `json:"mime_type,omitempty"`
	// Video thumbnail. Optional.
	Thumb *PhotoSize `json:"thumb,omitempty"`
}

// Voice object represents a voice note.
type Voice struct {
	MetaFile

	// Duration of the recording in seconds as defined by sender.
	Duration int `json:"duration"`
	// MIMEType of the file as defined by sender. Optional.
	MIMEType string `json:"mime_type,omitempty"`
}

// Contact object represents a phone contact of Telegram user
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`

	// UserID is a contact's user identifier in Telegram. Optional.
	UserID   int64  `json:"user_id,omitempty"`
	LastName string `json:"last_name,omitempty"`
}

// Values returns a url.Values representation of Contact object.
func (c Contact) Values() url.Values {
	v := url.Values{
		"phone_number": {c.PhoneNumber},
		"first_name":   {c.FirstName},
	}
	if c.UserID != 0 {
		v.Add("user_id", strconv.FormatInt(c.UserID, 10))
	}
	if c.LastName != "" {
		v.Add("last_name", c.LastName)
	}
	return v

}

// Location object represents geographic position.
type Location struct {
	// Longitude as defined by sender
	Longitude float64 `json:"longitude"`
	// Latitude as defined by sender
	Latitude float64 `json:"latitude"`
}

// Values returns a url.Values representation of Location object.
func (l Location) Values() url.Values {
	return url.Values{
		"latitude":  {strconv.FormatFloat(l.Latitude, 'f', -1, 64)},
		"longitude": {strconv.FormatFloat(l.Longitude, 'f', -1, 64)},
	}
}

// Venue object represents a venue.
type Venue struct {
	// Venue location
	Location Location `json:"location"`
	// Name of the venue
	Title string `json:"title"`
	// Address of the venue
	Address string `json:"address"`
	// Foursquare identifier of the venue. Optional.
	FoursquareID string `json:"foursquare_id,omitempty"`
}

// Values returns a url.Values representation of Venue object.
func (venue Venue) Values() url.Values {
	v := venue.Location.Values()
	v.Add("title", venue.Title)
	v.Add("address", venue.Address)
	if venue.FoursquareID != "" {
		v.Add("foursquare_id", venue.FoursquareID)
	}
	return v
}

// UserProfilePhotos contains a set of user profile pictures.
type UserProfilePhotos struct {
	// Total number of profile pictures the target user has
	TotalCount int `json:"total_count"`
	// Requested profile pictures (in up to 4 sizes each)
	Photos [][]PhotoSize `json:"photos"`
}

// CallbackQuery represents an incoming callback query
// from a callback button in an inline keyboard.
// If the button that originated the query
// was attached to a message sent by the bot,
// the field message will be presented.
// If the button was attached to a message
// sent via the bot (in inline mode),
// the field inline_message_id will be presented.
type CallbackQuery struct {
	// Unique identifier for this query
	ID string `json:"id"`
	// Sender
	From *User `json:"from"`
	// Message with the callback button that originated the query.
	// Note that message content and message date
	// will not be available if the message is too old. Optional.
	Message *Message `json:"message,omitempty"`
	// Identifier of the message sent via the bot in inline mode,
	// that originated the query. Optional.
	InlineMessageID string `json:"inline_message_id,omitempty"`
	// Data associated with the callback button.
	// Be aware that a bad client can send arbitrary data in this field
	Data string `json:"data"`
}

// ======= Markups

// A MarkReplyMarkup implements ReplyMarkup interface.
// You can mark your structures with this object.
type MarkReplyMarkup struct{}

// ReplyMarkup is a fake method that helps to identify implementations
func (MarkReplyMarkup) ReplyMarkup() {}

// KeyboardButton object represents one button of the reply keyboard.
// Optional fields are mutually exclusive.
//
// Note: request_contact and request_location options
// will only work in Telegram versions released after 9 April, 2016.
// Older clients will ignore them.
type KeyboardButton struct {
	// Text of the button. If none of the optional fields are used,
	// it will be sent to the bot as a message when the button is pressed
	Text string `json:"text"`
	// If true, the user's phone number will be sent as a contact
	// when the button is pressed. Available in private chats only.
	// Optional.
	RequestContact bool `json:"request_contact,omitempty"`
	// If true, the user's current location will be sent
	// when the button is pressed. Available in private chats only.
	// Optional.
	RequestLocation bool `json:"request_location,omitempty"`
}

// ReplyKeyboardMarkup represents a custom keyboard with reply options.
// Implements ReplyMarkup interface.
type ReplyKeyboardMarkup struct {
	MarkReplyMarkup

	// Array of button rows, each represented by an Array of Strings
	Keyboard [][]KeyboardButton `json:"keyboard"`
	// Requests clients to resize the keyboard vertically
	// for optimal fit (e.g., make the keyboard smaller
	// if there are just two rows of buttons).
	// Defaults to false, in which case the custom keyboard
	// is always of the same height as the app's standard keyboard.
	ResizeKeyboard bool `json:"resize_keyboard,omitempty"`
	// Requests clients to hide the keyboard as soon as it's been used.
	// The keyboard will still be available,
	// but clients will automatically display the usual
	// letter‐keyboard in the chat – the user can press
	// a special button in the input field to see the custom keyboard again.
	// Defaults to false.
	OneTimeKeyboard bool `json:"one_time_keyboard,omitempty"`
	// Use this parameter if you want to show the keyboard
	// to specific users only.
	// Targets:
	// 1) users that are @mentioned in the text of the Message object;
	// 2) if the bot's message is a reply (has reply_to_message_id),
	// sender of the original message.
	Selective bool `json:"selective,omitempty"`
}

// ReplyKeyboardHide tells Telegram clients to hide the current
// custom keyboard and display the default letter-keyboard.
// Implements ReplyMarkup interface.
type ReplyKeyboardHide struct {
	MarkReplyMarkup

	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"` // optional
}

// ForceReply allows the Bot to have users directly reply to it without
// additional interaction.
// Implements ReplyMarkup interface.
type ForceReply struct {
	MarkReplyMarkup

	ForceReply bool `json:"force_reply"`
	Selective  bool `json:"selective"` // optional
}

// InlineKeyboardButton object represents one button of an inline keyboard.
// You must use exactly one of the optional fields.
//
// Note: This will only work in Telegram versions
// released after 9 April, 2016.
// Older clients will display unsupported message.
type InlineKeyboardButton struct {
	// Label text on the button
	Text string `json:"text"`
	// HTTP url to be opened when button is pressed. Optional.
	URL string `json:"url,omitempty"`
	// Data to be sent in a callback query to the bot
	// when button is pressed. Optional.
	CallbackData string `json:"callback_data,omitempty"`
	// If set, pressing the button will prompt the user
	// to select one of their chats, open that chat
	// and insert the bot‘s username and the specified inline query
	// in the input field. Can be empty,
	// in which case just the bot’s username will be inserted.
	// Optional.
	//
	// Note: This offers an easy way for users
	// to start using your bot in inline mode
	// when they are currently in a private chat with it.
	// Especially useful when combined with switch_pm... actions
	// – in this case the user will be automatically returned to the chat
	// they switched from, skipping the chat selection screen.
	SwitchInlineQuery string `json:"switch_inline_query,omitempty"`
}

// InlineKeyboardMarkup object represents an inline keyboard
// that appears right next to the message it belongs to.
//
// Warning: Inline keyboards are currently being tested
// and are only available in one‐on‐one chats
// (i.e., user‐bot or user‐user in the case of inline bots).
// Note: This will only work in Telegram versions
// released after 9 April, 2016.
// Older clients will display unsupported message.
type InlineKeyboardMarkup struct {
	MarkReplyMarkup

	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}
