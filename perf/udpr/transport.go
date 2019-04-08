package udpr

import (
	"fmt"
	"github.com/adamb/netpupper/errors"
	"net"
)

/*
UDPTransport represents a state-ish machine for a UDP reliability test.

It is responsible for tracking the latency, jitter, and loss within the test by counting and validating sequence numbers.

*/
type UdpTransport struct {
	conn      net.Conn
	addr      *net.UDPAddr
	window    uint32
	maxlength uint64

	CurrentSequence uint64
	Buffer          []Datagram

	EffectiveLost []Datagram
}

/*
Create a UDP state machine
*/
func InitUdpSm(conn net.Conn, addr *net.UDPAddr, ac uint32, ml uint64) UdpTransport {
	u := UdpTransport{
		conn:      conn,
		addr:      addr,
		window:    ac,
		maxlength: ml,

		CurrentSequence: 0,
	}
	return u
}

/* Read from the connection
Will send an ACK every AckCount packets (ac)
Ends after max length
*/
func (u *UdpTransport) countedRead() {
	packet := make([]byte, 1500)

	// Read all incoming udp packets
	for u.CurrentSequence < u.maxlength {
		_, err := u.conn.Read(packet)
		errors.CheckError(err)

		h := ReadHeader(packet)
		if h.PacketType.Value == DG_TYPE {
			d := ReadDatagram(packet)
			fmt.Printf("Got a datagram. sequence %v Data: %vb\n", d.Sequence, d.Data)
			u.CheckSequence(d)
			u.CurrentSequence = d.Sequence
		}
	}
	fmt.Printf("finished read. Lost: %v, Eff Lost: %v\n", len(u.Buffer), len(u.EffectiveLost))

}

func (u *UdpTransport) countedSend() {
	for u.CurrentSequence <= u.maxlength {
		SendDatagram(u.conn, u.CurrentSequence, []byte{1, 1, 1, 1})
		u.CurrentSequence = u.CurrentSequence + 1
	}
	fmt.Printf("Got here")
}

func (u *UdpTransport) CheckSequence(d Datagram) {
	// If the datagram sequence number is ahead of what we expect
	// This is potential loss
	if d.Sequence > u.CurrentSequence+1 {
		u.Buffer = append(u.Buffer, d)
	}

	// If the datagram sn is less than what we expect -
	// This is a symptom of a lot of jitter
	// 1 2 4 5 3
	if d.Sequence < u.CurrentSequence {
		// Check to see if the datagram with the next sequence number exists in the buffer
		i := 0
		for _, datagram := range u.Buffer {
			if datagram.Sequence == d.Sequence+1 {
				u.EffectiveLost = append(u.EffectiveLost, datagram)
				// clear the previously stored datagram from the buffer
				u.Buffer = append(u.Buffer[:i], u.Buffer[i+1:]...)
			}
			i = i + 1
		}
	}
}
