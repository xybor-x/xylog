package xylog

import (
	"fmt"

	"github.com/xybor-x/xylock"
)

// Handler handles logging events. Do NOT instantiated directly this struct.
//
// Any Handler with a not-empty name will be associated with its name.
type Handler struct {
	f *filterer

	emitters  []Emitter
	level     int
	lock      xylock.RWLock
	formatter Formatter
}

// GetHandler gets a handler with the specified name, creating it if it doesn't
// yet exist.
//
// Leave the name as empty if you want to create an anonymous Handler.
func GetHandler(name string) *Handler {
	var h, ok = handlerManager[name]
	if ok {
		return h
	}

	h = &Handler{
		f:         newfilterer(),
		emitters:  nil,
		level:     NOTSET,
		lock:      xylock.RWLock{},
		formatter: defaultFormatter,
	}
	if name != "" {
		mapHandler(name, h)
	}

	return h
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

// AddEmitter adds a specified emitter.
func (h *Handler) AddEmitter(e Emitter) {
	h.emitters = append(h.emitters, e)
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
		h.lock.WLockFunc(func() {
			for i := range h.emitters {
				h.emitters[i].Emit(msg)
			}
		})
	}
}
