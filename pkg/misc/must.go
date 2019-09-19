package misc

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func MustReadAll(fromReader io.Reader) []byte {
	data, err := ioutil.ReadAll(fromReader)
	if err != nil {
		panic(err)
	}
	return data
}

func MustUnmarshalJSON(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}
