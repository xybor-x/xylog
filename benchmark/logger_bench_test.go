package benchmark

import (
	"testing"

	"github.com/xybor-x/xylog"
)

var DevnullEmitter *xylog.StreamEmitter

func BenchmarkLoggerWithoutLog(b *testing.B) {
	withBenchLogger(b, func(logger *xylog.Logger) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(getRandomMessage())
			}
		})
	})
}

func BenchmarkLoggerTextFormatter(b *testing.B) {
	withBenchLogger(b, func(logger *xylog.Logger) {
		var formatter = addFullMacros(xylog.NewTextFormatter())
		logger.Handlers()[0].SetFormatter(formatter)
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(getRandomMessage())
			}
		})
	})
}

func BenchmarkLoggerJSONFormatter(b *testing.B) {
	withBenchLogger(b, func(logger *xylog.Logger) {
		var formatter = addFullMacros(xylog.NewJSONFormatter())
		logger.Handlers()[0].SetFormatter(formatter)
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(getRandomMessage())
			}
		})
	})
}
