package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bot-api/telegram"
	"github.com/chzyer/readline"
	"golang.org/x/net/context"
)

var nl = []byte("\n")

var cmdUsage = map[string]string{
	"getMe":      "Returns information about bot. ",
	"getUpdates": "Get updates. Usage: getUpdates [offset, [limit, [timeout]]]",
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("help"),
	readline.PcItem("quit"),
	readline.PcItem("q"),
	readline.PcItem("getMe"),
	readline.PcItem("getUpdates"),
	readline.PcItem("getFile"),
	readline.PcItem("getUserProfilePhotos"),
	readline.PcItem("sendChatAction"),
	readline.PcItem("sendMessage"),
	readline.PcItem("sendLocation"),
	readline.PcItem("openChat"),
	readline.PcItem("setWebhook"),
	readline.PcItem("removeWebhook"),
)

func readToken(rl *readline.Instance) (string, error) {
	for {
		token, err := rl.ReadPassword("Give me the token: ")
		if err != nil {
			return "", err
		}
		if len(token) > 0 {
			return string(token), nil
		}
	}
}

func usage(rl *readline.Instance, completer *readline.PrefixCompleter) {
	w := rl.Stdout()
	io.WriteString(w, "Available commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

func help(rl *readline.Instance, cmd string) {
	fmt.Fprintln(rl.Stdout(), cmdUsage[cmd])
}

func writeJSON(rl *readline.Instance, method string, obj interface{}) {
	data, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		io.WriteString(
			rl.Stderr(),
			fmt.Sprintf("%s json error: %s\n", method, err.Error()),
		)
	}

	fmt.Fprintf(rl.Stdout(), "%s:", method)
	_, err = rl.Stdout().Write(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(rl.Stdout())
}

func getMe(rl *readline.Instance, ctx context.Context, cl *telegram.API) {
	me, err := cl.GetMe(ctx)
	if err != nil {
		log.Printf("getMe error: %s\n", err.Error())
		return
	}
	writeJSON(rl, "getMe", me)

}

func openChat(
	ctx context.Context,
	rl *readline.Instance,
	cl *telegram.API,
	args ...string) error {

	if len(args) != 1 {
		return fmt.Errorf("usage: openChat chatID")
	}
	chatID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	updateCh := make(chan *telegram.Update)
	upCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func(
		ctx context.Context,
		cfg telegram.UpdateCfg,
		out chan<- *telegram.Update) {
	loop:
		for {
			updates, err := cl.GetUpdates(ctx, cfg)
			if err != nil {
				if err == context.Canceled {
					close(out)
					return
				}
				log.Println(err)
				if telegram.IsForbiddenError(err) {
					close(out)
					return
				}
				log.Println("Failed to get updates, retrying in 3 seconds...")
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Second * 3):
					continue loop
				}
			}

			for _, update := range updates {
				if update.UpdateID >= cfg.Offset {
					cfg.Offset = update.UpdateID + 1
					select {
					case <-ctx.Done():
						return
					case out <- &update:
					}
				}
			}
		}
	}(upCtx, telegram.UpdateCfg{Timeout: 3}, updateCh)

	go func() {
		for {
			line, err := rl.Readline()
			if err != nil {
				// io.EOF, readline.ErrInterrupt
				break
			}
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "/q") {
				cancel()
				return
			}
			msg := telegram.NewMessage(chatID, line)
			msg.ParseMode = telegram.MarkdownMode
			_, err = cl.SendMessage(ctx, msg)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case update, ok := <-updateCh:
			if !ok {
				return nil
			}
			if update.Message.Chat.ID != chatID {
				continue
			}
			fmt.Fprintf(
				rl.Stdout(),
				"@%s: %s\n",
				update.Message.From.Username,
				update.Message.Text,
			)
			rl.Refresh()
		}
	}
}

func getUpdates(rl *readline.Instance, ctx context.Context, cl *telegram.API,
	cfg telegram.UpdateCfg) {

	updates, err := cl.GetUpdates(ctx, cfg)
	if err != nil {
		io.WriteString(rl.Stderr(),
			fmt.Sprintf("getUpdates error: %s\n", err.Error()))
		return
	}
	writeJSON(rl, "getUpdates", updates)
}

func main() {
	token := flag.String("token", "", "telegram bot token")
	debug := flag.Bool("debug", false, "show debug information")
	flag.Parse()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		HistoryFile:  "/tmp/tbot.history",
		AutoComplete: completer,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	log.SetOutput(rl.Stderr())
	if *token == "" {
		*token, err = readToken(rl)
		if err != nil {
			panic(err)
		}
	}
	if *debug {
		log.Printf("token: %s\n", *token)
	}

	cl := telegram.New(*token)
	cl.Debug(*debug)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	me, err := cl.GetMe(ctx)
	if err != nil {
		log.Printf("getMe error: %s\n", err.Error())
		log.Println("Exited")
		return
	}
	rl.SetPrompt(fmt.Sprintf("@%s> ", me.Username))

	fmt.Fprintln(rl.Stdout(), "Type 'help' for more information.")
loop:
	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}
		line = strings.TrimSpace(line)
		switch {
		case line == "quit" || line == "q":
			break loop
		case line == "help":
			usage(rl, completer)
		case strings.HasPrefix(line, "help "):
			help(rl, line[5:])
		case line == "getMe":
			getMe(rl, ctx, cl)
		case line == "getUpdates":
			getUpdates(rl, ctx, cl, telegram.UpdateCfg{})
		case strings.HasPrefix(line, "getUpdates"):
			args := strings.Split(line[len("getUpdates "):], " ")
			cfg := telegram.UpdateCfg{}
			switch len(args) {
			case 3:
				cfg.Timeout, err = strconv.Atoi(args[2])
				if err != nil {
					log.Println(err.Error())
					continue loop
				}
				fallthrough
			case 2:
				cfg.Limit, err = strconv.Atoi(args[1])
				if err != nil {
					log.Println(err.Error())
					continue loop
				}
				fallthrough
			case 1:
				cfg.Offset, err = strconv.ParseUint(args[0], 10, 64)
				if err != nil {
					log.Println(err.Error())
					continue loop
				}
			}
			getUpdates(rl, ctx, cl, cfg)

		case strings.HasPrefix(line, "sendChatAction "):
			args := strings.Split(line[len("sendChatAction "):], " ")
			if len(args) != 2 {
				log.Println("usage: sendChatAction @to action")
				continue loop
			}
			chatID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				log.Println(err.Error())
				continue loop
			}
			err = cl.SendChatAction(
				ctx,
				telegram.NewChatAction(chatID, args[1]),
			)
			if err != nil {
				log.Println(err.Error())
			}

		case strings.HasPrefix(line, "sendMessage "):
			args := strings.Split(line[len("sendMessage "):], " ")
			if len(args) < 2 {
				log.Println("usage: sendMessage @to text")
				continue loop
			}
			cfg := telegram.MessageCfg{
				Text: strings.Join(args[1:], " "),
			}
			if strings.HasPrefix(args[0], "@") {
				cfg.ChannelUsername = args[0]
			} else {
				chatID, err := strconv.Atoi(args[0])
				if err != nil {
					log.Println(err.Error())
					continue loop
				}
				cfg.ID = int64(chatID)
			}
			_, err = cl.SendMessage(ctx, cfg)
			if err != nil {
				log.Println(err.Error())
			}
		case strings.HasPrefix(line, "getUserProfilePhotos "):
			args := strings.Split(line[len("getUserProfilePhotos "):], " ")
			if len(args) != 1 {
				log.Println("usage: getUserProfilePhotos user_id")
				continue loop
			}
			userID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				log.Println(err.Error())
				continue loop
			}
			cfg := telegram.NewUserProfilePhotos(userID)
			photos, err := cl.GetUserProfilePhotos(ctx, cfg)
			if err != nil {
				log.Println(err.Error())
				continue loop
			}
			writeJSON(rl, "getUserProfilePhotos", photos)

		case strings.HasPrefix(line, "getFile "):
			args := strings.Split(line[len("getFile "):], " ")
			if len(args) != 1 {
				log.Println("usage: getFile file_id")
				continue loop
			}
			cfg := telegram.FileCfg{
				FileID: args[0],
			}
			file, err := cl.GetFile(ctx, cfg)
			if err != nil {
				log.Println(err.Error())
				continue loop
			}
			writeJSON(rl, "getFile", file)
		case line == "openChat":
			err := openChat(ctx, rl, cl)
			if err != nil {
				log.Println(err.Error())
				continue loop
			}
		case strings.HasPrefix(line, "openChat "):
			args := strings.Split(line[len("openChat "):], " ")

			err := openChat(ctx, rl, cl, args...)
			if err != nil {
				log.Println(err.Error())
				continue loop
			}

		case strings.HasPrefix(line, "sendLocation "):
			args := strings.Split(line[len("sendLocation "):], " ")
			if len(args) != 3 {
				log.Println("usage: sendLocation @to lat lon")
				continue loop
			}
			chatID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				log.Println(err.Error())
				continue loop
			}
			lat, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				log.Println(err.Error())
				continue loop
			}
			lon, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				log.Println(err.Error())
				continue loop
			}
			cfg := telegram.NewLocation(chatID, lat, lon)
			_, err = cl.Send(ctx, cfg)
			if err != nil {
				log.Println(err.Error())
				continue loop
			}
		case line == "removeWebhook":
			err := cl.SetWebhook(ctx, telegram.NewWebhook(""))
			if err != nil {
				log.Println(err.Error())
				continue loop
			}

		case strings.HasPrefix(line, "setWebhook "):
			args := strings.Split(line[len("setWebhook "):], " ")
			hookURL := ""
			if len(args) > 0 {
				hookURL = args[0]
			}
			var inputFile telegram.InputFile
			if len(args) > 1 {
				cert, err := os.Open(args[1])
				if err != nil {
					log.Println(err.Error())
					continue loop
				}
				inputFile = telegram.NewInputFile("cert.pem", cert)
			}
			cfg := telegram.NewWebhookWithCert(hookURL,
				inputFile)
			err = cl.SetWebhook(ctx, cfg)
			if inputFile != nil {
				if rc, ok := inputFile.Reader().(io.ReadCloser); ok {
					rc.Close()
				}
			}
			if err != nil {
				log.Println(err.Error())
				continue loop
			}
		}

	}
}
