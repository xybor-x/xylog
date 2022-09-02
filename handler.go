package xylog

import (
	"fmt"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylock"
)

// Handler handles logging events. Do NOT instantiated directly this struct.
//
// Any Handler with a not-empty name will be associated with its name.
type Handler struct {
	f *filterer
	e Emitter

	level     int
	lock      xylock.RWLock
	formatter Formatter
}

// NewHandler creates a Handler with a specified Emitter.
//
// Any Handler with a non-empty name will be associated with its name. Calling
// NewHandler twice with the same name will cause a panic. If you want to create
// an anonymous Handler, call this function with an empty name.
func NewHandler(name string, e Emitter) *Handler {
	var handler = GetHandler(name)
	xycond.AssertNil(handler)

	handler = &Handler{
		f:         newfilterer(),
		e:         e,
		level:     NOTSET,
		lock:      xylock.RWLock{},
		formatter: defaultFormatter,
	}

	if name != "" {
		mapHandler(name, handler)
	}

	return handler
}

// SetLevel sets the new logging level of handler. It is NOTSET by default.
func (h *Handler) SetLevel(level int) {
	h.lock.WLockFunc(func() { h.level = checkLevel(level) })
}

// SetFormatter sets the new formatter of handler.
func (h *Handler) SetFormatter(f Formatter) {
	h.lock.WLockFunc(func() { h.formatter = f })
}

// AddFilter adds a specified filter.
func (h *Handler) AddFilter(f Filter) {
	h.f.AddFilter(f)
}

// RemoveFilter removes an existed filter.
func (h *Handler) RemoveFilter(f Filter) {
	h.f.RemoveFilter(f)
}

// filter checks all filters in filterer, if there is any failed filter, it will
// returns false.
func (h *Handler) filter(r LogRecord) bool {
	return h.f.filter(r)
}

// handle handles a new record, it will check if the record should be logged or
// not, then call emit if it is.
func (h *Handler) handle(record LogRecord) {
	var level = h.lock.RLockFunc(func() any { return h.level }).(int)
	if h.filter(record) && record.LevelNo >= level {
		var msg, err = h.formatter.Format(record)
		if err != nil {
			msg = fmt.Sprint("An error occurred while formatting the message:",
				err)
		}
		h.lock.WLockFunc(func() { h.e.Emit(msg) })
	}
}
