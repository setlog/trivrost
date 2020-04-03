package gui

import (
	"testing"
	"time"
)

func TestRateString(t *testing.T){
	result := rateString(0, time.Second)
	if result != "0 B/s" {
		t.Errorf("TestTest")
	}
}
