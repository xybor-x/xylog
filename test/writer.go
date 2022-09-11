package test

import "github.com/xybor-x/xyerror"

// MockWriter is a LogWriter which all message will be captured.
type MockWriter struct {
	// Captured is the string MockWriter wrote.
	Captured string

	// Error decides if Write method returns error or not.
	Error bool
}

// Write append the byte slice to Captured string. It returns n error if Error
// attribute is true.
func (w *MockWriter) Write(b []byte) (int, error) {
	if w.Error {
		return 0, xyerror.BaseException.New("mockwriter raised an error")
	}

	w.Captured += string(b)
	return len(b), nil
}

// Close is a fake methods.
func (w *MockWriter) Close() error {
	return nil
}

// Reset sets the Captured to empty.
func (w *MockWriter) Reset() {
	w.Captured = ""
}
