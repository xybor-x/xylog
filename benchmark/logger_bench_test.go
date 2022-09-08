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
				logger.Debug(getRandomLoggerName())
			}
		})
	})
}

func BenchmarkLoggerTextFormatter(b *testing.B) {
	withBenchLogger(b, func(logger *xylog.Logger) {
		var formatter = xylog.NewTextFormatter(
			"%(asctime)s %(created)d %(filename)s %(funcname)s %(levelname)s " +
				"%(levelno)d %(lineno)d %(message)s %(module)s %(msecs)d " +
				"%(name)s %(pathname)s %(process)d %(relativeCreated)d")
		logger.Handlers()[0].SetFormatter(formatter)
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(getRandomLoggerName())
			}
		})
	})
}

func BenchmarkLoggerStructuredFormatter(b *testing.B) {
	withBenchLogger(b, func(logger *xylog.Logger) {
		var formatter = xylog.NewStructuredFormatter().
			AddField("asctime", "asctime").
			AddField("created", "created").
			AddField("filename", "filename").
			AddField("funcname", "funcname").
			AddField("levelname", "levelname").
			AddField("levelno", "levelno").
			AddField("lineno", "lineno").
			AddField("message", "message").
			AddField("module", "module").
			AddField("msecs", "msecs").
			AddField("pathname", "pathname").
			AddField("process", "process").
			AddField("relativeCreated", "relativeCreated")
		logger.Handlers()[0].SetFormatter(formatter)
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(getRandomLoggerName())
			}
		})
	})
}

func BenchmarkLoggerJSONFormatter(b *testing.B) {
	withBenchLogger(b, func(logger *xylog.Logger) {
		var formatter = xylog.NewJSONFormatter().
			AddField("asctime", "asctime").
			AddField("created", "created").
			AddField("filename", "filename").
			AddField("funcname", "funcname").
			AddField("levelname", "levelname").
			AddField("levelno", "levelno").
			AddField("lineno", "lineno").
			AddField("message", "message").
			AddField("module", "module").
			AddField("msecs", "msecs").
			AddField("pathname", "pathname").
			AddField("process", "process").
			AddField("relativeCreated", "relativeCreated")
		logger.Handlers()[0].SetFormatter(formatter)
		logger.SetLevel(xylog.DEBUG)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				logger.Debug(getRandomLoggerName())
			}
		})
	})
}
