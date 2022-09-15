package benchmark

import (
	"testing"

	"github.com/xybor-x/xylog"
	"github.com/xybor-x/xylog/encoding"
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
		logger.RemoveAllHandlers()
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(test.GetRandomMessage())
			}
		})
	})
}

func BenchmarkLoggerTextEncoding(b *testing.B) {
	test.WithBenchLogger(b, func(logger *xylog.Logger) {
		test.AddFullMacros(logger.Handlers()[0])
		logger.Handlers()[0].SetEncoding(encoding.NewTextEncoding())
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(test.GetRandomMessage())
			}
		})
	})
}

func BenchmarkLoggerJSONEncoding(b *testing.B) {
	test.WithBenchLogger(b, func(logger *xylog.Logger) {
		test.AddFullMacros(logger.Handlers()[0])
		logger.Handlers()[0].SetEncoding(encoding.NewJSONEncoding())
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(test.GetRandomMessage())
			}
		})
	})
}
