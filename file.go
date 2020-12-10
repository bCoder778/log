package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const splitStr = "*&^J%FGH&%"

type File struct {
	Day int
	log.Logger
	Level
}

func (f *File) log(level Level, subject string, messages []string) {
	var prefix = fmt.Sprintf("%-7s", level.ToString("[]"))
	f.SetPrefix(prefix)
	f.SetFlags(log.Lmicroseconds)
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}
	arr := strings.Split(file, "/")
	var msg string
	var printMsg string
	for i, item := range messages {
		if i > 0 {
			msg += " "
		}
		msg += item
	}
	fLn := fmt.Sprintf("%s:%d:", arr[len(arr)-1], line)
	printMsg = fmt.Sprintf("%-25s %s", fLn, msg)
	msgChan <- &Messgae{msg, subject, strings.TrimRight(printMsg, "\n"), level}
}

func (f *File) print(level Level, msg *Messgae) {
	var curTime = time.Now()
	if curTime.Day() != f.Day || !isExistLogFile(f.Level, curTime) {
		file, err := openLogFile(f.Level, curTime)
		if err != nil {
			fmt.Printf("Open log file %s failed! %s\n", curTime.String(), err)
		} else {
			f.Logger = *log.New(file, "file", log.LstdFlags)
			f.Day = curTime.Day()
		}
	}

	switch level {
	case MAIL:
		f.Println(msg.printMsg)

		err := email.SendEmail(msg.subject, msg.msg)
		if err != nil {
			Errorf("send email failed! error:%s, subject:%s, msg:%s", err, msg.subject, msg.msg)
		}
	default:
		f.Println(msg.printMsg)
	}
}

func openLogFile(level Level, curTime time.Time) (*os.File, error) {
	logDir := fmt.Sprintf(path+"/%04d%02d%02d", curTime.Year(), curTime.Month(), curTime.Day())
	if !Exist(path) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	if !Exist(logDir) {
		err := os.Mkdir(logDir, 0777)
		if err != nil {
			return nil, err
		}
	}
	logName := logDir + "/" + level.ToString("") + ".log"
	return os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}

func isExistLogFile(level Level, curTime time.Time) bool {
	fileName := fmt.Sprintf(path+"/%04d%02d%02d/%s.log", curTime.Year(), curTime.Month(), curTime.Day(), level.ToString(""))
	return Exist(fileName)
}

func Exist(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
