package xylog

import (
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/xybor-x/xycond"
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

	// Filename portion of pathname.
	FileName string

	// Function name logged the record.
	FuncName string

	// Text logging level for the message ("DEBUG", "INFO", "WARNING", "ERROR",
	// "CRITICAL").
	LevelName string

	// Numeric logging level for the message (DEBUG, INFO, WARNING, ERROR,
	// CRITICAL).
	LevelNo int

	// Source line number where the logging call was issued.
	LineNo int

	// The logging message.
	Message string

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

func (r LogRecord) mapIndex(i int) any {
	switch i {
	case 0:
		return r.Asctime
	case 1:
		return r.Created
	case 2:
		return r.FileName
	case 3:
		return r.FuncName
	case 4:
		return r.LevelName
	case 5:
		return r.LevelNo
	case 6:
		return r.LineNo
	case 7:
		return r.Message
	case 8:
		return r.Module
	case 9:
		return r.Msecs
	case 10:
		return r.Name
	case 11:
		return r.PathName
	case 12:
		return r.Process
	case 13:
		return r.RelativeCreated
	default:
		xycond.Panic("unknown index %d", i)
		return nil
	}
}

func (r LogRecord) mapName(name string) int {
	switch name {
	case "asctime":
		return 0
	case "created":
		return 1
	case "filename":
		return 2
	case "funcname":
		return 3
	case "levelname":
		return 4
	case "levelno":
		return 5
	case "lineno":
		return 6
	case "message":
		return 7
	case "module":
		return 8
	case "msecs":
		return 9
	case "name":
		return 10
	case "pathname":
		return 11
	case "process":
		return 12
	case "relativeCreated":
		return 13
	default:
		xycond.Panic("unknown name %s", name)
		return -1
	}
}

// makeRecord creates specialized LogRecords.
func makeRecord(
	name string, level int, pathname string, lineno int, msg string, pc uintptr,
) LogRecord {
	var created = time.Now()
	var module, funcname = extractFromPC(pc)

	return LogRecord{
		Asctime:         created.Format(timeLayout),
		Created:         created.Unix(),
		FileName:        filepath.Base(pathname),
		FuncName:        funcname,
		LevelName:       getLevelName(level),
		LevelNo:         level,
		LineNo:          lineno,
		Message:         msg,
		Module:          module,
		Msecs:           created.Nanosecond() / int(time.Millisecond),
		Name:            name,
		PathName:        pathname,
		Process:         processid,
		RelativeCreated: created.UnixMilli() - startTime,
	}
}

// extractFromPC returns module name and function name from program counter.
func extractFromPC(pc uintptr) (module, fname string) {
	var s = runtime.FuncForPC(pc).Name()

	// Split the funcname in the form of func with receiver.
	// E.g. module.(receiver).func
	var parts []string
	var sep = ".("
	parts = strings.Split(s, ".(")

	// If it is not the form of func with receiver, split it with normal func.
	// E.g. module.func
	if len(parts) <= 1 {
		sep = "."
		parts = strings.Split(s, ".")
	}

	// In case one of form is valid, remove the funcname from string.
	if len(parts) > 1 {
		var funcname = parts[len(parts)-1]
		module = strings.TrimSuffix(s, sep+funcname)
		fname = strings.TrimPrefix(sep+funcname, ".")
		return
	}

	// Otherwise, the string contains only funcname.
	return "unknown", s
}
