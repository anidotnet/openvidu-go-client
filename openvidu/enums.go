package openvidu

type OpenViduRole string

const (
	// Can subscribe to published Streams of other users
	SUBSCRIBER OpenViduRole = "SUBSCRIBER"

	// SUBSCRIBER permissions + can publish their own Streams (call
	// Session.publish())
	PUBLISHER OpenViduRole = "PUBLISHER"

	// SUBSCRIBER + PUBLISHER permissions + can force the unpublishing or
	// disconnection over a third-party Stream or Connection (call
	// Session.forceUnpublish() and
	// Session.forceDisconnect())
	MODERATOR OpenViduRole = "MODERATOR"
)

type MediaMode string

const (
	RELAYED MediaMode = "RELAYED"
	ROUTED  MediaMode = "ROUTED"
)

type RecordingMode string

const (
	ALWAYS RecordingMode = "ALWAYS"
	MANUAL RecordingMode = "MANUAL"
)

type OutputMode string

const (
	COMPOSED   OutputMode = "COMPOSED"
	INDIVIDUAL OutputMode = "INDIVIDUAL"
)

type RecordingLayout string

const (
	BEST_FIT                RecordingLayout = "BEST_FIT"
	PICTURE_IN_PICTURE      RecordingLayout = "PICTURE_IN_PICTURE"
	VERTICAL_PRESENTATION   RecordingLayout = "VERTICAL_PRESENTATION"
	HORIZONTAL_PRESENTATION RecordingLayout = "HORIZONTAL_PRESENTATION"
	CUSTOM                  RecordingLayout = "CUSTOM"
)

type RecordingStatus string

const (
	STARTING RecordingStatus = "starting"
	STARTED  RecordingStatus = "started"
	STOPPED  RecordingStatus = "stopped"
	READY    RecordingStatus = "ready"
	FAILED   RecordingStatus = "failed"
)
