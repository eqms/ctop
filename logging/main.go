package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
)

const (
	ringSize = 1024
)

var (
	Log *CTopLogger
)

type statusMsg struct {
	Text    string
	IsError bool
}

// ringBuffer is a fixed-size circular buffer for log messages
// that supports multiple concurrent readers via channels.
type ringBuffer struct {
	mu        sync.Mutex
	buf       []string
	pos       int
	full      bool
	listeners []chan string
}

func newRingBuffer(size int) *ringBuffer {
	return &ringBuffer{
		buf: make([]string, size),
	}
}

func (rb *ringBuffer) write(msg string) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.buf[rb.pos] = msg
	rb.pos = (rb.pos + 1) % len(rb.buf)
	if rb.pos == 0 {
		rb.full = true
	}

	// broadcast to all listeners
	for _, ch := range rb.listeners {
		select {
		case ch <- msg:
		default:
			// drop message if listener is too slow
		}
	}
}

// subscribe returns a channel that receives new log messages.
// Call unsubscribe with the returned channel when done.
func (rb *ringBuffer) subscribe() chan string {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	ch := make(chan string, 64)

	// send existing messages first
	if rb.full {
		for i := rb.pos; i < len(rb.buf); i++ {
			ch <- rb.buf[i]
		}
	}
	for i := 0; i < rb.pos; i++ {
		ch <- rb.buf[i]
	}

	rb.listeners = append(rb.listeners, ch)
	return ch
}

func (rb *ringBuffer) unsubscribe(ch chan string) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	for i, l := range rb.listeners {
		if l == ch {
			rb.listeners = append(rb.listeners[:i], rb.listeners[i+1:]...)
			close(ch)
			return
		}
	}
}

// ctopHandler is a custom slog.Handler that writes formatted messages
// to the ring buffer and optionally to a file.
type ctopHandler struct {
	level  *slog.LevelVar
	ring   *ringBuffer
	file   *os.File
	mu     sync.Mutex
	attrs  []slog.Attr
	groups []string
}

func (h *ctopHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level.Level()
}

func (h *ctopHandler) Handle(_ context.Context, r slog.Record) error {
	// format: "15:04:05.000 ▶ LEVL message"
	levelStr := r.Level.String()
	if len(levelStr) > 4 {
		levelStr = levelStr[:4]
	}
	msg := fmt.Sprintf("%s ▶ %s %s",
		r.Time.Format("15:04:05.000"),
		levelStr,
		r.Message,
	)

	h.ring.write(msg)

	h.mu.Lock()
	if h.file != nil {
		fmt.Fprintln(h.file, msg)
	}
	h.mu.Unlock()

	return nil
}

func (h *ctopHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ctopHandler{
		level:  h.level,
		ring:   h.ring,
		file:   h.file,
		attrs:  append(h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *ctopHandler) WithGroup(name string) slog.Handler {
	return &ctopHandler{
		level:  h.level,
		ring:   h.ring,
		file:   h.file,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}

// CTopLogger wraps slog.Logger with convenience methods matching
// the previous op/go-logging API, plus status message support.
type CTopLogger struct {
	*slog.Logger
	level   *slog.LevelVar
	ring    *ringBuffer
	logFile *os.File
	sLog    []statusMsg
	done    chan struct{}
}

// Convenience logging methods (matching the old op/go-logging API)

func (c *CTopLogger) Debugf(format string, args ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(format, args...))
}

func (c *CTopLogger) Infof(format string, args ...interface{}) {
	c.Logger.Info(fmt.Sprintf(format, args...))
}

func (c *CTopLogger) Info(msg string) {
	c.Logger.Info(msg)
}

func (c *CTopLogger) Warningf(format string, args ...interface{}) {
	c.Logger.Warn(fmt.Sprintf(format, args...))
}

func (c *CTopLogger) Errorf(format string, args ...interface{}) {
	c.Logger.Error(fmt.Sprintf(format, args...))
}

func (c *CTopLogger) Notice(msg string) {
	c.Logger.Info(msg)
}

func (c *CTopLogger) Noticef(format string, args ...interface{}) {
	c.Logger.Info(fmt.Sprintf(format, args...))
}

func (c *CTopLogger) IsDebugEnabled() bool {
	return c.level.Level() <= slog.LevelDebug
}

// Status message methods (unchanged API)

func (c *CTopLogger) FlushStatus() chan statusMsg {
	ch := make(chan statusMsg)
	go func() {
		for _, sm := range c.sLog {
			ch <- sm
		}
		close(ch)
		c.sLog = []statusMsg{}
	}()
	return ch
}

func (c *CTopLogger) StatusQueued() bool     { return len(c.sLog) > 0 }
func (c *CTopLogger) Status(s string)        { c.addStatus(statusMsg{s, false}) }
func (c *CTopLogger) StatusErr(err error)    { c.addStatus(statusMsg{err.Error(), true}) }
func (c *CTopLogger) addStatus(sm statusMsg) { c.sLog = append(c.sLog, sm) }

func (c *CTopLogger) Statusf(s string, a ...interface{}) { c.Status(fmt.Sprintf(s, a...)) }

func Init() *CTopLogger {
	if Log == nil {
		levelVar := &slog.LevelVar{}
		levelVar.Set(slog.LevelInfo)

		ring := newRingBuffer(ringSize)

		debugOn := debugMode()
		if debugOn {
			levelVar.Set(slog.LevelDebug)
		}

		handler := &ctopHandler{
			level: levelVar,
			ring:  ring,
		}

		logFilePath := debugModeFile()
		if logFilePath != "" {
			logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to create log file: %s\n", err)
			} else {
				handler.file = logFile
			}
		}

		logger := slog.New(handler)

		Log = &CTopLogger{
			Logger:  logger,
			level:   levelVar,
			ring:    ring,
			logFile: handler.file,
			sLog:    []statusMsg{},
			done:    make(chan struct{}),
		}

		if debugOn {
			StartServer()
		}
		Log.Notice("logger initialized")
	}
	return Log
}

func (c *CTopLogger) tail() chan string {
	return c.ring.subscribe()
}

func (c *CTopLogger) untail(ch chan string) {
	c.ring.unsubscribe(ch)
}

func (c *CTopLogger) Exit() {
	close(c.done)
	// give listeners a moment to drain
	time.Sleep(100 * time.Millisecond)
	if c.logFile != nil {
		_ = c.logFile.Close()
	}
	StopServer()
}

func debugMode() bool       { return os.Getenv("CTOP_DEBUG") == "1" }
func debugModeTCP() bool    { return os.Getenv("CTOP_DEBUG_TCP") == "1" }
func debugModeFile() string { return os.Getenv("CTOP_DEBUG_FILE") }
