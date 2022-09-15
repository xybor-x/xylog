package test

import "github.com/xybor-x/xylog"

// AddFullMacros adds all macros to a Formatter.
func AddFullMacros(h *xylog.Handler) {
	h.AddMacro("asctime", "asctime")
	h.AddMacro("created", "created")
	h.AddMacro("filename", "filename")
	h.AddMacro("funcname", "funcname")
	h.AddMacro("levelname", "levelname")
	h.AddMacro("levelno", "levelno")
	h.AddMacro("lineno", "lineno")
	h.AddMacro("module", "module")
	h.AddMacro("msecs", "msecs")
	h.AddMacro("pathname", "pathname")
	h.AddMacro("process", "process")
	h.AddMacro("relativeCreated", "relativeCreated")
}
