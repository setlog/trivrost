package gui

import (
	"testing"
	"time"
)

func TestRateString(t *testing.T){
	result := rateString(0, time.Second)
	if result != "0 B/s" {
		t.Errorf("rateString was incorrect, got: %s; want: %s", result, "0 B/s")
	}

	result = rateString(1024*200, time.Second)
	if result != "200 KiB/s" {
		t.Errorf("rateString was incorrect, got: %s; want: %s", result, "200 KiB/s")
	}

	result = rateString(1024*300, time.Second * 3)
	if result != "100 KiB/s" {
		t.Errorf("rateString was incorrect, got: %s; want: %s", result, "100 KiB/s")
	}

	result = rateString(1024*1024*10, time.Second)
	if result != "10.0 MiB/s" {
		t.Errorf("rateString was incorrect, got: %s; want: %s", result, "10.0 MiB/s")
	}

	result = rateString(1024*1024*10, time.Second / 2)
	if result != "20.0 MiB/s" {
		t.Errorf("rateString was incorrect, got: %s; want: %s", result, "20.0 MiB/s")
	}
}
