package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adamb/netpupper/errors"
	"github.com/adamb/netpupper/perf/tcpbw"
	"net/http"
)

type Client struct {
	Server string
	Addr   string
	Tags   map[string]string

	ApiAddr   string
	ByteCount string
}

/*
Start a bandwidth test from a client
*/
func (c *Client) StartbwTest() {
	cc := tcpbw.ClientConfig{
		Server: c.Server,
		Bytes:  c.ByteCount,
	}
	b, err := json.Marshal(cc)
	errors.CheckError(err)
	_, err = http.Post(fmt.Sprintf("http://%v/tcpbw", c.ApiAddr), "application/json", bytes.NewBuffer(b))
	errors.CheckError(err)

}

type Controller struct {
	Clients []Client
}

func (c *Controller) AddClient(client Client) {
	c.Clients = append(c.Clients, client)
}

func (c *Controller) GetFirstClient() Client {
	return c.Clients[0]
}
