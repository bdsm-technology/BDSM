package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/godbus/dbus"
)

type busLog struct {
	level        uint8
	tag, message string
}

var logTable = []string{"T", "D", "I", "N", "W", "E", "F"}

func (log busLog) write(out io.Writer) (err error) {
	_, err = fmt.Fprintf(out, "\033[0m%s [%v] %v\033[0m\n", logTable[log.level], log.tag, log.message)
	return
}

type bus struct {
	conn *dbus.Conn
	log  chan busLog
	obj  dbus.BusObject
}

func (b *bus) init(profile string, usermode bool) {
	var err error
	if usermode {
		b.conn, err = dbus.SessionBus()
	} else {
		b.conn, err = dbus.SystemBus()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	b.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		fmt.Sprintf("type='signal',path='/one/codehz/bedrockserver',interface='one.codehz.bedrockserver.core',sender='one.codehz.bedrockserver.%s'", profile))
	log := make(chan *dbus.Signal, 10)
	b.log = make(chan busLog, 10)
	b.conn.Signal(log)

	go func() {
		for v := range log {
			if v.Name == "one.codehz.bedrockserver.core.log" {
				b.log <- busLog{
					level:   v.Body[0].(uint8),
					tag:     v.Body[1].(string),
					message: v.Body[2].(string),
				}
			}
		}
	}()

	b.obj = b.conn.Object("one.codehz.bedrockserver."+profile, "/one/codehz/bedrockserver")
}

func (b bus) close() {
	b.conn.Close()
}

func (b bus) exec(cmd string) (rid string, err error) {
	err = b.obj.Call("one.codehz.bedrockserver.core.exec", 0, cmd).Store(&rid)
	return
}

func (b bus) ping() (result string, err error) {
	err = b.obj.Call("one.codehz.bedrockserver.core.ping", 0).Store(&result)
	return
}

func (b bus) complete(line string, pos uint) (result []string, err error) {
	err = b.obj.Call("one.codehz.bedrockserver.core.complete", 0, line, pos).Store(&result)
	return
}

// for readline completer
func (b bus) Do(line []rune, pos int) (newLine [][]rune, length int) {
	if len(line) == 0 {
		newLine = [][]rune{[]rune{'/'}}
		return
	}
	if len(line) == 0 || line[0] != '/' || len(line) != pos {
		return
	}
	r, err := b.complete(string(line), 10000)
	if err != nil {
		printWarn(err.Error())
		return
	}
	mx := strings.LastIndexAny(string(line), " /")
	pfx := string(line[mx+1:])
	for _, item := range r {
		if strings.HasPrefix(item, pfx) {
			newLine = append(newLine, []rune(item)[pos-mx-1:])
		}
	}
	length = pos
	return
}

type listItem struct {
	Name, UUID, XUID string
}

func (b bus) list() (result []listItem, err error) {
	err = b.obj.Call("one.codehz.bedrockserver.core.list", 0).Store(&result)
	return
}

func (b bus) stop() error {
	return b.obj.Call("one.codehz.bedrockserver.core.stop", 0).Err
}
