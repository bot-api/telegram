package main

// Simple proxy bot, that proxy messages from one user to another

import (
	"flag"
	"fmt"
	"log"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"golang.org/x/net/context"
)

func main() {
	token := flag.String("token", "", "telegram bot token")
	user1 := flag.Int("user1", 0, "first user")
	user2 := flag.Int("user2", 0, "second user")
	debug := flag.Bool("debug", false, "show debug information")
	flag.Parse()

	if *token == "" {
		log.Fatal("token flag is required")
	}
	if !telegram.IsValidToken(*token) {
		log.Fatal("token has wrong format")
	}

	api := telegram.New(*token)
	api.Debug(*debug)
	bot := telebot.NewWithAPI(api)

	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.Use(telebot.Recover())

	bot.HandleFunc(func(ctx context.Context) error {
		update := telebot.GetUpdate(ctx)
		if update.Message == nil {
			return nil
		}
		api := telebot.GetAPI(ctx)

		userTo := *user1
		if update.Message.Chat.ID == int64(userTo) {
			userTo = *user2
		}
		msg := telegram.CloneMessage(
			update.Message,
			&telegram.BaseMessage{
				BaseChat: telegram.BaseChat{
					ID: int64(userTo),
				},
				ReplyToMessageID: 0,
			},
		)
		if msg == nil {
			return fmt.Errorf("can't clone message")
		}
		_, err := api.Send(ctx, msg)
		return err

	})

	err := bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}

}
