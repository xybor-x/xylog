package xylog

// Filter instances are used to perform arbitrary filtering of LogRecord.
type Filter interface {
	Filter(record LogRecord) bool
}

// A base class for loggers and handlers which allows them to share common code.
type filterer struct {
	filters []Filter
}

// AddFilter adds a specified filter.
func (ftr *filterer) AddFilter(f Filter) {
	ftr.filters = append(ftr.filters, f)
}

// Filters returns all current filters.
func (ftr *filterer) Filters() []Filter {
	return ftr.filters
}

// RemoveFilter remove an existed filter.
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
	// Avoid calling locks.
	if len(ftr.filters) == 0 {
		return true
	}

	for i := range ftr.filters {
		if !ftr.filters[i].Filter(record) {
			return false
		}
	}
	return true
}
