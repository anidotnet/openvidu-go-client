package openvidu

type RecordingProperties struct {
	Name            string
	OutputMode      OutputMode
	RecordingLayout RecordingLayout
	CustomLayout    string
	Resolution      string
	HasVideo        bool
	HasAudio        bool
}

func (rp *RecordingProperties) Build() *RecordingProperties {
	if rp.OutputMode == COMPOSED {
		if len(rp.RecordingLayout) == 0 {
			rp.RecordingLayout = BEST_FIT
		}

		if len(rp.Resolution) == 0 {
			rp.Resolution = "1920x1080"
		}

		if rp.RecordingLayout == CUSTOM {
			if len(rp.CustomLayout) == 0 {
				rp.CustomLayout = ""
			}
		}
	}

	return rp
}
