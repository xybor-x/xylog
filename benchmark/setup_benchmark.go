// Package benchmark is used for benchmarking xylog.
package benchmark

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/xybor-x/xylog"
)

func init() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 1000; i++ {
		messages = append(messages,
			fmt.Sprintf("This is a long enough message %d", i))
	}

	var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 100000; i++ {
		var npart = rand.Intn(5)
		var name = make([]string, npart)
		for j := 0; j < npart; j++ {
			var p = rand.Intn(len(letters))
			name[j] = letters[p : p+1]
		}
		loggerNames = append(loggerNames, strings.Join(name, "."))
	}
}

var messages []string
var loggerNames []string

func getRandomLoggerName() string {
	return loggerNames[rand.Int()%len(loggerNames)]
}

func getRandomMessage() string {
	return messages[rand.Int()%len(messages)]
}

func withBenchLogger(b *testing.B, f func(logger *xylog.Logger)) {
	var devnull, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		b.Fail()
	}
	var emitter = xylog.NewStreamEmitter(devnull)
	var handler = xylog.GetHandler("")
	handler.AddEmitter(emitter)

	var logger = xylog.GetLogger(b.Name())
	logger.AddHandler(handler)

	f(logger)
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
