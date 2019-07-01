package openvidu

type Connection struct {
	ConnectionId string
	CreatedAt    int64
	Role         OpenViduRole
	Token        string
	Location     string
	Platform     string
	ServerData   string
	ClientData   string
	Publishers   map[string]*Publisher
	Subscribers  []string
}

func (c *Connection) GetPublishers() []*Publisher {
	v := make([]*Publisher, 0)
	for  _, value := range c.Publishers {
		v = append(v, value)
	}
	return v
}
