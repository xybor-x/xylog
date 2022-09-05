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

// NewTextFormatter creates a textFormatter which uses LogRecord macros and
// format string to contribute logging string, e.g. %(message)s or %(levelno)d.
// See LogRecord for more details.
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

// JSONFormatter allows logging message to be parsed as json format.
type JSONFormatter struct {
	fields []field
}

// NewJSONFormatter returns an empty JSONFormatter.
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{fields: make([]field, 0, 5)}
}

// AddField adds the macro value to the logging message under a name. It returns
// itself. For {message} macro and JSON event logger, you could leave the name
// as empty if you want to add all fields of message into the outer object
// directly.
func (js *JSONFormatter) AddField(name, macro string) *JSONFormatter {
	js.fields = append(js.fields, field{key: name, value: macro})
	return js
}

// Format creates the logging message of JSON format.
func (js JSONFormatter) Format(record LogRecord) (string, error) {
	var err error
	var data = make(map[string]any)
	for _, f := range js.fields {
		var attr, err = record.getAttributeByName(f.value.(string))
		if err != nil {
			return "", err
		}

		if mattr, ok := attr.(map[string]any); ok && f.key == "" {
			for k, v := range mattr {
				data[k] = v
			}
		} else {
			data[f.key] = attr
		}
	}

	var s []byte
	s, err = json.Marshal(data)
	if err != nil {
		return "", xyerror.ValueError.New(err)
	}
	return string(s), nil
}

// StructuredFormatter formats the logging message with the form of key=value.
type StructuredFormatter struct {
	fields []field
}

// NewStructuredFormatter creates an empty StructureFormatter.
func NewStructuredFormatter() *StructuredFormatter {
	return &StructuredFormatter{}
}

// AddField adds the macro to logging message under a name. It returns itself.
// If you leave the name as empty, it will adds the macro value without the name
// and equal character.
func (sf *StructuredFormatter) AddField(
	name, macro string,
) *StructuredFormatter {
	sf.fields = append(sf.fields, field{key: name, value: macro})
	return sf
}

// Format creates the logging message with the form of key=value.
func (sf StructuredFormatter) Format(record LogRecord) (string, error) {
	var msg string
	for _, f := range sf.fields {
		var attr, err = record.getAttributeByName(f.value.(string))
		if err != nil {
			return "", err
		}

		var value string
		switch attr.(type) {
		case map[string]any:
			var s []byte
			s, err = json.Marshal(attr)
			if err != nil {
				return "", err
			}
			value = string(s)
		default:
			value = fmt.Sprint(attr)
		}

		if f.key == "" {
			msg = prefixMessage(msg, value)
		} else {
			msg = prefixMessage(msg, fmt.Sprintf("%s=%s", f.key, value))
		}
	}

	return msg, nil
}
