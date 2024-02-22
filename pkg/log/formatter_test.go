package log_test

import (
	"testing"

	"github.com/tasshi-me/gh-iteration/pkg/log"
)

func TestPlainFormatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title   string
		message any
		output  string
	}{
		{"string single line", "Hello!", "[1995-07-10T01:23:45.000006789Z] WARN: Hello!\n"},
		{
			"string multi line",
			"First\nSecond\nThird",
			"[1995-07-10T01:23:45.000006789Z] WARN: First\n" +
				"[1995-07-10T01:23:45.000006789Z] WARN: Second\n" +
				"[1995-07-10T01:23:45.000006789Z] WARN: Third\n",
		},
		{"integer", 123, "[1995-07-10T01:23:45.000006789Z] WARN: 123\n"},
		{"float", 1.234, "[1995-07-10T01:23:45.000006789Z] WARN: 1.234\n"},
		{"bool", true, "[1995-07-10T01:23:45.000006789Z] WARN: true\n"},
		{"struct", struct {
			Public  int
			private float32
		}{1, 2.3}, "[1995-07-10T01:23:45.000006789Z] WARN: {Public:1 private:2.3}\n"},
		{"struct empty", struct{}{}, "[1995-07-10T01:23:45.000006789Z] WARN: {}\n"},
		{"nil", nil, "[1995-07-10T01:23:45.000006789Z] WARN: <nil>\n"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.title, func(t *testing.T) {
			t.Parallel()
			event := log.Event{Timestamp: nowMock(), Level: log.EventLevelWarn, Message: test.message}
			got := log.PlainFormatter(event)
			want := test.output
			if got != want {
				t.Errorf("Want %s, got %s", want, got)
			}
		})
	}
}

func TestJsonFormatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		title   string
		message any
		output  string
	}{
		{
			"string single line",
			"Hello!",
			"{\"timestamp\":\"1995-07-10T01:23:45.000006789Z\",\"level\":\"warn\",\"message\":\"Hello!\"}\n",
		},
		{
			"string multi line",
			"First\nSecond\nThird",
			"{\"timestamp\":\"1995-07-10T01:23:45.000006789Z\",\"level\":\"warn\",\"message\":\"First\\nSecond\\nThird\"}\n",
		},
		{"integer", 123, "{\"timestamp\":\"1995-07-10T01:23:45.000006789Z\",\"level\":\"warn\",\"message\":\"123\"}\n"},
		{"float", 1.234, "{\"timestamp\":\"1995-07-10T01:23:45.000006789Z\",\"level\":\"warn\",\"message\":\"1.234\"}\n"},
		{"bool", true, "{\"timestamp\":\"1995-07-10T01:23:45.000006789Z\",\"level\":\"warn\",\"message\":\"true\"}\n"},
		{
			"struct",
			struct {
				Public  int
				private float32
			}{1, 2.3},
			"{\"timestamp\":\"1995-07-10T01:23:45.000006789Z\",\"level\":\"warn\",\"message\":\"{Public:1 private:2.3}\"}\n",
		},
		{
			"struct empty",
			struct{}{},
			"{\"timestamp\":\"1995-07-10T01:23:45.000006789Z\",\"level\":\"warn\",\"message\":\"{}\"}\n",
		},
		{
			"nil",
			nil,
			"{\"timestamp\":\"1995-07-10T01:23:45.000006789Z\",\"level\":\"warn\",\"message\":\"\\u003cnil\\u003e\"}\n",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.title, func(t *testing.T) {
			t.Parallel()
			event := log.Event{Timestamp: nowMock(), Level: log.EventLevelWarn, Message: test.message}
			got := log.JSONFormatter(event)
			want := test.output
			if got != want {
				t.Errorf("Want %s, got %s", want, got)
			}
		})
	}
}
