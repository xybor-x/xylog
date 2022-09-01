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
	filters map[Filter]any
	lock    xylock.RWLock
}

func newfilterer() *filterer {
	return &filterer{
		filters: make(map[Filter]any),
		lock:    xylock.RWLock{},
	}
}

// AddFilter adds a specified filter.
func (ftr *filterer) AddFilter(f Filter) {
	ftr.lock.WLockFunc(func() {
		if _, ok := ftr.filters[f]; !ok {
			ftr.filters[f] = nil
		}
	})
}

// RemoveFilter removes an existed filter.
func (ftr *filterer) RemoveFilter(f Filter) {
	ftr.lock.WLockFunc(func() {
		delete(ftr.filters, f)
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
		for f := range ftr.filters {
			if !f.Filter(record) {
				return false
			}
		}
		return true
	}).(bool)
}
