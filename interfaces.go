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
	Field() string
	File() InputFile
	Exist() bool
}

// ReplyMarkup describes interface for reply_markup keyboards.
type ReplyMarkup interface {
	// Markup should return json string.
	Markup() (string, error)
}

type InputFile interface {
	Reader() io.Reader
	Name() string
}
