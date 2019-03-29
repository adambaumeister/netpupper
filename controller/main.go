package controller

type Client struct {
	Addr string
	Tags []Tag
}

type Tag struct {
	Name  string
	Value string
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
