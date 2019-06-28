package main

import (
	"fmt"
	"github.com/anidotnet/openvidu-go-client/openvidu"
)

func main() {
	connection := &openvidu.Connection{
		ConnectionId: "abc",
		CreatedAt: 12346,
		Role: openvidu.MODERATOR,
	}
	fmt.Printf("%+v", connection)
}
