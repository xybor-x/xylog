package xylog_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xylog"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 1000; i++ {
		messages = append(messages,
			fmt.Sprintf("This is a long enough message %d", i))
	}
}

func withLogger(t *testing.T, f func(logger *xylog.Logger, w *mockWriter)) {
	xylog.SetBufferSize(1)
	defer xylog.SetBufferSize(4096)
	var writer = &mockWriter{}
	var emitter = xylog.NewStreamEmitter(writer)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	var logger = xylog.GetLogger(t.Name())
	logger.AddHandler(handler)

	f(logger, writer)
}

func withStreamEmitter(
	t *testing.T, f func(e *xylog.StreamEmitter, w *mockWriter),
) {
	xylog.SetBufferSize(1)
	defer xylog.SetBufferSize(4096)
	var writer = &mockWriter{}
	var emitter = xylog.NewStreamEmitter(writer)
	f(emitter, writer)
}

type mockWriter struct {
	Captured string
	Error    bool
}

func (w *mockWriter) Write(b []byte) (int, error) {
	if w.Error {
		return 0, xyerror.BaseException.New("mockwriter raised an error")
	}

	w.Captured += string(b)
	return len(b), nil
}

func (w *mockWriter) Close() error {
	return nil
}

func (w *mockWriter) Reset() {
	w.Captured = ""
}

type LoggerNameFilter struct {
	name string
}

func (f *LoggerNameFilter) Filter(record xylog.LogRecord) bool {
	return record.Name == f.name
}

var messages []string

func getRandomMessage() string {
	return messages[rand.Int()%len(messages)]
}

func addFullMacros(f xylog.Formatter) xylog.Formatter {
	return f.AddMacro("asctime", "asctime").
		AddMacro("created", "created").
		AddMacro("filename", "filename").
		AddMacro("funcname", "funcname").
		AddMacro("levelname", "levelname").
		AddMacro("levelno", "levelno").
		AddMacro("lineno", "lineno").
		AddMacro("module", "module").
		AddMacro("msecs", "msecs").
		AddMacro("pathname", "pathname").
		AddMacro("process", "process").
		AddMacro("relativeCreated", "relativeCreated")
}
