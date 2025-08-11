package log

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockResponseWriter implements http.ResponseWriter for testing
type mockResponseWriter struct {
	statusCode int
	body       []byte
}

func (m *mockResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	m.body = append(m.body, data...)
	return len(data), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

// MockUser represents a user in our system
type MockUser struct {
	ID       string
	Email    string
	Username string
}

// MockDatabase simulates a database with logging
type MockDatabase struct {
	logger Logger
}

func NewMockDatabase(logger Logger) *MockDatabase {
	return &MockDatabase{
		logger: logger.With(String("component", "database")),
	}
}

func (db *MockDatabase) CreateUser(user *MockUser) error {
	db.logger.Debug("Creating user",
		String("user_id", user.ID),
		String("email", user.Email),
		String("username", user.Username),
	)

	// Simulate database operation
	time.Sleep(10 * time.Millisecond)

	// Simulate success
	db.logger.Info("User created successfully",
		String("user_id", user.ID),
		String("email", user.Email),
	)

	return nil
}

func (db *MockDatabase) GetUser(id string) (*MockUser, error) {
	db.logger.Debug("Fetching user", String("user_id", id))

	// Simulate database operation
	time.Sleep(5 * time.Millisecond)

	// Simulate user found
	user := &MockUser{
		ID:       id,
		Email:    "user@example.com",
		Username: "testuser",
	}

	db.logger.Info("User retrieved successfully",
		String("user_id", id),
		String("email", user.Email),
	)

	return user, nil
}

// MockEmailService simulates an email service with logging
type MockEmailService struct {
	logger Logger
}

func NewMockEmailService(logger Logger) *MockEmailService {
	return &MockEmailService{
		logger: logger.With(String("component", "email_service")),
	}
}

func (es *MockEmailService) SendWelcomeEmail(user *MockUser) error {
	es.logger.Info("Sending welcome email",
		String("user_id", user.ID),
		String("email", user.Email),
	)

	// Simulate email sending
	time.Sleep(20 * time.Millisecond)

	// Simulate success
	es.logger.Info("Welcome email sent successfully",
		String("user_id", user.ID),
		String("email", user.Email),
	)

	return nil
}

// MockAPIServer simulates an API server with request logging
type MockAPIServer struct {
	logger Logger
	db     *MockDatabase
	email  *MockEmailService
}

func NewMockAPIServer(logger Logger) *MockAPIServer {
	return &MockAPIServer{
		logger: logger.With(String("component", "api_server")),
		db:     NewMockDatabase(logger),
		email:  NewMockEmailService(logger),
	}
}

func (s *MockAPIServer) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Extract user data from request (simplified)
	user := &MockUser{
		ID:       "user_123",
		Email:    "newuser@example.com",
		Username: "newuser",
	}

	s.logger.Info("Create user request received",
		String("method", r.Method),
		String("path", r.URL.Path),
		String("user_id", user.ID),
		String("ip", r.RemoteAddr),
	)

	// Create user in database
	if err := s.db.CreateUser(user); err != nil {
		s.logger.Error("Failed to create user in database",
			String("user_id", user.ID),
			Error(err),
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send welcome email
	if err := s.email.SendWelcomeEmail(user); err != nil {
		s.logger.Error("Failed to send welcome email",
			String("user_id", user.ID),
			Error(err),
		)
		// Don't fail the request, just log the error
	}

	duration := time.Since(start)
	s.logger.Info("Create user request completed",
		String("method", r.Method),
		String("path", r.URL.Path),
		String("user_id", user.ID),
		Int("status_code", http.StatusOK),
		Float64("duration_ms", float64(duration.Microseconds())/1000),
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","user_id":"` + user.ID + `"}`))
}

func (s *MockAPIServer) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	userID := "user_123" // Simplified - would normally come from URL params

	s.logger.Info("Get user request received",
		String("method", r.Method),
		String("path", r.URL.Path),
		String("user_id", userID),
		String("ip", r.RemoteAddr),
	)

	user, err := s.db.GetUser(userID)
	if err != nil {
		s.logger.Error("Failed to get user from database",
			String("user_id", userID),
			Error(err),
		)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	duration := time.Since(start)
	s.logger.Info("Get user request completed",
		String("method", r.Method),
		String("path", r.URL.Path),
		String("user_id", userID),
		Int("status_code", http.StatusOK),
		Float64("duration_ms", float64(duration.Microseconds())/1000),
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"user_id":"%s","email":"%s","username":"%s"}`, user.ID, user.Email, user.Username)))
}

// TestIntegration_UserManagement demonstrates a complete user management flow
func TestIntegration_UserManagement(t *testing.T) {
	// Create a logger with JSON output for structured logging
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEncoder(NewJSONEncoder()),
	)

	// Create API server with all components
	server := NewMockAPIServer(logger)

	// Simulate concurrent requests
	var wg sync.WaitGroup
	requests := 5

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			// Create a request logger with request-specific context
			_ = logger.With(
				String("request_id", fmt.Sprintf("req_%d", requestID)),
				String("goroutine", fmt.Sprintf("goroutine_%d", requestID)),
			)

			// Simulate HTTP request
			req, _ := http.NewRequest("POST", "/api/users", nil)
			req.RemoteAddr = "192.168.1.100"

			// Handle create user request
			server.HandleCreateUser(&mockResponseWriter{}, req)

			// Simulate a small delay between requests
			time.Sleep(5 * time.Millisecond)

			// Handle get user request
			req, _ = http.NewRequest("GET", "/api/users/user_123", nil)
			req.RemoteAddr = "192.168.1.101"
			server.HandleGetUser(&mockResponseWriter{}, req)
		}(i)
	}

	wg.Wait()

	// Verify that we have log entries
	output := buf.String()
	logLines := strings.Split(strings.TrimSpace(output), "\n")

	// We should have multiple log entries from different components
	if len(logLines) < 10 {
		t.Errorf("Expected at least 10 log entries, got %d", len(logLines))
	}

	// Verify that different components are logging
	hasDatabaseLogs := false
	hasEmailLogs := false
	hasAPILogs := false

	for _, line := range logLines {
		if strings.Contains(line, `"component":"database"`) {
			hasDatabaseLogs = true
		}
		if strings.Contains(line, `"component":"email_service"`) {
			hasEmailLogs = true
		}
		if strings.Contains(line, `"component":"api_server"`) {
			hasAPILogs = true
		}
	}

	if !hasDatabaseLogs {
		t.Error("Expected database component logs")
	}
	if !hasEmailLogs {
		t.Error("Expected email service component logs")
	}
	if !hasAPILogs {
		t.Error("Expected API server component logs")
	}

	// Verify that we have different log levels
	hasDebug := false
	hasInfo := false

	for _, line := range logLines {
		if strings.Contains(line, `"level":"DEBUG"`) {
			hasDebug = true
		}
		if strings.Contains(line, `"level":"INFO"`) {
			hasInfo = true
		}
	}

	if !hasDebug {
		t.Error("Expected DEBUG level logs")
	}
	if !hasInfo {
		t.Error("Expected INFO level logs")
	}
	// Note: We don't expect ERROR logs in this test since all operations succeed
}

// TestIntegration_ErrorHandling demonstrates error handling with logging
func TestIntegration_ErrorHandling(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelError), // Only log errors for this test
		WithEncoder(NewTextEncoder()),
	)

	// Create a service that will generate errors
	service := struct {
		logger Logger
	}{
		logger: logger.With(String("component", "error_service")),
	}

	// Simulate various error scenarios
	errors := []struct {
		operation string
		err       error
	}{
		{"database_connection", fmt.Errorf("connection timeout")},
		{"file_upload", fmt.Errorf("disk full")},
		{"api_call", fmt.Errorf("rate limit exceeded")},
	}

	for _, scenario := range errors {
		service.logger.Error("Operation failed",
			String("operation", scenario.operation),
			Error(scenario.err),
			String("retry", "true"),
		)
	}

	output := buf.String()
	logLines := strings.Split(strings.TrimSpace(output), "\n")

	if len(logLines) != len(errors) {
		t.Errorf("Expected %d error log entries, got %d", len(errors), len(logLines))
	}

	// Verify error messages are included
	for _, line := range logLines {
		if !strings.Contains(line, "ERROR") {
			t.Error("Expected ERROR level in log line")
		}
		if !strings.Contains(line, "Operation failed") {
			t.Error("Expected error message in log line")
		}
	}
}

// TestIntegration_PerformanceMonitoring demonstrates performance logging
func TestIntegration_PerformanceMonitoring(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithOutput(&buf),
		WithLevel(LevelInfo),
		WithEncoder(NewJSONEncoder()),
	)

	// Simulate performance monitoring
	operations := []struct {
		name     string
		duration time.Duration
		success  bool
	}{
		{"database_query", 50 * time.Millisecond, true},
		{"api_request", 200 * time.Millisecond, true},
		{"file_processing", 100 * time.Millisecond, false},
	}

	for _, op := range operations {
		start := time.Now()
		time.Sleep(op.duration) // Simulate work
		actualDuration := time.Since(start)

		logger.Info("Operation completed",
			String("operation", op.name),
			Float64("duration_ms", float64(actualDuration.Microseconds())/1000),
			Bool("success", op.success),
			Int("status_code", func() int {
				if op.success {
					return 200
				}
				return 500
			}()),
		)
	}

	output := buf.String()
	logLines := strings.Split(strings.TrimSpace(output), "\n")

	if len(logLines) != len(operations) {
		t.Errorf("Expected %d performance log entries, got %d", len(operations), len(logLines))
	}

	// Verify performance data is included
	for _, line := range logLines {
		if !strings.Contains(line, `"duration_ms"`) {
			t.Error("Expected duration_ms field in performance log")
		}
		if !strings.Contains(line, `"operation"`) {
			t.Error("Expected operation field in performance log")
		}
	}
}
