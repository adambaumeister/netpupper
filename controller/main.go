package controller

type Client struct {
	TcpbwAddr string
	Tags      []Tag
}

type Tag struct {
	Name  string
	Value string
}

type Controller struct {
	Clients []Client
}
