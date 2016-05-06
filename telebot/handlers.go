package telebot

import "golang.org/x/net/context"

// Empty does nothing.
func Empty(context.Context) error {
	return nil
}

// EmptyHandler returns a handler that does nothing.
func EmptyHandler() Handler { return HandlerFunc(Empty) }
