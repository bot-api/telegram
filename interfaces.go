package telegram

import (
	"io"
	"net/url"
)

// Method describes interface for Telegram API request
//
// Every method is https://api.telegram.org/bot<token>/METHOD_NAME
// Values are passed as application/x-www-form-urlencoded for usual request
// and multipart/form-data when files are uploaded.
type Method interface {
	// method name
	Name() string
	// method params
	Values() (url.Values, error)
}

// Messenger is a virtual interface to distinct methods
// that return Message from others.BaseMessage
type Messenger interface {
	Method
	Message() *Message
}

// Filer is any config type that can be sent that includes a file.
type Filer interface {
	// Field name for file data
	Field() string
	// File data
	File() InputFile
	// Exist returns true if file exists on telegram servers
	Exist() bool
	// Reset removes FileID and sets new InputFile
	// Reset(InputFile)
	// GetFileID returns fileID if it's exist
	GetFileID() string
}

// ReplyMarkup describes interface for reply_markup keyboards.
type ReplyMarkup interface {
	// ReplyMarkup is a fake method that helps to identify implementations
	ReplyMarkup()
}

// InputFile describes interface for input files.
type InputFile interface {
	Reader() io.Reader
	Name() string
}

// InlineQueryResult interface represents one result of an inline query.
// Telegram clients currently support results of the following 19 types:
//
// - InlineQueryResultCachedAudio
// - InlineQueryResultCachedDocument
// - InlineQueryResultCachedGif
// - InlineQueryResultCachedMpeg4Gif
// - InlineQueryResultCachedPhoto
// - InlineQueryResultCachedSticker
// - InlineQueryResultCachedVideo
// - InlineQueryResultCachedVoice
// - InlineQueryResultArticle
// - InlineQueryResultAudio
// - InlineQueryResultContact
// - InlineQueryResultDocument
// - InlineQueryResultGif
// - InlineQueryResultLocation
// - InlineQueryResultMpeg4Gif
// - InlineQueryResultPhoto
// - InlineQueryResultVenue
// - InlineQueryResultVideo
// - InlineQueryResultVoice
//
type InlineQueryResult interface {
	// InlineQueryResult is a fake method that helps to identify implementations
	InlineQueryResult()
}

// InputMessageContent interface represents the content of a message
// to be sent as a result of an inline query.
// Telegram clients currently support the following 4 types:
//  - InputTextMessageContent
//  - InputLocationMessageContent
//  - InputVenueMessageContent
//  - InputContactMessageContent
type InputMessageContent interface {
	// MessageContent is a fake method that helps to identify implementations
	InputMessageContent()
}
