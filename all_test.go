package telegram_test

// all_test contains common test structures and interfaces

import (
	"errors"
	"net/url"

	"github.com/bot-api/telegram"
)

type valuesI interface {
	Values() (url.Values, error)
}

// cfgTT is a configs test table structure
type cfgTT struct {
	exp    url.Values
	expErr error
	cfg    valuesI
}

type replyBadMarkup struct {
	telegram.MarkReplyMarkup
}

var marshalError = errors.New("Can't be marshalled")

func (m replyBadMarkup) MarshalJSON() ([]byte, error) {
	return nil, marshalError
}

type badInlineQueryResult struct {
	telegram.MarkInlineQueryResult
}

func (m badInlineQueryResult) MarshalJSON() ([]byte, error) {
	return nil, marshalError
}
