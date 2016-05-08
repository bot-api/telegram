package main

// Simple echo bot, that responses with the same message

import (
	"flag"
	"log"

	"github.com/bot-api/telegram"
	telebot "github.com/bot-api/telegram/telebot"
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
	bot := telebot.NewWithApi(api)
	bot.Use(telebot.Recover()) // recover if handler panic

	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.HandleFunc(func(ctx context.Context) error {
		update := telebot.GetUpdate(ctx) // take update from context
		if update.Message == nil {
			return nil
		}
		api := telebot.GetAPI(ctx) // take api from context
		msg := telegram.CloneMessage(update.Message, nil)
		_, err := api.Send(ctx, msg)
		return err

	})
	err := bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}
}
