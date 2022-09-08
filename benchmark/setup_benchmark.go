// Package benchmark is used for benchmarking xylog.
package benchmark

import (
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/xybor-x/xylog"
)

func init() {
	rand.Seed(time.Now().UnixNano())

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

var loggerNames []string

func getRandomLoggerName() string {
	return loggerNames[rand.Int()%len(loggerNames)]
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
