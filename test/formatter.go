package test

import "github.com/xybor-x/xylog"

// AddFullMacros adds all macros to a Formatter.
func AddFullMacros(f xylog.Formatter) xylog.Formatter {
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
