// Inline example shows how to use inline bots
package main

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
	debug := flag.Bool("debug", false, "show debug information")
	flag.Parse()

	if *token == "" {
		log.Fatal("token flag is required")
	}

	api := telegram.New(*token)
	api.Debug(*debug)
	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover()) // recover if handler panic

	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.HandleFunc(func(ctx context.Context) error {
		update := telebot.GetUpdate(ctx) // take update from context
		if update.InlineQuery == nil {
			return nil
		}
		query := update.InlineQuery
		api := telebot.GetAPI(ctx) // take api from context

		api.AnswerInlineQuery(ctx, telegram.AnswerInlineQueryCfg{
			InlineQueryID: query.ID,
			Results: []telegram.InlineQueryResult{
				telegram.NewInlineQueryResultArticle(
					"10",
					fmt.Sprintf("one for %s", query.Query),
					fmt.Sprintf("result1: %s", query.Query),
				),
				telegram.NewInlineQueryResultArticle(
					"11",
					fmt.Sprintf("two for %s", query.Query),
					fmt.Sprintf("result2: %s", query.Query),
				),
			},
			CacheTime: 10, // cached for 10 seconds

		})
		return nil

	})
	err := bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}
}
