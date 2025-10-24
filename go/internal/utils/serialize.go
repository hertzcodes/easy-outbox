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

func Marshal(in any) ([]byte, error) {
	if in == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(in)
	if err != nil {
		return nil, ErrSerialization
	}
	return buf.Bytes(), nil
}

func Unmarshal[T any](data []byte, out *T) error {
    decoder := gob.NewDecoder(bytes.NewReader(data))
    if err := decoder.Decode(out); err != nil {
        return ErrDeserialization
    }
    return nil
}
