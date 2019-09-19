package misc

import "reflect"

// IsNil returns true if v holds a nil reference, even if it is a typed nil interface value.
func IsNil(v interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			_, isValueError := r.(*reflect.ValueError)
			if !isValueError {
				panic("IsNil() panicked with unexpected panic value.")
			}
		}
	}()
	return v == nil || reflect.ValueOf(v).IsNil()
}
