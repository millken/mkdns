package logger

import (
	"fmt"
	"io"
	"os"
	"time"
	"runtime"
)

var stdout io.Writer = os.Stdout


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

// This is the standard writer that prints to standard output.
type ConsoleLogWriter chan *LogRecord

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter() ConsoleLogWriter {
	records := make(ConsoleLogWriter, LogBufferLength)
	go records.run(stdout)
	return records
}

func (w ConsoleLogWriter) run(out io.Writer) {
	var timestr string
	var timestrAt int64

	for rec := range w {
		if at := rec.Created.UnixNano() / 1e9; at != timestrAt {
			timestr, timestrAt = rec.Created.Format("2006/01/02 15:04:05"), at
		}
		fmt.Fprintf(out, "[%s] [%s] (%s) %s\n",
			timestr,
			colorLevelString(levelStrings[rec.Level]),
			fmt.Sprintf("%s", Colorized(rec.Source, LIGHTGRAY)),
			fmt.Sprintf("%s", Colorized(rec.Message, BOLD_BLUE)),
		)
	}
}

// This is the ConsoleLogWriter's output method. This will block if the output
// buffer is full.
func (w ConsoleLogWriter) LogWrite(rec *LogRecord) {
	w <- rec
}

// Close stops the logger from sending messages to standard output. Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w ConsoleLogWriter) Close() {
	close(w)
	time.Sleep(50 * time.Millisecond) // Try to give console I/O time to complete
}
