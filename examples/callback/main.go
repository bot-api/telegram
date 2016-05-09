package main

// Callback example shows how to use callback query and how to edit bot message

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"golang.org/x/net/context"
)

func main() {
	token := flag.String("token", "", "telegram bot token")
	debug := flag.Bool("debug", false, "show debug information")
	flag.Parse()

	if *token == "" {
		log.Fatal("token flag is required")

	}

	api := telegram.New(*token)
	api.Debug(*debug)
	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover()) // recover if handler panics

	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.HandleFunc(func(ctx context.Context) error {
		update := telebot.GetUpdate(ctx) // take update from context
		api := telebot.GetAPI(ctx)       // take api from context

		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			if strings.HasPrefix(data, "sex:") {
				cfg := telegram.NewEditMessageText(
					update.Chat().ID,
					update.CallbackQuery.Message.MessageID,
					fmt.Sprintf("You sex: %s", data[4:]),
				)
				_, err := api.AnswerCallbackQuery(
					ctx,
					telegram.NewAnswerCallback(
						update.CallbackQuery.ID,
						"Your configs changed",
					),
				)
				if err != nil {
					return err
				}
				_, err = api.EditMessageText(ctx, cfg)
				return err
			}
		}

		msg := telegram.NewMessage(update.Chat().ID,
			"Your sex:")
		msg.ReplyMarkup = telegram.InlineKeyboardMarkup{
			InlineKeyboard: telegram.NewVInlineKeyboard(
				"sex:",
				[]string{"Female", "Male"},
				[]string{"female", "male"},
			),
		}
		_, err := api.SendMessage(ctx, msg)
		return err

	})

	err := bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}
}
