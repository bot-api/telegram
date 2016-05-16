package telegram

// InlineQuery is an incoming inline query. When the user sends
// an empty query, your bot could return some default or
// trending results.
type InlineQuery struct {
	// ID is a unique identifier for this query.
	ID string `json:"id"`
	// From is a sender.
	From User `json:"from"`
	// Sender location, only for bots that request user location.
	// Optional.
	Location *Location `json:"location,omitempty"`
	// Query is a text of the query.
	Query string `json:"query"`
	// Offset of the results to be returned, can be controlled by the bot.
	Offset string `json:"offset"`
}

// ChosenInlineResult represents a result of an inline query
// that was chosen by the user and sent to their chat partner
type ChosenInlineResult struct {
	// ResultID is a unique identifier for the result that was chosen.
	ResultID string `json:"result_id"`
	// From is a user that chose the result.
	From User `json:"from"`
	// Query is used to obtain the result.
	Query string `json:"query"`
}

// A MarkInlineQueryResult implements InlineQueryResult interface.
// You can mark your structures with this object.
type MarkInlineQueryResult struct{}

// InlineQueryResult is a fake method that helps to identify implementations
func (MarkInlineQueryResult) InlineQueryResult() {}

// BaseInlineQueryResult is a base class for InlineQueryResult
type BaseInlineQueryResult struct {
	Type                string              `json:"type"` // required
	ID                  string              `json:"id"`   // required
	InputMessageContent InputMessageContent `json:"input_message_content,omitempty"`
	// ReplyMarkup supports only InlineKeyboardMarkup for InlineQueryResult
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

// InlineThumb struct helps to describe thumbnail.
type InlineThumb struct {
	ThumbURL    string `json:"thumb_url,omitempty"`
	ThumbWidth  int    `json:"thumb_width,omitempty"`
	ThumbHeight int    `json:"thumb_height,omitempty"`
}

// InlineQueryResultArticle is an inline query response article.
type InlineQueryResultArticle struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	Title string `json:"title"` // required
	URL   string `json:"url,omitempty"`
	// Optional. Pass True, if you don't want the URL
	// to be shown in the message
	HideURL     bool   `json:"hide_url,omitempty"`
	Description string `json:"description,omitempty"`
}

// InlineQueryResultPhoto is an inline query response photo.
type InlineQueryResultPhoto struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	PhotoURL    string `json:"photo_url"` // required
	PhotoWidth  int    `json:"photo_width,omitempty"`
	PhotoHeight int    `json:"photo_height,omitempty"`
	MimeType    string `json:"mime_type,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Caption     string `json:"caption,omitempty"`
}

// InlineQueryResultGIF is an inline query response GIF.
type InlineQueryResultGIF struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	// A valid URL for the GIF file. File size must not exceed 1MB
	GifURL    string `json:"gif_url"` // required
	GifWidth  int    `json:"gif_width,omitempty"`
	GifHeight int    `json:"gif_height,omitempty"`
	Title     string `json:"title,omitempty"`
	Caption   string `json:"caption,omitempty"`
}

// InlineQueryResultMPEG4GIF is an inline query response MPEG4 GIF.
type InlineQueryResultMPEG4GIF struct {
	MarkInlineQueryResult
	BaseInlineQueryResult

	MPEG4URL    string `json:"mpeg4_url"` // required
	MPEG4Width  int    `json:"mpeg4_width,omitempty"`
	MPEG4Height int    `json:"mpeg4_height,omitempty"`
	Title       string `json:"title,omitempty"`
	Caption     string `json:"caption,omitempty"`
}

// InlineQueryResultVideo is an inline query response video.
type InlineQueryResultVideo struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	VideoURL      string `json:"video_url"` // required
	MimeType      string `json:"mime_type"` // required
	Title         string `json:"title,omitempty"`
	Caption       string `json:"caption,omitempty"`
	VideoWidth    int    `json:"video_width,omitempty"`
	VideoHeight   int    `json:"video_height,omitempty"`
	VideoDuration int    `json:"video_duration,omitempty"`
	Description   string `json:"description,omitempty"`
}

// InlineQueryResultAudio is an inline query response audio.
type InlineQueryResultAudio struct {
	MarkInlineQueryResult
	BaseInlineQueryResult

	AudioURL      string `json:"audio_url"` // required
	Title         string `json:"title"`     // required
	Performer     string `json:"performer"`
	AudioDuration int    `json:"audio_duration"`
}

// InlineQueryResultVoice is an inline query response voice.
type InlineQueryResultVoice struct {
	MarkInlineQueryResult
	BaseInlineQueryResult

	VoiceURL string `json:"voice_url"` // required
	Title    string `json:"title"`     // required
	Duration int    `json:"voice_duration,omitempty"`
}

// InlineQueryResultDocument is an inline query response document.
type InlineQueryResultDocument struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	DocumentURL string `json:"document_url"` // required
	// Mime type of the content of the file,
	// either “application/pdf” or “application/zip”
	MimeType string `json:"mime_type"` // required
	// Title for the result
	Title string `json:"title"` // required
	// Optional. Caption of the document to be sent, 0-200 characters
	Caption string `json:"caption,omitempty"`
	// Optional. Short description of the result
	Description string `json:"description,omitempty"`
}

// InlineQueryResultLocation is an inline query response location.
type InlineQueryResultLocation struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb

	Latitude  float64 `json:"latitude"`  // required
	Longitude float64 `json:"longitude"` // required
	Title     string  `json:"title"`     // required
}

// InlineQueryResultContact represents a contact with a phone number.
// By default, this contact will be sent by the user.
// Alternatively, you can use input_message_content
// to send a message with the specified content instead of the contact.
type InlineQueryResultContact struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb
	Contact
}

// InlineQueryResultVenue represents a venue.
// By default, the venue will be sent by the user.
// Alternatively, you can use input_message_content
// to send a message with the specified content instead of the venue.
type InlineQueryResultVenue struct {
	MarkInlineQueryResult
	BaseInlineQueryResult
	InlineThumb
	Venue
}

// A MarkInputMessageContent implements InputMessageContent interface.
// You can mark your structures with this object.
type MarkInputMessageContent struct{}

// InputMessageContent is a fake method that helps to identify implementations
func (MarkInputMessageContent) InputMessageContent() {}

// InputTextMessageContent represents the content of a text message
// to be sent as the result of an inline query.
type InputTextMessageContent struct {
	MarkInputMessageContent

	// Text of the message to be sent, 1‐4096 characters
	MessageText string `json:"message_text"`
	// Send Markdown or HTML, if you want Telegram apps to show
	// bold, italic, fixed‐width text or inline URLs in your bot's message.
	// Use Mode constants. Optional.
	ParseMode string `json:"parse_mode,omitempty"`
	// Disables link previews for links in this message.
	DisableWebPagePreview bool `json:"disable_web_page_preview,omitempty"`
}

// InputLocationMessageContent contains a location for displaying
// as an inline query result.
// Implements InputMessageContent
type InputLocationMessageContent struct {
	MarkInputMessageContent
	Location
}

// InputVenueMessageContent contains a venue for displaying
// as an inline query result.
// Implements InputMessageContent
type InputVenueMessageContent struct {
	MarkInputMessageContent
	Venue
}

// InputContactMessageContent contains a contact for displaying
// as an inline query result.
// Implements InputMessageContent
type InputContactMessageContent struct {
	MarkInputMessageContent
	Contact
}
