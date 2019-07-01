package openvidu

type Recording struct {
	Status              RecordingStatus
	Id                  string
	SessionId           string
	CreatedAt           int64
	Size                int64
	Duration            float64
	Url                 string
	RecordingProperties *RecordingProperties
}

type recordingJson struct {
	Status          RecordingStatus `json:"status"`
	Id              string          `json:"id"`
	SessionId       string          `json:"sessionId"`
	CreatedAt       int64           `json:"createdAt"`
	Size            int64           `json:"size"`
	Duration        float64         `json:"duration"`
	Url             string          `json:"url"`
	HasAudio        bool            `json:"hasAudio"`
	HasVideo        bool            `json:"hasVideo"`
	OutputMode      OutputMode      `json:"outputMode"`
	Name            string          `json:"name"`
	Resolution      string          `json:"resolution"`
	RecordingLayout RecordingLayout `json:"recordingLayout"`
	CustomLayout    string          `json:"customLayout"`
}

func NewRecording(rj *recordingJson) *Recording {
	r := &Recording{
		SessionId: rj.SessionId,
		Id:        rj.Id,
		CreatedAt: rj.CreatedAt,
		Url:       rj.Url,
		Duration:  rj.Duration,
		Size:      rj.Size,
		Status:    rj.Status,
	}

	outputMode := rj.OutputMode
	rp := &RecordingProperties{
		Name:       rj.Name,
		OutputMode: outputMode,
		HasAudio:   rj.HasAudio,
		HasVideo:   rj.HasVideo,
	}

	if outputMode == COMPOSED && rj.HasVideo {
		rp.Resolution = rj.Resolution
		rp.RecordingLayout = rj.RecordingLayout
		if len(rj.CustomLayout) > 0 {
			rp.CustomLayout = rj.CustomLayout
		}
	}

	r.RecordingProperties = rp.Build()
	return r
}

func (r *Recording) Name() string {
	return r.RecordingProperties.Name
}

func (r *Recording) OutputMode() OutputMode {
	return r.RecordingProperties.OutputMode
}

func (r *Recording) RecordingLayout() RecordingLayout {
	return r.RecordingProperties.RecordingLayout
}

func (r *Recording) CustomLayout() string {
	return r.RecordingProperties.CustomLayout
}

func (r *Recording) Resolution() string {
	return r.RecordingProperties.Resolution
}

func (r *Recording) HasAudio() bool {
	return r.RecordingProperties.HasAudio
}

func (r *Recording) HasVideo() bool {
	return r.RecordingProperties.HasVideo
}
