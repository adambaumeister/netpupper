package udpr

import (
	"encoding/binary"
)

/*
Base Field interface
All fields must satisfy these methods
*/
type Field interface {
	Read([]byte)
	Serialize() []byte
}

/*
IntField is a fixed-size (2 byte) field.
*/
type IntField struct {
	Length uint16
	Value  uint16
}

func (f *IntField) Read(b []byte) {
	f.Value = binary.BigEndian.Uint16(b)
}
func (f *IntField) Write(v uint16) {
	f.Value = v
}
func (f *IntField) Serialize() []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, f.Value)
	return b
}
