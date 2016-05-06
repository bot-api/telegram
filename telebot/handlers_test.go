package telebot_test

import (
	"testing"

	"github.com/bot-api/telegram/telebot"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestEmpty(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, telebot.Empty(ctx))
}

func TestEmptyHandler(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, telebot.EmptyHandler().Handle(ctx))
}
