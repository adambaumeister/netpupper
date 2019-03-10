package server

import "fmt"

const OPEN_TYPE = 1

const PACKETTYPE_LENGTH = 2
const PACKETLEN_LENGTH = 2

/*
Packet base is the common attributes that a slice of a packet has
*/
type PacketBase struct {
	Fields []Field
	Length uint16
}

func (pb *PacketBase) AddField(f Field) {
	pb.Fields = append(pb.Fields, f)
}

func (pb *PacketBase) GetFields() []Field {
	return pb.Fields
}

/*
Header is the packet header - it defines what type of packet it is and how long it is.
*/
type Header struct {
	Fields       []Field
	PacketType   *IntField
	PacketLength *IntField
}

func (h *Header) AddField(f Field) {
	h.Fields = append(h.Fields, f)
}

func (h *Header) Serialize() []byte {
	b := []byte{}
	for _, f := range h.Fields {
		b = append(b, f.Serialize()...)
	}

	return b
}
func ReadHeader(b []byte) Header {
	h := Header{}
	h.PacketType = &IntField{
		Length: PACKETTYPE_LENGTH,
	}
	h.PacketLength = &IntField{
		Length: PACKETLEN_LENGTH,
	}
	fmt.Printf("DEEBUG: %v\n", b)
	h.PacketType.Read(b[:2])
	h.PacketLength.Read(b[2:4])
	return h
}
