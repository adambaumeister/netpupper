package server

import (
	"encoding/binary"
	"fmt"
	"github.com/adamb/netpupper/errors"
	"net"
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

			data := make([]byte, o.DataLength)
			conn.Read(data)
			fmt.Printf("Value: %v\n", binary.BigEndian.Uint32(data))
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
		conn.Write([]byte{255, 255, 255, 255})
	}
}
