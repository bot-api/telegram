package telegram

import (
	"net/url"
	"regexp"
	"unicode"
)

var tokenRegex = regexp.MustCompile(`^[\d]{3,11}:[\w-]{35}$`)

// IsValidToken returns true if token is a valid telegram bot token
//
// Token format is like: 110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsawq
func IsValidToken(token string) bool {
	return tokenRegex.MatchString(token)
}

//SplitMessageString splits original message in case it's too big for telegram message and sends each piece of a split separately
func SplitMessageString(message string) <-chan string {
	out := make(chan string)
	originalMessage := []rune(message)
	originalMessageLength := len(originalMessage)

	go func() {
		previousChunkEnd, lastSuccessfulCheckpoint, curPosition := 0, 0, 0
		usedBytes, newBytes := 0, 0

		initializeChunk := func() bool {
			for ; previousChunkEnd < originalMessageLength && unicode.IsSpace(originalMessage[previousChunkEnd]); previousChunkEnd++ {
			}
			if previousChunkEnd >= originalMessageLength {
				return true
			}

			curPosition, lastSuccessfulCheckpoint = previousChunkEnd, previousChunkEnd
			usedBytes, newBytes = 0, 0
			return false
		}

		walkChars := func(checkFunc func(rune) bool) bool {
			for ; curPosition < originalMessageLength && checkFunc(originalMessage[curPosition]); curPosition++ {
				newBytes = usedBytes + len(string(originalMessage[curPosition]))
				if newBytes > MaximumMessageLength {
					return true
				}
				usedBytes = newBytes
			}
			return curPosition >= originalMessageLength
		}

		for curPosition < originalMessageLength {
			if initializeChunk() {
				break
			}

			for usedBytes <= MaximumMessageLength {
				if walkChars(unicode.IsSpace) {
					break
				}
				if walkChars(func(r rune) bool { return !unicode.IsSpace(r) }) {
					break
				}
				lastSuccessfulCheckpoint = curPosition
			}

			if lastSuccessfulCheckpoint == previousChunkEnd {
				out <- string(originalMessage[previousChunkEnd:curPosition])
				previousChunkEnd = curPosition
				continue
			}

			if curPosition >= originalMessageLength {
				lastSuccessfulCheckpoint = curPosition
			}

			out <- string(originalMessage[previousChunkEnd:lastSuccessfulCheckpoint])

			previousChunkEnd = lastSuccessfulCheckpoint
		}
		close(out)
	}()

	return out
}

// ========== Internal functions

func updateValues(to, from url.Values) {
	for key, values := range from {
		for _, value := range values {
			to.Add(key, value)
		}
	}
}

func updateValuesWithPrefix(to, from url.Values, prefix string) {
	for key, values := range from {
		for _, value := range values {
			to.Add(prefix+key, value)
		}
	}
}
