package telebot

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"golang.org/x/net/context"
)

type (
	// RecoverCfg defines the config for recover middleware.
	RecoverCfg struct {
		// StackSize is the stack size to be printed.
		// Optional, with default value as 4 KB.
		StackSize int

		// EnableStackAll enables formatting stack traces of all
		// other goroutines into buffer after the trace
		// for the current goroutine.
		// Optional, with default value as false.
		EnableStackAll bool

		// DisablePrintStack disables printing stack trace.
		// Optional, with default value as false.
		DisablePrintStack bool

		// LogFunc uses to write recover data to your own logger.
		// Please do not use stack arg after function execution,
		// because it will be freed.
		LogFunc func(ctx context.Context, cause error, stack []byte)
	}
)

var (
	// DefaultRecoverConfig is the default recover middleware config.
	DefaultRecoverConfig = RecoverCfg{
		StackSize:         4 << 10, // 4 KB
		DisablePrintStack: false,
	}
	// DefaultRecoverLogger is used to print recover information
	// if RecoverCfg.LogFunc is not set
	DefaultRecoverLogger = log.New(os.Stderr, "", log.LstdFlags)
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and returns nil error.
// It prints recovery information to standard log.
func Recover() MiddlewareFunc {
	return RecoverWithConfig(DefaultRecoverConfig)
}

// RecoverWithConfig returns a middleware which recovers from panics anywhere
// in the chain and returns nil error.
// It takes RecoverCfg to configure itself.
func RecoverWithConfig(cfg RecoverCfg) MiddlewareFunc {
	// Defaults
	if cfg.StackSize == 0 {
		cfg.StackSize = DefaultRecoverConfig.StackSize
	}

	var logFunc func(ctx context.Context, cause error, stack []byte)
	if cfg.LogFunc != nil {
		logFunc = cfg.LogFunc
	} else {
		logFunc = defaultLogFunc
	}

	var stackPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, cfg.StackSize)
		},
	}

	// getBuffer returns a buffer from the pool.
	getBuffer := func() (buf []byte) {
		return stackPool.Get().([]byte)
	}

	// putBuffer returns a buffer to the pool.
	// The buffer is reset before it is put back into circulation.
	putBuffer := func(buf []byte) {
		stackPool.Put(buf)
	}

	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context) error {
			defer func() {
				r := recover()
				if r == nil {
					return
				}
				var err error
				switch r := r.(type) {
				case error:
					err = r
				default:
					err = fmt.Errorf("%v", r)
				}
				var stackBuf []byte
				if !cfg.DisablePrintStack {
					stackBuf = getBuffer()
					length := runtime.Stack(
						stackBuf, cfg.EnableStackAll)

					stackBuf = stackBuf[:length]
				}
				logFunc(ctx, err, stackBuf)
				if stackBuf != nil {
					putBuffer(stackBuf)
				}

			}()
			return next.Handle(ctx)
		})
	}
}

func defaultLogFunc(ctx context.Context, cause error, stack []byte) {
	buf := bytes.NewBufferString("PANIC RECOVER")
	buf.WriteString("\ncause:")
	buf.WriteString(cause.Error())
	if stack != nil {
		buf.WriteString("\nstack:")
		buf.Write(stack)
	}
	DefaultRecoverLogger.Print(buf.String())
}
