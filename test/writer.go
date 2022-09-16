// MIT License
//
// Copyright (c) 2022 xybor-x
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
