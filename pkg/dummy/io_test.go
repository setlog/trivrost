package dummy_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/setlog/trivrost/pkg/dummy"
	"github.com/setlog/trivrost/pkg/misc"
)

func TestReadCloser(t *testing.T) {
	data := []byte("super amazing data")
	rc := &dummy.ReadCloser{Data: data}
	readData, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if !bytes.Equal(data, readData) {
		t.Fatalf("Data not equal. Got: %v", string(readData))
	}
}

func TestReadCloserBigFile(t *testing.T) {
	data := []byte("hyper amazing data" + misc.MustGetRandomHexString(10000))
	rc := &dummy.ReadCloser{Data: data}
	readData, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if !bytes.Equal(data, readData) {
		t.Fatalf("Data not equal. Got: %v", string(readData))
	}
}
