package server

import (
	"encoding/binary"
	"fmt"
)

/*
IntField is a fixed-size (2 byte) field.
*/
type IntField struct {
	Length uint16
	Value  uint16
}

func (f *IntField) Read(b []byte) {
	fmt.Printf("ptl: %v\n", b)
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

/*
Header is the packet header - it defines what type of packet it is and how long it is.
*/
type Header struct {
	PacketType   *IntField
	PacketLength *IntField
}

func (h *Header) Serialize() []byte {
	b := []byte{}

	b = append(b, h.PacketType.Serialize()...)
	b = append(b, h.PacketLength.Serialize()...)

	return b
}
func ReadHeader(b []byte) Header {
	h := Header{}
	h.PacketType = &IntField{
		Length: 2,
	}
	h.PacketType.Read(b[:2])
	return h
}
