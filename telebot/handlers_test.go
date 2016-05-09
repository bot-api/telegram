package telebot_test

import (
	"testing"

	"github.com/bot-api/telegram/telebot"
	"golang.org/x/net/context"
	"gopkg.in/stretchr/testify.v1/assert"
)

func TestEmpty(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, telebot.Empty(ctx))
}

func TestEmptyHandler(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, telebot.EmptyHandler().Handle(ctx))
}
