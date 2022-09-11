package test

import (
	"github.com/xybor-x/xylog"
)

// FullRecord is the record with all filled fields.
var FullRecord = xylog.LogRecord{
	Asctime:         "ASCTIME",
	Created:         1,
	FileName:        "FILENAME",
	FuncName:        "FUNCNAME",
	LevelName:       "LEVELNAME",
	LevelNo:         2,
	LineNo:          3,
	Module:          "MODULE",
	Msecs:           4,
	Name:            "NAME",
	PathName:        "PATHNAME",
	Process:         5,
	RelativeCreated: 6,
}
