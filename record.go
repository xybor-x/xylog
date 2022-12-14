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

package xylog

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/xybor-x/xyerror"
)

// A LogRecord instance represents an event being logged.
//
// LogRecord instances are created every time something is logged. They contain
// all the information pertinent to the event being logged. The main information
// passed in is Message. The record also includes information as when the record
// was created or the source line where the logging call was made.
type LogRecord struct {
	// Textual time when the LogRecord was created.
	Asctime string

	// Time when the LogRecord was created (time.Now().Unix() return value).
	Created int64

	// This a not a macro. Fields are always added to the logging message
	// without calling AddMacro.
	Fields []field

	// Filename is the portion of pathname.
	FileName string

	// Funcname is the name of function which logged the record.
	FuncName string

	// Text logging level for the message ("DEBUG", "INFO", "WARNING", "ERROR",
	// "CRITICAL").
	LevelName string

	// Numeric logging level for the message (DEBUG, INFO, WARNING, ERROR,
	// CRITICAL).
	LevelNo int

	// Source line number where the logging call was issued.
	LineNo int

	// The module called log method.
	Module string

	// Millisecond portion of the creation time.
	Msecs int

	// Name of the logger.
	Name string

	// Full pathname of the source file where the logging call was issued.
	PathName string

	// Process ID.
	Process int

	// Time in milliseconds when the LogRecord was created, relative to the time
	// the logging module was loaded (typically at application startup time).
	RelativeCreated int64
}

func (r LogRecord) getValue(name string) (any, error) {
	switch name {
	case "asctime":
		return r.Asctime, nil
	case "created":
		return r.Created, nil
	case "filename":
		return r.FileName, nil
	case "funcname":
		return r.FuncName, nil
	case "levelname":
		return r.LevelName, nil
	case "levelno":
		return r.LevelNo, nil
	case "lineno":
		return r.LineNo, nil
	case "module":
		return r.Module, nil
	case "msecs":
		return r.Msecs, nil
	case "name":
		return r.Name, nil
	case "pathname":
		return r.PathName, nil
	case "process":
		return r.Process, nil
	case "relativeCreated":
		return r.RelativeCreated, nil
	default:
		return nil, xyerror.ValueError.Newf("not found attribute %s", name)
	}
}

// makeRecord creates specialized LogRecords.
func makeRecord(name string, level int, fields ...field) LogRecord {
	var created = time.Now()
	var pc uintptr
	var lineno int
	var module, pathname, funcname = "unknown", "unknown", "unknown"
	if findCaller {
		var ok bool
		pc, pathname, lineno, ok = runtime.Caller(skipCall)
		if !ok {
			pathname = "unknown"
			lineno = -1
		} else {
			module, funcname = extractFromPC(pc)
		}
	}

	return LogRecord{
		Asctime:         created.Format(timeLayout),
		Created:         created.Unix(),
		Fields:          fields,
		FileName:        filepath.Base(pathname),
		FuncName:        funcname,
		LevelName:       GetLevelName(level),
		LevelNo:         level,
		LineNo:          lineno,
		Module:          module,
		Msecs:           created.Nanosecond() / int(time.Millisecond),
		Name:            name,
		PathName:        pathname,
		Process:         processid,
		RelativeCreated: created.UnixMilli() - startTime,
	}
}

// extractFromPC returns module name and function name from program counter.
func extractFromPC(pc uintptr) (string, string) {
	var s = runtime.FuncForPC(pc).Name()

	var moduleIdx = -1
	for i := range s {
		if s[i] == '.' && moduleIdx == -1 {
			moduleIdx = i
		}
		if s[i] == '/' {
			moduleIdx = -1
		}
	}

	if moduleIdx == -1 {
		return "unknown", s
	}

	return s[0:moduleIdx], s[moduleIdx+1:]

}
