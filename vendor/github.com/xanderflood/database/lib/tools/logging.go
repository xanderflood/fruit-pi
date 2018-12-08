package tools

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/mitchellh/colorstring"
)

//Logger logger
//go:generate counterfeiter . Logger
type Logger interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Detail(...interface{})
	Detailf(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
}

//LogLevel enum
type LogLevel int

const (
	//LogLevelError LogLevelError
	LogLevelError LogLevel = iota

	//LogLevelInfo LogLevelInfo
	LogLevelInfo LogLevel = iota

	//LogLevelDetail LogLevelDetail
	LogLevelDetail LogLevel = iota

	//LogLevelDebug LogLevelDebug
	LogLevelDebug LogLevel = iota
)

//NewStdoutLogger create a new logger that writes to stdout
func NewStdoutLogger(level LogLevel, tag string) Logger {
	return &LoggerImpl{
		LogLevel: level,
		Target:   os.Stdout,
		Tag:      tag,
	}
}

//LoggerImpl standard logger implementation
type LoggerImpl struct {
	LogLevel LogLevel
	Target   io.Writer
	Tag      string
}

func (li *LoggerImpl) print(origin string, color string, level string, vals ...interface{}) {
	line := fmt.Sprintln(vals...)
	line = fmt.Sprintf(
		"[light_red]%s [dark_gray][%s-%s] [%s]: [%s]%s",
		origin, li.Tag,
		time.Now().Format(time.RFC3339),
		level,
		color,
		line,
	)
	line = colorstring.Color(line)
	li.Target.Write([]byte(line))
}

func (li *LoggerImpl) printf(origin string, color string, level string, format string, vals ...interface{}) {
	li.print(origin, color, level, fmt.Sprintf(format, vals...))
}

//Error Error
func (li *LoggerImpl) Error(vals ...interface{}) {
	if li.LogLevel >= LogLevelError {
		li.print(li.originInfoString(), "red", "ERROR", vals...)
	}
}

//Errorf Errorf
func (li *LoggerImpl) Errorf(format string, vals ...interface{}) {
	if li.LogLevel >= LogLevelError {
		li.printf(li.originInfoString(), "red", "ERROR", format, vals...)
	}
}

//Info Info
func (li *LoggerImpl) Info(vals ...interface{}) {
	if li.LogLevel >= LogLevelInfo {
		li.print(li.originInfoString(), "blue", "INFO", vals...)
	}
}

//Infof Infof
func (li *LoggerImpl) Infof(format string, vals ...interface{}) {
	if li.LogLevel >= LogLevelInfo {
		li.printf(li.originInfoString(), "blue", "INFO", format, vals...)
	}
}

//Detail Detail
func (li *LoggerImpl) Detail(vals ...interface{}) {
	if li.LogLevel >= LogLevelDetail {
		li.print(li.originInfoString(), "green", "DETAIL", vals...)
	}
}

//Detailf Detailf
func (li *LoggerImpl) Detailf(format string, vals ...interface{}) {
	if li.LogLevel >= LogLevelDetail {
		li.printf(li.originInfoString(), "green", "DETAIL", format, vals...)
	}
}

//Debug Debug
func (li *LoggerImpl) Debug(vals ...interface{}) {
	if li.LogLevel >= LogLevelDebug {
		li.print(li.originInfoString(), "cyan", "DEBUG", vals...)
	}
}

//Debugf Debugf
func (li *LoggerImpl) Debugf(format string, vals ...interface{}) {
	if li.LogLevel >= LogLevelDebug {
		li.printf(li.originInfoString(), "cyan", "DEBUG", format, vals...)
	}
}

func (li *LoggerImpl) originInfoString() string {
	//the ignored value is filepath, but Func.Name
	//includes the (GOPATH-relative) file anyways
	ptr, _, line, _ := runtime.Caller(2)
	funcNameWithPath := runtime.FuncForPC(ptr).Name()

	return fmt.Sprintf("%s:%d", funcNameWithPath, line)
}
