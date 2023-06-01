package misc

import (
	"encoding/json"
	"io"
)

func MustReadAll(fromReader io.Reader) []byte {
	data, err := io.ReadAll(fromReader)
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
