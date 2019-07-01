package main

import (
	"fmt"
	"github.com/anidotnet/openvidu-go-client/openvidu"
)

func main() {
	c1 := &openvidu.Connection{
		ConnectionId: "abc",
		CreatedAt:    12346,
		Role:         openvidu.MODERATOR,
	}
	c2 := &openvidu.Connection{
		ConnectionId: "xyz",
		CreatedAt:    91635,
		Role:         openvidu.PUBLISHER,
	}

	ac := make(map[string]*openvidu.Connection, 2)
	ac["abc"] = c1
	ac["xyz"] = c2

	p := &openvidu.SessionProperties{
		DefaultCustomLayout:    "abcd",
		DefaultRecordingLayout: openvidu.BEST_FIT,
		DefaultOutputMode:      openvidu.COMPOSED,
		RecordingMode:          openvidu.MANUAL,
		MediaMode:              openvidu.RELAYED,
		CustomSessionId:        "abcd",
	}

	s := &openvidu.Session{
		SessionId:         "13abcd",
		CreatedAt:         15648793,
		Recording:         true,
		ActiveConnections: ac,
		Properties:        p,
	}

	json, e := s.ToJson()
	if e != nil {
		fmt.Printf("Error %v", e)
	}

	fmt.Printf("%s", json)
}
