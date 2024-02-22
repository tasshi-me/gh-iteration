package log_test

import (
	"bytes"
	"strconv"
	"testing"
	"time"

	"github.com/tasshi-me/gh-iteration/pkg/log"
)

func nowMock() time.Time { return time.Date(1995, 0o7, 10, 0o1, 23, 45, 6789, time.UTC) }

func TestLogger_SetFormat(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	logger := log.NewWithOptions(log.FormatPlain, log.ConfigLevelDebug, output, nowMock)
	logger.Info("Hello in plain text!")
	want := "[1995-07-10T01:23:45.000006789Z] INFO: Hello in plain text!\n"
	if output.String() != want {
		t.Errorf("Want %s, got %s", want, output)
	}

	output.Reset()
	logger.SetFormat(log.FormatJSON)
	logger.Info("Hello in JSON text!")
	want = "{\"timestamp\":\"1995-07-10T01:23:45.000006789Z\",\"level\":\"info\",\"message\":\"Hello in JSON text!\"}\n"
	if output.String() != want {
		t.Errorf("Want %s, got %s", want, output)
	}

	output.Reset()
	logger.SetFormat(log.FormatPlain)
	logger.Info("Hello in plain text again!")
	want = "[1995-07-10T01:23:45.000006789Z] INFO: Hello in plain text again!\n"
	if output.String() != want {
		t.Errorf("Want %s, got %s", want, output)
	}
}

func TestLogger_SetLevel(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	logger := log.NewWithOptions(log.FormatPlain, log.ConfigLevelDebug, output, nowMock)
	logger.Info("Hello in debug level!")
	want := "[1995-07-10T01:23:45.000006789Z] INFO: Hello in debug level!\n"
	if output.String() != want {
		t.Errorf("Want %s, got %s", want, output)
	}

	output.Reset()
	logger.SetLevel(log.ConfigLevelWarn)
	logger.Info("Hello in warn level!")
	want = ""
	if output.String() != want {
		t.Errorf("Want %s, got %s", want, output)
	}

	output.Reset()
	logger.SetLevel(log.ConfigLevelInfo)
	logger.Info("Hello in info level!")
	want = "[1995-07-10T01:23:45.000006789Z] INFO: Hello in info level!\n"
	if output.String() != want {
		t.Errorf("Want %s, got %s", want, output)
	}
}

func TestLogger_Trace(t *testing.T) { //nolint:dupl
	t.Parallel()

	tests := []struct {
		configLevel log.ConfigLevel
		hasOutput   bool
	}{
		{log.ConfigLevelTrace, true},
		{log.ConfigLevelDebug, false},
		{log.ConfigLevelInfo, false},
		{log.ConfigLevelWarn, false},
		{log.ConfigLevelError, false},
		{log.ConfigLevelNone, false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(strconv.Itoa(int(test.configLevel)), func(t *testing.T) {
			t.Parallel()
			output := &bytes.Buffer{}
			logger := log.NewWithOptions(log.FormatPlain, test.configLevel, output, nowMock)
			logger.Trace("This is a trace message")
			var want string
			if test.hasOutput {
				want = "[1995-07-10T01:23:45.000006789Z] TRACE: This is a trace message\n"
			} else {
				want = ""
			}
			if output.String() != want {
				t.Errorf("Want %s, got %s", want, output)
			}
		})
	}
}

func TestLogger_Debug(t *testing.T) { //nolint:dupl
	t.Parallel()

	tests := []struct {
		configLevel log.ConfigLevel
		hasOutput   bool
	}{
		{log.ConfigLevelTrace, true},
		{log.ConfigLevelDebug, true},
		{log.ConfigLevelInfo, false},
		{log.ConfigLevelWarn, false},
		{log.ConfigLevelError, false},
		{log.ConfigLevelNone, false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(strconv.Itoa(int(test.configLevel)), func(t *testing.T) {
			t.Parallel()
			output := &bytes.Buffer{}
			logger := log.NewWithOptions(log.FormatPlain, test.configLevel, output, nowMock)
			logger.Debug("This is a debug message")
			var want string
			if test.hasOutput {
				want = "[1995-07-10T01:23:45.000006789Z] DEBUG: This is a debug message\n"
			} else {
				want = ""
			}
			if output.String() != want {
				t.Errorf("Want %s, got %s", want, output)
			}
		})
	}
}

func TestLogger_Info(t *testing.T) { //nolint:dupl
	t.Parallel()
	tests := []struct {
		configLevel log.ConfigLevel
		hasOutput   bool
	}{
		{log.ConfigLevelTrace, true},
		{log.ConfigLevelDebug, true},
		{log.ConfigLevelInfo, true},
		{log.ConfigLevelWarn, false},
		{log.ConfigLevelError, false},
		{log.ConfigLevelNone, false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(strconv.Itoa(int(test.configLevel)), func(t *testing.T) {
			t.Parallel()
			output := &bytes.Buffer{}
			logger := log.NewWithOptions(log.FormatPlain, test.configLevel, output, nowMock)
			logger.Info("This is a info message")
			var want string
			if test.hasOutput {
				want = "[1995-07-10T01:23:45.000006789Z] INFO: This is a info message\n"
			} else {
				want = ""
			}

			if output.String() != want {
				t.Errorf("Want %s, got %s", want, output)
			}
		})
	}
}

func TestLogger_Warn(t *testing.T) { //nolint:dupl
	t.Parallel()
	tests := []struct {
		configLevel log.ConfigLevel
		hasOutput   bool
	}{
		{log.ConfigLevelTrace, true},
		{log.ConfigLevelDebug, true},
		{log.ConfigLevelInfo, true},
		{log.ConfigLevelWarn, true},
		{log.ConfigLevelError, false},
		{log.ConfigLevelNone, false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(strconv.Itoa(int(tt.configLevel)), func(t *testing.T) {
			t.Parallel()
			output := &bytes.Buffer{}
			logger := log.NewWithOptions(log.FormatPlain, test.configLevel, output, nowMock)
			logger.Warn("This is a warn message")
			var want string
			if test.hasOutput {
				want = "[1995-07-10T01:23:45.000006789Z] WARN: This is a warn message\n"
			} else {
				want = ""
			}
			if output.String() != want {
				t.Errorf("Want %s, got %s", want, output)
			}
		})
	}
}

func TestLogger_Error(t *testing.T) { //nolint:dupl
	t.Parallel()
	tests := []struct {
		configLevel log.ConfigLevel
		hasOutput   bool
	}{
		{log.ConfigLevelTrace, true},
		{log.ConfigLevelDebug, true},
		{log.ConfigLevelInfo, true},
		{log.ConfigLevelWarn, true},
		{log.ConfigLevelError, true},
		{log.ConfigLevelNone, false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(strconv.Itoa(int(test.configLevel)), func(t *testing.T) {
			t.Parallel()
			output := &bytes.Buffer{}
			logger := log.NewWithOptions(log.FormatPlain, test.configLevel, output, nowMock)
			logger.Error("This is a error message")
			var want string
			if test.hasOutput {
				want = "[1995-07-10T01:23:45.000006789Z] ERROR: This is a error message\n"
			} else {
				want = ""
			}

			if output.String() != want {
				t.Errorf("Want %s, got %s", want, output)
			}
		})
	}
}
