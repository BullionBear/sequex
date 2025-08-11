package log

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLogger_StructuredLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEncoder(NewJSONEncoder()),
	)

	logger.Info("User logged in", String("user_id", "123"), String("ip", "192.168.1.1"))

	output := buf.String()
	if !strings.Contains(output, `"message":"User logged in"`) {
		t.Errorf("Expected message in output, got: %s", output)
	}
	if !strings.Contains(output, `"user_id":"123"`) {
		t.Errorf("Expected user_id field in output, got: %s", output)
	}
	if !strings.Contains(output, `"ip":"192.168.1.1"`) {
		t.Errorf("Expected ip field in output, got: %s", output)
	}
	if !strings.Contains(output, `"level":"INFO"`) {
		t.Errorf("Expected INFO level in output, got: %s", output)
	}
}

func TestLogger_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEncoder(NewTextEncoder()),
	)

	logger.Warn("Database connection failed", String("db", "postgres"), Int("retries", 3))

	output := buf.String()
	if !strings.Contains(output, "WARN") {
		t.Errorf("Expected WARN level in output, got: %s", output)
	}
	if !strings.Contains(output, "Database connection failed") {
		t.Errorf("Expected message in output, got: %s", output)
	}
	if !strings.Contains(output, "db=postgres") {
		t.Errorf("Expected db field in output, got: %s", output)
	}
	if !strings.Contains(output, "retries=3") {
		t.Errorf("Expected retries field in output, got: %s", output)
	}
}

func TestLogger_FormatStringSupport(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEncoder(NewTextEncoder()),
	)

	logger.Infof("Processing %d items for user %s", 42, "alice")

	output := buf.String()
	if !strings.Contains(output, "Processing 42 items for user alice") {
		t.Errorf("Expected formatted message in output, got: %s", output)
	}
}

func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEncoder(NewJSONEncoder()),
	)

	userLogger := logger.With(String("user_id", "123"), String("session_id", "abc"))
	userLogger.Info("Action performed", String("action", "click"))

	output := buf.String()
	if !strings.Contains(output, `"user_id":"123"`) {
		t.Errorf("Expected user_id from With() in output, got: %s", output)
	}
	if !strings.Contains(output, `"session_id":"abc"`) {
		t.Errorf("Expected session_id from With() in output, got: %s", output)
	}
	if !strings.Contains(output, `"action":"click"`) {
		t.Errorf("Expected action field in output, got: %s", output)
	}
}

func TestLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelWarn), // Only WARN and above
		WithEncoder(NewTextEncoder()),
	)

	logger.Debug("Debug message") // Should be filtered out
	logger.Info("Info message")   // Should be filtered out
	logger.Warn("Warn message")   // Should be included
	logger.Error("Error message") // Should be included

	output := buf.String()
	if strings.Contains(output, "Debug message") {
		t.Errorf("Debug message should be filtered out, got: %s", output)
	}
	if strings.Contains(output, "Info message") {
		t.Errorf("Info message should be filtered out, got: %s", output)
	}
	if !strings.Contains(output, "Warn message") {
		t.Errorf("Warn message should be included, got: %s", output)
	}
	if !strings.Contains(output, "Error message") {
		t.Errorf("Error message should be included, got: %s", output)
	}
}

func TestLogger_FieldHelpers(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEncoder(NewJSONEncoder()),
	)

	logger.Info("Test fields",
		String("string_field", "value"),
		Int("int_field", 42),
		Int64("int64_field", 123456789),
		Float64("float_field", 3.14),
		Bool("bool_field", true),
		Error(fmt.Errorf("test error")),
		Any("any_field", map[string]string{"key": "value"}),
	)

	output := buf.String()
	if !strings.Contains(output, `"string_field":"value"`) {
		t.Errorf("Expected string_field in output, got: %s", output)
	}
	if !strings.Contains(output, `"int_field":42`) {
		t.Errorf("Expected int_field in output, got: %s", output)
	}
	if !strings.Contains(output, `"int64_field":123456789`) {
		t.Errorf("Expected int64_field in output, got: %s", output)
	}
	if !strings.Contains(output, `"float_field":3.140000`) {
		t.Errorf("Expected float_field in output, got: %s", output)
	}
	if !strings.Contains(output, `"bool_field":true`) {
		t.Errorf("Expected bool_field in output, got: %s", output)
	}
	if !strings.Contains(output, `"error":"test error"`) {
		t.Errorf("Expected error field in output, got: %s", output)
	}
}

func TestTimeRotateWriter_Basic(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	writer, err := NewTimeRotateWriter(filename, 24*time.Hour, 3)
	if err != nil {
		t.Fatalf("Failed to create TimeRotateWriter: %v", err)
	}
	defer writer.Close()

	// Write some data
	data := []byte("test log entry\n")
	n, err := writer.Write(data)
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}

	// Check that file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Log file was not created")
	}
}

func TestTimeRotateWriter_Rotation(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	// Create writer with very short rotation interval for testing
	writer, err := NewTimeRotateWriter(filename, 1*time.Millisecond, 2)
	if err != nil {
		t.Fatalf("Failed to create TimeRotateWriter: %v", err)
	}
	defer writer.Close()

	// Write initial data
	writer.Write([]byte("first entry\n"))

	// Wait for rotation
	time.Sleep(2 * time.Millisecond)

	// Write more data (should trigger rotation)
	writer.Write([]byte("second entry\n"))

	// Check that rotated file exists
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	rotatedPattern := base + "-*" + ext
	matches, err := filepath.Glob(rotatedPattern)
	if err != nil {
		t.Fatalf("Failed to glob rotated files: %v", err)
	}
	if len(matches) == 0 {
		t.Errorf("Expected rotated file to exist")
	}
}

func TestLogger_ContextualInformation(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEncoder(NewJSONEncoder()),
	)

	logger.Info("Test contextual info")

	output := buf.String()
	if !strings.Contains(output, `"file":"logger_test.go"`) {
		t.Errorf("Expected file information in output, got: %s", output)
	}
	if !strings.Contains(output, `"function":"github.com/BullionBear/sequex/pkg/log.TestLogger_ContextualInformation"`) {
		t.Errorf("Expected function information in output, got: %s", output)
	}
	if !strings.Contains(output, `"line":`) {
		t.Errorf("Expected line number in output, got: %s", output)
	}
	if !strings.Contains(output, `"timestamp"`) {
		t.Errorf("Expected timestamp in output, got: %s", output)
	}
}

func TestLogger_ThreadSafety(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEncoder(NewTextEncoder()),
	)

	// Test concurrent logging
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				logger.Info("Concurrent log", Int("goroutine", id), Int("iteration", j))
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Count log entries
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	expectedLines := 1000 // 10 goroutines * 100 iterations each
	if len(lines) != expectedLines {
		t.Errorf("Expected %d log entries, got %d", expectedLines, len(lines))
	}
}

func TestLogger_SetLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelInfo),
		WithEncoder(NewTextEncoder()),
	)

	// Initially, debug should be filtered
	logger.Debug("Debug message")
	if strings.Contains(buf.String(), "Debug message") {
		t.Errorf("Debug message should be filtered initially")
	}

	// Change level to debug
	logger.SetLevel(LevelDebug)
	buf.Reset()

	// Now debug should be included
	logger.Debug("Debug message")
	if !strings.Contains(buf.String(), "Debug message") {
		t.Errorf("Debug message should be included after level change")
	}
}

func TestLogger_SetOutput(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	logger := New(
		WithOutput(&buf1),
		WithLevel(LevelInfo),
		WithEncoder(NewTextEncoder()),
	)

	logger.Info("Message to first output")
	if !strings.Contains(buf1.String(), "Message to first output") {
		t.Errorf("Message should be in first output")
	}

	// Change output
	logger.SetOutput(&buf2)
	logger.Info("Message to second output")

	if strings.Contains(buf2.String(), "Message to first output") {
		t.Errorf("First message should not be in second output")
	}
	if !strings.Contains(buf2.String(), "Message to second output") {
		t.Errorf("Second message should be in second output")
	}
}

// Benchmark tests
func BenchmarkLogger_StructuredLogging(b *testing.B) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEncoder(NewJSONEncoder()),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message",
			String("key1", "value1"),
			Int("key2", i),
			Float64("key3", 3.14),
		)
	}
}

func BenchmarkLogger_TextFormat(b *testing.B) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEncoder(NewTextEncoder()),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message",
			String("key1", "value1"),
			Int("key2", i),
			Float64("key3", 3.14),
		)
	}
}

func BenchmarkLogger_FormatString(b *testing.B) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEncoder(NewTextEncoder()),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("Benchmark message %d with value %s", i, "test")
	}
}
