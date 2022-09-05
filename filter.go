package xylog

import (
	"github.com/xybor-x/xylock"
)

// Filter instances are used to perform arbitrary filtering of LogRecord.
type Filter interface {
	Filter(record LogRecord) bool
}

// A base class for loggers and handlers which allows them to share common code.
type filterer struct {
	filters []Filter
	lock    xylock.RWLock
}

func newfilterer() *filterer {
	return &filterer{
		filters: nil,
		lock:    xylock.RWLock{},
	}
}

// AddFilter adds a specified filter.
func (ftr *filterer) AddFilter(f Filter) {
	ftr.lock.WLockFunc(func() {
		ftr.filters = append(ftr.filters, f)
	})
}

// filter checks all filters in filterer, if there is any failed filter, it will
// returns false.
func (ftr *filterer) filter(record LogRecord) bool {
	// Avoid calling locks.
	if len(ftr.filters) == 0 {
		return true
	}

	return ftr.lock.RLockFunc(func() any {
		for i := range ftr.filters {
			if !ftr.filters[i].Filter(record) {
				return false
			}
		}
		return true
	}).(bool)
}
