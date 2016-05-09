package main

// Command example shows how to handle commands

import (
	"flag"
	"log"

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

	// Use command middleware, that helps to work with commands
	bot.Use(telebot.Commands(map[string]telebot.Commander{
		"start": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {

				api := telebot.GetAPI(ctx)
				update := telebot.GetUpdate(ctx)
				_, err := api.SendMessage(ctx,
					telegram.NewMessagef(update.Chat().ID,
						"received start with arg %s", arg,
					))
				return err
			}),
		"": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {

				api := telebot.GetAPI(ctx)
				update := telebot.GetUpdate(ctx)
				command, arg := update.Message.Command()
				_, err := api.SendMessage(ctx,
					telegram.NewMessagef(update.Chat().ID,
						"received unrecognized"+
							" command %s with arg %s",
						command, arg,
					))
				return err
			}),
	}))

	err := bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}
}
