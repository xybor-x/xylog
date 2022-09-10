package xylog

import (
	"encoding/json"
	"fmt"

	"github.com/xybor-x/xyerror"
	"github.com/xybor-x/xylog/encoding"
)

// Formatter instances are used to convert a LogRecord to text.
//
// Formatter need to know how a LogRecord is constructed. They are responsible
// for converting a LogRecord to a string which can be interpreted by either a
// human or an external system.
type Formatter interface {
	Format(LogRecord) (string, error)
	AddMacro(string, string) Formatter
	AddField(string, any) Formatter
}

// TextFormatter formats the logging message with the form of key=value.
type TextFormatter struct {
	macros []macroField
	fixed  *encoding.Buffer
}

// NewTextFormatter creates an empty TextFormatter.
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		macros: make([]macroField, 0, 10),
		fixed:  encoding.NewBuffer(),
	}
}

// AddMacro adds the macro to logging message under a name. It returns itself.
func (tf *TextFormatter) AddMacro(name, macro string) Formatter {
	tf.macros = append(tf.macros, macroField{key: name, macro: macro})
	return tf
}

// AddField adds a fixed field to logging message. It returns itself.
func (tf *TextFormatter) AddField(name string, value any) Formatter {
	tf.fixed.AppendSeperator()
	tf.fixed.AppendString(name)
	tf.fixed.AppendByte('=')
	tf.fixed.AppendQuotedString(fmt.Sprint(value))
	return tf
}

// Format creates the logging message with the form of key=value.
func (tf TextFormatter) Format(record LogRecord) (string, error) {
	var buf = tf.fixed.Copy()

	for _, m := range tf.macros {
		var attr, err = record.getValue(m.macro)
		if err != nil {
			return "", err
		}

		buf.AppendSeperator()
		buf.AppendString(m.key)
		buf.AppendByte('=')
		buf.AppendQuotedString(fmt.Sprint(attr))
	}

	for _, f := range record.Fields {
		buf.AppendSeperator()
		buf.AppendString(f.key)
		buf.AppendByte('=')
		buf.AppendQuotedString(fmt.Sprint(f.value))
	}

	return buf.String(), nil
}

// JSONFormatter allows logging message to be parsed as json format.
type JSONFormatter struct {
	macros []macroField
	fields map[string]any
}

// NewJSONFormatter returns an empty JSONFormatter.
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		macros: make([]macroField, 0, 5),
		fields: make(map[string]any),
	}
}

// AddMacro adds the macro value to the logging message under a name. It returns
// itself.
func (js *JSONFormatter) AddMacro(name, macro string) Formatter {
	js.macros = append(js.macros, macroField{key: name, macro: macro})
	return js
}

// AddField adds a fixed field to the logging message. It returns itself.
func (js *JSONFormatter) AddField(name string, value any) Formatter {
	js.fields[name] = value
	return js
}

// Format creates the logging message of JSON format.
func (js JSONFormatter) Format(record LogRecord) (string, error) {
	// Copy the predefined fields to the new map.
	var data = make(map[string]any)
	for k, v := range js.fields {
		data[k] = v
	}

	for _, m := range js.macros {
		var attr, err = record.getValue(m.macro)
		if err != nil {
			return "", err
		}
		data[m.key] = attr
	}

	for _, f := range record.Fields {
		data[f.key] = f.value
	}

	var s, err = json.Marshal(data)
	if err != nil {
		return "", xyerror.ValueError.New(err)
	}
	return string(s), nil
}
