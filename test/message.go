// MIT License
//
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

package test

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
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

// GetRandomMessage returns a random long message.
func GetRandomMessage() string {
	return messages[rand.Int()%len(messages)]
}

// GetRandomLoggerName returns a random logger name with dot-seperated.
func GetRandomLoggerName() string {
	return loggerNames[rand.Int()%len(loggerNames)]
}
