package openvidu

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
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
	url := o.hostName + API_RECORDINGS + API_RECORDINGS_START
	rj := &recordingJson{
		SessionId:  sessionId,
		Name:       properties.Name,
		OutputMode: properties.OutputMode,
		HasAudio:   properties.HasAudio,
		HasVideo:   properties.HasVideo,
	}

	if properties.OutputMode == COMPOSED && properties.HasVideo {
		rj.Resolution = properties.Resolution
		rj.RecordingLayout = properties.RecordingLayout
		if properties.RecordingLayout == CUSTOM {
			rj.CustomLayout = properties.CustomLayout
		}
	}

	reqString, err := json.Marshal(rj)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqString))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+o.basicAuth)
	response, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var rj *recordingJson
		err = json.Unmarshal(body, &rj)
		if err != nil {
			return nil, err
		}

		r := NewRecording(rj)

		activeSession := o.activeSessions[r.SessionId]
		if activeSession != nil {
			activeSession.Recording = true
		}
		return r, nil
	} else {
		return nil, newOpenViduError(statusCode)
	}
}

func (o *OpenVidu) StartRecordingByName(sessionId string, name string) (*Recording, error) {
	rp := &RecordingProperties{
		Name:       name,
		OutputMode: COMPOSED,
		HasAudio:   true,
		HasVideo:   true,
	}

	return o.StartRecording(sessionId, rp.Build())
}

func (o *OpenVidu) StartRecordingById(sessionId string) (*Recording, error) {
	return o.StartRecordingByName(sessionId, "")
}

func (o *OpenVidu) StopRecording(recordingId string) (*Recording, error)  {
	url := o.hostName + API_RECORDINGS + API_RECORDINGS_STOP + "/" + recordingId
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+o.basicAuth)
	response, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var rj *recordingJson
		err = json.Unmarshal(body, &rj)
		if err != nil {
			return nil, err
		}

		r := NewRecording(rj)
		activeSession := o.activeSessions[r.SessionId]
		if activeSession != nil {
			activeSession.Recording = false
		}
		return r, nil
	} else {
		return nil, newOpenViduError(statusCode)
	}
}

func (o *OpenVidu) GetRecording(recordingId string) (*Recording, error) {
	url := o.hostName + API_RECORDINGS + "/" + recordingId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+o.basicAuth)
	response, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var rj *recordingJson
		err = json.Unmarshal(body, &rj)
		if err != nil {
			return nil, err
		}

		r := NewRecording(rj)
		return r, nil
	} else {
		return nil, newOpenViduError(statusCode)
	}
}

func (o *OpenVidu) ListRecording() ([]*Recording, error) {
	url := o.hostName + API_RECORDINGS
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+o.basicAuth)
	response, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		obj := struct {
			Items	[]*recordingJson	`json:"items"`
		}{}

		err = json.Unmarshal(body, &obj)
		if err != nil {
			return nil, err
		}

		var recordings []*Recording
		for _, item := range obj.Items {
			r := NewRecording(item)
			recordings = append(recordings, r)
		}

		return recordings, nil
	} else {
		return nil, newOpenViduError(statusCode)
	}
}

func (o *OpenVidu) DeleteRecording(recordingId string) error {
	url := o.hostName + API_RECORDINGS + "/" + recordingId
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Basic "+o.basicAuth)
	response, err := o.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode != http.StatusNoContent {
		return newOpenViduError(statusCode)
	} else {
		return nil
	}
}

func (o *OpenVidu) GetActiveSessions() []*Session {
	var sessions []*Session
	for _, v := range o.activeSessions {
		sessions = append(sessions, v)
	}
	return sessions
}

func (o *OpenVidu) Fetch() (bool, error) {
	url := o.hostName + API_SESSIONS
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Basic "+o.basicAuth)
	response, err := o.httpClient.Do(req)
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

		var sas *serverActiveSessions
		err = json.Unmarshal(body, &sas)
		if err != nil {
			return false, err
		}

		var fetchedSessionIds []string
		hasChanged := false
		for _, session := range sas.Content {
			sessionId := session.SessionId
			fetchedSessionIds = append(fetchedSessionIds, sessionId)
			computeIfPresent(o.activeSessions, sessionId, func(sId string, s *Session) *Session {
				beforeJSON, _ := s.ToJson()
				s.resetSessionWithJson(session)
				afterJSON, _ := s.ToJson()
				changed := strings.Compare(beforeJSON, afterJSON) != 0
				hasChanged = hasChanged || changed
				return s
			})

			computeIfAbsent(o.activeSessions, sessionId, func(sId string) *Session {
				hasChanged = true
				s, _ := NewSession2(o, session)
				return s
			})
		}

		newActiveSessions := make(map[string]*Session, 0)
		for k, v := range o.activeSessions {
			if contains(fetchedSessionIds, k) {
				newActiveSessions[k] = v
			} else {
				hasChanged = true
			}
		}
		o.activeSessions = newActiveSessions
		return hasChanged, nil
	} else {
		return false, newOpenViduError(statusCode)
	}
}

func (o *OpenVidu) FetchSessions() ([]*Session, error) {
	url := o.hostName + API_SESSIONS
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+o.basicAuth)
	response, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	statusCode := response.StatusCode
	if statusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var sas *serverActiveSessions
		err = json.Unmarshal(body, &sas)
		if err != nil {
			return nil, err
		}

		var sessions []*Session
		for _, ss := range sas.Content {
			session, err := NewSession2(o, ss)
			if err != nil {
				return nil, err
			}

			sessions = append(sessions, session)
		}
		return sessions, nil
	} else {
		return nil, newOpenViduError(statusCode)
	}
}

func computeIfPresent(m map[string]*Session, key string, fn func(sId string, s *Session) *Session) *Session {
	oldValue := m[key]
	if oldValue != nil {
		newValue := fn(key, oldValue)
		if newValue != nil {
			m[key] = newValue
			return newValue
		} else {
			delete(m, key)
			return nil
		}
	} else {
		return nil
	}
}

func computeIfAbsent(m map[string]*Session, key string, fn func(sId string) *Session) *Session {
	v := m[key]
	if v == nil {
		newValue := fn(key)
		if newValue != nil {
			m[key] = newValue
			return newValue
		}
	}
	return v
}

func contains(set []string, item string) bool {
	for _, v := range set {
		if strings.Compare(v, item) == 0 {
			return true
		}
	}
	return false
}
