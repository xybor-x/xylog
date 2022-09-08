package benchmark

import (
	"testing"

	"github.com/xybor-x/xylog"
)

func BenchmarkTextFormatter(b *testing.B) {
	var record = xylog.LogRecord{}
	var formatter = xylog.NewTextFormatter(
		"%(asctime)s %(created)d %(filename)s %(funcname)s %(levelname)s " +
			"%(levelno)d %(lineno)d %(message)s %(module)s %(msecs)d " +
			"%(name)s %(pathname)s %(process)d %(relativeCreated)d",
	)
	for i := 0; i < b.N; i++ {
		formatter.Format(record)
	}
}

func BenchmarkStructuredFormatter(b *testing.B) {
	var record = xylog.LogRecord{}
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

	for i := 0; i < b.N; i++ {
		formatter.Format(record)
	}
}

func BenchmarkJSONFormatter(b *testing.B) {
	var record = xylog.LogRecord{}
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

	for i := 0; i < b.N; i++ {
		formatter.Format(record)
	}
}
