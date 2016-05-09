package telegram

import (
	"net/url"
	"regexp"
)

var tokenRegex = regexp.MustCompile(`^[\d]{3,11}:[\w-]{35}$`)

// IsValidToken returns true if token is a valid telegram bot token
//
// Token format is like: 110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsawq
func IsValidToken(token string) bool {
	return tokenRegex.MatchString(token)
}

// ========== Internal functions

func updateValues(to, from url.Values) {
	for key, values := range from {
		for _, value := range values {
			to.Add(key, value)
		}
	}
}

func updateValuesWithPrefix(to, from url.Values, prefix string) {
	for key, values := range from {
		for _, value := range values {
			to.Add(prefix+key, value)
		}
	}
}
