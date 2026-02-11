package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Logger struct {
	Enabled map[string]bool
	Output  io.Writer
}

var logTypes = map[string]string{
	"debug":   "DEBUG",
	"info":    "INFO",
	"warning": "WARN",
	"warn":    "WARN",
	"error":   "ERROR",
}

func NewLogger(types []string, logFile string) (*Logger, error) {
	var out io.Writer = os.Stderr
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		out = f
	}

	enabled := make(map[string]bool)
	for _, t := range types {
		t = strings.ToLower(strings.TrimSpace(t))
		if logType, ok := logTypes[t]; ok {
			enabled[logType] = true
		}
	}

	return &Logger{
		Enabled: enabled,
		Output:  out,
	}, nil
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.Enabled["DEBUG"] {
		l.log("DEBUG", format, args...)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.Enabled["INFO"] {
		l.log("INFO", format, args...)
	}
}

func (l *Logger) Warn(format string, args ...interface{}) {
	if l.Enabled["WARN"] {
		l.log("WARN", format, args...)
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.Enabled["ERROR"] {
		l.log("ERROR", format, args...)
	}
}

func (l *Logger) log(logType string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.Output, "[%s] [%s] %s\n", timestamp, logType, msg)
}
