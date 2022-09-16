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

package xylog

// Filter instances are used to perform arbitrary filtering of LogRecord.
type Filter interface {
	Filter(record LogRecord) bool
}

// A base class for loggers and handlers which allows them to share common code.
type filterer struct {
	filters []Filter
}

// Filters returns all current filters.
func (ftr *filterer) Filters() []Filter {
	return ftr.filters
}

// AddFilter adds a specified filter.
func (ftr *filterer) AddFilter(f Filter) {
	ftr.filters = append(ftr.filters, f)
}

// RemoveFilter removes an existed filter.
func (ftr *filterer) RemoveFilter(f Filter) {
	for i := range ftr.filters {
		if ftr.filters[i] == f {
			ftr.filters = append(ftr.filters[:i], ftr.filters[i+1:]...)
		}
	}
}

// filter checks all filters in filterer, if there is any failed filter, it will
// returns false.
func (ftr *filterer) filter(record LogRecord) bool {
	for i := range ftr.filters {
		if !ftr.filters[i].Filter(record) {
			return false
		}
	}
	return true
}
