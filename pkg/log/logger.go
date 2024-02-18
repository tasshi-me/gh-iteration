package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

type EventLevel int

const (
	EventLevelTrace EventLevel = 1 << iota
	EventLevelDebug
	EventLevelInfo
	EventLevelWarn
	EventLevelError
)

type ConfigLevel int

const (
	ConfigLevelTrace ConfigLevel = 1 << iota
	ConfigLevelDebug
	ConfigLevelInfo
	ConfigLevelWarn
	ConfigLevelError
	ConfigLevelNone
)

type Format int

const (
	FormatPlain Format = 1 << iota
	FormatJSON
)

type Event struct {
	Timestamp time.Time
	Level     EventLevel
	Message   any
}

type (
	Printer io.Writer
	Now     func() time.Time
)

type Logger struct {
	format  Format
	level   ConfigLevel
	printer Printer
	now     Now
}

var (
	defaultPrinter = os.Stderr //nolint:gochecknoglobals
	defaultNow     = time.Now  //nolint:gochecknoglobals
)

func New() Logger {
	return NewWithOptions(FormatPlain, ConfigLevelWarn, defaultPrinter, defaultNow)
}

func NewWithOptions(format Format, level ConfigLevel, printer Printer, now Now) Logger {
	return Logger{
		format:  format,
		level:   level,
		printer: printer,
		now:     now,
	}
}

func (logger *Logger) SetLevel(level ConfigLevel) {
	logger.level = level
}

func (logger *Logger) SetFormat(format Format) {
	logger.format = format
}

func (logger *Logger) Trace(message any) {
	event := Event{
		Timestamp: logger.now(),
		Level:     EventLevelTrace,
		Message:   message,
	}
	logger.print(event)
}

func (logger *Logger) Debug(message any) {
	event := Event{
		Timestamp: logger.now(),
		Level:     EventLevelDebug,
		Message:   message,
	}
	logger.print(event)
}

func (logger *Logger) Info(message any) {
	event := Event{
		Timestamp: logger.now(),
		Level:     EventLevelInfo,
		Message:   message,
	}
	logger.print(event)
}

func (logger *Logger) Warn(message any) {
	event := Event{
		Timestamp: logger.now(),
		Level:     EventLevelWarn,
		Message:   message,
	}
	logger.print(event)
}

func (logger *Logger) Error(message any) {
	event := Event{
		Timestamp: logger.now(),
		Level:     EventLevelError,
		Message:   message,
	}
	logger.print(event)
}

func (logger *Logger) print(event Event) {
	eventFilterMap := map[ConfigLevel]map[EventLevel]bool{
		ConfigLevelTrace: {EventLevelTrace: true, EventLevelDebug: true, EventLevelInfo: true, EventLevelWarn: true, EventLevelError: true},
		ConfigLevelDebug: {EventLevelTrace: false, EventLevelDebug: true, EventLevelInfo: true, EventLevelWarn: true, EventLevelError: true},
		ConfigLevelInfo:  {EventLevelTrace: false, EventLevelDebug: false, EventLevelInfo: true, EventLevelWarn: true, EventLevelError: true},
		ConfigLevelWarn:  {EventLevelTrace: false, EventLevelDebug: false, EventLevelInfo: false, EventLevelWarn: true, EventLevelError: true},
		ConfigLevelError: {EventLevelTrace: false, EventLevelDebug: false, EventLevelInfo: false, EventLevelWarn: false, EventLevelError: true},
		ConfigLevelNone:  {EventLevelTrace: false, EventLevelDebug: false, EventLevelInfo: false, EventLevelWarn: false, EventLevelError: false},
	}

	if !eventFilterMap[logger.level][event.Level] {
		return
	}

	var content string
	switch logger.format {
	case FormatPlain:
		content = plainFormatter(event)
	case FormatJSON:
		content = jsonFormatter(event)
	}

	_, _ = fmt.Fprint(logger.printer, content)
}
