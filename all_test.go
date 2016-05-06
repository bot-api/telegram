package telegram_test

// all_test contains common test structures and interfaces

import (
	"errors"
	"net/url"
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

type replyBadMarkup struct{}

var marshalError = errors.New("Can't be marshalled")

func (m replyBadMarkup) Markup() (string, error) {
	return "", marshalError
}
