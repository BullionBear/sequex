package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level represents the log level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Helper functions for creating fields
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err.Error()}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Logger interface defines the logging methods
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	With(fields ...Field) Logger
	SetLevel(level Level)
	SetOutput(w io.Writer)
}

// Encoder interface for different output formats
type Encoder interface {
	Encode(entry *LogEntry) ([]byte, error)
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	Level     Level
	Message   string
	File      string
	Function  string
	Line      int
	Fields    []Field // Changed from map[string]interface{} to []Field to maintain order
}

// JSONEncoder implements JSON format encoding
type JSONEncoder struct{}

func NewJSONEncoder() Encoder {
	return &JSONEncoder{}
}

func (e *JSONEncoder) Encode(entry *LogEntry) ([]byte, error) {
	// Simple JSON encoding - in production you might want to use encoding/json
	fields := make([]string, 0, len(entry.Fields)+6)

	fields = append(fields, fmt.Sprintf(`"timestamp":"%s"`, entry.Timestamp.Format(time.RFC3339)))
	fields = append(fields, fmt.Sprintf(`"level":"%s"`, entry.Level.String()))
	fields = append(fields, fmt.Sprintf(`"message":"%s"`, entry.Message))
	fields = append(fields, fmt.Sprintf(`"file":"%s"`, entry.File))
	fields = append(fields, fmt.Sprintf(`"function":"%s"`, entry.Function))
	fields = append(fields, fmt.Sprintf(`"line":%d`, entry.Line))

	// Add fields in order to maintain consistent output
	for _, field := range entry.Fields {
		switch val := field.Value.(type) {
		case string:
			fields = append(fields, fmt.Sprintf(`"%s":"%s"`, field.Key, val))
		case int, int64:
			fields = append(fields, fmt.Sprintf(`"%s":%v`, field.Key, val))
		case float64:
			fields = append(fields, fmt.Sprintf(`"%s":%f`, field.Key, val))
		case bool:
			fields = append(fields, fmt.Sprintf(`"%s":%t`, field.Key, val))
		default:
			fields = append(fields, fmt.Sprintf(`"%s":"%v"`, field.Key, val))
		}
	}

	return []byte("{" + strings.Join(fields, ",") + "}\n"), nil
}

// TextEncoder implements text format encoding
type TextEncoder struct{}

func NewTextEncoder() Encoder {
	return &TextEncoder{}
}

func (e *TextEncoder) Encode(entry *LogEntry) ([]byte, error) {
	// Format: timestamp level file:line function > message key=value
	parts := []string{
		entry.Timestamp.Format(time.RFC3339),
		entry.Level.String(),
		fmt.Sprintf("%s:%d", entry.File, entry.Line),
		entry.Function,
		">",
		entry.Message,
	}

	// Add fields as key=value pairs in order to maintain consistent output
	for _, field := range entry.Fields {
		parts = append(parts, fmt.Sprintf("%s=%v", field.Key, field.Value))
	}

	return []byte(strings.Join(parts, " ") + "\n"), nil
}

// TimeRotateWriter handles time-based log rotation
type TimeRotateWriter struct {
	filename    string
	current     *os.File
	lastRotate  time.Time
	rotateEvery time.Duration
	maxBackups  int
	mu          sync.Mutex
}

func NewTimeRotateWriter(filename string, rotateEvery time.Duration, maxBackups int) (*TimeRotateWriter, error) {
	w := &TimeRotateWriter{
		filename:    filename,
		rotateEvery: rotateEvery,
		maxBackups:  maxBackups,
	}
	if err := w.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *TimeRotateWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotateIfNeeded(); err != nil {
		return 0, err
	}
	return w.current.Write(p)
}

func (w *TimeRotateWriter) rotateIfNeeded() error {
	now := time.Now()
	if w.current == nil || now.Sub(w.lastRotate) >= w.rotateEvery {
		if w.current != nil {
			w.current.Close()
		}

		// Rename old log if exists
		if _, err := os.Stat(w.filename); err == nil {
			ext := filepath.Ext(w.filename)
			base := strings.TrimSuffix(w.filename, ext)
			newName := fmt.Sprintf("%s-%s%s", base, w.lastRotate.Format("20060102"), ext)
			os.Rename(w.filename, newName)
		}

		// Open new log file
		f, err := os.OpenFile(w.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		w.current = f
		w.lastRotate = now

		// Clean up old logs
		w.cleanupOldLogs()
	}
	return nil
}

func (w *TimeRotateWriter) cleanupOldLogs() {
	if w.maxBackups <= 0 {
		return
	}

	dir := filepath.Dir(w.filename)
	base := strings.TrimSuffix(filepath.Base(w.filename), filepath.Ext(w.filename))
	ext := filepath.Ext(w.filename)

	pattern := filepath.Join(dir, base+"-*"+ext)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	// Sort by modification time (oldest first)
	type fileInfo struct {
		path    string
		modTime time.Time
	}

	var files []fileInfo
	for _, match := range matches {
		if info, err := os.Stat(match); err == nil {
			files = append(files, fileInfo{match, info.ModTime()})
		}
	}

	// Remove oldest files if we have too many
	if len(files) > w.maxBackups {
		// Sort by modification time (oldest first)
		for i := 0; i < len(files)-w.maxBackups; i++ {
			os.Remove(files[i].path)
		}
	}
}

func (w *TimeRotateWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.current != nil {
		return w.current.Close()
	}
	return nil
}

// logger implements the Logger interface
type logger struct {
	level      Level
	output     io.Writer
	encoder    Encoder
	mu         sync.Mutex
	callerSkip int
	fields     []Field // Changed from map[string]interface{} to []Field to maintain order
}

// Option is a function that configures a logger
type Option func(*logger)

// New creates a new logger with the given options
func New(opts ...Option) Logger {
	l := &logger{
		level:      LevelInfo,
		output:     os.Stdout,
		encoder:    NewJSONEncoder(),
		callerSkip: 2,
		fields:     make([]Field, 0),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// WithLevel sets the log level
func WithLevel(level Level) Option {
	return func(l *logger) {
		l.level = level
	}
}

// WithOutput sets the output writer
func WithOutput(w io.Writer) Option {
	return func(l *logger) {
		l.output = w
	}
}

// WithEncoder sets the encoder
func WithEncoder(encoder Encoder) Option {
	return func(l *logger) {
		l.encoder = encoder
	}
}

// WithTimeRotation enables time-based log rotation
func WithTimeRotation(dir, filename string, rotateEvery time.Duration, maxBackups int) Option {
	return func(l *logger) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(fmt.Sprintf("failed to create log directory: %v", err))
		}
		fullPath := filepath.Join(dir, filename)
		writer, err := NewTimeRotateWriter(fullPath, rotateEvery, maxBackups)
		if err != nil {
			panic(fmt.Sprintf("failed to create time rotate writer: %v", err))
		}
		l.output = writer
	}
}

// WithCallerSkip sets the number of call frames to skip
func WithCallerSkip(skip int) Option {
	return func(l *logger) {
		l.callerSkip = skip
	}
}

func (l *logger) log(level Level, msg string, fields ...Field) {
	if level < l.level {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(l.callerSkip)
	if !ok {
		file = "unknown"
		line = 0
	}

	// Extract function name
	pc, _, _, ok := runtime.Caller(l.callerSkip)
	var function string
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			function = fn.Name()
		}
	}
	if function == "" {
		function = "unknown"
	}

	// Combine fields
	allFields := make([]Field, 0, len(l.fields)+len(fields))
	for _, field := range l.fields {
		allFields = append(allFields, field)
	}
	for _, field := range fields {
		allFields = append(allFields, field)
	}

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		File:      filepath.Base(file),
		Function:  function,
		Line:      line,
		Fields:    allFields,
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := l.encoder.Encode(entry)
	if err != nil {
		// Fallback to simple output
		fmt.Fprintf(l.output, "%s %s %s:%d > %s\n",
			entry.Timestamp.Format(time.RFC3339),
			entry.Level.String(),
			entry.File, entry.Line, entry.Message)
		return
	}

	l.output.Write(data)

	if level == LevelFatal {
		os.Exit(1)
	}
}

func (l *logger) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

func (l *logger) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

func (l *logger) Warn(msg string, fields ...Field) {
	l.log(LevelWarn, msg, fields...)
}

func (l *logger) Error(msg string, fields ...Field) {
	l.log(LevelError, msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...Field) {
	l.log(LevelFatal, msg, fields...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.log(LevelDebug, fmt.Sprintf(format, args...))
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.log(LevelInfo, fmt.Sprintf(format, args...))
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.log(LevelWarn, fmt.Sprintf(format, args...))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log(LevelError, fmt.Sprintf(format, args...))
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.log(LevelFatal, fmt.Sprintf(format, args...))
}

func (l *logger) With(fields ...Field) Logger {
	newLogger := &logger{
		level:      l.level,
		output:     l.output,
		encoder:    l.encoder,
		callerSkip: l.callerSkip,
		fields:     make([]Field, 0, len(l.fields)+len(fields)),
	}

	// Copy existing fields
	for _, field := range l.fields {
		newLogger.fields = append(newLogger.fields, field)
	}

	// Add new fields
	for _, field := range fields {
		newLogger.fields = append(newLogger.fields, field)
	}

	return newLogger
}

func (l *logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
}
