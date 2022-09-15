package xylog

import (
	"fmt"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xylock"
	"github.com/xybor-x/xylog/encoding"
)

// Handler handles logging events. Do NOT instantiated directly this struct.
//
// Any Handler with a not-empty name will be associated with its name.
type Handler struct {
	f *filterer

	name     string
	emitters []Emitter
	level    int
	lock     *xylock.RWLock
	encoder  *encoding.Encoder
	macros   []macroField
	fields   []field
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
		f:       &filterer{},
		name:    name,
		level:   NOTSET,
		lock:    &xylock.RWLock{},
		encoder: encoding.NewEncoder(encoding.NewTextEncoding()),
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
	h.lock.WLockFunc(func() { h.level = level })
}

// SetEncoding sets a new Encoding.
func (h *Handler) SetEncoding(e encoding.Encoding) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.encoder = encoding.NewEncoder(e)
	for i := range h.fields {
		h.encoder.Add(h.fields[i].key, h.fields[i].value)
	}
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

// AddMacro adds the macro value to the logging message under a name.
func (h *Handler) AddMacro(name, macro string) {
	h.lock.WLockFunc(func() {
		h.macros = append(h.macros, macroField{key: name, macro: macro})
	})
}

// AddField adds a fixed field to the logging message.
func (h *Handler) AddField(name string, value any) {
	h.lock.WLockFunc(func() {
		h.fields = append(h.fields, makeField(name, value))
		h.encoder.Add(name, value)
	})
}

// Handle checks if a record should be logged or not, then calls Emitters if it
// is.
func (h *Handler) Handle(record LogRecord) {
	if h.filter(record) && record.LevelNo >= h.Level() {
		var msg, err = h.format(record)
		if err != nil {
			msg = []byte(fmt.Sprintf(
				"An error occurred while formatting the message (%s)", err))
		}
		var emitters = h.Emitters()
		for i := range emitters {
			emitters[i].Emit(msg)
		}
	}
}

// format creates the logging message based on the encoding.
func (h Handler) format(record LogRecord) ([]byte, error) {
	var encoder = h.encoder.Clone()
	defer encoder.Free()

	for i := range h.macros {
		var attr, err = record.getValue(h.macros[i].macro)
		if err != nil {
			return nil, err
		}
		encoder.Add(h.macros[i].key, attr)
	}

	for _, f := range record.Fields {
		encoder.Add(f.key, f.value)
	}

	return encoder.Encode(), nil
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
