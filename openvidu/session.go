package openvidu

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Session struct {
	openVidu          *OpenVidu
	SessionId         string
	CreatedAt         int64
	Properties        *SessionProperties
	ActiveConnections map[string]*Connection
	Recording         bool
}

type sessionJson struct {
	SessionId              string          `json:"sessionId"`
	CreatedAt              int64           `json:"createdAt"`
	CustomSessionId        string          `json:"customSessionId"`
	Recording              bool            `json:"recording"`
	MediaMode              MediaMode       `json:"mediaMode"`
	RecordingMode          RecordingMode   `json:"recordingMode"`
	DefaultOutputMode      OutputMode      `json:"defaultOutputMode"`
	DefaultRecordingLayout RecordingLayout `json:"defaultRecordingLayout"`
	DefaultCustomLayout    string          `json:"defaultCustomLayout"`
	Connections            *connections    `json:"connections"`
}

type connections struct {
	NumberOfElements int               `json:"numberOfElements"`
	Content          []*connectionJson `json:"content"`
}

type connectionJson struct {
	ConnectionId string       `json:"connectionId"`
	Role         OpenViduRole `json:"role"`
	Token        string       `json:"token"`
	ClientData   string       `json:"clientData"`
	ServerData   string       `json:"serverData"`
	Publishers   []*Publisher `json:"publishers"`
	Subscribers  []string     `json:"subscribers"`
}

type tokenRequest struct {
	Session        string          `json:"session"`
	Role           OpenViduRole    `json:"role"`
	Data           string          `json:"data"`
	KurentoOptions *KurentoOptions `json:"kurentoOptions"`
}

func NewSession0(o *OpenVidu) (*Session, error) {
	session := &Session{
		openVidu: o,
		Properties: &SessionProperties{
			CustomSessionId:        "",
			MediaMode:              ROUTED,
			RecordingMode:          MANUAL,
			DefaultOutputMode:      COMPOSED,
			DefaultRecordingLayout: BEST_FIT,
			DefaultCustomLayout:    "",
		},
	}
	err := session.getSessionIdHttp()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func NewSession1(ov *OpenVidu, properties *SessionProperties) (*Session, error) {
	session := &Session{
		openVidu:   ov,
		Properties: properties,
	}
	err := session.getSessionIdHttp()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func NewSession2(ov *OpenVidu, json *serverSession) (*Session, error) {
	session := &Session{
		openVidu: ov,
	}
	session.resetSessionWithJson(json)
	return session, nil
}

func (s *Session) GenerateToken(to *TokenOptions) (string, error) {
	if to == nil {
		to = &TokenOptions{
			Data: "",
			Role: PUBLISHER,
		}
	}

	obj := &tokenRequest{
		Session: s.SessionId,
		Role:    to.Role,
		Data:    to.Data,
	}

	if to.KurentoOptions != nil {
		if to.KurentoOptions.VideoMaxRecvBandwidth != nil {
			obj.KurentoOptions.VideoMaxRecvBandwidth = to.KurentoOptions.VideoMaxRecvBandwidth
		}
		if to.KurentoOptions.VideoMinRecvBandwidth != nil {
			obj.KurentoOptions.VideoMinRecvBandwidth = to.KurentoOptions.VideoMinRecvBandwidth
		}
		if to.KurentoOptions.VideoMaxSendBandwidth != nil {
			obj.KurentoOptions.VideoMaxSendBandwidth = to.KurentoOptions.VideoMaxSendBandwidth
		}
		if to.KurentoOptions.VideoMinSendBandwidth != nil {
			obj.KurentoOptions.VideoMinSendBandwidth = to.KurentoOptions.VideoMinSendBandwidth
		}
		if len(to.KurentoOptions.AllowedFilters) > 0 {
			obj.KurentoOptions.AllowedFilters = to.KurentoOptions.AllowedFilters
		}
	}

	reqString, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	url := s.openVidu.hostName + API_TOKENS
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqString))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+s.openVidu.basicAuth)
	response, err := s.openVidu.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		res := struct {
			Id string `json:"id"`
		}{}
		err = json.Unmarshal(body, &res)
		if err != nil {
			return "", err
		}

		return res.Id, nil
	}
	return "", newOpenViduError(statusCode)
}

func (s *Session) GetActiveConnections() []*Connection {
	v := make([]*Connection, 0, len(s.ActiveConnections))
	for _, value := range s.ActiveConnections {
		v = append(v, value)
	}
	return v
}

func (s *Session) Close() error {
	url := s.openVidu.hostName + API_SESSIONS + "/" + s.SessionId
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+s.openVidu.basicAuth)
	response, err := s.openVidu.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode == http.StatusNoContent {
		delete(s.openVidu.activeSessions, s.SessionId)
	} else {
		return newOpenViduError(statusCode)
	}
	return nil
}

func (s *Session) Fetch() (bool, error) {
	beforeJson, err := s.ToJson()
	if err != nil {
		return false, err
	}

	url := s.openVidu.hostName + API_SESSIONS + "/" + s.SessionId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+s.openVidu.basicAuth)
	response, err := s.openVidu.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return false, err
		}

		var response serverSession
		err = json.Unmarshal(body, &response)
		if err != nil {
			return false, err
		}

		s.resetSessionWithJson(&response)
		afterJson, err := s.ToJson()
		if err != nil {
			return false, err
		}

		if strings.Compare(beforeJson, afterJson) != 0 {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, newOpenViduError(statusCode)
	}
}

func (s *Session) ForceDisconnect(c *Connection) error {
	return s.ForceDisconnectById(c.ConnectionId)
}

func (s *Session) ForceDisconnectById(connectionId string) error {
	url := s.openVidu.hostName + API_SESSIONS + "/" + s.SessionId + "/connection/" + connectionId
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+s.openVidu.basicAuth)
	response, err := s.openVidu.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode == http.StatusNoContent {
		connectionClosed := s.ActiveConnections[connectionId]
		delete(s.ActiveConnections, connectionId)

		if connectionClosed != nil {
			for _, publisher := range connectionClosed.Publishers {
				streamId := publisher.StreamId
				for _, connection := range s.ActiveConnections {
					var newSubscribers []string
					for _, subscriber := range connection.Subscribers {
						if strings.Compare(streamId, subscriber) != 0 {
							newSubscribers = append(newSubscribers, subscriber)
						}
					}
					connection.Subscribers = newSubscribers
				}
			}
		}
	} else {
		return newOpenViduError(statusCode)
	}

	return nil
}

func (s *Session) ForceUnpublish(pub *Publisher) error {
	return s.ForceUnpublishById(pub.StreamId)
}

func (s *Session) ForceUnpublishById(streamId string) error {
	url := s.openVidu.hostName + API_SESSIONS + "/" + s.SessionId + "/stream/" + streamId
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+s.openVidu.basicAuth)
	response, err := s.openVidu.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode == http.StatusNoContent {
		for _, connection := range s.ActiveConnections {
			if connection.Publishers[streamId] != nil {
				delete(connection.Publishers, streamId)
				continue
			}

			var newSubscribers []string
			for _, subscriber := range connection.Subscribers {
				if strings.Compare(streamId, subscriber) != 0 {
					newSubscribers = append(newSubscribers, subscriber)
				}
			}
			connection.Subscribers = newSubscribers
		}
	} else {
		return newOpenViduError(statusCode)
	}
	return nil
}

func (s *Session) String() string {
	return s.SessionId
}

func (s *Session) ToJson() (string, error) {
	ac := s.GetActiveConnections()
	var content []*connectionJson
	for _, con := range ac {
		cJson := &connectionJson{
			ConnectionId: con.ConnectionId,
			Role:         con.Role,
			Token:        con.Token,
			ClientData:   con.ClientData,
			ServerData:   con.ServerData,
			Publishers:   con.GetPublishers(),
			Subscribers:  con.Subscribers,
		}
		content = append(content, cJson)
	}

	sJson := &sessionJson{
		SessionId:              s.SessionId,
		CreatedAt:              s.CreatedAt,
		CustomSessionId:        s.Properties.CustomSessionId,
		Recording:              s.Recording,
		MediaMode:              s.Properties.MediaMode,
		RecordingMode:          s.Properties.RecordingMode,
		DefaultOutputMode:      s.Properties.DefaultOutputMode,
		DefaultRecordingLayout: s.Properties.DefaultRecordingLayout,
		DefaultCustomLayout:    s.Properties.DefaultCustomLayout,
		Connections: &connections{
			NumberOfElements: len(ac),
			Content:          content,
		},
	}

	b, err := json.Marshal(sJson)
	if err != nil {
		return "{}", err
	}
	return string(b), nil
}

func (s *Session) getSessionIdHttp() error {
	if len(s.SessionId) > 0 {
		return nil
	}

	url := s.openVidu.hostName + API_SESSIONS
	obj := &SessionProperties{
		MediaMode:              s.Properties.MediaMode,
		RecordingMode:          s.Properties.RecordingMode,
		DefaultOutputMode:      s.Properties.DefaultOutputMode,
		DefaultRecordingLayout: s.Properties.DefaultRecordingLayout,
		DefaultCustomLayout:    s.Properties.DefaultCustomLayout,
		CustomSessionId:        s.Properties.CustomSessionId,
	}

	reqString, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqString))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+s.openVidu.basicAuth)
	response, err := s.openVidu.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		res := struct {
			Id        string `json:"id"`
			CreatedAt int64  `json:"createdAt"`
		}{}
		err = json.Unmarshal(body, &res)
		if err != nil {
			return err
		}

		s.SessionId = res.Id
		s.CreatedAt = res.CreatedAt
	} else if statusCode == http.StatusConflict {
		s.SessionId = s.Properties.CustomSessionId
	} else {
		return newOpenViduError(statusCode)
	}
	return nil
}

func (s *Session) resetSessionWithJson(sj *serverSession) {
	s.SessionId = sj.SessionId
	s.CreatedAt = sj.CreatedAt
	s.Recording = sj.Recording

	sp := &SessionProperties{
		MediaMode:         sj.MediaMode,
		RecordingMode:     sj.RecordingMode,
		DefaultOutputMode: sj.DefaultOutputMode,
	}
	if len(sj.DefaultRecordingLayout) > 0 {
		sp.DefaultRecordingLayout = sj.DefaultRecordingLayout
	}
	if len(sj.DefaultCustomLayout) > 0 {
		sp.DefaultCustomLayout = sj.DefaultCustomLayout
	}
	if s.Properties != nil && len(s.Properties.CustomSessionId) > 0 {
		sp.CustomSessionId = s.Properties.CustomSessionId
	} else if len(sj.CustomSessionId) > 0 {
		sp.CustomSessionId = sj.CustomSessionId
	}
	s.Properties = sp

	connArray := sj.Connections.Content
	s.ActiveConnections = make(map[string]*Connection, 0)
	for _, con := range connArray {
		publishers := con.Publishers
		pubMap := make(map[string]*Publisher, 0)
		for _, publisher := range publishers {
			mediaOptions := publisher.MediaOptions
			p := &Publisher{
				CreatedAt:       publisher.CreatedAt,
				AudioActive:     mediaOptions.AudioActive,
				FrameRate:       mediaOptions.FrameRate,
				HasAudio:        mediaOptions.HasAudio,
				HasVideo:        mediaOptions.HasVideo,
				StreamId:        publisher.StreamID,
				TypeOfVideo:     mediaOptions.TypeOfVideo,
				VideoActive:     mediaOptions.VideoActive,
				VideoDimensions: mediaOptions.VideoDimensions,
			}
			pubMap[p.StreamId] = p
		}

		subscribers := make([]string, 0)
		for _, subscriber := range con.Subscribers {
			subscribers = append(subscribers, subscriber.StreamID)
		}

		s.ActiveConnections[con.ConnectionId] = &Connection{
			CreatedAt:    con.CreatedAt,
			Subscribers:  subscribers,
			ServerData:   con.ServerData,
			ClientData:   con.ClientData,
			Token:        con.Token,
			Role:         con.Role,
			ConnectionId: con.ConnectionId,
			Publishers:   pubMap,
			Location:     con.Location,
			Platform:     con.Platform,
		}
	}
}

//TODO: fetch all activeSessions from server GET /api/sessions, response serverActiveSessions
