package telebot_test

import (
	"fmt"
	"testing"
	"bytes"

	"github.com/bot-api/telegram/telebot"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestRecoverWithConfig(t *testing.T) {
	var (
		actCause error
		actStack string
	)
	ctx := context.Background()
	expErr := fmt.Errorf("expected error")

	buf := &bytes.Buffer{}
	telebot.DefaultRecoverLogger.SetOutput(buf)

	errHandler := telebot.HandlerFunc(func(ctx context.Context) error {
		return expErr
	})

	logFunc := func(ctx context.Context, cause error, stack []byte) {
		actCause = cause
		actStack = string(stack) // we must copy stack
	}

	m := telebot.Recover()
	{
		err := m(errHandler).Handle(ctx)
		// no panic just passed
		assert.Equal(t, expErr, err)
		assert.Equal(t, 0, buf.Len())
	}
	{
		err := m(telebot.HandlerFunc(func(ctx context.Context) error {
			panic("whatever")
		})).Handle(ctx)
		// no panic just passed, recovery info printed to default logger
		assert.Nil(t, err)
		assert.NotEqual(t, 0, buf.Len())
		stack := buf.String()
		assert.Contains(t, stack, "whatever")
		assert.Contains(t, stack, "bot.RecoverWithConfig")
		assert.Contains(t, stack, "recover_middleware.go")
	}

	m = telebot.RecoverWithConfig(telebot.RecoverCfg{
		LogFunc: logFunc,
	})
	{
		err := m(telebot.HandlerFunc(func(ctx context.Context) error {
			panic(expErr)
		})).Handle(ctx)
		// no panic just passed
		assert.Nil(t, err)
		assert.Equal(t, expErr, actCause)
		assert.Contains(t, actStack, "bot.RecoverWithConfig")
		assert.Contains(t, actStack, "recover_middleware.go")

	}

}
