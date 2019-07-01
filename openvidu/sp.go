package openvidu

type SessionProperties struct {
	MediaMode              MediaMode
	RecordingMode          RecordingMode
	DefaultOutputMode      OutputMode
	DefaultRecordingLayout RecordingLayout
	DefaultCustomLayout    string
	CustomSessionId        string
}
