package telebot_test

import (
	"fmt"
	"testing"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestHandlerFunc_Handle(t *testing.T) {
	ctx := context.Background()
	err := fmt.Errorf("expected error")
	f := telebot.HandlerFunc(func(ctx2 context.Context) error {
		assert.Equal(t, ctx, ctx2)
		return err
	})
	assert.Equal(t, err, f.Handle(ctx))
}

func TestCommandFunc_Command(t *testing.T) {
	ctx := context.Background()
	err := fmt.Errorf("expected error")
	f := telebot.CommandFunc(func(ctx2 context.Context, arg string) error {
		assert.Equal(t, ctx, ctx2)
		assert.Equal(t, "arg", arg)
		return err
	})
	assert.Equal(t, err, f.Command(ctx, "arg"))
}

func TestCommands(t *testing.T) {
	err := fmt.Errorf("expected error")
	hErr := fmt.Errorf("handler error")
	defErr := fmt.Errorf("default error")

	f := telebot.HandlerFunc(func(context.Context) error {
		return hErr
	})
	cmdOne := false
	cmdTwo := false
	cmdDef := false

	c := telebot.Commands(map[string]telebot.Commander{
		"one": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {
				cmdOne = true
				assert.Equal(t, "oneArg", arg)
				return err
			}),
		"two": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {
				cmdTwo = true
				assert.Equal(t, "", arg)
				return nil
			}),
		// this command pass execution to a handler
		"three": nil,
		"": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {
				cmdDef = true
				assert.Equal(t, "def", arg)
				return defErr
			}),
	})
	{
		ctx := telebot.WithUpdate(context.Background(), &telegram.Update{
			Message: &telegram.Message{
				Text: "/one oneArg",
			},
		})
		assert.Equal(t, c(f).Handle(ctx), err)
	}
	{
		ctx := telebot.WithUpdate(context.Background(), &telegram.Update{
			Message: &telegram.Message{
				Text: "/two",
			},
		})
		assert.Equal(t, c(f).Handle(ctx), nil)
	}
	{
		ctx := telebot.WithUpdate(context.Background(), &telegram.Update{
			Message: &telegram.Message{
				Text: "/three",
			},
		})
		assert.Equal(t, c(f).Handle(ctx), hErr)
	}
	{
		ctx := telebot.WithUpdate(context.Background(), &telegram.Update{
			Message: &telegram.Message{
				Text: "not a command",
			},
		})
		// non commands passed directly to a handler
		assert.Equal(t, c(f).Handle(ctx), hErr)
	}
	{
		ctx := telebot.WithUpdate(context.Background(), &telegram.Update{})
		// non message update passed directly to a handler
		assert.Equal(t, c(f).Handle(ctx), hErr)
	}
	{
		ctx := telebot.WithUpdate(context.Background(), &telegram.Update{
			Message: &telegram.Message{
				Text: "/four def",
			},
		})
		// unknown commands passed to an empty commander
		assert.Equal(t, c(f).Handle(ctx), defErr)
	}
	assert.True(t, cmdOne, "command one hasn't been executed")
	assert.True(t, cmdTwo, "command two hasn't been executed")
	assert.True(t, cmdDef, "default command hasn't been executed")

	c = telebot.Commands(map[string]telebot.Commander{
		"one": telebot.CommandFunc(
			func(ctx context.Context, arg string) error {
				assert.Equal(t, "oneArg twoArg", arg)
				return err
			}),
	})
	{
		ctx := telebot.WithUpdate(context.Background(), &telegram.Update{
			Message: &telegram.Message{
				Text: "/one oneArg twoArg",
			},
		})
		assert.Equal(t, c(f).Handle(ctx), err)
	}
	{
		ctx := telebot.WithUpdate(context.Background(), &telegram.Update{
			Message: &telegram.Message{
				Text: "/two oneArg twoArg",
			},
		})
		// two commands passed directly to a handler
		assert.Equal(t, c(f).Handle(ctx), hErr)
	}
}
