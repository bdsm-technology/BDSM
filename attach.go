package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"

	"github.com/chzyer/readline"
	"github.com/valyala/fasttemplate"
)

func attach(profile string, usermode, keep bool, prompt *fasttemplate.Template) {
	var bus bus
	bus.init(profile, usermode)
	defer bus.close()
	vs, err := bus.ping()
	if err != nil {
		printWarn("Server is not running!")
		if !keep {
			return
		}
	} else {
		printPair("Server Version", vs)
	}

	username := "nobody"
	hostname := "bedrockserver"
	{
		u, err := user.Current()
		if err == nil {
			username = u.Username
		}
		hn, err := os.Hostname()
		if err == nil {
			hostname = hn
		}
	}

	rl, _ := readline.NewEx(&readline.Config{
		Prompt: prompt.ExecuteString(map[string]interface{}{
			"username": username,
			"hostname": hostname,
			"esc":      "\033",
		}),
		HistoryFile:     ".readline-history",
		AutoComplete:    bus,
		InterruptPrompt: "^C",
		EOFPrompt:       ":quit",

		HistorySearchFold: true,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			if r == readline.CharCtrlZ {
				return r, false
			}
			return r, true
		},
	})
	defer rl.Close()
	lw := rl.Stdout()
	execFn := func(src, cmd string) {
		ncmd := strings.TrimSpace(cmd)
		if len(ncmd) == 0 {
			return
		}
		result, err := bus.exec(ncmd)
		if err != nil {
			fmt.Fprintf(lw, "\033[0m%v\033[0m\n", err)
		} else {
			if len(result) > 0 {
				fmt.Fprintf(lw, "\033[0m%s\033[0m", replacer.Replace(result))
			}
		}
	}
	go func() {
		for v := range bus.log {
			v.write(lw)
		}
	}()
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		execFn("console", line)
	}
}
