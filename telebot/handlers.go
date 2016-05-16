package telebot

import (
	"github.com/bot-api/telegram"
	"golang.org/x/net/context"
)

// Empty does nothing.
func Empty(context.Context) error {
	return nil
}

// EmptyHandler returns a handler that does nothing.
func EmptyHandler() Handler { return HandlerFunc(Empty) }

// StringHandler sends user a text
func StringHandler(text string) Handler {
	return HandlerFunc(func(ctx context.Context) error {
		update := GetUpdate(ctx)
		api := GetAPI(ctx)
		chat := update.Chat()
		if chat == nil {
			return nil
		}
		_, err := api.SendMessage(ctx, telegram.NewMessage(chat.ID, text))
		return err
	})
}
