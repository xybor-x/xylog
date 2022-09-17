// Copyright (c) 2022 xybor-x
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package xylog

import (
	"bufio"
	"fmt"
	"io"
	"runtime/debug"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylock"
)

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
	w    *bufio.Writer
	lock *xylock.Lock
}

// NewStreamEmitter creates a StreamEmitter which writes message to a LogWriter
// (os.Stderr by default).
func NewStreamEmitter(w io.Writer) *StreamEmitter {
	xycond.AssertNotNil(w)

	var size = globalLock.RLockFunc(func() any { return bufferSize }).(int)
	var e = &StreamEmitter{
		lock: &xylock.Lock{},
		w:    bufio.NewWriterSize(w, size),
	}

	globalLock.WLockFunc(func() {
		emitterManager = append(emitterManager, e)
	})
	return e
}

// Emit will be called after a record was decided to log.
func (e *StreamEmitter) Emit(msg []byte) {
	e.lock.Lock()
	defer e.lock.Unlock()

	var _, err = e.w.Write(msg)
	if err == nil {
		_, err = e.w.WriteRune('\n')
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

	e.w.Flush()
}
