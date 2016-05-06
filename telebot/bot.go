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

//type botContext struct {
//	context.Context
//}

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

func WithUpdate(ctx context.Context, u *telegram.Update) context.Context {
	return context.WithValue(ctx, updateKey{}, u)
}

func WithAPI(ctx context.Context, api *telegram.API) context.Context {
	return context.WithValue(ctx, apiKey{}, api)
}

// Bot
// Bot is not thread safe
type Bot struct {
	api *telegram.API
	me  *telegram.User

	handler    Handler
	middleware []MiddlewareFunc
	errFunc    ErrorFunc
}

// NewWithApi returns bot with custom api client
func NewWithApi(api *telegram.API) *Bot {
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
	return NewWithApi(telegram.New(token))
}

func (b *Bot) Use(middleware ...MiddlewareFunc) {
	b.middleware = append(b.middleware, middleware...)
}

func (b *Bot) Handle(handler Handler) {
	b.handler = handler
}

func (b *Bot) HandleFunc(handler HandlerFunc) {
	b.handler = handler
}

// ErrorFunc set a ErrorFunc, that handles error returned
// from handlers/middlewares.
func (b *Bot) ErrorFunc(errFunc ErrorFunc) {
	b.errFunc = errFunc
}

func (b *Bot) Serve(ctx context.Context) error {
	if err := b.updateMe(ctx); err != nil {
		return err
	}

	var rErr error
	errCh := make(chan error, 1)
	updatesCh := make(chan telegram.Update)
	lastUpdate := uint64(0)
	go func() {
		errCh <- telegram.GetUpdates(
			ctx,
			b.api,
			telegram.NewUpdate(lastUpdate),
			updatesCh)
	}()
loop:
	for {
		select {
		case rErr = <-errCh:
			break loop
		case update := <-updatesCh:
			b.handleUpdate(ctx, &update)

		}
	}
	return rErr
}

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
				b.handleUpdate(ctx, &update)
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
