package main

import (
	"log"
	"flag"

	"github.com/bot-api/telegram"
	"golang.org/x/net/context"
)


func main() {
	token := flag.String("token", "", "telegram bot token")
	debug := flag.Bool("debug", false, "show debug information")
	flag.Parse()

	if *token == "" {
		log.Fatal("token flag required")
	}

	api := telegram.New(*token)
	api.Debug(*debug)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if user, err := api.GetMe(ctx); err != nil {
		log.Panic(err)
	} else {
		log.Printf("bot info: %#v", user)
	}

	updatesCh := make(chan telegram.Update)

	go telegram.GetUpdates(ctx, api, telegram.UpdateCfg{
		Timeout: 10, 	// Timeout in seconds for long polling.
		Offset: 0, 	// Start with the oldest update
	}, updatesCh)

	for update := range updatesCh {
		log.Printf("got update from %s", update.Message.From.Username)
		if update.Message == nil {
			continue
		}
		msg := telegram.CloneMessage(update.Message, nil)
		// echo with the same message
		if _, err := api.Send(ctx, msg); err != nil {
			log.Print("send error: %v", err)
		}
	}
}