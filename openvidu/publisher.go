package openvidu

type Publisher struct {
	StreamId        string		`json:"streamId"`
	CreatedAt       int64		`json:"_"`
	HasVideo        bool		`json:"hasVideo"`
	HasAudio        bool		`json:"hasAudio"`
	AudioActive     bool		`json:"audioActive"`
	VideoActive     bool		`json:"videoActive"`
	FrameRate       int32		`json:"frameRate"`
	TypeOfVideo     string		`json:"typeOfVideo"`
	VideoDimensions string		`json:"videoDimensions"`
}
