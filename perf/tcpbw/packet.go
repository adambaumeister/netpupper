package tcpbw

import (
	"encoding/binary"
	"github.com/adamb/netpupper/errors"
	"net"
)

const OPEN_TYPE = 1
const CONFIRM_TYPE = 2

const PACKETTYPE_LENGTH = 2
const PACKETLEN_LENGTH = 2

const OPEN_LENGTH = 10

type Packet struct {
	Header  *Header
	Message Message
}

func (p *Packet) Serialize() []byte {
	b := []byte{}
	b = append(b, p.Header.Serialize()...)
	b = append(b, p.Message.Serialize()...)
	return b
}

/*
Message methods
*/
type Message interface {
	Serialize() []byte
}

/*
Open message, used to initiate a transfer
	Datalength: Length of data in this transfer
	Reverse: Request the transfer to operate in reverse (server to client)
*/
type Open struct {
	DataLength uint64
	Reverse    uint16
}

func (m *Open) Serialize() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, m.DataLength)
	binary.BigEndian.PutUint16(b, m.Reverse)
	return b
}
func (m *Open) Write(i uint64, r uint16) {
	m.DataLength = i
	m.Reverse = r
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
func ReadHeader(c net.Conn) Header {
	h := Header{}
	var b = make([]byte, 4)
	_, err := c.Read(b)
	errors.CheckError(err)
	h.PacketType = &IntField{
		Length: PACKETTYPE_LENGTH,
	}
	h.PacketLength = &IntField{
		Length: PACKETLEN_LENGTH,
	}
	h.PacketLength.Read(b[:2])
	h.PacketType.Read(b[2:4])
	return h
}

func ReadOpen(c net.Conn) Open {
	var b = make([]byte, OPEN_LENGTH)
	_, err := c.Read(b)
	errors.CheckError(err)
	o := Open{}
	o.DataLength = binary.BigEndian.Uint64(b[:8])
	o.Reverse = binary.BigEndian.Uint16(b[8:10])
	return o
}

func ReadData(conn net.Conn, l uint64) []byte {
	var b = make([]byte, l)
	_, err := conn.Read(b)
	errors.CheckError(err)
	return b
}

func ReadPacket(conn net.Conn) Packet {
	h := ReadHeader(conn)

	p := Packet{}
	p.Header = &h
	switch {
	case p.Header.PacketType.Value == OPEN_TYPE:
		o := ReadOpen(conn)
		p.Message = &o
	}

	return p
}

/*
SendOpen
	conn: Net.conn instance
	dl:	Data length
	r: Reverse option
*/
func SendOpen(conn net.Conn, dl uint64, r uint16) {
	h := Header{}
	pl := IntField{}
	pl.Write(4)

	h.AddField(&pl)
	pt := IntField{}
	pt.Write(OPEN_TYPE)
	h.AddField(&pt)
	p := Packet{}

	p.Header = &h
	msg := Open{}
	msg.Write(dl, r)

	p.Message = &msg
	b := p.Serialize()
	_, err := conn.Write(b)
	errors.CheckError(err)
}

func SendConfirm(conn net.Conn) {
	p := Header{}
	pl := IntField{}
	pl.Write(4)

	p.AddField(&pl)
	pt := IntField{}
	pt.Write(CONFIRM_TYPE)
	p.AddField(&pt)

	b := p.Serialize()
	_, err := conn.Write(b)
	errors.CheckError(err)
}
