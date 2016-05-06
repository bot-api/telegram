package telebot

import (
	"golang.org/x/net/context"
)

// A Handler takes update message.
type Handler interface {
	Handle(context.Context) error
}

// A Commander takes command message.
type Commander interface {
	Command(ctx context.Context, arg string) error
}

type (

	// MiddlewareFunc defines a function to process middleware.
	MiddlewareFunc func(next Handler) Handler

	// HandlerFunc defines a function to serve on update message.
	// Implements Handler interface.
	HandlerFunc func(context.Context) error

	//MessageFunc func(context.Context, *telegram.Message) error

	// ErrorFunc handles error, if
	ErrorFunc func(ctx context.Context, err error)

	// CommandFunc defines a function to handle commands.
	// Implements Commander interface.
	CommandFunc func(ctx context.Context, arg string) error
)

// Handle method handles message update.
func (h HandlerFunc) Handle(c context.Context) error {
	return h(c)
}

// Command method handles command on message update.
func (c CommandFunc) Command(ctx context.Context, arg string) error {
	return c(ctx, arg)
}

// Commands middleware takes map of commands.
// It runs associated Commander if update messages has a command message.
// Empty command (e.x. "": Commander) used as a default Commander.
// Nil command (e.x. "cmd": nil) used as an EmptyHandler
// Take a look on examples/commands/main.go to know more.
func Commands(commands map[string]Commander) MiddlewareFunc {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context) error {
			update := GetUpdate(ctx)
			if update.Message == nil {
				return next.Handle(ctx)
			}
			command, arg := update.Message.Command()
			if command == "" {
				return next.Handle(ctx)
			}
			cmd, ok := commands[command]
			if !ok {
				if cmd = commands[""]; cmd == nil {
					return next.Handle(ctx)
				}
			}
			if cmd == nil {
				return next.Handle(ctx)
			}
			return cmd.Command(ctx, arg)
		})
	}
}
