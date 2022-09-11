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
