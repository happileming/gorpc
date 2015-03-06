package gorpc

import (
	"log"
	"sync"
	"time"
)

const (
	// DefaultRequestTimeout is the default timeout for client request.
	DefaultRequestTimeout = 20 * time.Second

	// DefaultPendingMessages is the default number of pending messages
	// handled by Client and Server.
	DefaultPendingMessages = 32 * 1024

	// DefaultFlushDelay is the default delay between message flushes
	// on Client and Server.
	DefaultFlushDelay = 5 * time.Millisecond

	// DefaultBufferSize is the default size for Client and Server buffers.
	DefaultBufferSize = 64 * 1024
)

// LoggerFunc is an error logging function to pass to gorpc.SetErrorLogger().
type LoggerFunc func(format string, args ...interface{})

var errorLogger = LoggerFunc(log.Printf)

// SetErrorLogger sets the given error logger to use in gorpc.
//
// By default log.Printf is used for error logging.
func SetErrorLogger(f LoggerFunc) {
	errorLogger = f
}

func logError(format string, args ...interface{}) {
	errorLogger(format, args...)
}

var timerPool sync.Pool

func acquireTimer(timeout time.Duration) *time.Timer {
	tv := timerPool.Get()
	if tv == nil {
		return time.NewTimer(timeout)
	}

	t := tv.(*time.Timer)
	if t.Reset(timeout) {
		panic("BUG: Active timer trapped into acquireTimer()")
	}
	return t
}

func releaseTimer(t *time.Timer) {
	if !t.Stop() {
		// Collect possibly added time from the channel
		// if timer has been stopped and nobody collected its' value.
		select {
		case <-t.C:
		default:
		}
	}

	timerPool.Put(t)
}
