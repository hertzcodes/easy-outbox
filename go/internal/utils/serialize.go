package utils

import (
	"bytes"
	"encoding/gob"
	"errors"
)

var (
	ErrSerialization   = errors.New("failed to serialize object to bytes")
	ErrDeserialization = errors.New("failed to deserialize bytes to object")
)

func MarshalGob(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(v)
	if err != nil {
		return nil, ErrSerialization
	}
	return buf.Bytes(), nil
}

func UnmarshalGob(data []byte, v interface{}) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if decoder.Decode(v) != nil {
		return ErrDeserialization
	}
	return nil
}
