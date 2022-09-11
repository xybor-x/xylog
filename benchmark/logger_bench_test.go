package benchmark

import (
	"testing"

	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/test"
)

var DevnullEmitter *xylog.StreamEmitter

func BenchmarkLoggerDisable(b *testing.B) {
	test.WithBenchLogger(b, func(logger *xylog.Logger) {
		logger.SetLevel(xylog.ERROR)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(test.GetRandomMessage())
			}
		})
	})
}

func BenchmarkLoggerWithoutHandler(b *testing.B) {
	test.WithBenchLogger(b, func(logger *xylog.Logger) {
		logger.RemoveHandler(logger.Handlers()[0])
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(test.GetRandomMessage())
			}
		})
	})
}

func BenchmarkLoggerTextFormatter(b *testing.B) {
	test.WithBenchLogger(b, func(logger *xylog.Logger) {
		var formatter = test.AddFullMacros(xylog.NewTextFormatter())
		logger.Handlers()[0].SetFormatter(formatter)
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(test.GetRandomMessage())
			}
		})
	})
}

func BenchmarkLoggerJSONFormatter(b *testing.B) {
	test.WithBenchLogger(b, func(logger *xylog.Logger) {
		var formatter = test.AddFullMacros(xylog.NewJSONFormatter())
		logger.Handlers()[0].SetFormatter(formatter)
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(test.GetRandomMessage())
			}
		})
	})
}
