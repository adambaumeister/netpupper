package udpr

import (
	"context"
	"fmt"
	"github.com/adamb/netpupper/perf/stats"
	"golang.org/x/time/rate"
	"net"
	"time"
)

/*
UDPTransport represents a state-ish machine for a UDP reliability test.

It is responsible for tracking the latency, jitter, and loss within the test by counting and validating sequence numbers.

*/
type UdpTransport struct {
	conn            net.Conn
	addr            *net.UDPAddr
	window          uint32
	maxlength       uint64
	timeout         time.Duration
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
	u.timeout = 3 * time.Second

	return u
}

/* Read from the connection
Will send an ACK every AckCount packets (ac)
Ends after max length
*/
func (u *UdpTransport) countedRead(uc *net.UDPConn, test *stats.Test) {
	packet := make([]byte, 1500)
	count := uint32(0)
	// Read all incoming udp packets
	for u.CurrentSequence < u.maxlength {
		u.conn.SetDeadline(time.Now().Local().Add(u.timeout))
		_, addr, err := uc.ReadFromUDP(packet)
		if err != nil {
			fmt.Printf("Failed read(%v): %v, %v\n", err, u.CurrentSequence, u.maxlength)
		} else {

			h := ReadHeader(packet)
			if h.PacketType.Value == DG_TYPE {
				d := ReadDatagram(packet)
				u.CheckSequence(d)
				u.CurrentSequence = d.Sequence
				//fmt.Printf("Seq: %v\n", u.CurrentSequence)
			}
			if count == u.window {
				loss := len(u.Buffer)
				ef := len(u.EffectiveLost)
				fmt.Printf("Loss count: %v\n", loss)
				r := stats.ReliabilityResult{
					Loss:          loss,
					EffectiveLoss: ef,
				}
				test.InRelTests <- r
				count = 0

				SendAck(uc, addr, uint32(loss), uint32(ef))

			}
		}
		count = count + 1
	}
}

func (u *UdpTransport) countedSend(test *stats.Test, ratel int) {
	ctx := context.Background()
	// 10000 events p/s goal
	// 10000 events p/s * 1000 bytes/event = 10MBps
	limit := CalcLimiter(ratel)
	count := uint32(0)
	for u.CurrentSequence <= u.maxlength {
		u.conn.SetDeadline(time.Now().Local().Add(u.timeout))
		limit.Wait(ctx)

		SendDatagram(u.conn, u.CurrentSequence, []byte{1, 1, 1, 1})
		u.CurrentSequence = u.CurrentSequence + 1
		if count == u.window {
			packet := make([]byte, 1500)
			_, err := u.conn.Read(packet)
			if err != nil {
				fmt.Printf("Failed read(%v): %v, %v\n", err, u.CurrentSequence, u.maxlength)
			}
			h := ReadHeader(packet)
			if h.PacketType.Value == ACK_TYPE {
				ack := ReadAck(packet)
				r := stats.ReliabilityResult{
					Loss:          int(ack.Loss),
					EffectiveLoss: int(ack.EffectiveLoss),
				}
				test.InRelTests <- r
			}
			count = 0
		}
		count = count + 1
	}
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

func CalcLimiter(cir int) *rate.Limiter {
	// bc = cir *tc /1000
	// below chops one second up into 100 time intervals
	bc := cir * 10 / 1000
	limit := rate.Limit(cir)
	l := rate.NewLimiter(limit, bc)
	return l
}
