package xylog_test

import (
	"os"
	"testing"

	"github.com/xybor-x/xylog"
)

func BenchmarkLoggerWithoutLog(b *testing.B) {
	var logger = xylog.GetLogger(b.Name())
	logger.SetLevel(xylog.CRITICAL)
	for i := 0; i < b.N; i++ {
		logger.Debug("msg")
	}
}

func BenchmarkLoggerWithOneHandler(b *testing.B) {
	var handler = xylog.NewHandler("", &CapturedEmitter{})
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
		var handler = xylog.NewHandler("", &CapturedEmitter{})
		handler.SetLevel(xylog.DEBUG)
		logger.AddHandler(handler)
	}
	for i := 0; i < b.N; i++ {
		logger.Critical("msg")
	}
}

func benchEmitter(b *testing.B, logger *xylog.Logger, emitter xylog.Emitter) {
	var handler = xylog.NewHandler("", emitter)
	handler.SetLevel(xylog.DEBUG)
	logger.AddHandler(handler)
	logger.SetLevel(xylog.DEBUG)
	for i := 0; i < b.N; i++ {
		logger.Critical("msg")
	}
}

func BenchmarkLoggerStreamEmitter(b *testing.B) {
	var logger = xylog.GetLogger(b.Name())
	var devnull, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		b.Skipf("can not open /dev/null: %v", err)
	}
	benchEmitter(b, logger, xylog.NewStreamEmitter(devnull))
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
