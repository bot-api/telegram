package telegram_test

import (
	"testing"

	"github.com/bot-api/telegram"
	"gopkg.in/stretchr/testify.v1/assert"
)

func TestIsValidToken(t *testing.T) {
	testTable := []struct {
		token  string
		result bool
	}{
		{"110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsawq", true},
		{"110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaq", false},
		{"113:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsawq", true},
		{"12345678901:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsawq", true},
		{"1234567890123:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsawq", false},
		{"110201543:AAHdqTcvCH1vGWJxf-eofSAs0K5PALDsawq", true},
	}
	for i, tt := range testTable {
		t.Logf("test #%d", i)
		assert.Equal(t, tt.result, telegram.IsValidToken(tt.token))
	}
}
