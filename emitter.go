package xylog

import (
	"bufio"
	"fmt"
	"io"
	"runtime/debug"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylock"
)

// LogWriter instances define a writer using to log.
type LogWriter interface {
	io.Writer
	Close() error
}

// Emitter instances dispatch logging events to specific destinations.
type Emitter interface {
	// Emit will be called after a record was decided to log.
	Emit([]byte)

	// Flush writes unflushed buffered data to destination, then closes the
	// Emitter.
	Flush()
}

// StreamEmitter writes logging message to a stream.
type StreamEmitter struct {
	w      LogWriter
	stream *bufio.Writer
	lock   *xylock.Lock
}

// NewStreamEmitter creates a StreamEmitter which writes message to a LogWriter
// (os.Stderr by default).
func NewStreamEmitter(w LogWriter) *StreamEmitter {
	xycond.AssertNotNil(w)
	var e = &StreamEmitter{
		w: w, lock: &xylock.Lock{},
		stream: bufio.NewWriterSize(w, bufferSize),
	}
	return e
}

// Emit will be called after a record was decided to log.
func (e *StreamEmitter) Emit(msg []byte) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.stream == nil {
		return
	}

	var _, err = e.stream.Write(msg)
	if err == nil {
		_, err = e.stream.WriteRune('\n')
	}
	if err != nil {
		fmt.Println("------------ Logging error ------------")
		fmt.Printf("An error occurs when logging: %v\n", err)
		fmt.Println(string(debug.Stack()))
	}
}

// Flush writes unflushed buffered data to destination.
func (e *StreamEmitter) Flush() {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.stream == nil {
		return
	}

	e.stream.Flush()
}

// Close writes unflushed buffered data to destination, then closes the stream.
func (e *StreamEmitter) Close() {
	e.Flush()

	e.lock.Lock()
	defer e.lock.Unlock()
	e.stream = nil
	e.w.Close()
}
