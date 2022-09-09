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

	name      string
	emitters  []Emitter
	level     int
	lock      *xylock.RWLock
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
		f:         &filterer{},
		name:      name,
		emitters:  nil,
		level:     NOTSET,
		lock:      &xylock.RWLock{},
		formatter: defaultFormatter,
	}
	if name != "" {
		mapHandler(name, h)
	}

	return h
}

// Name returns the current name. An anonymous Handler returns the empty name.
func (h *Handler) Name() string {
	return h.lock.RLockFunc(func() any { return h.name }).(string)
}

// Level returns the current logging level.
func (h *Handler) Level() int {
	return h.lock.RLockFunc(func() any { return h.level }).(int)
}

// SetLevel sets the new logging level. It is NOTSET by default.
func (h *Handler) SetLevel(level int) {
	h.lock.WLockFunc(func() { h.level = CheckLevel(level) })
}

// Formatter returns the current Formatter.
func (h *Handler) Formatter() Formatter {
	return h.lock.RLockFunc(func() any { return h.formatter }).(Formatter)
}

// SetFormatter sets a new Formatter.
func (h *Handler) SetFormatter(f Formatter) {
	h.lock.WLockFunc(func() { h.formatter = f })
}

// Filters returns all current Filters.
func (h *Handler) Filters() []Filter {
	return h.lock.RLockFunc(func() any { return h.f.Filters() }).([]Filter)
}

// AddFilter adds a specified Filter.
func (h *Handler) AddFilter(f Filter) {
	h.lock.WLockFunc(func() { h.f.AddFilter(f) })
}

// RemoveFilter remove an existed Filter.
func (h *Handler) RemoveFilter(f Filter) {
	h.lock.WLockFunc(func() { h.f.RemoveFilter(f) })
}

// Emitters returns all current Emitters.
func (h *Handler) Emitters() []Emitter {
	return h.lock.RLockFunc(func() any { return h.emitters }).([]Emitter)
}

// AddEmitter adds a specified Emitter.
func (h *Handler) AddEmitter(e Emitter) {
	h.lock.WLockFunc(func() { h.emitters = append(h.emitters, e) })
}

// RemoveEmitter remove an existed Emitter.
func (h *Handler) RemoveEmitter(e Emitter) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for i := range h.emitters {
		if h.emitters[i] == e {
			h.emitters = append(h.emitters[:i], h.emitters[i+1:]...)
		}
	}
}

// handle checks if a record should be logged or not, then call Emitter if it
// is.
func (h *Handler) handle(record LogRecord) {
	if h.filter(record) && record.LevelNo >= h.Level() {
		var msg, err = h.Formatter().Format(record)
		if err != nil {
			msg = fmt.Sprintf(
				"An error occurred while formatting the message (%s)", err)
		}
		var emitters = make([]Emitter, len(h.Emitters()))
		copy(emitters, h.Emitters())
		for i := range emitters {
			emitters[i].Emit(msg)
		}
	}
}

// filter checks all Filters, if there is any failed one, it will returns false.
func (h *Handler) filter(r LogRecord) bool {
	return h.lock.RLockFunc(func() any { return h.f.filter(r) }).(bool)
}

// mapHandler associates a name with a handler.
func mapHandler(name string, h *Handler) {
	xycond.AssertNotIn(name, handlerManager)
	handlerManager[name] = h
}
