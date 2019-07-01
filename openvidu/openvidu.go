package openvidu

import (
	"crypto/tls"
	"encoding/base64"
	"net/http"
	"strings"
	"time"
)

const (
	API_SESSIONS         = "api/sessions"
	API_TOKENS           = "api/tokens"
	API_RECORDINGS       = "api/recordings"
	API_RECORDINGS_START = "/start"
	API_RECORDINGS_STOP  = "/stop"
)

type OpenVidu struct {
	hostName       string
	secret         string
	activeSessions map[string]*Session
	httpClient     *http.Client
	basicAuth      string
}

type serverActiveSessions struct {
	NumberOfElements int              `json:"numberOfElements"`
	Content          []*serverSession `json:"content"`
}

type serverSession struct {
	SessionId              string           `json:"sessionId"`
	CreatedAt              int64            `json:"createdAt"`
	MediaMode              MediaMode        `json:"mediaMode"`
	RecordingMode          RecordingMode    `json:"recordingMode"`
	DefaultOutputMode      OutputMode       `json:"defaultOutputMode"`
	DefaultRecordingLayout RecordingLayout  `json:"defaultRecordingLayout"`
	DefaultCustomLayout    string           `json:"defaultCustomLayout"`
	CustomSessionId        string           `json:"customSessionId"`
	Connections            *connectionsInfo `json:"connections"`
	Recording              bool             `json:"recording"`
}
type connectionsInfo struct {
	NumberOfElements int                  `json:"numberOfElements"`
	Content          []*connectionContent `json:"content"`
}
type connectionContent struct {
	ConnectionId string        `json:"connectionId"`
	CreatedAt    int64         `json:"createdAt"`
	Location     string        `json:"location"`
	Platform     string        `json:"platform"`
	Token        string        `json:"token"`
	Role         OpenViduRole  `json:"role"`
	ServerData   string        `json:"serverData"`
	ClientData   string        `json:"clientData"`
	Publishers   []*publisher  `json:"publishers"`
	Subscribers  []*subscriber `json:"subscribers"`
}
type subscriber struct {
	CreatedAt int64  `json:"createdAt"`
	StreamID  string `json:"streamId"`
	Publisher string `json:"publisher"`
}
type publisher struct {
	CreatedAt    int64         `json:"createdAt"`
	StreamID     string        `json:"streamId"`
	MediaOptions *mediaOptions `json:"mediaOptions"`
}

type mediaOptions struct {
	HasAudio        bool                   `json:"hasAudio"`
	AudioActive     bool                   `json:"audioActive"`
	HasVideo        bool                   `json:"hasVideo"`
	VideoActive     bool                   `json:"videoActive"`
	TypeOfVideo     string                 `json:"typeOfVideo"`
	FrameRate       int32                  `json:"frameRate"`
	VideoDimensions string                 `json:"videoDimensions"`
	Filter          map[string]interface{} `json:"filter"`
}

func NewOpenVidu(hostName string, secret string) *OpenVidu {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	openVidu := &OpenVidu{
		hostName:   hostName,
		secret:     secret,
		httpClient: &http.Client{Transport: tr, Timeout: 30 * time.Second},
		basicAuth:  base64.StdEncoding.EncodeToString([]byte("OPENVIDUAPP:" + secret)),
	}

	if !strings.HasSuffix(openVidu.hostName, "/") {
		openVidu.hostName = openVidu.hostName + "/"
	}

	return openVidu
}

func (o *OpenVidu) CreateSession0() (*Session, error) {
	session, err := NewSession0(o)
	if err != nil {
		return nil, err
	}
	o.activeSessions[session.SessionId] = session
	return session, nil
}

func (o *OpenVidu) CreateSession1(properties *SessionProperties) (*Session, error) {
	session, err := NewSession1(o, properties)
	if err != nil {
		return nil, err
	}
	o.activeSessions[session.SessionId] = session
	return session, nil
}

func (o *OpenVidu) StartRecording(sessionId string, properties *RecordingProperties) (*Recording, error) {

}
