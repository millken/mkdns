package main

/*
import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

type Color string

const (
	// https://github.com/git/git/blob/master/color.h
	NORMAL       = Color("")
	RESET        = Color("\033[0m")
	BOLD         = Color("\033[1m")
	RED          = Color("\033[31m")
	GREEN        = Color("\033[32m")
	YELLOW       = Color("\033[33m")
	BLUE         = Color("\033[34m")
	MAGENTA      = Color("\033[35m")
	CYAN         = Color("\033[36m")
	LIGHTGRAY    = Color("\033[37m")
	BOLD_RED     = Color("\033[1;31m")
	BOLD_GREEN   = Color("\033[1;32m")
	BOLD_YELLOW  = Color("\033[1;33m")
	BOLD_BLUE    = Color("\033[1;34m")
	BOLD_MAGENTA = Color("\033[1;35m")
	BOLD_CYAN    = Color("\033[1;36m")
	BG_RED       = Color("\033[41m")
	BG_GREEN     = Color("\033[42m")
	BG_YELLOW    = Color("\033[43m")
	BG_BLUE      = Color("\033[44m")
	BG_MAGENTA   = Color("\033[45m")
	BG_CYAN      = Color("\033[46m")
)

type Level int

const (
	FINEST Level = iota
	FINE
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
)

// Logging level strings
var (
	levelStrings = [...]string{"FNST", "FINE", "DEBG", "TRAC", "INFO", "WARN", "EROR", "CRIT"}
)

func (l Level) String() string {
	if l < 0 || int(l) > len(levelStrings) {
		return "UNKNOWN"
	}
	return levelStrings[int(l)]
}

func Colorized(s string, c Color) string {
	if runtime.GOOS == "windows" {
		return s
	}
	return string(c) + s + string(RESET)
}

//Color the level string
type colorLevelString string

func (c colorLevelString) String() (str string) {
	switch c {
	case "FNST":
		str = fmt.Sprintf("%s", Colorized(string(c), BOLD_GREEN))
	case "FINE":
		str = fmt.Sprintf("%s", Colorized(string(c), GREEN))
	case "DEBG":
		str = fmt.Sprintf("%s", Colorized(string(c), MAGENTA))
	case "TRAC":
		str = fmt.Sprintf("%s", Colorized(string(c), LIGHTGRAY))
	case "INFO":
		str = fmt.Sprintf("%s", Colorized(string(c), BLUE))
	case "WARN":
		str = fmt.Sprintf("%s", Colorized(string(c), YELLOW))
	case "EROR":
		str = fmt.Sprintf("%s", Colorized(string(c), RED))
	case "CRIT":
		str = fmt.Sprintf("%s", Colorized(string(c), BOLD_RED))
	default:
		str = string(c)
	}
	return

}

type Logger struct {
	out    io.Writer
	prefix string
	level  Level
	logger *log.Logger
}

func NewLogger(out io.Writer, prefix string, level Level) *Logger {
	logger := log.New(out, prefix, log.Ldate|log.Ltime)
	return &Logger{out: out, level: level, logger: logger}
}

func (l *Logger) logr(lvl Level, arg0 interface{}, args ...interface{}) {
	switch first := arg0.(type) {
	case string:
		l.logf(lvl, first, args...)
	case func() string:
		l.logc(lvl, first)
	default:
		l.logf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func (l *Logger) Finest(arg0 interface{}, args ...interface{}) {
	l.logr(FINEST, arg0, args...)
}

func (l *Logger) Fine(arg0 interface{}, args ...interface{}) {
	l.logr(FINE, arg0, args...)
}

func (l *Logger) Debug(arg0 interface{}, args ...interface{}) {
	l.logr(DEBUG, arg0, args...)
}

func (l *Logger) Trace(arg0 interface{}, args ...interface{}) {
	l.logr(TRACE, arg0, args...)
}

func (l *Logger) Info(arg0 interface{}, args ...interface{}) {
	l.logr(INFO, arg0, args...)
}

func (l *Logger) Warn(arg0 interface{}, args ...interface{}) {
	l.logr(WARNING, arg0, args...)
}

func (l *Logger) Error(arg0 interface{}, args ...interface{}) {
	l.logr(ERROR, arg0, args...)
}

func (l *Logger) Exit(arg0 interface{}, args ...interface{}) {
	l.logr(ERROR, arg0, args...)
	os.Exit(0)
}

func (l *Logger) logf(lvl Level, format string, args ...interface{}) {
	if lvl < l.level {
		return
	}
	pc, _, lineno, ok := runtime.Caller(3)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	msg := fmt.Sprintf("[%s] (%s) %s",
		colorLevelString(levelStrings[lvl]),
		fmt.Sprintf("%s", Colorized(src, LIGHTGRAY)),
		fmt.Sprintf("%s", Colorized(format, BOLD_BLUE)),
	)
	if len(args) > 0 {
		l.logger.Printf(msg, args...)
	} else {
		l.logger.Println(msg)
	}
}

func (l *Logger) logc(lvl Level, closure func() string) {
	if lvl < l.level {
		return
	}
	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	msg := fmt.Sprintf("[%s] (%s) %s",
		colorLevelString(levelStrings[lvl]),
		fmt.Sprintf("%s", Colorized(src, LIGHTGRAY)),
		fmt.Sprintf("%s", Colorized(closure(), BOLD_BLUE)),
	)
	l.logger.Println(msg)
}
*/
