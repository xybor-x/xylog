package xylog

import (
	"path/filepath"
	"runtime"
	"strings"
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

	// This is not a macro. Extra provides possibility of using custom macros.
	Extra map[string]any

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
		if attr, ok := r.Extra[name]; ok {
			return attr, nil
		}
		return nil, xyerror.ValueError.Newf("not found attribute %s", name)
	}
}

// makeRecord creates specialized LogRecords.
func makeRecord(
	name string, level int, pathname string, lineno int, pc uintptr,
	extra map[string]any, fields ...field,
) LogRecord {
	var created = time.Now()
	var module, funcname = extractFromPC(pc)

	return LogRecord{
		Asctime:         created.Format(timeLayout),
		Created:         created.Unix(),
		Extra:           extra,
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
func extractFromPC(pc uintptr) (module, fname string) {
	var s = runtime.FuncForPC(pc).Name()

	// Split the funcname in the form of func with receiver.
	// module.(receiver).func for example.
	var parts []string
	var sep = ".("
	parts = strings.Split(s, ".(")

	// If it is not the form of func with receiver, split it with normal func.
	// module.func for example.
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
