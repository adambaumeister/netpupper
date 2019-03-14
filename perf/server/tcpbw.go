package server

import (
	"fmt"
	"github.com/adamb/netpupper/errors"
	"net"
	"time"
)

func Server() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		errors.RaiseError("Failed to open socket.")
	}
	for {
		conn, err := ln.Accept()
		errors.CheckError(err)

		h := ReadHeader(conn)
		switch {
		case h.PacketType.Value == OPEN_TYPE:
			var o = Open{}
			o = ReadOpen(conn)
			fmt.Printf("Got a connection from: %v, Packet Type: %v Data to follow: %v bytes\n", conn.RemoteAddr(),
				h.PacketType.Value, o.DataLength)
			SendConfirm(conn)

			timedRead(conn, o.DataLength)
		}
	}
}

func Client() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		errors.RaiseError("Failed to open connection!")
	}
	fmt.Printf("Succesfully connected to: %v\n", conn.RemoteAddr())

	// Send the open message, request to start
	SendOpen(conn)
	// Wait for a confirmation
	h := ReadHeader(conn)
	switch {
	case h.PacketType.Value == CONFIRM_TYPE:
		fmt.Printf("OPEN Request confirmed. Sending data...\n")
		// Test by splitting up the data
		conn.Write([]byte{255, 255})
		time.Sleep(2 * time.Second)
		conn.Write([]byte{255, 255})
		time.Sleep(1 * time.Second)
	}
}

func timedRead(conn net.Conn, rl uint64) {
	start := time.Now().UnixNano()

	// Chunk size is how much we read at each time interval
	chunk := rl / 4
	data := make([]byte, rl)

	currentChunk := 1
	// Read each chunk until we've read the entire thing
	for currentChunk <= 4 {
		chunkData := make([]byte, chunk)
		conn.Read(chunkData)
		fmt.Printf("Chunkdat: %v\n", chunkData)
		// Append each read chunk to the full data array
		data = append(data, chunkData...)
		currentChunk = currentChunk + 1
	}
	t := time.Now().UnixNano()
	elapsed := t - start
	fmt.Printf("Read took %v ns\n", elapsed)

}
