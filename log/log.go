package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Log struct {
	channel  chan LogMsg
	json     bool
	writer   io.Writer
	errorNum int
}

type LogOptions struct {
	ChannelSize  int
	FormatJson   bool
	Writer       io.Writer
	FailHookFunc []func(msg LogMsg)
}

func New(options ...LogOptions) *Log {
	var size int = 500
	var writer io.Writer = os.Stdout
	var ff = []func(msg LogMsg){}
	var formateJson bool
	if len(options) > 0 {
		if options[0].ChannelSize > 0 {
			size = options[0].ChannelSize
		}
		if options[0].Writer != nil {
			writer = options[0].Writer
		}
		formateJson = options[0].FormatJson
		ff = options[0].FailHookFunc
	}
	l := Log{
		channel: make(chan LogMsg, size),
		json:    formateJson,
		writer:  writer,
	}
	go l.run(ff...)
	return &l
}

type LogLevel int

const (
	DEBUG = iota + 1
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"

	default:
		return "UNKNOWN"
	}
}

type LogMsg struct {
	Level          LogLevel
	ErrorWithStack error
	Msg            string
}

func (l LogMsg) MarshalJSON() ([]byte, error) {
	stack := fmt.Sprintf("%+v", l.ErrorWithStack)

	if l.ErrorWithStack == nil {
		stack = ""
	}
	var msg = map[string]any{
		"level":            l.Level.String(),
		"error_with_stack": stack,
		"message":          l.Msg,
	}
	return json.Marshal(msg)
}

func (l *Log) SendMsg(msg LogMsg) {
	l.channel <- msg
}

func (l *Log) Send(msg string, err error, level LogLevel) {
	l.SendMsg(LogMsg{
		Level:          level,
		ErrorWithStack: err,
		Msg:            msg,
	})
}

func (l *Log) run(failedHook ...func(msg LogMsg)) {
	for msg := range l.channel {
		var marshal []byte
		var err error
		if l.json {
			marshal, err = json.Marshal(msg)
			if err != nil {
				l.channel <- msg
				l.errorNum++
			}
		} else {
			marshal = []byte(fmt.Sprintf("%+v", msg))
		}
		if l.errorNum > 5 {
			if len(failedHook) > 0 {
				for _, f := range failedHook {
					f(msg)
				}
			}
			l.writer = os.Stdout
		}
		_, err = l.writer.Write(marshal)
		if err != nil {
			log.Println(err)
			l.channel <- msg
			l.errorNum++
		}
	}
}
