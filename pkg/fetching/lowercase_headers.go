package fetching

import (
	"net/http"
	"strings"
)

// LowercaseHeaders models http.Header with each header name being lower-case.
type LowercaseHeaders map[string][]string

// NewLowercaseHeaders returns a map of headers for the given http.Header with all
// header names in lower-case. Header names which overlap after becoming lower-case
// have their values merged into one header in an unspecified order.
func NewLowercaseHeaders(headers http.Header) LowercaseHeaders {
	newHeaders := make(LowercaseHeaders)
	for oldKey, oldValues := range headers {
		newKey := strings.ToLower(oldKey)
		if newValues, ok := newHeaders[newKey]; ok {
			newHeaders[newKey] = append(newValues, oldValues...)
		} else {
			newValues := make([]string, 0, 0)
			newHeaders[newKey] = append(newValues, oldValues...)
		}
	}
	return newHeaders
}

// Get is analogous to http.Header.Get()
func (h LowercaseHeaders) Get(header string) string {
	header = strings.ToLower(header)
	if values, ok := h[header]; ok {
		if len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

// Values is analogous to http.Header.Values()
func (h LowercaseHeaders) Values(header string) []string {
	header = strings.ToLower(header)
	if values, ok := h[header]; ok {
		return values
	}
	return nil
}
