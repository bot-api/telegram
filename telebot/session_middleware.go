package telebot

import (
	"bytes"

	"encoding/json"

	"golang.org/x/net/context"
)

type sessionKey struct{}

// SessionData describes key:value data
type SessionData map[string]interface{}

// UpdateFunc describes a func to update session data
type UpdateFunc func(data []byte) error

// SessionConfig helps to configure Session Middleware
type SessionConfig struct {
	// Encode takes SessionData and return encoded version
	Encode func(interface{}) ([]byte, error)
	// Decode takes encoded session data and fill dst struct
	Decode func(data []byte, dst interface{}) error
	// GetSession function should receive request context and return
	// []bytes of current session and UpdateFunc
	// that is invoked if session is modified
	GetSession func(context.Context) ([]byte, UpdateFunc, error)
}

// GetSession returns SessionData or nil for current context
func GetSession(ctx context.Context) SessionData {
	if item, ok := ctx.Value(sessionKey{}).(SessionData); ok {
		return item
	}
	return nil
}

// WithSession returns a new context with SessionData inside.
func WithSession(ctx context.Context, item SessionData) context.Context {
	return context.WithValue(ctx, sessionKey{}, item)
}

// Session is a default middleware to work with sessions.
// getSession function should receive request context and return
// []bytes of current session, UpdateFunc that is invoked if session is modified
// error if something goes wrong.
func Session(getSession func(context.Context) ([]byte, UpdateFunc, error)) MiddlewareFunc {
	return SessionWithConfig(SessionConfig{
		GetSession: getSession,
	})
}

// SessionWithConfig takes SessionConfig and returns SessionMiddleware
func SessionWithConfig(cfg SessionConfig) MiddlewareFunc {
	encode := cfg.Encode
	if encode == nil {
		encode = json.Marshal
	}
	decode := cfg.Decode
	if decode == nil {
		decode = json.Unmarshal
	}

	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context) error {
			sessionBytes, update, err := cfg.GetSession(ctx)
			if err != nil {
				return err
			}
			sessionData := SessionData{}
			if sessionBytes != nil && len(sessionBytes) > 0 {
				err = decode(sessionBytes, &sessionData)
				if err != nil {
					return err
				}
			}

			ctx = WithSession(ctx, sessionData)
			defer func() {
				data, err := encode(sessionData)
				if err != nil {
					return
				}
				if !bytes.Equal(data, sessionBytes) {
					err := update(data)
					if err != nil {
						return
					}
				}
			}()
			err = next.Handle(ctx)
			return err
		})
	}
}
