package server

const PACKETTYPE_LENGTH = 2

/*
Packet base is the common attributes that a slice of a packet has
*/
type PacketBase struct {
	Fields []*Field
	Length uint16
}

func (pb *PacketBase) AddField(f Field) {
	pb.Fields = append(pb.Fields, &f)
}

/*
Header is the packet header - it defines what type of packet it is and how long it is.
*/
type Header struct {
	Base         PacketBase
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
		Length: PACKETTYPE_LENGTH,
	}
	h.PacketType.Read(b[:2])
	return h
}
