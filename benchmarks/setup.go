// Copyright (c) 2022 xybor-x
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package benchmarks

import (
	"errors"
	"time"

	"github.com/xybor-x/xylog"
)

type user struct {
	Name      string
	Email     string
	CreatedAt time.Time
}

var (
	errExample = errors.New("fail")

	_tenInts    = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	_tenStrings = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	_tenTimes   = []time.Time{
		time.Unix(0, 0),
		time.Unix(1, 0),
		time.Unix(2, 0),
		time.Unix(3, 0),
		time.Unix(4, 0),
		time.Unix(5, 0),
		time.Unix(6, 0),
		time.Unix(7, 0),
		time.Unix(8, 0),
		time.Unix(9, 0),
	}
	_oneUser = &user{
		Name:      "Jane Doe",
		Email:     "jane@test.com",
		CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	_tenUsers = []*user{
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
	}
)

// add10FieldsToHandler adds 10 fields like zap benchmark.
func add10FieldsToHandler(handler *xylog.Handler) {
	handler.AddField("int", _tenInts[0])
	handler.AddField("ints", _tenInts)
	handler.AddField("string", _tenStrings[0])
	handler.AddField("strings", _tenStrings)
	handler.AddField("time", _tenTimes[0])
	handler.AddField("times", _tenTimes)
	handler.AddField("user1", _oneUser)
	handler.AddField("user2", _oneUser)
	handler.AddField("users", _tenUsers)
	handler.AddField("error", errExample)
}

// add10EventLoggerField create an EventLogger with 10 fields like zap benchmark.
func add10EventLoggerField(logger *xylog.Logger) *xylog.EventLogger {
	var elogger = logger.Event("event")
	elogger.Field("int", _tenInts[0])
	elogger.Field("ints", _tenInts)
	elogger.Field("string", _tenStrings[0])
	elogger.Field("strings", _tenStrings)
	elogger.Field("time", _tenTimes[0])
	elogger.Field("times", _tenTimes)
	elogger.Field("user1", _oneUser)
	elogger.Field("user2", _oneUser)
	elogger.Field("users", _tenUsers)
	elogger.Field("error", errExample)
	return elogger
}
