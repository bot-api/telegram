package telebot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bot-api/telegram"
	"golang.org/x/net/context"
)

type apiKey struct{}
type updateKey struct{}
type webhookKey struct{}

// GetAPI takes telegram API from context.
// Raises panic if context doesn't have API instance.
func GetAPI(ctx context.Context) *telegram.API {
	return ctx.Value(apiKey{}).(*telegram.API)
}

// GetUpdate takes telegram Update from context.
// Raises panic if context doesn't have Update instance.
func GetUpdate(ctx context.Context) *telegram.Update {
	return ctx.Value(updateKey{}).(*telegram.Update)
}

// WithUpdate returns context with telegram Update inside.
// Use GetUpdate to take Update from context.
func WithUpdate(ctx context.Context, u *telegram.Update) context.Context {
	return context.WithValue(ctx, updateKey{}, u)
}

// WithAPI returns context with telegram api inside.
// Use GetAPI to take api from context.
func WithAPI(ctx context.Context, api *telegram.API) context.Context {
	return context.WithValue(ctx, apiKey{}, api)
}

// IsWebhook returns true if update received by webhook
func IsWebhook(ctx context.Context) bool {
	return ctx.Value(webhookKey{}) != nil
}

// A Bot object helps to work with telegram bot api using handlers.
//
// Bot initialization is not thread safe.
type Bot struct {
	api *telegram.API
	me  *telegram.User

	handler    Handler
	middleware []MiddlewareFunc
	errFunc    ErrorFunc
}

// NewWithAPI returns bot with custom API client
func NewWithAPI(api *telegram.API) *Bot {
	return &Bot{
		api:        api,
		middleware: []MiddlewareFunc{},
		errFunc: func(ctx context.Context, err error) {
			log.Printf("update error: %s", err.Error())
		},
	}
}

// New returns bot with default api client
func New(token string) *Bot {
	return NewWithAPI(telegram.New(token))
}

// Use adds middleware to a middleware chain.
func (b *Bot) Use(middleware ...MiddlewareFunc) {
	b.middleware = append(b.middleware, middleware...)
}

// Handle setups handler to handle telegram updates.
func (b *Bot) Handle(handler Handler) {
	b.handler = handler
}

// HandleFunc takes HandlerFunc and sets handler.
func (b *Bot) HandleFunc(handler HandlerFunc) {
	b.handler = handler
}

// ErrorFunc set a ErrorFunc, that handles error returned
// from handlers/middlewares.
func (b *Bot) ErrorFunc(errFunc ErrorFunc) {
	b.errFunc = errFunc
}

// ServeWithConfig runs update cycle with custom update config.
func (b *Bot) ServeWithConfig(ctx context.Context, cfg telegram.UpdateCfg) error {
	if err := b.updateMe(ctx); err != nil {
		return err
	}

	var rErr error
	errCh := make(chan error, 1)
	updatesCh := make(chan telegram.Update)
	go func() {
		errCh <- telegram.GetUpdates(
			ctx,
			b.api,
			cfg,
			updatesCh)
	}()
loop:
	for {
		select {
		case rErr = <-errCh:
			break loop
		case update, ok := <-updatesCh:
			if !ok {
				// update channel was closed, wait for error
				select {
				case rErr = <-errCh:
					break loop
				}
			}
			b.handleUpdate(ctx, &update)

		}
	}
	return rErr
}

// Serve runs update cycle with default update config.
// Offset is zero and timeout is 30 seconds.
func (b *Bot) Serve(ctx context.Context) error {
	cfg := telegram.NewUpdate(0)
	return b.ServeWithConfig(ctx, cfg)
}

// ServeByWebhook returns webhook handler,
// that can handle incoming telegram webhook messages.
//
// Use IsWebhook function to identify webhook updates.
func (b *Bot) ServeByWebhook(ctx context.Context) (http.HandlerFunc, error) {
	if err := b.updateMe(ctx); err != nil {
		return nil, err
	}

	updatesCh := make(chan telegram.Update)
	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case update := <-updatesCh:
				b.handleUpdate(
					context.WithValue(
						ctx,
						webhookKey{},
						struct{}{}),
					&update,
				)
			}
		}
	}()
	return b.getWebhookHandler(ctx, updatesCh), nil
}

// ============== Internal ================================================== //

func (b *Bot) updateMe(ctx context.Context) (err error) {
	b.me, err = b.api.GetMe(ctx)
	return err
}

func (b *Bot) handleUpdate(ctx context.Context, update *telegram.Update) {
	ctx = WithAPI(ctx, b.api)
	ctx = WithUpdate(ctx, update)
	ctx = context.WithValue(ctx, "update.id", update.UpdateID)
	h := b.handler
	if h == nil {
		h = EmptyHandler()
	}
	for i := len(b.middleware) - 1; i >= 0; i-- {
		h = b.middleware[i](h)
	}

	err := h.Handle(ctx)
	if err != nil && b.errFunc != nil {
		eh := b.errFunc
		eh(ctx, err)
	}
}

func (b *Bot) getWebhookHandler(
	ctx context.Context,
	out chan<- telegram.Update) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}

		var update telegram.Update
		err = json.Unmarshal(bytes, &update)
		if err != nil {
			log.Println(err)
			return
		}
		select {
		case out <- update:
		case <-ctx.Done():
			return
		}
	}
}
