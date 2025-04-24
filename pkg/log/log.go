package main

import (
	"bytes"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Level int

const (
	InfoLevel Level = iota
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelStrings = map[Level]string{
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

type Logger struct {
	level     Level
	writers   []io.Writer
	entryChan chan *logEntry
	wg        sync.WaitGroup
	done      chan struct{}
	once      sync.Once
}

type logEntry struct {
	timestamp time.Time
	level     Level
	file      string
	line      int
	msg       string
}

func NewLogger(level Level, outputs ...string) (*Logger, error) {
	var writers []io.Writer
	for _, output := range outputs {
		switch output {
		case "stdout":
			writers = append(writers, os.Stdout)
		default:
			f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, err
			}
			writers = append(writers, f)
		}
	}

	logger := &Logger{
		level:     level,
		writers:   writers,
		entryChan: make(chan *logEntry, 1000),
		done:      make(chan struct{}),
	}

	logger.wg.Add(1)
	go logger.processEntries()

	return logger, nil
}

func (l *Logger) Close() {
	l.once.Do(func() {
		close(l.done)
		l.wg.Wait()
		for _, w := range l.writers {
			if f, ok := w.(*os.File); ok && f != os.Stdout {
				f.Sync() // Ensure data is flushed to disk
				f.Close()
			}
		}
	})
}

func (l *Logger) processEntries() {
	defer l.wg.Done()
	for {
		select {
		case entry := <-l.entryChan:
			buf := l.formatEntry(entry)
			for _, w := range l.writers {
				w.Write(buf)
			}
		case <-l.done:
			// Drain remaining entries
			for {
				select {
				case entry := <-l.entryChan:
					buf := l.formatEntry(entry)
					for _, w := range l.writers {
						w.Write(buf)
					}
				default:
					return
				}
			}
		}
	}
}

func (l *Logger) formatEntry(entry *logEntry) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 128))

	// Timestamp with milliseconds
	buf.WriteString(entry.timestamp.Format("2006-01-02T15:04:05.000Z07:00"))
	buf.WriteByte(' ')

	// Level
	buf.WriteString(levelStrings[entry.level])
	buf.WriteByte(' ')

	// Message
	buf.WriteString(entry.msg)
	buf.WriteByte(' ')

	// File:Line
	buf.WriteString(entry.file)
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(entry.line))
	buf.WriteByte('\n')

	return buf.Bytes()
}

func (l *Logger) log(level Level, msg string) {
	if level < l.level {
		return
	}

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	file = path.Base(file)

	select {
	case l.entryChan <- &logEntry{
		timestamp: time.Now(),
		level:     level,
		file:      file,
		line:      line,
		msg:       msg,
	}:
	case <-l.done:
	}
}

func (l *Logger) Info(msg string) {
	l.log(InfoLevel, msg)
}

func (l *Logger) Warn(msg string) {
	l.log(WarnLevel, msg)
}

func (l *Logger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

func (l *Logger) Fatal(msg string) {
	if FatalLevel < l.level {
		return
	}

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	file = path.Base(file)

	buf := l.formatEntry(&logEntry{
		timestamp: time.Now(),
		level:     FatalLevel,
		file:      file,
		line:      line,
		msg:       msg,
	})

	for _, w := range l.writers {
		w.Write(buf)
	}
	os.Exit(1)
}
