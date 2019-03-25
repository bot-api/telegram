package telegram_test

import (
	"bytes"
	"testing"

	"github.com/bot-api/telegram"
	"gopkg.in/stretchr/testify.v1/assert"
)

func getValuesFromChannel(in <-chan string) []string {
	var result []string
	for s := range in {
		result = append(result, s)
	}
	return result
}

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

func TestSplittingMessageEndgeOnSpace(t *testing.T) {
	var buffer bytes.Buffer
	for i := 0; i < telegram.MaximumMessageLength-1; i++ {
		buffer.WriteString("a")
	}
	firstPart := buffer.String()
	buffer.WriteString("      ")
	secondPart := "bbb ccc ddd"
	buffer.WriteString(secondPart)

	originalMessage := buffer.String()
	messages := getValuesFromChannel(telegram.SplitMessageString(originalMessage))
	assert.Len(t, messages, 2)
	assert.Equal(t, messages[0], firstPart)
	assert.Equal(t, messages[1], secondPart)
}

func TestSplittingMessageEndgeOnLetter(t *testing.T) {
	var buffer bytes.Buffer
	for i := 0; i < telegram.MaximumMessageLength-3; i++ {
		buffer.WriteString("a")
	}
	firstPart := buffer.String()
	buffer.WriteString(" ")
	secondPart := "bbb ccc ddd"
	buffer.WriteString(secondPart)

	originalMessage := buffer.String()
	messages := getValuesFromChannel(telegram.SplitMessageString(originalMessage))
	assert.Len(t, messages, 2)
	assert.Equal(t, messages[0], firstPart)
	assert.Equal(t, messages[1], secondPart)
}

func TestSplittingMessageLongWordNoSpace(t *testing.T) {
	var buffer bytes.Buffer
	for i := 0; i < telegram.MaximumMessageLength; i++ {
		buffer.WriteString("a")
	}
	firstPart := buffer.String()
	secondPart := "bbb"

	buffer.WriteString(secondPart)
	originalMessage := buffer.String()

	messages := getValuesFromChannel(telegram.SplitMessageString(originalMessage))
	assert.Len(t, messages, 2)
	assert.Equal(t, messages[0], firstPart)
	assert.Equal(t, messages[1], secondPart)
}

func TestSplittingMessageLongWordNoSpace0(t *testing.T) {
	var buffer bytes.Buffer
	for i := 0; i < telegram.MaximumMessageLength; i++ {
		buffer.WriteString("a")
	}
	firstPart := buffer.String()
	secondPart := "a"

	buffer.WriteString(secondPart)
	originalMessage := buffer.String()

	messages := getValuesFromChannel(telegram.SplitMessageString(originalMessage))
	assert.Len(t, messages, 2)
	assert.Equal(t, messages[0], firstPart)
	assert.Equal(t, messages[1], secondPart)
}

func TestSplittingMessageLongWordNoSpace2(t *testing.T) {
	var originalMessageBuffer, buffer bytes.Buffer
	for i := 0; i < telegram.MaximumMessageLength/2; i++ {
		buffer.WriteString("к")
	}
	firstPart := buffer.String()
	originalMessageBuffer.WriteString(firstPart)
	buffer.Reset()

	for i := 0; i < telegram.MaximumMessageLength; i++ {
		buffer.WriteString("a")
	}
	secondPart := buffer.String()
	originalMessageBuffer.WriteString(secondPart)
	buffer.Reset()

	for i := 0; i < telegram.MaximumMessageLength/2; i++ {
		buffer.WriteString("к")
	}
	thirdPart := buffer.String()
	originalMessageBuffer.WriteString(thirdPart)

	originalMessage := originalMessageBuffer.String()
	messages := getValuesFromChannel(telegram.SplitMessageString(originalMessage))
	assert.Len(t, messages, 3)
	assert.Equal(t, messages[0], firstPart)
	assert.Equal(t, messages[1], secondPart)
	assert.Equal(t, messages[2], thirdPart)

}

func TestSplittingMessageLongWordNoSpace3(t *testing.T) {
	var buffer bytes.Buffer

	for i := 0; i < telegram.MaximumMessageLength-1; i++ {
		buffer.WriteString("a")
	}
	firstPart := buffer.String()

	buffer.WriteString("к")
	secondPart := "к"

	originalMessage := buffer.String()
	messages := getValuesFromChannel(telegram.SplitMessageString(originalMessage))
	assert.Len(t, messages, 2)
	assert.Equal(t, messages[0], firstPart)
	assert.Equal(t, messages[1], secondPart)

}

func TestSpacesInFron(t *testing.T) {
	var buffer bytes.Buffer
	for i := 0; i < telegram.MaximumMessageLength-1; i++ {
		buffer.WriteString(" ")
	}
	firstPart := "aaaa"

	buffer.WriteString(firstPart)
	originalMessage := buffer.String()

	messages := getValuesFromChannel(telegram.SplitMessageString(originalMessage))
	assert.Len(t, messages, 1)
	assert.Equal(t, messages[0], firstPart)
}

func TestSpacesOnly(t *testing.T) {
	var buffer bytes.Buffer
	for i := 0; i < telegram.MaximumMessageLength+1; i++ {
		buffer.WriteString(" ")
	}
	originalMessage := buffer.String()

	messages := getValuesFromChannel(telegram.SplitMessageString(originalMessage))
	assert.Len(t, messages, 0)
}
