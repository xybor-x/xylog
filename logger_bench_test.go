package xylog_test

import (
	"os"
	"testing"

	"github.com/xybor-x/xylog"
)

var DevnullEmitter *xylog.StreamEmitter

func init() {
	var devnull, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err == nil {
		DevnullEmitter = xylog.NewStreamEmitter(devnull)
	}
}

func BenchmarkLoggerWithoutLog(b *testing.B) {
	var logger = xylog.GetLogger(b.Name())
	logger.SetLevel(xylog.CRITICAL)
	for i := 0; i < b.N; i++ {
		logger.Debug("msg")
	}
}

func BenchmarkEventLoggerWithoutLog(b *testing.B) {
	var logger = xylog.GetLogger(b.Name())
	logger.SetLevel(xylog.CRITICAL)
	for i := 0; i < b.N; i++ {
		logger.Event("foo").Field("foo", "bar").Field("bar", "foo").Debug()
	}
}

func BenchmarkLoggerWithOneHandler(b *testing.B) {
	var handler = xylog.GetHandler("")
	handler.AddEmitter(&CapturedEmitter{})
	handler.SetLevel(xylog.DEBUG)
	var logger = xylog.GetLogger(b.Name())
	logger.SetLevel(xylog.DEBUG)
	logger.AddHandler(handler)
	for i := 0; i < b.N; i++ {
		logger.Critical("msg")
	}
}

func BenchmarkLoggerWithMultiHandler(b *testing.B) {
	var logger = xylog.GetLogger(b.Name())
	logger.SetLevel(xylog.DEBUG)
	for i := 0; i < 100; i++ {
		var handler = xylog.GetHandler("")
		handler.AddEmitter(&CapturedEmitter{})
		handler.SetLevel(xylog.DEBUG)
		logger.AddHandler(handler)
	}
	for i := 0; i < b.N; i++ {
		logger.Critical("msg")
	}
}

func benchEmitter(b *testing.B, logger *xylog.Logger, emitter xylog.Emitter) {
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)
	handler.SetLevel(xylog.DEBUG)
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	for i := 0; i < b.N; i++ {
		logger.Critical("msg")
	}
}

func BenchmarkLoggerStreamEmitter(b *testing.B) {
	var logger = xylog.GetLogger(b.Name())
	if DevnullEmitter == nil {
		b.Skipf("can not open /dev/null")
	}
	benchEmitter(b, logger, DevnullEmitter)
}
func BenchmarkLoggerFileEmitter(b *testing.B) {
	var logger = xylog.GetLogger(b.Name())
	var emitter = xylog.NewFileEmitter("example.log")
	benchEmitter(b, logger, emitter)
}

func BenchmarkLoggerSizeEmitter(b *testing.B) {
	var logger = xylog.GetLogger(b.Name())
	var emitter = xylog.NewSizeRotatingFileEmitter(
		"example.log", 1024*1024, 3)
	benchEmitter(b, logger, emitter)
}
