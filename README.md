# Telegram Bot Api  [![GoDoc](https://godoc.org/github.com/bot-api/telegram?status.svg)](http://godoc.org/github.com/bot-api/telegram) [![Build Status](https://travis-ci.org/bot-api/telegram.svg?branch=master)](https://travis-ci.org/bot-api/telegram)  [![Coverage Status](https://coveralls.io/repos/github/bot-api/telegram/badge.svg?branch=master)](https://coveralls.io/github/bot-api/telegram?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/bot-api/telegram)](https://goreportcard.com/report/github.com/bot-api/telegram)

Supported go version: 1.5, 1.6, tip


Implementation of the telegram bot API, inspired by github.com/go-telegram-bot-api/telegram-bot-api.

The main difference between telegram-bot-api and this version is supporting net/context.
Also, this library handles errors more correctly at this time (telegram-bot-api v4).


## Package contains:

1. Client for telegram bot api.
2. Bot with:
    1. Middleware support
        1. Command middleware to handle commands.
        2. Recover middleware to recover on panics.
    2. Webhook support



# Get started

Get last telegram api:
 
`go get github.com/bot-api/telegram`

## If you want to use telegram bot api directly:

`go run ./examples/api/main.go -debug -token BOT_TOKEN`

```go
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
			log.Printf("send error: %v", err)
		}
	}
}
```

## If you want to use bot

`go run ./examples/echo/main.go -debug -token BOT_TOKEN`

```go
package main
// Simple echo bot, that responses with the same message

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
	}))


	err := bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
	}
}
```

Use callback query and edit bot's message

`go run ./examples/callback/main.go -debug -token BOT_TOKEN`


```go

bot.HandleFunc(func(ctx context.Context) error {
    update := telebot.GetUpdate(ctx) // take update from context
    api := telebot.GetAPI(ctx) // take api from context

    if update.CallbackQuery != nil {
        data := update.CallbackQuery.Data
        if strings.HasPrefix(data, "sex:") {
            cfg := telegram.NewEditMessageText(
                update.Chat().ID,
                update.CallbackQuery.Message.MessageID,
                fmt.Sprintf("You sex: %s", data[4:]),
            )
            api.AnswerCallbackQuery(
                ctx,
                telegram.NewAnswerCallback(
                    update.CallbackQuery.ID,
                    "Your configs changed",
                ),
            )
            _, err := api.EditMessageText(ctx, cfg)
            return err
        }
    }

    msg := telegram.NewMessage(update.Chat().ID,
        "Your sex:")
    msg.ReplyMarkup = telegram.InlineKeyboardMarkup{
        InlineKeyboard: telegram.NewVInlineKeyboard(
            "sex:",
            []string{"Female", "Male",},
            []string{"female", "male",},
        ),
    }
    _, err := api.SendMessage(ctx, msg)
    return err

})


```

Take a look at `./examples/` to know more how to use bot and telegram api.


# TODO:

- [x] Handlers 
- [x] Middleware
- [x] Command middleware
- [ ] Session middleware
- [ ] Log middleware
- [ ] Menu middleware
- [ ] Examples
    - [x] Command
    - [x] CallbackAnswer
    - [x] Inline
    - [x] Proxy
    - [ ] Menu
- [x] Add travis-ci integration
- [x] Add coverage badge
- [x] Add integration tests


- [ ] Add gopkg version
- [ ] Improve documentation
- [ ] Benchmark ffjson and easyjson.
- [ ] Add GAE example. 


# Supported API methods:
- [x] getMe
- [x] sendMessage
- [x] forwardMessage
- [x] sendPhoto
- [x] sendAudio
- [x] sendDocument
- [x] sendSticker
- [x] sendVideo
- [x] sendVoice
- [x] sendLocation
- [x] sendChatAction
- [x] getUserProfilePhotos
- [x] getUpdates
- [x] setWebhook
- [x] getFile
- [x] answerInlineQuery inline bots

#  Supported API v2 methods:
- [x] sendVenue
- [x] sendContact
- [x] editMessageText
- [x] editMessageCaption
- [x] editMessageReplyMarkup
- [x] kickChatMember
- [x] unbanChatMember
- [x] answerCallbackQuery
- [x] getChat
- [x] getChatMember
- [x] getChatMembersCount
- [x] getChatAdministrators
- [x] leaveChat

# Supported Inline modes


- [x] InlineQueryResultArticle
- [x] InlineQueryResultAudio
- [x] InlineQueryResultContact
- [x] InlineQueryResultDocument
- [x] InlineQueryResultGif
- [x] InlineQueryResultLocation
- [x] InlineQueryResultMpeg4Gif
- [x] InlineQueryResultPhoto
- [x] InlineQueryResultVenue
- [x] InlineQueryResultVideo
- [x] InlineQueryResultVoice
- [ ] InlineQueryResultCachedAudio
- [ ] InlineQueryResultCachedDocument
- [ ] InlineQueryResultCachedGif
- [ ] InlineQueryResultCachedMpeg4Gif
- [ ] InlineQueryResultCachedPhoto
- [ ] InlineQueryResultCachedSticker
- [ ] InlineQueryResultCachedVideo
- [ ] InlineQueryResultCachedVoice
- [ ] InputTextMessageContent
- [ ] InputLocationMessageContent


Other bots:
I like this handler system
https://bitbucket.org/master_groosha/telegram-proxy-bot/src/07a6b57372603acae7bdb78f771be132d063b899/proxy_bot.py?fileviewer=file-view-default

