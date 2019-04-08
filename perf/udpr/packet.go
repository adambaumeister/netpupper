package udpr

import (
	"encoding/binary"
	"github.com/adamb/netpupper/errors"
	"net"
)

const OPEN_TYPE = 1
const CONFIRM_TYPE = 2
const CLOSE_TYPE = 3
const DG_TYPE = 4
const ACK_TYPE = 5

const PACKETTYPE_LENGTH = 2
const PACKETLEN_LENGTH = 2

const HEADER_LENGTH = 4
const OPEN_LENGTH = 12
const DG_MIN_LENGTH = 10

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
	AckCount   uint32
}

func (m *Open) Serialize() []byte {
	b := make([]byte, OPEN_LENGTH)
	binary.BigEndian.PutUint64(b[:8], m.DataLength)
	return b
}
func (m *Open) Write(i uint64, ac uint32) {
	m.DataLength = i
	m.AckCount = ac
}

/*
Datagram message, sends data
	Datalength: Length of data in this datagram
	Sequence: Sequence number of this UDP packet
	Data: Actual data contained within this datagram
*/
type Datagram struct {
	DataLength uint16
	Sequence   uint64
	Data       []byte
}

func (m *Datagram) Serialize() []byte {
	b := make([]byte, 10)
	binary.BigEndian.PutUint16(b[:2], m.DataLength)
	binary.BigEndian.PutUint64(b[2:10], m.Sequence)
	b = append(b, m.Data...)
	return b
}
func (m *Datagram) Write(b []byte, sn uint64) {
	m.DataLength = uint16(len(b))
	m.Sequence = sn
	m.Data = b
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

func HeaderFromBytes(b []byte) Header {
	h := Header{}
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

func ReadHeader(packet []byte) Header {
	h := HeaderFromBytes(packet[:HEADER_LENGTH])

	return h
}

func ReadOpen(b []byte) Open {
	o := Open{}
	o.DataLength = binary.BigEndian.Uint64(b[HEADER_LENGTH:12])
	return o
}

func ReadDatagram(b []byte) Datagram {
	d := Datagram{}
	// Cut out the header
	b = b[HEADER_LENGTH:]

	// First two bytes are the length
	d.DataLength = binary.BigEndian.Uint16(b[:2])
	// Next 8 are sequence number
	d.Sequence = binary.BigEndian.Uint64(b[2:10])
	d.Data = b[10 : 10+d.DataLength]
	return d
}

func ReadPacket(c *net.UDPConn) (Packet, *net.UDPAddr) {
	packet := make([]byte, 1500)
	_, addr, err := c.ReadFromUDP(packet)

	errors.CheckError(err)
	h := ReadHeader(packet)

	p := Packet{}
	p.Header = &h
	return p, addr
}

/*
SendOpen
	conn: Net.conn instance
	dl:	Total length of UDP trans, 0 for infinite
	ac: Ack count - count of packets to require an acknowelegement
*/
func SendOpen(conn net.Conn, dl uint64, ac uint32) {
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

	msg.Write(dl, ac)

	p.Message = &msg
	b := p.Serialize()
	_, err := conn.Write(b)
	errors.CheckError(err)
}

func SendConfirm(conn *net.UDPConn, a *net.UDPAddr) {
	p := Header{}
	pl := IntField{}
	pl.Write(4)

	p.AddField(&pl)
	pt := IntField{}
	pt.Write(CONFIRM_TYPE)
	p.AddField(&pt)

	b := p.Serialize()
	_, err := conn.WriteToUDP(b, a)
	errors.CheckError(err)
}

func SendDatagram(conn net.Conn, sn uint64, data []byte) {
	h := Header{}
	pt := IntField{}
	pt.Write(DG_TYPE)

	dg := Datagram{}
	dg.Write(data, sn)
	pl := IntField{}
	pl.Write(4 + DG_MIN_LENGTH + dg.DataLength)

	h.AddField(&pl)
	h.AddField(&pt)

	p := Packet{}
	p.Header = &h
	p.Message = &dg

	b := p.Serialize()
	_, err := conn.Write(b)
	errors.CheckError(err)
}

func SendClose(conn net.Conn) {
	p := Header{}
	pl := IntField{}
	pl.Write(4)

	p.AddField(&pl)
	pt := IntField{}
	pt.Write(CLOSE_TYPE)
	p.AddField(&pt)

	b := p.Serialize()
	_, err := conn.Write(b)
	errors.CheckError(err)
}
