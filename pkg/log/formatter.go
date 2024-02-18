package log

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

func plainFormatter(event Event) string {
	label := map[EventLevel]string{
		EventLevelTrace: "TRACE",
		EventLevelDebug: "DEBUG",
		EventLevelInfo:  "INFO",
		EventLevelWarn:  "WARN",
		EventLevelError: "ERROR",
	}[event.Level]

	// Because trailing zero is truncated in time.RFC3339Nano, we use custom layout
	layout := "2006-01-02T15:04:05.000000000Z07:00"
	timestamp := event.Timestamp.Format(layout)

	message := fmt.Sprintf("%+v", event.Message)

	output := ""
	sc := bufio.NewScanner(strings.NewReader(message))
	for sc.Scan() {
		output += fmt.Sprintf("[%s] %s: %s\n", timestamp, label, sc.Text())
	}

	return output
}

type JSONFormattedEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

func jsonFormatter(event Event) string {
	label := map[EventLevel]string{
		EventLevelDebug: "debug",
		EventLevelInfo:  "info",
		EventLevelWarn:  "warn",
		EventLevelError: "error",
	}[event.Level]

	message := fmt.Sprintf("%+v", event.Message)

	formatted := JSONFormattedEvent{
		Timestamp: event.Timestamp,
		Level:     label,
		Message:   message,
	}

	str, err := json.Marshal(formatted)
	if err != nil {
		log.Fatal(err)
	}
	return string(str) + "\n"
}
