package openvidu

type KurentoOptions struct {
	VideoMaxRecvBandwidth *int32
	VideoMinRecvBandwidth *int32
	VideoMaxSendBandwidth *int32
	VideoMinSendBandwidth *int32
	AllowedFilters        []string
}
