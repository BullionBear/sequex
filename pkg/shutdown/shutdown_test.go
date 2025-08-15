package shutdown

import (
	"os"
	"testing"
	"time"

	"github.com/BullionBear/sequex/pkg/log"
)

func TestShutdownWithTimeout(t *testing.T) {
	logger := log.New(
		log.WithLevel(log.LevelDebug),
		log.WithEncoder(log.NewTextEncoder()),
		log.WithOutput(os.Stdout),
		log.WithCallerSkip(2),
	)
	shutdown := NewShutdown(logger)

	// Track completion status
	quickCompleted := false
	slowCompleted := false
	timeoutOccurred := false

	// Hook a quick callback
	shutdown.HookShutdownCallback("quick", func() {
		time.Sleep(50 * time.Millisecond)
		quickCompleted = true
	}, 1*time.Second)

	// Hook a slow callback that will timeout
	shutdown.HookShutdownCallback("slow", func() {
		time.Sleep(2 * time.Second) // This will timeout
		slowCompleted = true
	}, 100*time.Millisecond)

	// Hook a callback to detect timeout
	shutdown.HookShutdownCallback("timeout-detector", func() {
		time.Sleep(200 * time.Millisecond)
		timeoutOccurred = true
	}, 50*time.Millisecond)

	// Trigger shutdown
	shutdown.ShutdownNow()

	// Verify results
	if !quickCompleted {
		t.Error("Quick callback should have completed")
	}

	if slowCompleted {
		t.Error("Slow callback should not have completed due to timeout")
	}

	if timeoutOccurred {
		t.Error("Timeout detector should not have completed due to timeout")
	}
}

func TestShutdownWithoutTimeout(t *testing.T) {
	logger := log.New(
		log.WithLevel(log.LevelDebug),
		log.WithEncoder(log.NewTextEncoder()),
		log.WithOutput(os.Stdout),
		log.WithCallerSkip(2),
	)
	shutdown := NewShutdown(logger)

	completed := false

	// Hook a callback without timeout
	shutdown.HookShutdownCallback("no-timeout", func() {
		time.Sleep(100 * time.Millisecond)
		completed = true
	}, 0) // No timeout

	// Trigger shutdown
	shutdown.ShutdownNow()

	if !completed {
		t.Error("Callback without timeout should have completed")
	}
}
