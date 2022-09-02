package xylog

import (
	"encoding/json"
	"fmt"

	"github.com/xybor-x/xycond"
	"github.com/xybor-x/xyerror"
)

// Formatter instances are used to convert a LogRecord to text.
//
// Formatter need to know how a LogRecord is constructed. They are responsible
// for converting a LogRecord to a string which can be interpreted by either a
// human or an external system.
type Formatter interface {
	Format(LogRecord) (string, error)
}

// The TextFormatter can be initialized with a format string which makes use of
// knowledge of the LogRecord attributes - e.g. %(message)s or %(levelno)d. See
// LogRecord for more details.
type TextFormatter struct {
	fmtstr string
	names  []string
}

// NewTextFormatter creates a textFormatter which uses LogRecord attributes to
// contribute logging string, e.g. %(message)s or %(levelno)d. See LogRecord for
// more details.
func NewTextFormatter(s string) TextFormatter {
	var tf = TextFormatter{}
	var i, n = 0, len(s)
	for i < n {
		tf.fmtstr += string(s[i])
		if s[i] == '%' {
			xycond.AssertLessThan(i+1, n)
			i++
			switch s[i] {
			case '%':
			case '(':
				i++
				var token = ""
				for {
					xycond.AssertLessThan(i, n)
					if s[i] == ')' {
						break
					}
					token += string(s[i])
					i++
				}
				tf.names = append(tf.names, token)
			default:
				xycond.Panicf("unexpected token: %s", s[i-2:i])
			}
		}
		i++
	}

	return tf
}

// Format creates a logging string by combining format string and logging
// record.
func (tf TextFormatter) Format(record LogRecord) (string, error) {
	var err error
	var attrs = make([]any, len(tf.names))
	for i := range tf.names {
		attrs[i], err = record.getAttributeByName(tf.names[i])
		if err != nil {
			return "", err
		}

		switch attrs[i].(type) {
		case map[string]any:
			var s []byte
			s, err = json.Marshal(attrs[i])
			if err != nil {
				return "", err
			}
			attrs[i] = string(s)
		}
	}

	return fmt.Sprintf(tf.fmtstr, attrs...), nil
}

// JSONFormatter allows logging message to be parsed as json format. It uses the
// same macros as other Formatter.
type JSONFormatter struct {
	fields []field
}

// NewJSONFormatter returns an empty JSONFormatter.
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{fields: make([]field, 0, 5)}
}

// AddField adds the macro to logging message under key name. It returns itself.
func (js *JSONFormatter) AddField(name, macro string) *JSONFormatter {
	js.fields = append(js.fields, field{key: name, value: macro})
	return js
}

// Format creates the logging message as the json format.
func (js JSONFormatter) Format(record LogRecord) (string, error) {
	var err error
	var attr = make(map[string]any)
	for _, f := range js.fields {
		attr[f.key], err = record.getAttributeByName(f.value.(string))
		if err != nil {
			return "", err
		}
	}

	var s []byte
	s, err = json.Marshal(attr)
	if err != nil {
		return "", xyerror.ValueError.New(err)
	}
	return string(s), nil
}
