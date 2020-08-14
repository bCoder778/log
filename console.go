package log

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type Console struct {
	LevelStyle Style
	log.Logger
}

func (c *Console) log(level Level, subject string, messages []string) {
	isWin := runtime.GOOS == "windows"
	var prefix string
	if isWin {
		prefix = fmt.Sprintf("%-7s", level.ToString("[]"))
	} else {
		preCol := level.ToForeColor()
		prefix = fmt.Sprintf("%c[0;%d;%dm%-7s%c[0m %c[0;0;%dm", 0x1B, c.LevelStyle, preCol, level.ToString("[]"), 0x1B, 0x1B, ForeGreen)
	}
	c.SetPrefix(prefix)
	c.SetFlags(log.Lmicroseconds)
	/*	c.Mu.Lock()
		defer c.Mu.Unlock()*/
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}
	arr := strings.Split(file, "/")
	var msg string
	for i, item := range messages {
		if i > 0 {
			//msg += "\n"
			msg += " "
		}
		msg += item
	}
	if subject != "" {
		msg = subject + splitStr + msg
	}
	fLn := fmt.Sprintf("%s:%d:", arr[len(arr)-1], line)
	if isWin {
		msg = fmt.Sprintf("%-25s %s", fLn, msg)
	} else {
		msg = fmt.Sprintf("%c[4;%d;%dm%-25s%c[0m %s", 0x1B, StyleDefault, ForeBlue, fLn, 0x1B, msg)
	}
	msgChan <- &Messgae{strings.TrimRight(msg, "\n"), level}
}

func (c *Console) print(level Level, msg string) {
	switch level {
	case MAIL:
		c.Println(msg)

		strs := strings.Split(msg, splitStr)
		email.SendEmail(strs[0], strs[1])
	default:
		c.Println(msg)
	}
}
