package fetching

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewLowercaseHeaders(t *testing.T) {
	tests := []struct {
		headers         http.Header
		expectedHeaders LowercaseHeaders
	}{
		{nil, LowercaseHeaders{}},
		{http.Header{}, LowercaseHeaders{}},
		{http.Header{"content-type": {"foo"}}, LowercaseHeaders{"content-type": {"foo"}}},
		{http.Header{"Content-Type": {"foo"}}, LowercaseHeaders{"content-type": {"foo"}}},
		{http.Header{"Content-Type": {"foo"}, "Content-Length": {"3"}}, LowercaseHeaders{"content-type": {"foo"}, "content-length": {"3"}}},
	}
	for i, test := range tests {
		lowerCaseHeaders := NewLowercaseHeaders(test.headers)
		if !reflect.DeepEqual(lowerCaseHeaders, test.expectedHeaders) {
			t.Errorf("Test #%d: lowerCaseHeaders = %v. Expected %v.", i+1, lowerCaseHeaders, test.expectedHeaders)
		}
	}
}

func TestLowercaseHeadersGet(t *testing.T) {
	headers := NewLowercaseHeaders(http.Header{
		"Content-Type":   {"foo"},
		"Content-Length": {"3"},
	})
	if headers.Get("Content-Type") != "foo" {
		t.Errorf(`headers.Get("Content-Type") != "foo"`)
	}
	if headers.Get("Content-Length") != "3" {
		t.Errorf(`headers.Get("Content-Length") != "3"`)
	}
	if headers.Get("content-type") != "foo" {
		t.Errorf(`headers.Get("content-type") != "foo"`)
	}
	if headers.Get("content-length") != "3" {
		t.Errorf(`headers.Get("content-length") != "3"`)
	}
}

func TestLowercaseValues(t *testing.T) {
	headers := NewLowercaseHeaders(http.Header{
		"Content-Type":   {"foo"},
		"Content-Length": {"3"},
		"content-type":   {"bar"},
		"content-length": {"4"},
	})
	tests := []struct {
		headerName     string
		expectedValues map[string]int
	}{
		{"cOnTeNt-TyPe", map[string]int{"foo": 1, "bar": 1}},
		{"cOnTeNt-LeNgTh", map[string]int{"3": 1, "4": 1}},
	}
	for _, test := range tests {
		for _, value := range headers.Values(test.headerName) {
			_, ok := test.expectedValues[value]
			if !ok {
				t.Errorf("unexpected or duplicate value '%s' for header name '%s'", value, test.headerName)
			} else {
				delete(test.expectedValues, value)
			}
		}
		if len(test.expectedValues) > 0 {
			t.Errorf("missing values for header name '%s': %v", test.headerName, test.expectedValues)
		}
	}
}
