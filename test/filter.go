package test

import "github.com/xybor-x/xylog"

// LoggerNameFilter is a Filter allows to log Log Records whose name is the same
// with the predefined one.
type LoggerNameFilter struct {
	Name string
}

// Filter returns true if record's nname is equal to filter's name.
func (f *LoggerNameFilter) Filter(record xylog.LogRecord) bool {
	return record.Name == f.Name
}
