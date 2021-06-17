package websocket

type client struct {
	Ip string
	Id int
	HeartBreath uint64
}

func NewClient(ip string, id int, heartBreath uint64) *client {
	return &client{
		Ip:ip,
		Id:id,
		HeartBreath:heartBreath,
	}
}