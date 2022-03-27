package superLog

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	debug                                        = false
	infoLog, debugLog, warnLog, errLog, fatalLog *log.Logger
)

func init() {
	infoLog = log.New(os.Stdout, color.GreenString("[INFO]"), log.LstdFlags|log.LUTC)
	debugLog = log.New(os.Stdout, color.CyanString("[DEBUG]"), log.LstdFlags|log.Lshortfile|log.LUTC)
	errLog = log.New(os.Stderr, color.RedString("[ERROR]"), log.LstdFlags|log.Lshortfile|log.LUTC)
	warnLog = log.New(os.Stdout, color.YellowString("[WARN]"), log.LstdFlags|log.LUTC)
	fatalLog = log.New(os.Stderr, color.RedString("[FATAL]"), log.LstdFlags|log.Llongfile|log.LUTC)
}

func EnableDebug() {
	debug = true
}

func Error(v ...interface{}) {
	if v[0] != nil {
		errLog.Output(2, fmt.Sprintln(v...))
	}
}

func PanicError(err error, msg ...interface{}) {
	if err != nil {
		if len(msg) > 0 {
			errLog.Output(2, err.Error()+":"+fmt.Sprint(msg...))
		} else {
			errLog.Output(2, err.Error())
		}
		panic(err)
	}
}

func Warn(v ...interface{}) {
	if v[0] != nil {
		warnLog.Output(2, fmt.Sprintln(v...))
	}
}

func Info(v ...interface{}) {
	if v[0] != nil {
		infoLog.Output(2, fmt.Sprintln(v...))
	}
}

func Debug(v ...interface{}) {
	if debug && v[0] != nil {
		debugLog.Output(2, fmt.Sprintln(v...))
	}
}

func Infof(msg string, v ...interface{}) {
	strMsg := fmt.Sprintf(msg, v...)
	Info(strMsg)
}

func Debugf(msg string, v ...interface{}) {
	strMsg := fmt.Sprintf(msg, v...)
	Debug(strMsg)
}

func Warnf(msg string, v ...interface{}) {
	strMsg := fmt.Sprintf(msg, v...)
	Warn(strMsg)
}

func Errorf(msg string, v ...interface{}) {
	strMsg := fmt.Sprintf(msg, v...)
	Error(strMsg)
}

func Fatal(v ...interface{}) {
	var msg []string
	for _, i := range v {
		msg = append(msg, fmt.Sprintf("%v", i))
	}
	fatalLog.Output(2, strings.Join(msg, " "))
	os.Exit(1)
}

func Fatalf(msg string, v ...interface{}) {
	Fatal(msg, v)
}
