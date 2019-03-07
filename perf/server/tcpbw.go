package server

import (
	"../../errors"
	"fmt"
	"net"
)

func Server() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		errors.RaiseError("Failed to open socket.")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			errors.RaiseError("Failed to receive connection")
		}
		var b = make([]byte, 4)
		_, err = conn.Read(b)
		h := ReadHeader(b)
		errors.CheckError(err)
		fmt.Printf("Got a connection from: %v, Packet Type: %v\n", conn.RemoteAddr(), h.PacketType.Value)
	}
}

func Client() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		errors.RaiseError("Failed to open connection!")
	}
	fmt.Printf("Succesfully connected to: %v\n", conn.RemoteAddr())
	p := Header{}
	pl := IntField{}
	pl.Write(4)

	p.Base.AddField(pl)
	pt := IntField{}
	pt.Write(1)
	p.PacketLength = &pl
	p.PacketType = &pt

	b := p.Serialize()
	_, err = conn.Write(b)
	errors.CheckError(err)
}
