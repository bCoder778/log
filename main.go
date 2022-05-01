package log

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
)

type Style int

const (
	StyleDefault Style = 0 + iota
	StyleLight
	StyleUnderline Style = 4 + iota
	StyleFlash
	StyleInverse Style = 7 + iota
	StyleHide
)

type Color int

const (
	ForeBlack Color = 30 + iota
	ForeRed
	ForeGreen
	ForeYellow
	ForeBlue
	ForePurple
	ForeLightBlue
	ForeWhite
	BackBlack Color = 40 + iota
	BackRed
	BackGreen
	BackYellow
	BackBlue
	BackPurple
	BackLightBlue
	BackWhite
)

type Level int
type Mode int

const (
	DEBUG Level = iota
	INFO
	WARN
	FAIL
	ERROR
	MAIL
)

const (
	Mode_File        Mode = 11
	Mode_Console     Mode = 12
	Mode_FileConsole Mode = 13
)

type ILog interface {
	log(level Level, subject string, messages []string)
	print(level Level, msg *Messgae)
}

type Log struct {
	Level Level
	Debug ILog
	Info  ILog
	Warn  ILog
	Fail  ILog
	Error ILog
	Mail  ILog
}

type Messgae struct {
	msg      string
	subject  string
	printMsg string
	Level
}

var msgChan chan *Messgae
var msgChanFile chan *Messgae
var _conf *Log
var _conf2 *Log
var stop chan bool
var email *EMail
var emailOpt = &EMailOption{}

var (
	logLevel = DEBUG
	mode     = Mode_Console
	path     = "logs"
)

func init() {
	stop = make(chan bool)
	msgChan = make(chan *Messgae, 1000)
	msgChanFile = make(chan *Messgae, 1000)
	start()
}

func start() {
	email = NewEmail(emailOpt)
	level := logLevel
	switch mode {
	case Mode_File:
		var curTime = time.Now()
		_conf = &Log{
			Level: level,
			Debug: &File{curTime.Day(), *createLevelLogger(DEBUG, curTime), DEBUG},
			Info:  &File{curTime.Day(), *createLevelLogger(INFO, curTime), INFO},
			Warn:  &File{curTime.Day(), *createLevelLogger(WARN, curTime), WARN},
			Fail:  &File{curTime.Day(), *createLevelLogger(FAIL, curTime), FAIL},
			Error: &File{curTime.Day(), *createLevelLogger(ERROR, curTime), ERROR},
			Mail:  &File{curTime.Day(), *createLevelLogger(MAIL, curTime), MAIL},
		}
		dealFatal()
		_conf2 = nil
	case Mode_FileConsole:
		var curTime = time.Now()
		_conf = &Log{
			Level: level,
			Debug: &File{curTime.Day(), *createLevelLogger(DEBUG, curTime), DEBUG},
			Info:  &File{curTime.Day(), *createLevelLogger(INFO, curTime), INFO},
			Warn:  &File{curTime.Day(), *createLevelLogger(WARN, curTime), WARN},
			Fail:  &File{curTime.Day(), *createLevelLogger(FAIL, curTime), FAIL},
			Error: &File{curTime.Day(), *createLevelLogger(ERROR, curTime), ERROR},
			Mail:  &File{curTime.Day(), *createLevelLogger(MAIL, curTime), MAIL},
		}
		dealFatal()

		cliLogger := log.New(os.Stdout, "", log.Ltime)
		_conf2 = &Log{
			Level: level,
			Debug: &Console{StyleFlash, *cliLogger},
			Info:  &Console{StyleFlash, *cliLogger},
			Warn:  &Console{StyleFlash, *cliLogger},
			Fail:  &Console{StyleFlash, *cliLogger},
			Error: &Console{StyleFlash, *cliLogger},
			Mail:  &Console{StyleFlash, *cliLogger},
		}

	default:
		cliLogger := log.New(os.Stdout, "", log.Ltime)
		_conf2 = &Log{
			Level: level,
			Debug: &Console{StyleFlash, *cliLogger},
			Info:  &Console{StyleFlash, *cliLogger},
			Warn:  &Console{StyleFlash, *cliLogger},
			Fail:  &Console{StyleFlash, *cliLogger},
			Error: &Console{StyleFlash, *cliLogger},
			Mail:  &Console{StyleFlash, *cliLogger},
		}
	}
	go OutputMsg()
}

type Option struct {
	LogLevel Level
	Mode     Mode
	Path     string
	Email    *EMailOption
}

func SetOption(opt *Option) {
	stop <- true
	logLevel = opt.LogLevel
	mode = opt.Mode
	emailOpt = opt.Email
	path = opt.Path
	start()
}

func Debug(msg ...string) {
	if _conf.Level > DEBUG {
		return
	}
	if _conf != nil {
		_conf.Debug.log(DEBUG, "", msg)
	}
	if _conf2 != nil {
		_conf2.Debug.log(DEBUG, "", msg)
	}
}

func Debugf(format string, a ...interface{}) {
	Debug(fmt.Sprintf(format, a...))
}

func Info(msg ...string) {
	if _conf.Level > INFO {
		return
	}
	if _conf != nil {
		_conf.Info.log(INFO, "", msg)
	}
	if _conf2 != nil {
		_conf2.Info.log(INFO, "", msg)
	}
}
func Infof(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a...))
}

func Warn(msg ...string) {
	if _conf.Level > WARN {
		return
	}
	if _conf != nil && _conf2 != nil {
		_conf.Warn.log(WARN, "", msg)
	} else if _conf != nil {
		_conf.Warn.log(WARN, "", msg)
	} else if _conf2 != nil {
		_conf2.Warn.log(WARN, "", msg)
	}
}
func Warnf(format string, a ...interface{}) {
	Warn(fmt.Sprintf(format, a...))
}

func Fail(msg ...string) {
	if _conf.Level > FAIL {
		return
	}
	if _conf != nil {
		_conf.Fail.log(FAIL, "", msg)
	}
	if _conf2 != nil {
		_conf2.Fail.log(FAIL, "", msg)
	}
}
func Failf(format string, a ...interface{}) {
	Fail(fmt.Sprintf(format, a...))
}

func Error(msg ...string) {
	if _conf.Level > ERROR {
		return
	}
	if _conf != nil {
		_conf.Error.log(ERROR, "", msg)
	}
	if _conf2 != nil {
		_conf2.Error.log(ERROR, "", msg)
	}
}

func Errorf(format string, a ...interface{}) {
	Error(fmt.Sprintf(format, a...))
}

func Mail(subject string, msg ...string) {
	if _conf.Level > MAIL {
		return
	}

	if _conf != nil {
		_conf.Mail.log(MAIL, subject, msg)
	}
	if _conf2 != nil {
		_conf2.Mail.log(MAIL, subject, msg)
	}
}

func Mailf(suject string, format string, a ...interface{}) {
	Mail(suject, fmt.Sprintf(format, a...))
}

func (lv Level) ToString(flag string) string {
	rs := ""
	switch lv {
	case DEBUG:
		rs = "DEBUG"
		break
	case INFO:
		rs = "INFO"
		break
	case WARN:
		rs = "WARN"
		break
	case FAIL:
		rs = "FAIL"
		break
	case ERROR:
		rs = "ERROR"
		break
	case MAIL:
		rs = "MAIL"
		break
	}

	switch flag {
	case "<>":
		rs = fmt.Sprintf("<%s>", rs)
		break
	case "[]":
		rs = fmt.Sprintf("[%s]", rs)
		break
	case "#":
		rs = fmt.Sprintf("#%s#", rs)
		break
	}

	return rs
}

func (lv Level) ToBackColor() Color {
	switch lv {
	case DEBUG:
		return BackGreen
	case INFO:
		return BackLightBlue
	case WARN:
		return BackYellow
	case FAIL:
		return BackPurple
	case ERROR:
		return BackRed
	case MAIL:
		return BackBlue
	}
	return 0
}

func (lv Level) ToForeColor() Color {
	switch lv {
	case DEBUG:
		return ForeGreen
	case INFO:
		return ForeLightBlue
	case WARN:
		return ForeYellow
	case FAIL:
		return ForePurple
	case ERROR:
		return ForeRed
	case MAIL:
		return ForeBlue
	}
	return 0
}

func TranslateToLevel(l string) Level {
	switch strings.ToUpper(l) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "FAIL":
		return FAIL
	case "ERROR":
		return ERROR
	case "MAIL":
		return MAIL
	}
	return DEBUG
}

func OutputMsg() {
	var msg *Messgae
	for {
		select {
		case <-stop:
			return
		case msg = <-msgChan:
			switch msg.Level {
			case DEBUG:
				_conf2.Debug.print(msg.Level, msg)
			case INFO:
				_conf2.Info.print(msg.Level, msg)
			case WARN:
				_conf2.Warn.print(msg.Level, msg)
			case FAIL:
				_conf2.Fail.print(msg.Level, msg)
			case ERROR:
				_conf2.Error.print(msg.Level, msg)
			case MAIL:
				_conf2.Mail.print(msg.Level, msg)
			}
		case msg = <-msgChanFile:
			switch msg.Level {
			case DEBUG:
				_conf.Debug.print(msg.Level, msg)
			case INFO:
				_conf.Info.print(msg.Level, msg)
			case WARN:
				_conf.Warn.print(msg.Level, msg)
			case FAIL:
				_conf.Fail.print(msg.Level, msg)
			case ERROR:
				_conf.Error.print(msg.Level, msg)
			case MAIL:
				_conf.Mail.print(msg.Level, msg)
			}
		}
	}
}

func createLevelLogger(level Level, t time.Time) *log.Logger {
	if level >= logLevel {
		file, err := openLogFile(level, t)
		if err != nil {
			fmt.Printf("create level logger failed! %s/n", err)
			os.Exit(1)
		}
		fileLogger := log.New(file, "file", log.LstdFlags)
		return fileLogger
	}
	return &log.Logger{}
}

func dealFatal() {
	if !Exist(path) {
		err := os.Mkdir(path, 0666)
		if err != nil {
			fmt.Println("deal fatal failed! ", err)
			os.Exit(1)
		}
	}
	logFile, err := os.OpenFile(path+"/fatal.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("deal fatal failed! ", err)
		os.Exit(1)
	}
	syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))
}
