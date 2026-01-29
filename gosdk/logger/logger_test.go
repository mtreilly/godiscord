package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{Level(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.level.String()
			if got != tt.want {
				t.Errorf("Level(%d).String() = %q, want %q", tt.level, got, tt.want)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input string
		want  Level
	}{
		{"debug", DebugLevel},
		{"DEBUG", InfoLevel}, // case sensitive, falls through to default
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"error", ErrorLevel},
		{"", InfoLevel},        // default
		{"unknown", InfoLevel}, // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseLevel(tt.input)
			if got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	log := New(DebugLevel, "json", &buf)

	if log.level != DebugLevel {
		t.Errorf("level = %v, want %v", log.level, DebugLevel)
	}
	if log.format != "json" {
		t.Errorf("format = %q, want %q", log.format, "json")
	}
	if log.writer != &buf {
		t.Error("writer not set correctly")
	}
}

func TestNewNilWriterDefaultsToStderr(t *testing.T) {
	log := New(InfoLevel, "json", nil)
	if log.writer != os.Stderr {
		t.Error("expected nil writer to default to os.Stderr")
	}
}

func TestDefault(t *testing.T) {
	log := Default()
	if log.level != InfoLevel {
		t.Errorf("Default() level = %v, want %v", log.level, InfoLevel)
	}
	if log.format != "json" {
		t.Errorf("Default() format = %q, want %q", log.format, "json")
	}
}

func TestIsDebug(t *testing.T) {
	tests := []struct {
		level Level
		want  bool
	}{
		{DebugLevel, true},
		{InfoLevel, false},
		{WarnLevel, false},
		{ErrorLevel, false},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			log := New(tt.level, "json", &bytes.Buffer{})
			if got := log.IsDebug(); got != tt.want {
				t.Errorf("IsDebug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebugLogging(t *testing.T) {
	var buf bytes.Buffer
	log := New(DebugLevel, "json", &buf)
	log.Debug("test message", "key", "value")

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["message"] != "test message" {
		t.Errorf("message = %q, want %q", entry["message"], "test message")
	}
	if entry["level"] != "debug" {
		t.Errorf("level = %q, want %q", entry["level"], "debug")
	}
	if entry["key"] != "value" {
		t.Errorf("key = %v, want %v", entry["key"], "value")
	}
}

func TestDebugBelowLevelNotLogged(t *testing.T) {
	var buf bytes.Buffer
	log := New(InfoLevel, "json", &buf)
	log.Debug("should not appear")

	if buf.Len() > 0 {
		t.Error("Debug message logged when level is Info")
	}
}

func TestInfoLogging(t *testing.T) {
	var buf bytes.Buffer
	log := New(InfoLevel, "json", &buf)
	log.Info("info message")

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["message"] != "info message" {
		t.Errorf("message = %q, want %q", entry["message"], "info message")
	}
	if entry["level"] != "info" {
		t.Errorf("level = %q, want %q", entry["level"], "info")
	}
}

func TestWarnLogging(t *testing.T) {
	var buf bytes.Buffer
	log := New(WarnLevel, "json", &buf)
	log.Warn("warn message", "count", 42)

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["message"] != "warn message" {
		t.Errorf("message = %q, want %q", entry["message"], "warn message")
	}
	if entry["level"] != "warn" {
		t.Errorf("level = %q, want %q", entry["level"], "warn")
	}
	if entry["count"] != float64(42) { // JSON numbers are float64
		t.Errorf("count = %v, want %v", entry["count"], float64(42))
	}
}

func TestErrorLogging(t *testing.T) {
	var buf bytes.Buffer
	log := New(ErrorLevel, "json", &buf)
	log.Error("error message", "err", "something failed")

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["message"] != "error message" {
		t.Errorf("message = %q, want %q", entry["message"], "error message")
	}
	if entry["level"] != "error" {
		t.Errorf("level = %q, want %q", entry["level"], "error")
	}
}

func TestMultipleFields(t *testing.T) {
	var buf bytes.Buffer
	log := New(DebugLevel, "json", &buf)
	log.Info("multi field", "a", 1, "b", "two", "c", true)

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["a"] != float64(1) {
		t.Errorf("a = %v, want %v", entry["a"], float64(1))
	}
	if entry["b"] != "two" {
		t.Errorf("b = %v, want %v", entry["b"], "two")
	}
	if entry["c"] != true {
		t.Errorf("c = %v, want %v", entry["c"], true)
	}
}

func TestOddNumberOfFields(t *testing.T) {
	var buf bytes.Buffer
	log := New(DebugLevel, "json", &buf)
	// Odd number of fields - last one should be ignored
	log.Info("odd fields", "key", "value", "orphan")

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["key"] != "value" {
		t.Errorf("key = %v, want %v", entry["key"], "value")
	}
	// "orphan" should not appear as it has no value
	if _, exists := entry["orphan"]; exists {
		t.Error("orphan field should not exist")
	}
}

func TestTextFormat(t *testing.T) {
	var buf bytes.Buffer
	log := New(InfoLevel, "text", &buf)
	log.Info("text message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "text message") {
		t.Errorf("output missing message: %q", output)
	}
	if !strings.Contains(output, "info") {
		t.Errorf("output missing level: %q", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("output missing field: %q", output)
	}
}

func TestTextFormatMultipleFields(t *testing.T) {
	var buf bytes.Buffer
	log := New(InfoLevel, "text", &buf)
	log.Info("test", "a", 1, "b", 2)

	output := buf.String()
	// Text format should contain both fields
	if !strings.Contains(output, "a=1") {
		t.Errorf("output missing field a: %q", output)
	}
	if !strings.Contains(output, "b=2") {
		t.Errorf("output missing field b: %q", output)
	}
}

func TestTimestampFormat(t *testing.T) {
	var buf bytes.Buffer
	log := New(InfoLevel, "json", &buf)
	log.Info("timestamp test")

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	timestamp, ok := entry["timestamp"].(string)
	if !ok {
		t.Fatalf("timestamp not a string: %v", entry["timestamp"])
	}

	// Verify it's a valid RFC3339 timestamp
	if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
		t.Errorf("timestamp %q is not valid RFC3339: %v", timestamp, err)
	}
}

func TestLevelFiltering(t *testing.T) {
	tests := []struct {
		level       Level
		debugLogged bool
		infoLogged  bool
		warnLogged  bool
		errorLogged bool
	}{
		{DebugLevel, true, true, true, true},
		{InfoLevel, false, true, true, true},
		{WarnLevel, false, false, true, true},
		{ErrorLevel, false, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			var buf bytes.Buffer
			log := New(tt.level, "json", &buf)

			log.Debug("debug")
			hasDebug := strings.Contains(buf.String(), "debug")
			buf.Reset()

			log.Info("info")
			hasInfo := strings.Contains(buf.String(), "info")
			buf.Reset()

			log.Warn("warn")
			hasWarn := strings.Contains(buf.String(), "warn")
			buf.Reset()

			log.Error("error")
			hasError := strings.Contains(buf.String(), "error")

			if hasDebug != tt.debugLogged {
				t.Errorf("Debug logged = %v, want %v", hasDebug, tt.debugLogged)
			}
			if hasInfo != tt.infoLogged {
				t.Errorf("Info logged = %v, want %v", hasInfo, tt.infoLogged)
			}
			if hasWarn != tt.warnLogged {
				t.Errorf("Warn logged = %v, want %v", hasWarn, tt.warnLogged)
			}
			if hasError != tt.errorLogged {
				t.Errorf("Error logged = %v, want %v", hasError, tt.errorLogged)
			}
		})
	}
}
