package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLevelFiltering(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "level.log")

	logger, err := NewLogger(WarnLevel, logFile)
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	logger.Info("should not appear")
	logger.Warn("should appear")
	logger.Error("should appear")

	// Add small delay to ensure logs are processed
	time.Sleep(100 * time.Millisecond)

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	logs := string(content)
	t.Logf("Log file content:\n%s", logs) // Debug output

	if strings.Contains(logs, "INFO") {
		t.Error("Info level message was not filtered")
	}
	if !strings.Contains(logs, "WARN") {
		t.Error("Warn level message missing")
	}
	if !strings.Contains(logs, "ERROR") {
		t.Error("Error level message missing")
	}
}

func TestOutputDestinations(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "multi.log")

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	logger, err := NewLogger(InfoLevel, "stdout", logFile)
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	logger.Info("multi-output test")
	logger.Close()

	w.Close()
	var stdoutBuf bytes.Buffer
	io.Copy(&stdoutBuf, r)

	// Check stdout
	if !strings.Contains(stdoutBuf.String(), "multi-output test") {
		t.Error("Message missing from stdout")
	}

	// Check file
	fileContent, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(fileContent), "multi-output test") {
		t.Error("Message missing from log file")
	}
}

func TestFatalExits(t *testing.T) {
	if os.Getenv("TEST_FATAL") == "1" {
		logger, _ := NewLogger(InfoLevel, "stdout")
		logger.Fatal("fatal test")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatalExits")
	cmd.Env = append(os.Environ(), "TEST_FATAL=1")
	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && e.ExitCode() == 1 {
		if !strings.Contains(stdoutBuf.String(), "FATAL fatal test") {
			t.Error("Fatal message missing from output")
		}
		return
	}
	t.Fatalf("Process exited with %v, want exit status 1", err)
}

func TestFileLineFormat(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "fileline.log")

	logger, err := NewLogger(InfoLevel, logFile)
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	_, file, line, _ := runtime.Caller(0)
	logger.Info("line test")
	logger.Close()

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedFile := filepath.Base(file)
	expectedLine := line + 1 // Caller(0) + 1 line for log call
	expectedPattern := expectedFile + ":" + strconv.Itoa(expectedLine)

	if !strings.Contains(string(content), expectedPattern) {
		t.Errorf("Missing file:line pattern %q in log: %q", expectedPattern, string(content))
	}
}

func TestTimestampFormat(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "timestamp.log")

	logger, err := NewLogger(InfoLevel, logFile)
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	// Use UTC for consistency
	before := time.Now().UTC()
	logger.Info("timestamp test")
	logger.Close()
	after := time.Now().UTC()

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	// Extract timestamp from log
	logLine := string(content)
	tsStr := strings.Split(logLine, " ")[0]
	t.Logf("Parsing timestamp string: %s", tsStr) // Debug output

	// Parse with location set to UTC
	ts, err := time.Parse("2006-01-02T15:04:05.000Z07:00", tsStr)
	if err != nil {
		t.Fatalf("Failed to parse timestamp: %v", err)
	}

	// Convert parsed time to UTC for comparison
	tsUTC := ts.UTC()
	t.Logf("Parsed timestamp (UTC): %v", tsUTC) // Debug output

	// Allow small buffer for test execution time (1 second)
	adjustedBefore := before.Add(-1 * time.Second)
	adjustedAfter := after.Add(1 * time.Second)

	if tsUTC.Before(adjustedBefore) || tsUTC.After(adjustedAfter) {
		t.Errorf("Timestamp %v (UTC) not in expected range [%v - %v]", tsUTC, adjustedBefore, adjustedAfter)
	}
}

func TestCloseDrainsLogs(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "drain.log")

	logger, err := NewLogger(InfoLevel, logFile)
	if err != nil {
		t.Fatal(err)
	}

	const numLogs = 1000
	for i := 0; i < numLogs; i++ {
		logger.Info("drain test")
	}

	logger.Close()

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Count(string(content), "\n")
	if lines != numLogs {
		t.Errorf("Expected %d lines, got %d", numLogs, lines)
	}
}

func TestConcurrentLogging(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "concurrent.log")

	logger, err := NewLogger(InfoLevel, logFile)
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	var wg sync.WaitGroup
	const numLogs = 1000
	for i := 0; i < numLogs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.Info("concurrent test")
		}()
	}
	wg.Wait()
	logger.Close()

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Count(string(content), "\n")
	if lines != numLogs {
		t.Errorf("Expected %d lines, got %d", numLogs, lines)
	}
}

func TestInvalidFile(t *testing.T) {
	_, err := NewLogger(InfoLevel, "/invalid/path/file.log")
	if err == nil {
		t.Error("Expected error for invalid file path")
	}
}
