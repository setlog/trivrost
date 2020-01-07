package main

import "fmt"

type report struct {
	message string
	isError bool
}

func errorReport(info string, args ...interface{}) *report {
	return &report{message: fmt.Sprintf(info, args...), isError: true}
}

func statusReport(info string, args ...interface{}) *report {
	return &report{message: fmt.Sprintf(info, args...), isError: false}
}

type reports []*report

func (r reports) HaveError() bool {
	for _, report := range r {
		if report != nil && report.isError {
			return true
		}
	}
	return false
}
